package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

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
	Sizes       string         `json:"sizes" db:"Sizes"`
}

type ProductFormReq struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int32   `json:"quantity"`
	Image    string  `json:"image"`
	Sizes    string  `json:"sizes"`
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
		Sizes:    product.Sizes,
	}
}

func SubmitProduct(product *ProductFormReq) error {
	// Define the SQL query with named parameters
	query := `
        INSERT INTO "Product" ("Id", "Name", "Price", "Quantity", "Image", "Sizes") 
        VALUES (:id, :name, :price, :quantity, :image, :sizes)
    `

	// Generate a new UUID for the product
	newId, _ := uuid.NewV4()

	// Print the product details for debugging
	fmt.Println("Inserting product:", newId.String(), product.Name, product.Price, product.Quantity, product.Image, product.Sizes)

	// Create a map for the parameters
	params := map[string]interface{}{
		"id":       newId.String(),
		"name":     product.Name,
		"price":    product.Price,
		"quantity": product.Quantity,
		"image":    product.Image,
		"sizes":    product.Sizes,
	}

	// Execute the query with the parameters
	_, err := db.DB.NamedExec(query, params)
	if err != nil {
		return err
	}

	return nil
}

func UpdateProduct(product *Product) error {
	query := `update "Product" set "Name"=$1, "Price"=$2, "Quantity"=$3, "Image"=$4, "Sizes"=$5 where "Id"=$6`

	_, err := db.DB.Exec(query, product.Name, product.Price, product.Quantity, product.Image, product.Sizes, product.Id)

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
	err := row.Scan(&p.Id, &p.Name, &p.Description, &p.Category, &p.Image, &p.Price, &p.Quantity, &p.Sizes)

	if err != nil {
		return nil, err
	}

	return &p, nil
}
