package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

const (
	toolPrefix             = "//iwrapper:"
	targetDirectivePrefix  = toolPrefix + "target"
	requireDirectivePrefix = toolPrefix + "require"
)

type ParseResult struct {
	FuncName                               string
	RequiredInterfaces, OptionalInterfaces []*Interface
}

var (
	ErrNoPkgName = errors.New("no package name")
)

func ParseTarget(r io.Reader) (string, []*ParseResult, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", r, parser.ParseComments)
	if err != nil {
		return "", nil, err
	}

	if f.Name == nil {
		return "", nil, ErrNoPkgName
	}

	pkgName := f.Name.Name

	pkgMap := make(map[string]*Package, len(f.Imports))
	for _, impt := range f.Imports {
		if impt.Path == nil {
			continue
		}
		path, err := strconv.Unquote(impt.Path.Value)
		if err != nil {
			return "", nil, fmt.Errorf("invalid import path(%s): %w", impt.Path.Value, err)
		}

		if impt.Name == nil {
			pkgs, err := packages.Load(&packages.Config{
				Mode: packages.NeedName,
			}, path)
			if err != nil {
				log.Printf("failed to load package(%s): %v", path, err)
				continue
			}

			if len(pkgs) < 1 {
				log.Printf("failed to load package(%s): no packages", path)
				continue
			}

			for _, pkg := range pkgs {
				name := pkg.Name
				pkgMap[name] = NewPackage(name, path)
			}
		} else {
			name := impt.Name.Name
			pkgMap[name] = NewPackage(name, path)
		}
	}

	var results []*ParseResult
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl == nil || genDecl.Tok != token.TYPE || genDecl.Doc == nil {
			continue
		}

		if len(genDecl.Specs) == 1 {
			typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
			if !ok || typeSpec == nil || typeSpec.Name == nil {
				continue
			}

			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok || interfaceType == nil || interfaceType.Methods == nil {
				continue
			}

			var docs []*ast.Comment
			if genDecl.Doc != nil {
				docs = append(docs, genDecl.Doc.List...)
			}
			if typeSpec.Doc != nil {
				docs = append(docs, typeSpec.Doc.List...)
			}

			var funcName string
			for _, comment := range docs {
				if !strings.HasPrefix(comment.Text, targetDirectivePrefix) {
					continue
				}

				annotationTagText := strings.TrimPrefix(comment.Text, targetDirectivePrefix)
				annotationTag := reflect.StructTag(annotationTagText)

				funcName, ok = annotationTag.Lookup("func")
				if !ok {
					funcName = fmt.Sprintf("%sWrapper", typeSpec.Name.Name)
				}
				break
			}

			requireInterfaces := []*Interface{}
			optionalInterfaces := []*Interface{}
			for _, field := range interfaceType.Methods.List {
				if field.Type == nil {
					strField, err := astToString(fset, field)
					if err != nil {
						return "", nil, errors.New("invalid interface field: no field type")
					}
					return "", nil, fmt.Errorf("invalid interface field(%s): no field type", strField)
				}

				var interfaceValue *Interface
				switch expr := field.Type.(type) {
				case *ast.Ident:
					interfaceValue = NewInterface(nil, expr.Name)
				case *ast.SelectorExpr:
					pkgIdent, ok := expr.X.(*ast.Ident)
					if !ok {
						strField, err := astToString(fset, field)
						if err != nil {
							return "", nil, errors.New("invalid interface field: invalid field type")
						}
						return "", nil, fmt.Errorf("invalid interface field(%s): invalid field type", strField)
					}

					pkg, ok := pkgMap[pkgIdent.Name]
					if !ok {
						strField, err := astToString(fset, field)
						if err != nil {
							return "", nil, errors.New("invalid interface field: invalid field type")
						}
						return "", nil, fmt.Errorf("invalid interface field(%s): invalid field type", strField)
					}

					if expr.Sel == nil {
						strField, err := astToString(fset, field)
						if err != nil {
							return "", nil, errors.New("invalid interface field: invalid field type")
						}
						return "", nil, fmt.Errorf("invalid interface field(%s): invalid field type", strField)
					}

					interfaceValue = NewInterface(pkg, expr.Sel.Name)
				default:
					strField, err := astToString(fset, field)
					if err != nil {
						return "", nil, errors.New("invalid interface field: invalid field type")
					}
					return "", nil, fmt.Errorf("invalid interface field(%s): invalid field type", strField)
				}

				required := false
				if field.Doc != nil {
					for _, docs := range field.Doc.List {
						if strings.HasPrefix(docs.Text, requireDirectivePrefix) {
							required = true
							break
						}
					}
				}

				if required {
					requireInterfaces = append(requireInterfaces, interfaceValue)
				} else {
					optionalInterfaces = append(optionalInterfaces, interfaceValue)
				}
			}

			results = append(results, &ParseResult{
				FuncName:           funcName,
				RequiredInterfaces: requireInterfaces,
				OptionalInterfaces: optionalInterfaces,
			})
		}
	}

	return pkgName, results, nil
}

func astToString(fset *token.FileSet, node ast.Node) (string, error) {
	var sb strings.Builder
	err := format.Node(&sb, fset, node)
	if err != nil {
		return "", fmt.Errorf("failed to format node: %w", err)
	}

	return sb.String(), nil
}
