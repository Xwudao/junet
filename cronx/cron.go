package cronx

import (
	"fmt"

	"github.com/robfig/cron/v3"

	"github.com/Xwudao/junet/shutdown"
)

var config = Config{
	cn:          cron.New(),
	minuteFunc:  map[int][]func(){},
	secondsFunc: map[int][]func(){},
	hourFunc:    map[int][]func(){},
}

type Config struct {
	//Run once a year, midnight, Jan. 1st
	yearlyFunc []func()
	//Run once a month, midnight, first of month
	monthlyFunc []func()
	//Run once a week, midnight between Sat/Sun
	weeklyFunc []func()
	//Run once a day, midnight
	dailyFunc []func()
	//Run once an hour, beginning of hour
	hourlyFunc []func()
	//Run on minute
	minuteFunc map[int][]func()
	//Run on seconds
	secondsFunc map[int][]func()
	//Run on hour
	hourFunc map[int][]func()

	cn *cron.Cron
}

func (c Config) Close() error {
	c.cn.Stop()
	return nil
}

func AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	return config.cn.AddFunc(spec, cmd)
}

type Opt func(*Config)

func AddSecondsFunc(s int, f ...func()) Opt {
	return func(c *Config) {
		c.secondsFunc[s] = append(c.secondsFunc[s], f...)
	}
}
func AddMinuteFunc(m int, f ...func()) Opt {
	return func(c *Config) {
		c.minuteFunc[m] = append(c.minuteFunc[m], f...)
	}
}
func AddHourly(f func()) Opt {
	return func(c *Config) {
		c.hourlyFunc = append(c.hourlyFunc, f)
	}
}
func AddDaily(f func()) Opt {
	return func(c *Config) {
		c.dailyFunc = append(c.dailyFunc, f)
	}
}
func AddMonthly(f func()) Opt {
	return func(c *Config) {
		c.monthlyFunc = append(c.monthlyFunc, f)
	}
}
func AddYearly(f func()) Opt {
	return func(c *Config) {
		c.yearlyFunc = append(c.yearlyFunc, f)
	}
}

func Init(opts ...Opt) error {
	for _, opt := range opts {
		opt(&config)
	}
	var err error
	for _, f := range config.yearlyFunc {
		_, err = config.cn.AddFunc("@yearly", f)
		return err
	}
	for _, f := range config.monthlyFunc {
		_, err = config.cn.AddFunc("@monthly", f)
		return err
	}
	for _, f := range config.weeklyFunc {
		_, err = config.cn.AddFunc("@weekly", f)
		return err
	}
	for _, f := range config.dailyFunc {
		_, err = config.cn.AddFunc("@daily", f)
		return err
	}
	for _, f := range config.hourlyFunc {
		_, err = config.cn.AddFunc("@hourly", f)
		return err
	}

	for sec, funcs := range config.secondsFunc {
		for _, f := range funcs {
			_, err = config.cn.AddFunc(fmt.Sprintf("@every %ds", sec), f)
			if err != nil {
				return err
			}
		}
	}
	for min, funcs := range config.minuteFunc {
		for _, f := range funcs {
			_, err = config.cn.AddFunc(fmt.Sprintf("@every %dm", min), f)
			if err != nil {
				return err
			}
		}
	}
	for hour, funcs := range config.hourFunc {
		for _, f := range funcs {
			_, err = config.cn.AddFunc(fmt.Sprintf("@every %dh", hour), f)
			if err != nil {
				return err
			}
		}
	}

	config.cn.Start()
	shutdown.Add(config)
	return nil
}
