package gen

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/Xwudao/junet/cmd/junet/utils"
	"github.com/Xwudao/junet/internal/parser"
	"github.com/Xwudao/junet/internal/queryset/generator"
)

type DB struct {
	PackageName string

	ModelName string

	RootPath string
	timeout  time.Duration

	CurrentFileName string
	FilenameSuffix  string
	FilePath        string
}

func (r *DB) Cmd() *cobra.Command {
	var c = &cobra.Command{
		Use:   "db",
		Short: "generate gorm model",
		Run: func(cmd *cobra.Command, args []string) {
			r.init()
			r.generate()
		},
	}
	//c.Flags().StringVarP(&r.ModelName, "model", "m", "", "the model entity's name")
	c.Flags().DurationVarP(&r.timeout, "timeout", "t", time.Minute, "timeout of generator")
	//_ = c.MarkFlagRequired("model")
	return c
}

func (r *DB) init() {
	r.FilenameSuffix = "_gen.go"
	r.RootPath = utils.CurrentDir()
	r.CurrentFileName = os.Getenv("GOFILE")
	r.PackageName = os.Getenv("GOPACKAGE")
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
