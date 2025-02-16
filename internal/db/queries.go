package db

import (
	"database/sql"
	"fmt"
	"log"
)

func GetProducts() []Product {
	var products []Product

	resp, err := DB.Query("select Id, Name, Description, Category, Image, Price, Quantity, Sizes from Product")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Close()
	for resp.Next() {
		var p Product
		resp.Scan(&p.Id, &p.Name, &p.Description, &p.Category, &p.Image, &p.Price, &p.Quantity, &p.Sizes)
		fmt.Println(p)
		products = append(products, p)
	}

	return products
}

func GetProduct(id string) *Product {
	record := DB.QueryRow("select Id, Name, Description, Category, Image, Price, Quantity, Sizes from Product where Id = ?", id)

	if record == nil {
		return nil
	}

	var p Product
	err := record.Scan(&p.Id, &p.Name, &p.Description, &p.Category, &p.Image, &p.Price, &p.Quantity, &p.Sizes)
	if err == sql.ErrNoRows {
		return nil
	}

	return &p

}

func (o *Order) Insert() error {
	resp, err := DB.Exec(`
		INSERT INTO "Order" (
			Id, SubTotal, Tax, GrandTotal, ShippingTotal, ShippingLine1, ShippingLine2, 
			ShippingCity, ShippingState, ShippingPostalCode, CustomerName, CustomerEmail, 
			CreatedAt, Status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		o.Id, o.SubTotal, o.Tax, o.GrandTotal, o.ShippingTotal, o.ShippingLine1, o.ShippingLine2,
		o.ShippingCity, o.ShippingState, o.ShippingPostalCode, o.CustomerName, o.CustomerEmail,
		o.CreatedAt, o.Status,
	)
	if err != nil {
		return err
	}

	lastInsertedId, err := resp.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Println(lastInsertedId)

	return nil
}

func GetOrder(id string) *Order {
	row := DB.QueryRow(`select Id, PaymentIntentId, SubTotal, Tax, GrandTotal, ShippingTotal, ShippingLine1, ShippingLine2, ShippingCity, ShippingState, ShippingPostalCode, CustomerName, CustomerEmail, CreatedAt, Status from "Order" where Id = ?`, id)
	if row == nil {
		return nil
	}

	var order Order
	err := row.Scan(&order.Id, &order.PaymentIntentId, &order.SubTotal, &order.Tax, &order.GrandTotal, &order.ShippingTotal, &order.ShippingLine1, &order.ShippingLine2, &order.ShippingCity, &order.ShippingState, &order.ShippingPostalCode, &order.CustomerName, &order.CustomerEmail, &order.CreatedAt, &order.Status)
	if err != nil {
		return nil
	}

	return &order
}

func (o *Order) Save() error {
	_, err := DB.Exec(`update "Order" set Id = ?, PaymentIntentId = ?, SubTotal = ?, Tax = ?, GrandTotal = ?, ShippingTotal = ?, ShippingLine1 = ?, ShippingLine2 = ?, ShippingCity = ?, ShippingState = ?, ShippingPostalCode = ?, CustomerName = ?, CustomerEmail = ?, CreatedAt = ?, Status = ? where Id = ?`, o.Id, o.PaymentIntentId, o.SubTotal, o.Tax, o.GrandTotal, o.ShippingTotal, o.ShippingLine1, o.ShippingLine2, o.ShippingCity, o.ShippingState, o.ShippingPostalCode, o.CustomerName, o.CustomerEmail, o.CreatedAt, o.Status, o.Id)
	if err != nil {
		return err
	}

	return nil

}
