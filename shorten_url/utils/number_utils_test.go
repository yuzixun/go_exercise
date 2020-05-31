package utils

import (
	"fmt"
	"testing"
)

func TestDecimalTo62(t *testing.T) {
	fmt.Println(decimalTo62(61))
	fmt.Println(decimalTo62(100))
	fmt.Println(decimalTo62(10000))
}

func TestDecimalFrom62(t *testing.T) {
	fmt.Println(decimalFrom62("C1"))
	fmt.Println(decimalFrom62("iB2"))
	fmt.Println(decimalFrom62("Z"))
}
