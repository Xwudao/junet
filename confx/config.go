package confx

import (
	"github.com/spf13/viper"
)

var config = Config{
	name: "config",
	ext:  "yaml",
	path: []string{".", "../config", "./config"},
}

type Opt func(*Config)
type Config struct {
	name string
	ext  string
	path []string

	run func(*viper.Viper)
}

func SetDefaultValue(run func(*viper.Viper)) Opt {
	return func(c *Config) {
		c.run = run
	}
}

func SetPath(p []string) Opt {
	return func(c *Config) {
		c.path = p
	}
}
func SetFileExt(e string) Opt {
	return func(c *Config) {
		c.ext = e
	}
}
func SetFileName(n string) Opt {
	return func(c *Config) {
		c.name = n
	}
}

func Init(opts ...Opt) {
	for _, opt := range opts {
		opt(&config)
	}

	viper.SetConfigName(config.name)
	viper.SetConfigType(config.ext)
	for _, p := range config.path {
		viper.AddConfigPath(p)
	}

	if config.run != nil {
		config.run(viper.GetViper())
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("config file not found")
		} else {
			panic("init config failed: err")
		}
	}

	viper.WatchConfig()
}

func SaveConfigAs(file string) error {
	return viper.WriteConfigAs(file)
}
func SaveConfig() error {
	return viper.WriteConfig()
}

//func init() {
//initOnceViper.Do(func() {
//	v = viper.New()
//	v.SetConfigName("config")
//	v.SetConfigType("yaml")
//	v.AddConfigPath(".")
//	v.AddConfigPath("..")
//	v.AddConfigPath("../config")
//	v.AddConfigPath("./config")
//	configPathFile, err := getConfigPathFile()
//	if err == nil {
//		v.AddConfigPath(configPathFile)
//	}
//
//	if err := v.ReadInConfig(); err != nil {
//		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
//			panic(errors.New("config file not found"))
//		} else {
//			panic(errors.New("init config failed: err"))
//		}
//	}
//	v.WatchConfig()
//})
//}

// our functions
//
//func GetViper() *viper.Viper {
//	return v
//}
//
//func SetDefaultValue(f func(v *viper.Viper)) {
//	f(v)
//}
//
//// Viper functions
//
//// AddConfigPath adds a path for Viper to search for the config file in.
//// Can be called multiple times to define multiple search paths.
//func AddConfigPath(in string) { v.AddConfigPath(in) }
//
//// SetDefault sets the default value for this key.
//// SetDefault is case-insensitive for a key.
//// Default only used when no value is provided by the user via flag, config or ENV.
//func SetDefault(key string, value interface{}) { v.SetDefault(key, value) }
//
//// GetString returns the value associated with the key as a string.
//func GetString(key string) string { return v.GetString(key) }
//
//// GetBool returns the value associated with the key as a boolean.
//func GetBool(key string) bool { return v.GetBool(key) }
//
//// GetInt returns the value associated with the key as an integer.
//func GetInt(key string) int { return v.GetInt(key) }
//
//// GetInt32 returns the value associated with the key as an integer.
//func GetInt32(key string) int32 { return v.GetInt32(key) }
//
//// GetInt64 returns the value associated with the key as an integer.
//func GetInt64(key string) int64 { return v.GetInt64(key) }
//
//// GetUint returns the value associated with the key as an unsigned integer.
//func GetUint(key string) uint { return v.GetUint(key) }
//
//// GetUint32 returns the value associated with the key as an unsigned integer.
//func GetUint32(key string) uint32 { return v.GetUint32(key) }
//
//// GetUint64 returns the value associated with the key as an unsigned integer.
//func GetUint64(key string) uint64 { return v.GetUint64(key) }
//
//// GetFloat64 returns the value associated with the key as a float64.
//func GetFloat64(key string) float64 { return v.GetFloat64(key) }
//
//// GetTime returns the value associated with the key as time.
//func GetTime(key string) time.Time { return v.GetTime(key) }
//
//// GetDuration returns the value associated with the key as a duration.
//func GetDuration(key string) time.Duration { return v.GetDuration(key) }
//
//// GetIntSlice returns the value associated with the key as a slice of int values.
//func GetIntSlice(key string) []int { return v.GetIntSlice(key) }
//
//// GetStringSlice returns the value associated with the key as a slice of strings.
//func GetStringSlice(key string) []string { return v.GetStringSlice(key) }
//
//// GetStringMap returns the value associated with the key as a map of interfaces.
//func GetStringMap(key string) map[string]interface{} { return v.GetStringMap(key) }
//
//// GetStringMapString returns the value associated with the key as a map of strings.
//func GetStringMapString(key string) map[string]string { return v.GetStringMapString(key) }
//
//// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
//func GetStringMapStringSlice(key string) map[string][]string { return v.GetStringMapStringSlice(key) }
