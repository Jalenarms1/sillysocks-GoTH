package models

import (
	"database/sql"
	"encoding/json"

	"github.com/Jalenarms1/sillysocks-GoTH/db"
	"github.com/gofrs/uuid"
)

type Product struct {
	Id          uuid.UUID      `json:"id" db:"Id"`
	Name        string         `json:"name" db:"Name"`
	Description sql.NullString `json:"description" db:"Description"`
	Category    sql.NullString `json:"category" db:"Category"`
	Image       string         `json:"image" db:"Image"`
	Price       float64        `json:"price" db:"Price"`
	Quantity    int32          `json:"quantity" db:"Quantity"`
}

func (product *Product) ToJson() string {
	bytes, _ := json.Marshal(product)

	return string(bytes)
}

func NewProduct(product Product) *Product {
	if product.Id == uuid.Nil {
		newId := generateUUIDv4()

		product.Id = newId
	}

	return &Product{
		Id:       product.Id,
		Name:     product.Name,
		Price:    product.Price,
		Quantity: product.Quantity,
		Image:    product.Image,
	}
}

func SubmitProduct(product *Product) error {
	query := `insert into "Product" ("Id", "Name", "Price", "Quantity", "Image") values ($1, $2, $3, $4, $5)`

	_, err := db.DB.Exec(query, product.Id, product.Name, product.Price, product.Quantity, product.Image)

	if err != nil {
		return err
	}

	return nil

}

func UpdateProduct(product *Product) error {
	query := `update "Product" set "Name"=$1, "Price"=$2, "Quantity"=$3, "Image"=$4 where "Id"=$5`

	_, err := db.DB.Exec(query, product.Name, product.Price, product.Quantity, product.Image, product.Id)

	if err != nil {
		return err
	}

	return nil
}

func DeleteProduct(productId string) error {
	query := `delete from "Product" where "Id"=$1`

	_, err := db.DB.Exec(query, productId)

	if err != nil {
		return err
	}

	return nil
}

func GetProducts() ([]Product, error) {
	var products []Product
	query := `
		select * from "Product"
	`
	err := db.DB.Select(&products, query)
	if err != nil {
		return nil, err
	}

	return products, nil

}

func GetProductById(productId string) (*Product, error) {
	row := db.DB.QueryRow(`select * from "Product" where "Id"=$1`, productId)
	if row.Err() == sql.ErrNoRows {
		return nil, row.Err()
	}

	var p Product
	err := row.Scan(&p.Id, &p.Name, &p.Description, &p.Category, &p.Image, &p.Price, &p.Quantity)

	if err != nil {
		return nil, err
	}

	return &p, nil
}
