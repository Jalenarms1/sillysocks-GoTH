package models

import (
	"log"
	"time"
)

type Order struct {
	Id            string       `json:"id"`
	CustomerEmail string       `json:"customerEmail"`
	CreatedAt     time.Time    `json:"createdAt"`
	SubTotal      int64        `json:"total"`
	Tax           int64        `json:"tax"`
	Shipped       bool         `json:"shipped"`
	OrderItems    *[]OrderItem `json:"orderItems"`
}

type OrderItem struct {
	Id        string `json:"id"`
	OrderId   string `json:"orderId"`
	ProductId string `json:"productId"`
}

func NewOrder(order Order) *Order {
	if order.Id == "" {
		newId, err := generateUUIDv4()
		if err != nil {
			log.Fatal(err)
		}
		order.Id = newId
	}

	return &Order{
		Id:            order.Id,
		CustomerEmail: order.CustomerEmail,
		SubTotal:      order.SubTotal,
		Tax:           order.Tax,
		Shipped:       order.Shipped,
		CreatedAt:     order.CreatedAt,
	}
}

// func HandleOrderSubmit(order Order) {
// 	query := `insert into Order ("Id", "CustomerEmail", "SubTotal", "Tax", "Shipped", "CreatedAt") values ($1,$2,$3,$4,$5, $6)`

// 	res, err := db.DB.Exec(query, order.Id, order.CustomerEmail, order.SubTotal, order.Tax, order.Shipped, order.CreatedAt)
// }
