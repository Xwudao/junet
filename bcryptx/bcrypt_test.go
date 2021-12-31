package bcryptx

import (
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	password, err := GeneratePassword("hello")
	if err != nil {
		t.FailNow()
		return
	}
	t.Logf("password: %s\n", password)
}

func TestValidatePassword(t *testing.T) {
	ok := ValidatePassword("$2a$10$5EUcgy8KKNKXoj/uJYm5Ye8YLRpS5cLlhd2hfZNcpQbItZLgjVm8K", "hello")
	if !ok {
		t.FailNow()
	}
	t.Logf("ok!\n")
}
