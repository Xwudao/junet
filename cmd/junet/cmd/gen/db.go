package gen

import (
	"bytes"
	"context"
	_ "embed"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"github.com/vetcher/go-astra"
	"gopkg.in/errgo.v2/fmt/errors"

	"github.com/Xwudao/junet/cmd/junet/utils"
	"github.com/Xwudao/junet/internal/parser"
	"github.com/Xwudao/junet/internal/queryset/generator"
)

//go:embed dto.tpl
var dtoTplCnt string

type DB struct {
	PackageName string

	ModelName string

	RunFilePath string //which file run cmd?

	RootPath string
	timeout  time.Duration

	WithDto bool

	CurrentFileName string
	FilenameSuffix  string
	FilePath        string
}
type DtoData struct {
	Name   string
	Fields []DtoFields
}
type DtoFields struct {
	Name     string
	Type     string
	JsonName string
}

func (r *DB) Cmd() *cobra.Command {
	var c = &cobra.Command{
		Use:   "db",
		Short: "generate gorm model",
		Run: func(cmd *cobra.Command, args []string) {
			r.init()
			r.generate()
			r.generateDto()
		},
	}
	//c.Flags().StringVarP(&r.ModelName, "model", "m", "", "the model entity's name")
	c.Flags().DurationVarP(&r.timeout, "timeout", "t", time.Minute, "timeout of generator")
	c.Flags().BoolVarP(&r.WithDto, "dto", "d", false, "generate with dto struct")
	//_ = c.MarkFlagRequired("model")
	return c
}

func (r *DB) init() {
	r.FilenameSuffix = "_gen.go"
	r.RootPath = utils.CurrentDir()
	r.CurrentFileName = os.Getenv("GOFILE")
	r.PackageName = os.Getenv("GOPACKAGE")

	if r.PackageName == "" || r.CurrentFileName == "" {
		utils.CheckErrWithStatus(errors.Newf("please run with //go:generate"))
		return
	}
	r.RunFilePath = filepath.Join(r.RootPath, os.Getenv("GOFILE"))
}

func (r *DB) generate() {

	g := generator.Generator{
		StructsParser: &parser.Structs{},
	}

	ctx, finish := context.WithTimeout(context.Background(), r.timeout)
	defer finish()

	inFile := filepath.Join(r.RootPath, r.CurrentFileName)
	outFile := filepath.Join(r.RootPath, utils.RemoveExt(r.CurrentFileName)+r.FilenameSuffix)

	utils.Info(inFile)
	utils.Info(outFile)
	if err := g.Generate(ctx, inFile, outFile); err != nil {
		log.Fatalf("can't generate query sets: %s", err)
	}
}

func (r *DB) generateDto() {
	if !r.WithDto {
		return
	}
	//fset := token.NewFileSet()
	//
	//f, err := parser.ParseFile(fset, r.RunFilePath, nil, parser.ParseComments|parser.AllErrors)
	//if err != nil {
	//	utils.CheckErrWithStatus(err)
	//}
	file, err := astra.ParseFile(r.RunFilePath)
	if err != nil {
		utils.CheckErrWithStatus(err)
	}
	var allDto strings.Builder
	for _, structure := range file.Structures {
		if strings.Contains(strings.Join(structure.Docs, "|"), "//gen:qs") {
			parse, err := template.New("dto").Parse(dtoTplCnt)
			utils.CheckErrWithStatus(err)

			var fields []DtoFields
			for _, field := range structure.Fields {
				fields = append(fields, DtoFields{
					Name:     field.Name,
					Type:     field.Type.String(),
					JsonName: strcase.ToLowerCamel(field.Name),
				})
			}
			var data = &DtoData{
				Name:   structure.Name,
				Fields: fields,
			}

			var buf bytes.Buffer
			err = parse.Execute(&buf, data)
			utils.CheckErrWithStatus(err)
			allDto.WriteString(buf.String())
			allDto.WriteRune('\n')
		}
	}
	if allDto.String() != "" {
		readFile, err := ioutil.ReadFile(r.RunFilePath)
		utils.CheckErrWithStatus(err)

		buffer := bytes.NewBuffer(readFile)

		buffer.WriteString(allDto.String())
		source, err := format.Source(buffer.Bytes())
		utils.CheckErrWithStatus(err)
		err = ioutil.WriteFile(r.RunFilePath, source, os.ModePerm)
		utils.CheckErrWithStatus(err)

		utils.Info("generate dto struct")
	}
}
