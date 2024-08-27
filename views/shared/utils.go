package shared

import (
	"fmt"
)

func ConvStr[T int | float64](v T) string {
	return fmt.Sprintf("%v", v)

}

func GetJSON[T any | []any](v T) string {
	str := fmt.Sprintf("%+v", v)
	fmt.Printf("%s\n", str)
	return str

}
