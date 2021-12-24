package captchax

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	generate, s, err := Generate(30, 400, 4)
	if err != nil {
		return
	}
	fmt.Println(generate, s)
}
