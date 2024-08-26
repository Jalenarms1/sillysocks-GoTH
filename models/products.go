package models

import (
	"database/sql"
	"log"

	"github.com/Jalenarms1/sillysocks-GoTH/db"
)

type Product struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

func NewProduct(product Product) *Product {
	if product.Id == "" {
		newId, err := generateUUIDv4()
		if err != nil {
			log.Fatal(err)
		}
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
	rows, err := db.DB.Query(`select * from "Product"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product

		if err := rows.Scan(&p.Id, &p.Name, &p.Description, &p.Category, &p.Image, &p.Price, &p.Quantity); err != nil {
			return nil, err
		}

		products = append(products, p)

	}

	if err := rows.Err(); err != nil {
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
