package main

import (
	"go/ast"
	"go/token"
	"strconv"
)

type Package struct {
	name string
	path string
}

func NewPackage(name, path string) *Package {
	return &Package{
		name: name,
		path: path,
	}
}

func (p *Package) ID() string {
	return p.name
}

func (p *Package) Expr() ast.Expr {
	return ast.NewIdent(p.name)
}

func (p *Package) ImportSpec() *ast.ImportSpec {
	return &ast.ImportSpec{
		Name: ast.NewIdent(p.name),
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote(p.path),
		},
	}
}

type AnonymousInterface struct {
	interfaces []*Interface
}

func NewAnonymousInterface(interfaces []*Interface) *AnonymousInterface {
	return &AnonymousInterface{
		interfaces: interfaces,
	}
}

func (ai *AnonymousInterface) Expr() ([]*Package, ast.Expr) {
	switch len(ai.interfaces) {
	// if the interface has only one interface, return that interface
	case 1:
		pkg, expr := ai.interfaces[0].Expr()

		return []*Package{pkg}, expr
	// if the interface has more than one interface, return the named interface
	default:
		pkgs := make([]*Package, 0, len(ai.interfaces))
		fieldList := make([]*ast.Field, 0, len(ai.interfaces))
		for _, intrfc := range ai.interfaces {
			pkg, expr := intrfc.Expr()
			pkgs = append(pkgs, pkg)
			fieldList = append(fieldList, &ast.Field{
				Type: expr,
			})
		}

		return pkgs, &ast.InterfaceType{
			Methods: &ast.FieldList{
				List: fieldList,
			},
		}
	}
}

type NamedInterface struct {
	name       string
	interfaces []*Interface
	declared   bool
}

func NewNamedInterface(name string, interfaces []*Interface, declared bool) *NamedInterface {
	return &NamedInterface{
		name:       name,
		interfaces: interfaces,
		declared:   declared,
	}
}

func (ni *NamedInterface) Expr() ast.Expr {
	// if the interface is already declared, return the named interface
	if ni.declared {
		return ast.NewIdent(ni.name)
	}

	switch len(ni.interfaces) {
	// if the interface is empty, return an empty interface
	case 0:
		return &ast.InterfaceType{
			Methods: &ast.FieldList{
				List: []*ast.Field{},
			},
		}
	// if the interface has only one interface, return that interface
	case 1:
		_, expr := ni.interfaces[0].Expr()

		return expr
	// if the interface has more than one interface, return the named interface
	default:
		return ast.NewIdent(ni.name)
	}
}

func (n *NamedInterface) Decl() ([]*Package, *ast.GenDecl) {
	// if declared or empty, no declaration is needed
	if n.declared || len(n.interfaces) == 0 {
		return nil, nil
	}

	// if the interface has only one interface, import of the dependency is needed
	if len(n.interfaces) == 1 {
		pkg, _ := n.interfaces[0].Expr()
		if pkg == nil {
			return nil, nil
		}

		return []*Package{pkg}, nil
	}

	pkgs := make([]*Package, 0, len(n.interfaces))
	fieldList := make([]*ast.Field, 0, len(n.interfaces))
	for _, intrfc := range n.interfaces {
		pkg, expr := intrfc.Expr()
		pkgs = append(pkgs, pkg)
		fieldList = append(fieldList, &ast.Field{
			Type: expr,
		})
	}

	return pkgs, &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{&ast.TypeSpec{
			Name: ast.NewIdent(n.name),
			Type: &ast.InterfaceType{
				Methods: &ast.FieldList{
					List: fieldList,
				},
			},
		}},
	}
}

type Interface struct {
	pkg  *Package
	name string
}

func NewInterface(pkg *Package, name string) *Interface {
	return &Interface{
		pkg:  pkg,
		name: name,
	}
}

func (i *Interface) Expr() (*Package, ast.Expr) {
	if i.pkg == nil {
		return nil, ast.NewIdent(i.name)
	}

	return i.pkg, &ast.SelectorExpr{
		X:   i.pkg.Expr(),
		Sel: ast.NewIdent(i.name),
	}
}
