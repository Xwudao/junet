package captchax

import (
	"image/color"

	cp "github.com/mojocn/base64Captcha"
)

var config = Config{
	source: "abcdefghjklmnpqrsxyz23456789",
	store:  cp.DefaultMemStore,
}

type Opt func(*Config)
type Config struct {
	source string
	store  cp.Store
}

func SetSource(s string) Opt {
	return func(c *Config) {
		if s != "" {
			c.source = s
		}
	}
}
func SetStore(s cp.Store) Opt {
	return func(c *Config) {
		if s != nil {
			c.store = s
		}
	}
}

func Generate(height, width, length int) (string, string, error) {
	driver := cp.NewDriverString(
		height,
		width, 0,
		cp.OptionShowHollowLine,
		length,
		config.source,
		&color.RGBA{R: 0, G: 0, B: 0, A: 50},
		nil,
		[]string{"RitaSmith.ttf", "3Dumb.ttf", "actionj.ttf", "DENNEthree-dee.ttf"},
	)
	cape := cp.NewCaptcha(driver, config.store)
	return cape.Generate()
}

func DefaultGenerate() (string, string, error) {
	return Generate(32, 100, 4)
}

func Verify(id, answer string, clear bool) bool {
	return config.store.Verify(id, answer, clear)
}
