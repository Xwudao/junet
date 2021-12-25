package utils

import (
	"io/ioutil"
	"os"

	"gopkg.in/errgo.v2/fmt/errors"
)

func SaveToFile(p string, cnt []byte, cover bool) (err error) {
	if string(cnt) == "" {
		return errors.Newf("write file: empty content")
	}
	_, err = os.Stat(p)
	if err == nil {
		if !cover {
			return errors.Newf("file [%s] existed, please rename or remove it.", p)
		}
	}
	err = ioutil.WriteFile(p, cnt, os.ModePerm)
	if err != nil {
		return
	}
	return nil
}
