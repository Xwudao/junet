package goose

import (
	"fmt"
	"github.com/Xwudao/junet/confx"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type Down struct {
	Type string // mysql, postgres
	Dir  string // migration direction

	ConfigPath string // path to config file
}

func (u *Down) Cmd() *cobra.Command {
	var c = &cobra.Command{
		Use:   "down",
		Short: "down rollback the DB to the most recent version available",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			dir, _ := os.Getwd()
			confx.Init(confx.SetPath([]string{dir}))
			dsn := fmt.Sprintf(
				"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",

				confx.GetString("mysql.username"),
				confx.GetString("mysql.password"),
				confx.GetString("mysql.host"),
				confx.GetInt("mysql.port"),
				confx.GetString("mysql.dbName"),
				confx.GetString("mysql.charset"),
			)
			fullPath := filepath.Join(dir, u.Dir)
			//var cd *exec.Cmd
			//if runtime.GOOS == "windows" {
			//	cd = exec.Command("cmd", "/c", fmt.Sprintf("--dir %s %s %s up", fullPath, u.Type, strconv.Quote(dsn)))
			//} else {
			//	cd = exec.Command("bash", "-c", fmt.Sprintf("--dir %s %s %s up", fullPath, u.Type, strconv.Quote(dsn)))
			//}
			//cd := exec.Command("goose", fmt.Sprintf("--dir %s %s %s up", fullPath, u.Type, strconv.Quote(dsn)))
			cd := exec.Command("goose", "--dir", fullPath, u.Type, strconv.Quote(dsn), "down")
			fmt.Println(cd.String())
			//var out bytes.Buffer
			//cd.Stdout = &out
			//err := cd.Run()
			//if err != nil {
			//	panic(err)
			//}
			//fmt.Println(out.String())
		},
	}

	c.Flags().StringVarP(&u.Type, "type", "t", "mysql", "database driver to use")
	c.Flags().StringVarP(&u.Dir, "dir", "d", "./internal/db/migrations", "directory with migration files")
	c.Flags().StringVarP(&u.ConfigPath, "config", "c", "./config.yml", "path to config file")

	return c
}
