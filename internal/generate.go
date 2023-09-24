package iwrapper

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"os"
	"strconv"
)

type GenerateConfig struct {
	FuncName           string
	RequireInterface   *AnonymousInterface
	WrappedInterface   *NamedInterface
	OptionalInterfaces []*Interface
}

func Generate(w io.Writer, pkgName string, confs []*GenerateConfig) error {
	fset := token.NewFileSet()

	decls := []ast.Decl{}
	importPkgMap := map[string]*Package{}

	for _, conf := range confs {
		requireDepPkgs, valueType := conf.RequireInterface.Expr()
		for _, pkg := range requireDepPkgs {
			importPkgMap[pkg.ID()] = pkg
		}

		wrappedDepPkgs, wrappedDecl := conf.WrappedInterface.Decl()
		if wrappedDecl != nil {
			decls = append(decls, wrappedDecl)
		}
		for _, pkg := range wrappedDepPkgs {
			importPkgMap[pkg.ID()] = pkg
		}

		var (
			valueIdent      = ast.NewIdent("v")
			wrapFuncIdent   = ast.NewIdent("wrapper")
			wrappedTypeExpr = conf.WrappedInterface.Expr()
		)

		bodyDepPkgs, bodyStmts := getBody(valueIdent, wrapFuncIdent, valueType, conf.OptionalInterfaces)
		for _, pkg := range bodyDepPkgs {
			importPkgMap[pkg.ID()] = pkg
		}

		decls = append(decls, &ast.FuncDecl{
			Name: ast.NewIdent(conf.FuncName),
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{{
						Names: []*ast.Ident{
							valueIdent,
						},
						Type: valueType,
					}, {
						Names: []*ast.Ident{
							wrapFuncIdent,
						},
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{{
									Type: valueType,
								}},
							},
							Results: &ast.FieldList{
								List: []*ast.Field{{
									Type: wrappedTypeExpr,
								}},
							},
						},
					}},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{{
						Type: valueType,
					}},
				},
			},
			Body: &ast.BlockStmt{
				List: bodyStmts,
			},
		})
	}

	importSpecs := make([]ast.Spec, 0, len(importPkgMap))
	for _, pkg := range importPkgMap {
		importSpecs = append(importSpecs, pkg.ImportSpec())
	}
	decls = append(decls, &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: importSpecs,
	})

	_, err := io.WriteString(os.Stdout, "// Code generated by internal/tools/interface-wrapper.go; DO NOT EDIT.\n")
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	f := ast.File{
		Name:  ast.NewIdent(pkgName),
		Decls: decls,
	}

	format.Node(os.Stdout, fset, &f)

	return nil
}

func getBody(valueIdent, wrapFuncIdent *ast.Ident, valueType ast.Expr, optionalInterfaces []*Interface) ([]*Package, []ast.Stmt) {
	if len(optionalInterfaces) == 0 {
		return nil, []ast.Stmt{&ast.ReturnStmt{
			Results: []ast.Expr{valueIdent},
		}}
	}

	wrappedValueIdent := ast.NewIdent("wrapped")
	indexIdent := ast.NewIdent("i")
	bodyStmts := []ast.Stmt{
		&ast.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []ast.Expr{wrappedValueIdent},
			Rhs: []ast.Expr{&ast.CallExpr{
				Fun:  wrapFuncIdent,
				Args: []ast.Expr{valueIdent},
			}},
		},
		&ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{&ast.ValueSpec{
					Names: []*ast.Ident{indexIdent},
					Type:  ast.NewIdent("uint64"),
				}},
			},
		},
	}

	depPkgs := make([]*Package, 0, len(optionalInterfaces))
	constSpecs := make([]ast.Spec, 0, len(optionalInterfaces))
	checkStmts := make([]ast.Stmt, 0, len(optionalInterfaces))
	optionalInterfaceExprs := make([]ast.Expr, 0, len(optionalInterfaces))
	for i, intrfc := range optionalInterfaces {
		ident := ast.NewIdent(fmt.Sprintf("i%d", i))
		var values []ast.Expr
		if i == 0 {
			values = []ast.Expr{&ast.BinaryExpr{
				X: &ast.BasicLit{
					Kind:  token.INT,
					Value: "1",
				},
				Op: token.SHL,
				Y:  ast.NewIdent("iota"),
			}}
		}

		constSpecs = append(constSpecs, &ast.ValueSpec{
			Names:  []*ast.Ident{ident},
			Values: values,
		})

		pkg, expr := intrfc.Expr()
		depPkgs = append(depPkgs, pkg)
		optionalInterfaceExprs = append(optionalInterfaceExprs, expr)

		okIdent := ast.NewIdent("ok")
		checkStmts = append(checkStmts, &ast.IfStmt{
			Init: &ast.AssignStmt{
				Lhs: []ast.Expr{ast.NewIdent("_"), okIdent},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{&ast.TypeAssertExpr{
					X:    valueIdent,
					Type: expr,
				}},
			},
			Cond: okIdent,
			Body: &ast.BlockStmt{
				List: []ast.Stmt{&ast.AssignStmt{
					Lhs: []ast.Expr{indexIdent},
					Tok: token.OR_ASSIGN,
					Rhs: []ast.Expr{ident},
				}},
			},
		})
	}

	bodyStmts = append(bodyStmts, &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok:   token.CONST,
			Specs: constSpecs,
		},
	})
	bodyStmts = append(bodyStmts, checkStmts...)

	caseClauseStmts := make([]ast.Stmt, 0, 1<<len(optionalInterfaces))
	var i uint64 = 0
	for ; i < 1<<len(optionalInterfaces); i++ {
		typeFields := []*ast.Field{{
			Type: valueType,
		}}
		elementsExprs := []ast.Expr{wrappedValueIdent}

		tmpI := i
		for j := 0; j < len(optionalInterfaces); j++ {
			if tmpI&1 != 0 {
				typeFields = append(typeFields, &ast.Field{
					Type: optionalInterfaceExprs[j],
				})
				elementsExprs = append(elementsExprs, wrappedValueIdent)
			}
		}

		caseClauseStmts = append(caseClauseStmts, &ast.CaseClause{
			List: []ast.Expr{&ast.BasicLit{
				Kind:  token.INT,
				Value: "0b" + strconv.FormatUint(i, 2),
			}},
			Body: []ast.Stmt{&ast.ReturnStmt{
				Results: []ast.Expr{&ast.CompositeLit{
					Type: &ast.StructType{
						Fields: &ast.FieldList{
							List: typeFields,
						},
					},
					Elts: elementsExprs,
				}},
			}},
		})
	}

	bodyStmts = append(bodyStmts, &ast.SwitchStmt{
		Tag: indexIdent,
		Body: &ast.BlockStmt{
			List: caseClauseStmts,
		},
	}, &ast.ReturnStmt{
		Results: []ast.Expr{valueIdent},
	})

	return depPkgs, bodyStmts
}
