package models

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"reflect"

	"github.com/gofrs/uuid"
)

func generateUUIDv4() uuid.UUID {
	newId, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}
	return newId
}

func generateOrderNbr() string {

	return fmt.Sprintf("%d-%d", rand.Int31(), rand.Int31())
}

func DataGateway[T any](f func() ([]T, error)) []T {
	data, err := f()
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func ScanStruct(rows *sql.Rows, dest interface{}) error {
	val := reflect.ValueOf(&dest).Elem()
	typ := val.Type()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	columnNames := make(map[string]int)
	for i, col := range columns {
		columnNames[col] = i
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return err
		}

		elem := reflect.New(typ).Elem()
		for i, col := range columns {
			if _, ok := columnNames[col]; ok {
				field := elem.FieldByName(col)
				if field.IsValid() && field.CanSet() {
					val := reflect.ValueOf(*(values[i].(*interface{})))
					if val.Type().ConvertibleTo(field.Type()) {
						field.Set(val.Convert(field.Type()))
					}
				}
			}
		}

		val.Set(elem)
	}

	return nil
}
