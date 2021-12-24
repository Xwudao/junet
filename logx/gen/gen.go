package main

import (
	"bytes"
	"fmt"
	"go/ast"
	gobuild "go/build"
	"go/format"
	goparser "go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	dir, err := getImportPkg("go.uber.org/zap")
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	log.Printf("dir: %+v", dir)

	pkg, err := parseDir(dir, "zap")
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	funcs, err := walkAst(pkg)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	saveDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	saveDir = filepath.Dir(saveDir)
	f := filepath.Join(saveDir, "log.go")
	ff, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer ff.Close()
	err = writeGoFile(os.Stdout, funcs)
	err = writeGoFile(ff, funcs)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}

func getImportPkg(pkg string) (string, error) {
	p, err := gobuild.Import(pkg, "", gobuild.FindOnly)
	if err != nil {
		return "", err
	}

	return p.Dir, err

}

func parseDir(dir, pkgName string) (*ast.Package, error) {
	pkgMap, err := goparser.ParseDir(
		token.NewFileSet(),
		dir,
		func(info os.FileInfo) bool {
			// skip go-test
			return !strings.Contains(info.Name(), "_test.go")
		},
		goparser.Mode(0), // no comment
	)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	pkg, ok := pkgMap[pkgName]
	if !ok {
		err := errors.New("not found")
		return nil, fmt.Errorf(err.Error())
	}

	return pkg, nil
}

type visitor struct {
	funcs []*ast.FuncDecl
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if n.Recv == nil ||
			!n.Name.IsExported() ||
			len(n.Recv.List) != 1 {
			return nil
		}
		t, ok := n.Recv.List[0].Type.(*ast.StarExpr)
		if !ok {
			return nil
		}

		if t.X.(*ast.Ident).String() != "SugaredLogger" {
			return nil
		}

		log.Printf("func name: %s", n.Name.String())

		v.funcs = append(v.funcs, rewriteFunc(n))

	}
	return v
}

func walkAst(node ast.Node) ([]ast.Decl, error) {
	v := &visitor{}
	ast.Walk(v, node)

	log.Printf("funcs len: %d", len(v.funcs))

	var decls []ast.Decl
	for _, v := range v.funcs {
		decls = append(decls, v)
	}

	return decls, nil
}

func rewriteFunc(fn *ast.FuncDecl) *ast.FuncDecl {
	fn.Recv = nil

	fnName := fn.Name.String()

	var args []string
	for _, field := range fn.Type.Params.List {
		for _, id := range field.Names {
			idStr := id.String()
			_, ok := field.Type.(*ast.Ellipsis)
			if ok {
				// Ellipsis args
				idStr += "..."
			}
			args = append(args, idStr)
		}
	}

	exprStr := fmt.Sprintf(`logger.%s(%s)`, fnName, strings.Join(args, ","))
	expr, err := goparser.ParseExpr(exprStr)
	if err != nil {
		panic(err)
	}

	var body []ast.Stmt
	if fn.Type.Results != nil {
		body = []ast.Stmt{
			&ast.ReturnStmt{
				// Return:
				Results: []ast.Expr{expr},
			},
		}
	} else {
		body = []ast.Stmt{
			&ast.ExprStmt{
				X: expr,
			},
		}
	}

	fn.Body.List = body

	return fn
}

func astToGo(dst *bytes.Buffer, node interface{}) error {
	addNewline := func() {
		err := dst.WriteByte('\n') // add newline
		if err != nil {
			log.Panicln(err)
		}
	}

	addNewline()

	err := format.Node(dst, token.NewFileSet(), node)
	if err != nil {
		return err
	}

	addNewline()

	return nil
}

// Output Go code
func writeGoFile(wr io.Writer, funcs []ast.Decl) error {
	header := `// Code generated by log-gen. DO NOT EDIT.
package wzap

import (
	"go.uber.org/zap"
)

type SugaredLogger = zap.SugaredLogger
type Logger = zap.Logger
`
	buffer := bytes.NewBufferString(header)

	for _, fn := range funcs {
		err := astToGo(buffer, fn)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
	}
	source, err := format.Source(buffer.Bytes())
	if err != nil {
		return err
	}

	_, err = wr.Write(source)
	return err
}