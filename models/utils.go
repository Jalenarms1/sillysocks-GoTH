package models

import (
	"database/sql"
	"log"
	"reflect"

	"github.com/gofrs/uuid"
)

func generateUUIDv4() (uuid.UUID, error) {
	// var uuid [16]byte
	// _, err := rand.Read(uuid[:])
	// if err != nil {
	// 	return uuid.Nil, err
	// }

	// // Set version (4 bits) to 0100
	// uuid[6] = (uuid[6] & 0x0f) | 0x40
	// // Set variant (2 bits) to 10
	// uuid[8] = (uuid[8] & 0x3f) | 0x80

	// // Format UUID as a string
	// return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
	// 		uuid[0:4],
	// 		uuid[4:6],
	// 		uuid[6:8],
	// 		uuid[8:10],
	// 		uuid[10:]),
	// 	nil

	return uuid.NewV4()
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
