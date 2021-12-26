# JuNet

An `composition` framework. It contains `gin`, `zap`, `jwt-go`, `viper`

And, out of box to use.

## Get Started

### install cli

```bash
## for go1.16+
go install github.com/Xwudao/junet/cmd/junet@latest
## for go1.15-
go get -u github.com/Xwudao/junet/cmd/junet
```

### init project

```bash
junet init demo -m github.com/Xwudao/my-new-demo
```

- project name: demo
- go module name: github.com/Xwudao/my-new-demo

## more

1、generate an new route

```bash
// in `pkg/routes/v1/gen.go`
//go:generate junet gen route -n Home
```

2、generate new db helper (gorm enhance)

```bash
//go:generate junet gen db

// the "gen:qs" is required, which tell junet cli generate funcs for `User`
//gen:qs
type User struct {
	ID       int64  `gorm:"column:id;primaryKey;AUTO_INCREMENT"`
	Username string `gorm:"column:username"`
	Email    string `gorm:"column:email"`
}
```

//this function original author: https://github.com/jirfag/go-queryset

//`junet` modify it to support gorm v2, and combine with `junet cli`