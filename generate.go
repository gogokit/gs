package gs

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

const (
	LockFiledName = "lock"
	RecvObjName   = "obj"
	ParamObjName  = "param"
)

func AddGetAndSet(src string) (newSrc string, err error) {
	file, err := parser.ParseFile(token.NewFileSet(), "", src, 0)
	if err != nil {
		return "", fmt.Errorf("parse src to ast fail, err=%v", err)
	}

	structSpecs := getAllStruct(file)
	for name, ts := range structSpecs {
		if err = insertLockFiled(ts); err != nil {
			delete(structSpecs, name)
			_, _ = fmt.Fprintf(os.Stderr, "insert lock field fail, err=%v\n", err)
		}
	}

	if len(structSpecs) == 0 {
		return src, nil
	}

	insertSyncImport(file) // 为结构体增加sync导入

	objMethods := getAllObjMethod(file)
	for _, ss := range structSpecs {
		for _, field := range ss.Type.(*ast.StructType).Fields.List {
			insertGetAndSetMethod(file, ss, field, objMethods)
		}
	}

	buf := bytes.NewBufferString("")
	if err = format.Node(buf, token.NewFileSet(), file); err != nil {
		return "", fmt.Errorf("exec format.Node fail, err=%v", err)
	}
	newSrcBytes, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("exec format.Source fail, err=%v", err)
	}
	return string(newSrcBytes), nil
}

// 返回所有全局定义的有名结构体, 返回结果为 结构体类型名->结构体对应的ast.TypeSpec对象
func getAllStruct(file *ast.File) (ret map[string]*ast.TypeSpec) {
	ret = make(map[string]*ast.TypeSpec)
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd == nil {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok || ts == nil {
				continue
			}
			if ty, ok := ts.Type.(*ast.StructType); !ok || ty == nil {
				continue
			}
			name := ts.Name
			if name == nil {
				continue
			}
			ret[name.Name] = ts
		}
	}
	return ret
}

func insertLockFiled(ts *ast.TypeSpec) error {
	st := ts.Type.(*ast.StructType)
	for _, field := range st.Fields.List {
		for _, name := range field.Names {
			if name.Name == LockFiledName {
				return fmt.Errorf("struct [%s] already had filed [%s]", ts.Name.Name, LockFiledName)
			}
		}
	}
	field := &ast.Field{
		Type: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "sync",
			},
			Sel: &ast.Ident{
				Name: "RWMutex",
			},
		},
	}
	field.Names = append(field.Names, &ast.Ident{
		Name: LockFiledName,
		Obj: &ast.Object{
			Kind: ast.Var,
			Name: LockFiledName,
			Decl: field,
		},
	})
	st.Fields.List = append(st.Fields.List, field)
	return nil
}

func insertSyncImport(file *ast.File) {
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd == nil {
			continue
		}
		for _, spec := range gd.Specs {
			is, ok := spec.(*ast.ImportSpec)
			if !ok || is == nil {
				continue
			}
			path := is.Path
			if path == nil {
				continue
			}
			if path.Value == `"sync"` {
				return
			}
		}
	}

	importSpec := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: `"sync"`,
		},
	}
	file.Decls = append([]ast.Decl{&ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: []ast.Spec{importSpec},
	}}, file.Decls...)
	file.Imports = append(file.Imports, importSpec)
}

// 返回对象已有的方法名称, 对象类型名称->method名称列表
func getAllObjMethod(file *ast.File) (ret map[string][]string) {
	ret = make(map[string][]string)
	for _, decl := range file.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd == nil || fd.Recv == nil || fd.Recv.NumFields() != 1 || fd.Name == nil || fd.Name.Name == "" {
			continue
		}

		typeExpr := fd.Recv.List[0].Type
		for {
			if se, ok := typeExpr.(*ast.StarExpr); ok && se != nil {
				typeExpr = se.X
				continue
			}
			break
		}

		ident, ok := typeExpr.(*ast.Ident)
		if !ok || ident == nil {
			continue
		}
		typeName := ident.Name
		ret[typeName] = append(ret[typeName], fd.Name.Name)
	}
	return ret
}

func inSlice(val string, slice []string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func insertGetAndSetMethod(f *ast.File, ss *ast.TypeSpec, field *ast.Field, objMethods map[string][]string) {
	// 检查结构体是否已经有名为${LockFiledName}的字段
	for _, name := range field.Names {
		if name.Name == LockFiledName || name.Name == "" {
			continue
		}

		recvField := &ast.Field{
			Type: &ast.StarExpr{
				X: &ast.Ident{
					Name: ss.Name.Name,
					Obj:  ss.Name.Obj,
				},
			},
		}
		recvField.Names = []*ast.Ident{
			{
				Name: RecvObjName,
				Obj: &ast.Object{
					Kind: ast.Var,
					Name: RecvObjName,
					Decl: recvField,
				},
			},
		}
		recvObj := recvField.Names[0].Obj

		if getMethodName := "Get" + strings.ToUpper(name.Name[:1]) + name.Name[1:]; !inSlice(getMethodName, objMethods[ss.Name.Name]) {
			// 为结构体增加Get方法
			f.Decls = append(f.Decls, &ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{recvField},
				},
				Name: &ast.Ident{
					Name: getMethodName,
				},
				Type: &ast.FuncType{
					Results: &ast.FieldList{
						List: []*ast.Field{
							{Type: field.Type},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: RecvObjName,
											Obj:  recvObj,
										},
										Sel: &ast.Ident{
											Name: LockFiledName,
										},
									},
									Sel: &ast.Ident{
										Name: "RLock",
									},
								},
							},
						},
						&ast.DeferStmt{
							Call: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: RecvObjName,
											Obj:  recvObj,
										},
										Sel: &ast.Ident{
											Name: LockFiledName,
										},
									},
									Sel: &ast.Ident{
										Name: "RUnlock",
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: RecvObjName,
										Obj:  recvObj,
									},
									Sel: &ast.Ident{
										Name: name.Name,
									},
								},
							},
						},
					},
				},
			})
		}

		if setMethodName := "Set" + strings.ToUpper(name.Name[:1]) + name.Name[1:]; !inSlice(setMethodName, objMethods[ss.Name.Name]) {
			// 为结构体增加Set方法
			paramField := &ast.Field{
				Type: field.Type,
			}
			paramField.Names = []*ast.Ident{
				{
					Name: ParamObjName,
					Obj: &ast.Object{
						Kind: ast.Var,
						Decl: paramField,
					},
				},
			}
			f.Decls = append(f.Decls, &ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{recvField},
				},
				Name: &ast.Ident{
					Name: setMethodName,
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{paramField},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: recvField.Type,
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.Ident{
									Name: RecvObjName,
									Obj:  recvObj,
								},
								Op: token.EQL,
								Y: &ast.Ident{
									Name: "nil",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.Ident{
												Name: "nil",
											},
										},
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: RecvObjName,
											Obj:  recvObj,
										},
										Sel: &ast.Ident{
											Name: LockFiledName,
										},
									},
									Sel: &ast.Ident{
										Name: "Lock",
									},
								},
							},
						},
						&ast.DeferStmt{
							Call: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: RecvObjName,
											Obj:  recvObj,
										},
										Sel: &ast.Ident{
											Name: LockFiledName,
										},
									},
									Sel: &ast.Ident{
										Name: "Unlock",
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   recvField.Names[0],
									Sel: name,
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{paramField.Names[0]},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{recvField.Names[0]},
						},
					},
				},
			})
		}
	}
}
