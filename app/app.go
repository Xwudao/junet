package app

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/Xwudao/junet"
)

type H gin.H
type InfoOpt func(*Info)

func SetMode(s string) InfoOpt {
	return func(info *Info) {
		info.Mode = s
	}
}
func SetShort(s string) InfoOpt {
	return func(info *Info) {
		info.Short = s
	}
}
func SetLong(l string) InfoOpt {
	return func(info *Info) {
		info.Long = l
	}
}
func SetUse(u string) InfoOpt {
	return func(info *Info) {
		info.Use = u
	}
}

type App struct {
	*Engine
	rootCmd *cobra.Command

	info *Info
}

func init() {
}
func NewApp(opts ...InfoOpt) *App {
	info := &Info{Mode: junet.Mode}
	for _, opt := range opts {
		opt(info)
	}
	cmd := &cobra.Command{
		Use:   info.Use,
		Short: info.Short,
		Long:  info.Long,
	}
	app := &App{
		rootCmd: cmd,
		info:    info,
	}
	return app
}

func (a *App) AddCommand(cmd ...*cobra.Command) {
	a.rootCmd.AddCommand(cmd...)
}

func (a *App) Mdw(f func(app *App)) *App {
	if a.Engine == nil {
		a.Engine = NewEngine()
	}
	f(a)
	return a
}
func (a *App) Mount(f func(app *App)) *App {
	gin.SetMode(a.info.Mode)
	if a.Engine == nil {
		a.Engine = NewEngine()
	}
	f(a)
	return a
}

func (a *App) Execute(run func(cmd *cobra.Command, args []string)) {
	a.rootCmd.Run = run
	err := a.rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
