package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Order struct {
	Id                uuid.UUID    `json:"id" db:"Id"`
	CstId             string       `json:"cstId" db:"CstId"`
	PmtIntId          string       `json:"pmtIntId" db:"PmtIntId"`
	CustomerEmail     string       `json:"customerEmail" db:"CustomerEmail"`
	CustomerName      string       `json:"customerName" db:"CustomerName"`
	CreatedAt         time.Time    `json:"createdAt" db:"CreatedAt"`
	SubTotal          float64      `json:"total" db:"Total"`
	Tax               float64      `json:"tax" db:"Tax"`
	Shipping          float64      `json:"shipping" db:"Shipping"`
	ShippingAddrLine1 string       `json:"shippingAddrLine1" db:"ShippingAddrLine1"`
	ShippingAddrLine2 string       `json:"shippingAddrLine2" db:"ShippingAddrLine2"`
	ShippingAddrCity  string       `json:"shippingAddrCity" db:"ShippingAddrCity"`
	ShippingAddrState string       `json:"shippingAddrState" db:"ShippingAddrState"`
	ShippingAddrZip   string       `json:"shippingAddrZip" db:"ShippingAddrZip"`
	Shipped           bool         `json:"shipped" db:"Shipped"`
	OrderItems        *[]OrderItem `json:"orderItems"`
}

type OrderItem struct {
	Id        uuid.UUID `json:"id" db:"Id"`
	OrderId   string    `json:"orderId" db:"OrderId"`
	ProductId string    `json:"productId" db:"ProductId"`
}

func NewOrder(cstId, pmtId, cstEmail, cstName, shpAddr1, shpAddr2, shpAddrCity, shpAddrState, shpAddrZip string, subTotal, tax, shipping float64) *Order {
	newId, _ := generateUUIDv4()

	return &Order{
		Id:                newId,
		CstId:             cstId,
		PmtIntId:          pmtId,
		CustomerEmail:     cstEmail,
		CustomerName:      cstName,
		SubTotal:          subTotal,
		Tax:               tax,
		Shipping:          shipping,
		ShippingAddrLine1: shpAddr1,
		ShippingAddrLine2: shpAddr2,
		ShippingAddrCity:  shpAddrCity,
		ShippingAddrState: shpAddrState,
		ShippingAddrZip:   shpAddrZip,
		CreatedAt:         time.Now(),
	}
}

// func HandleOrderSubmit(order Order) {
// 	query := `insert into Order ("Id", "CustomerEmail", "SubTotal", "Tax", "Shipped", "CreatedAt") values ($1,$2,$3,$4,$5, $6)`

// 	res, err := db.DB.Exec(query, order.Id, order.CustomerEmail, order.SubTotal, order.Tax, order.Shipped, order.CreatedAt)
// }
