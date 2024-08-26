package shared

import (
	"encoding/json"
	"fmt"
	"log"
)

func ConvStr[T int | float64](v T) string {
	return fmt.Sprintf("%v", v)

}

func GetJSON[T any | []any](v T) string {
	var keeper T

	jsonData, err := json.Marshal(keeper)
	if err != nil {
		log.Fatal(err)
	}

	escapedProductJSON := fmt.Sprint(string(jsonData))

	return escapedProductJSON

}
