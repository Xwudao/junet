package gen

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"golang.org/x/tools/go/ast/astutil"
	"gopkg.in/errgo.v2/fmt/errors"

	"github.com/Xwudao/junet/cmd/junet/utils"
)

//go:embed route.tpl
var tplCnt string

//go:embed service.tpl
var serviceTplCnt string

type Route struct {
	PackageName string

	RootPath  string
	ModName   string
	RouteName string

	FilePath string

	WithService     bool
	ServiceFilePath string

	FilenameSuffix  string
	RouteNameSuffix string
}

func (r *Route) Cmd() *cobra.Command {
	var c = &cobra.Command{
		Use:   "route",
		Short: "gen route",
		Run: func(cmd *cobra.Command, args []string) {
			r.init()
			r.generate()
			r.updateRoot()
			r.generateService()
		},
	}
	c.Flags().StringVarP(&r.RouteName, "name", "n", "", "the route of name")
	c.Flags().BoolVarP(&r.WithService, "service", "s", false, "generate service file")
	_ = c.MarkFlagRequired("name")
	return c
}
func (r *Route) init() {
	r.RootPath = utils.CurrentDir()
	r.FilenameSuffix = "_routes.go"
	r.RouteNameSuffix = "Routes"
	r.PackageName = os.Getenv("GOPACKAGE")
	r.ModName = utils.GetModName()
	if r.PackageName == "" {
		utils.CheckErrWithStatus(errors.Newf("please run with //go:generate"))
		return
	}
}
func (r *Route) generate() {
	r.FilePath = filepath.Join(r.RootPath, fmt.Sprintf("%s%s", strcase.ToSnake(r.RouteName), r.FilenameSuffix))
	parse, err := template.New("route").Parse(tplCnt)
	utils.CheckErrWithStatus(err)

	buffer := bytes.NewBuffer([]byte{})
	err = parse.Execute(buffer, r)
	utils.CheckErrWithStatus(err)

	source, err := format.Source(buffer.Bytes())
	utils.CheckErrWithStatus(err)

	err = utils.SaveToFile(r.FilePath, source, false)
	utils.CheckErrWithStatus(err)
}
func (r *Route) updateRoot() {
	utils.Info("updating root.go")
	rootFilePath := filepath.Join(filepath.Dir(r.RootPath), "root.go")
	exist := utils.CheckExist(rootFilePath)
	if !exist {
		utils.CheckErrWithStatus(errors.Newf("can't find root.go file [%s]", rootFilePath))
	}

	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, rootFilePath, nil, 0)
	if err != nil {
		utils.CheckErrWithStatus(err)
	}

	var importedName []string
	for _, decl := range f.Decls {
		if x, ok := decl.(*ast.FuncDecl); ok && x.Name.String() == "Setup" {
			x.Body.List = append(x.Body.List, r.genUse2())
			x.Body.List = append(x.Body.List, r.genUse1())
		}

		if x, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range x.Specs {
				if y, ok := spec.(*ast.ImportSpec); ok {
					if y.Name != nil && y.Name.Name != "" {
						importedName = append(importedName, y.Name.Name)
					}
				}
			}
		}
	}
	if !utils.InStrArr(importedName, r.PackageName) {
		utils.Info("no package")
		added := astutil.AddNamedImport(fset, f,
			r.PackageName,
			fmt.Sprintf(`%s/pkg/routes/%s`, r.ModName, r.PackageName),
		)
		if !added {
			utils.Error(errors.Newf("can't add import sec for root.go file"))
		}

	}

	buf := &bytes.Buffer{}
	err = format.Node(buf, fset, f)
	utils.CheckErrWithStatus(err)
	err = utils.SaveToFile(rootFilePath, buf.Bytes(), true)
	utils.CheckErrWithStatus(err)

	utils.Info("updated root.go")
}

func (r *Route) generateService() {
	if r.WithService {
		utils.Info("generate with service file")
		servicePath := filepath.Join(r.FilePath, "../../../services")
		err := os.MkdirAll(servicePath, os.ModePerm)
		utils.CheckErrWithStatus(err)
		serviceFilePath := filepath.Join(servicePath, strcase.ToLowerCamel(r.RouteName)+"_service.go")

		parse, err := template.New("service").Parse(serviceTplCnt)
		utils.CheckErrWithStatus(err)
		var buf bytes.Buffer
		err = parse.Execute(&buf, r)
		utils.CheckErrWithStatus(err)

		err = utils.SaveToFile(serviceFilePath, buf.Bytes(), false)
		utils.CheckErrWithStatus(err)
		utils.Info("generate service file in: " + serviceFilePath)
	}
}

func (r *Route) genUse1() *ast.ExprStmt {
	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: strcase.ToLowerCamel(r.PackageName + r.RouteName + r.RouteNameSuffix)},
				Sel: &ast.Ident{Name: "SetUpRoutes"},
			},
		},
	}
}

func (r *Route) genUse2() *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.Ident{
				Name: strcase.ToLowerCamel(r.PackageName + r.RouteName + r.RouteNameSuffix),
				Obj: &ast.Object{
					Kind: ast.Var,
					Name: strcase.ToLowerCamel(r.RouteName + r.RouteNameSuffix),
				},
			},
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CompositeLit{
				Type: &ast.SelectorExpr{
					X:   &ast.Ident{Name: r.PackageName},
					Sel: &ast.Ident{Name: r.RouteName + r.RouteNameSuffix},
				},
				Elts: []ast.Expr{
					&ast.KeyValueExpr{
						Key:   &ast.Ident{Name: "Engine"},
						Value: &ast.Ident{Name: "a"},
					},
				},
			},
		},
	}
}

//tpl fun
func (r Route) ToSnake(p string) string {
	return strcase.ToSnake(p)
}
