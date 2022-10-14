package config

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	SetConfigPath("../../")
	fmt.Println(Config().Hsf)
}
