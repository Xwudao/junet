package logx

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	rotatelogs "github.com/iproj/file-rotatelogs"
)

var z = zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), os.Stdout, zap.DebugLevel))
var logger = z.Sugar()

var config = Config{
	encoder:        encoder(),
	maxAge:         time.Hour * 24 * 30,
	rotationTime:   time.Hour * 24,
	disableDefault: false,
	logPath:        "logs",
	logSuffix:      ".%Y-%m-%d.log",
}

type Opt func(*Config)
type Config struct {
	disableDefault bool

	rotationTime time.Duration
	maxAge       time.Duration
	encoder      zapcore.Encoder

	logPath   string
	logSuffix string

	core []zapcore.Core
}

func SetLogSuffix(s string) Opt {
	return func(c *Config) {
		c.logSuffix = s
	}
}
func AddCore(core ...zapcore.Core) Opt {
	return func(c *Config) {
		c.core = append(c.core, core...)
	}
}
func SetLogPath(p string) Opt {
	return func(c *Config) {
		c.logPath = p
	}
}
func SetRotationTime(tme time.Duration) Opt {
	return func(c *Config) {
		c.rotationTime = tme
	}
}
func SetMaxAge(age time.Duration) Opt {
	return func(c *Config) {
		c.maxAge = age
	}
}
func SetDisableDefault(b bool) Opt {
	return func(c *Config) {
		c.disableDefault = b
	}
}
func SetEncorder(en zapcore.Encoder) Opt {
	return func(c *Config) {
		c.encoder = en
	}
}

func encoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func Init(opts ...Opt) {
	for _, opt := range opts {
		opt(&config)
	}
	var core []zapcore.Core
	if !config.disableDefault {
		c, _ := getDefaultCore()
		core = append(core, c...)
	}
	core = append(core, config.core...)
	z = zap.New(zapcore.NewTee(core...))
	logger = z.Sugar()
}

func getDefaultCore() ([]zapcore.Core, zapcore.Encoder) {
	dir, _ := os.Getwd()
	logsDir := filepath.Join(dir, config.logPath)
	//日志级别
	allPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.DebugLevel
	})
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error级别
		return lev >= zap.WarnLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info和debug级别,debug级别是最低的
		return lev < zap.WarnLevel && lev >= zap.DebugLevel
	})

	allLogSync := zapcore.AddSync(getWriter(filepath.Join(logsDir, "all"+config.logSuffix), "all"))
	errorLogSync := zapcore.AddSync(getWriter(filepath.Join(logsDir, "error"+config.logSuffix), "error"))
	infoLogSync := zapcore.AddSync(getWriter(filepath.Join(logsDir, "info"+config.logSuffix), "info"))

	var cores []zapcore.Core

	cores = append(cores,
		zapcore.NewCore(
			config.encoder,
			zapcore.NewMultiWriteSyncer(allLogSync, zapcore.AddSync(os.Stdout)),
			allPriority,
		),
		zapcore.NewCore(
			config.encoder,
			zapcore.NewMultiWriteSyncer(errorLogSync),
			highPriority,
		),
		zapcore.NewCore(
			config.encoder,
			zapcore.NewMultiWriteSyncer(infoLogSync),
			lowPriority,
		),
	)
	return cores, config.encoder
}

//filename: %Y%m%d%H%M
func getWriter(filename string, lev string) io.Writer {
	// 保存30天内的日志，每24小时(整点)分割一次日志
	win := runtime.GOOS == "windows"
	var hook *rotatelogs.RotateLogs
	var err error
	var opts []rotatelogs.Option
	opts = append(opts,
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if !win {
		opts = append(opts,
			rotatelogs.WithLinkName(fmt.Sprintf("current.%s.log", lev)),
		)
	}
	hook, err = rotatelogs.New(
		filename,
		opts...,
	)
	if err != nil {
		panic(err)
	}
	return hook
}
