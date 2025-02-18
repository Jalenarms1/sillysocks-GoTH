package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gofrs/uuid"
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

func (o *Order) Insert(tx *sql.Tx) error {
	resp, err := tx.Exec(`
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

func GetOrder(id string) (*Order, error) {
	row := DB.QueryRow(`select Id, PaymentIntentId, SubTotal, Tax, GrandTotal, ShippingTotal, ShippingLine1, ShippingLine2, ShippingCity, ShippingState, ShippingPostalCode, CustomerName, CustomerEmail, CreatedAt, Status from "Order" where Id = ?`, id)
	if row == nil {
		return nil, errors.New("order not found")
	}

	var order Order
	err := row.Scan(&order.Id, &order.PaymentIntentId, &order.SubTotal, &order.Tax, &order.GrandTotal, &order.ShippingTotal, &order.ShippingLine1, &order.ShippingLine2, &order.ShippingCity, &order.ShippingState, &order.ShippingPostalCode, &order.CustomerName, &order.CustomerEmail, &order.CreatedAt, &order.Status)
	if err != nil {
		return nil, errors.New("error scanning data into order")
	}

	err = order.LoadCartItems()
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (o *Order) Save() error {
	_, err := DB.Exec(`update "Order" set Id = ?, PaymentIntentId = ?, SubTotal = ?, Tax = ?, GrandTotal = ?, ShippingTotal = ?, ShippingLine1 = ?, ShippingLine2 = ?, ShippingCity = ?, ShippingState = ?, ShippingPostalCode = ?, CustomerName = ?, CustomerEmail = ?, CreatedAt = ?, Status = ? where Id = ?`, o.Id, o.PaymentIntentId, o.SubTotal, o.Tax, o.GrandTotal, o.ShippingTotal, o.ShippingLine1, o.ShippingLine2, o.ShippingCity, o.ShippingState, o.ShippingPostalCode, o.CustomerName, o.CustomerEmail, o.CreatedAt, o.Status, o.Id)
	if err != nil {
		return err
	}

	return nil

}

func (o *Order) LoadCartItems() error {

	var cartItems []CartItem
	res, err := DB.Query("select Id, ProductId, ProductName, ProductImage, ProductPrice, OrderId, Quantity, SubTotal, Size from CartItem where OrderId = ?", o.Id)
	if err != nil {
		return err
	}

	for res.Next() {
		var ci CartItem
		err := res.Scan(&ci.Id, &ci.ProductId, &ci.ProductName, &ci.ProductImage, &ci.ProductPrice, &ci.OrderId, &ci.Quantity, &ci.SubTotal, &ci.Size)
		if err != nil {
			return err
		}

		cartItems = append(cartItems, ci)
	}

	o.CartItems = cartItems

	return nil
}

func InsertCartItems(tx *sql.Tx, cartItems []CartItem, orderId string) error {
	queryStrings := make([]string, len(cartItems))
	queryArgs := []interface{}{}

	for i, item := range cartItems {
		queryStrings[i] = "(?, ?, ?, ?, ?, ?, ?, ?, ?)"
		fmt.Println(item.Id)
		fmt.Println(item.Product.Id)
		fmt.Println(orderId)
		fmt.Println(item.Quantity)
		fmt.Println(item.SubTotal)

		uid, _ := uuid.NewV4()

		queryArgs = append(queryArgs, uid.String())
		queryArgs = append(queryArgs, item.Product.Id)
		queryArgs = append(queryArgs, item.Product.Name)
		queryArgs = append(queryArgs, item.Product.Image)
		queryArgs = append(queryArgs, item.Product.Price)
		queryArgs = append(queryArgs, orderId)
		queryArgs = append(queryArgs, item.Quantity)
		queryArgs = append(queryArgs, item.SubTotal)
		queryArgs = append(queryArgs, item.Size)

	}

	fmt.Println(queryArgs)

	_, err := tx.Exec("insert into CartItem (Id, ProductId, ProductName, ProductImage, ProductPrice, OrderId, Quantity, SubTotal, Size) values "+strings.Join(queryStrings, ", "), queryArgs...)

	return err

}
