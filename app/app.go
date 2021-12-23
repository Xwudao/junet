package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/Xwudao/junet"
	"github.com/Xwudao/junet/shutdown"
)

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
	info := &Info{Mode: junet.Debug}
	for _, opt := range opts {
		opt(info)
	}
	gin.SetMode(info.Mode)
	cmd := &cobra.Command{
		Use:   info.Use,
		Short: info.Short,
		Long:  info.Long,
	}
	app := &App{
		Engine:  NewEngine(),
		rootCmd: cmd,
		info:    info,
	}
	return app
}

func (a App) AddCommand(cmd ...*cobra.Command) {
	a.rootCmd.AddCommand(cmd...)
}
func (a App) Start(port int) error {
	a.rootCmd.Run = func(cmd *cobra.Command, args []string) {
		go func() {
			err := a.Run(fmt.Sprintf(":%d", port))
			if err != nil {
				panic(err)
			}
		}()

		shutdown.Wait()
	}

	return a.rootCmd.Execute()
}
