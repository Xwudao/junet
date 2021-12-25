package cmd

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/errgo.v2/fmt/errors"
)

const (
	gitUrl = "https://gitee.com/xwudao/junet-template.git"
)

type initProject struct {
	modPath     string
	rootPath    string
	projectName string

	originModName string
	newModName    string
}

func (p *initProject) cmd() *cobra.Command {
	var c = &cobra.Command{
		Use:   "init",
		Short: "init junet project",
		Run: func(cmd *cobra.Command, args []string) {
			p.init(args)
			p.cloneProject()
			p.rewriteMod()

			Info("finished, happy hacking!")
		},
	}

	c.Flags().StringVarP(&p.newModName, "mode", "m", "", "the mod name, eg: github.com/Xwudao/junet-example")
	return c
}

func (p *initProject) init(args []string) {
	//Info(strings.Join(args, ","))
	if len(args) == 0 {
		CheckErrWithStatus(errors.Newf("please input project name"))
		return
	}
	p.projectName = args[0]

	dir, _ := os.Getwd()
	p.rootPath = filepath.Join(dir, p.projectName)
	p.modPath = filepath.Join(p.rootPath, "go.mod")

	_, err := os.Stat(p.rootPath)
	if err == nil {
		CheckErrWithStatus(errors.Newf("maybe %s path existed, please rename or remove it.", p.rootPath))
	}
}
func (p *initProject) cloneProject() {
	Info("cloning project....")
	Info(p.projectName)
	Info(gitUrl)
	cmd := exec.Command("git", "clone", gitUrl, p.projectName)
	err := cmd.Run()
	CheckErrWithStatus(err)
	Info("cloned project....")
}
func (p *initProject) rewriteMod() {
	if p.newModName == "" {
		p.newModName = p.projectName
	}
	var err error
	p.originModName, err = p.getOriginName()
	CheckErrWithStatus(err)
	cwd, _ := os.Getwd()
	files := LoadFiles(cwd, func(filename string) bool {
		return path.Ext(filename) == ".go" && !strings.Contains(filename, "/vendor/")
	})
	Info("changing mod name...")
	for _, f := range files {
		node, fset, err := p.parse(f)
		if err != nil {
			Error(err)
			continue
		}
		err = p.write(f, node, fset)
		if err != nil {
			Error(err)
			continue
		}
	}
	err = p.setModName()
	CheckErrWithStatus(err)
	Info("changed mod name")
}
func (p *initProject) parse(filename string) (*ast.File, *token.FileSet, error) {

	fileSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fileSet, filename, nil, parser.ParseComments)

	if err != nil {
		return nil, nil, err
	}

	fset := fileSet

	for _, importSpec := range astFile.Imports {
		originPath := importSpec.Path.Value
		importSpec.Path.Value = strings.Replace(originPath, p.originModName, p.newModName, 1)
	}

	return astFile, fset, nil
}
func (p *initProject) getOriginName() (name string, err error) {
	_, err = os.Stat(p.modPath)
	if err != nil {
		return
	}

	cnt, err := ioutil.ReadFile(p.modPath)
	if err != nil {
		return
	}

	compile := regexp.MustCompile("(?m)module\\s([^\\s]+)")
	matches := compile.FindStringSubmatch(string(cnt))
	if len(matches) >= 2 {
		return matches[1], nil
	}
	return
}
func (p *initProject) write(filename string, node *ast.File, fset *token.FileSet) error {

	var buf bytes.Buffer

	err := format.Node(&buf, fset, node)
	if err != nil {
		return err
	}

	if filename == "" {
		return errors.Newf("no file name")
	}

	err = ioutil.WriteFile(filename, buf.Bytes(), os.ModePerm)
	if err != nil {
		return errors.Newf("write file err: %s", err.Error())
	}

	return nil
}

func (p *initProject) setModName() (err error) {
	_, err = os.Stat(p.modPath)
	if err != nil {
		return
	}

	cnt, err := ioutil.ReadFile(p.modPath)
	if err != nil {
		return
	}
	nCnt := strings.Replace(string(cnt), p.originModName, p.newModName, 1)
	err = ioutil.WriteFile(p.modPath, []byte(nCnt), os.ModePerm)
	if err != nil {
		return
	}
	return nil
}
func init() {
	rootCmd.AddCommand((&initProject{}).cmd())
}
