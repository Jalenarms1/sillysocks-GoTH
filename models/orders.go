package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type Order struct {
	Id                uuid.UUID    `json:"id" db:"Id"`
	OrderNbr          string       `json:"orderNbr" db:"OrderNbr"`
	CstId             string       `json:"cstId" db:"CstId"`
	PmtIntId          string       `json:"pmtIntId" db:"PmtIntId"`
	CustomerEmail     string       `json:"customerEmail" db:"CustomerEmail"`
	CustomerName      string       `json:"customerName" db:"CustomerName"`
	CreatedAt         int64        `json:"createdAt" db:"CreatedAt"`
	SubTotal          float64      `json:"total" db:"SubTotal"`
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
	OrderId   uuid.UUID `json:"orderId" db:"OrderId"`
	Total     float64   `json:"total" db:"Total"`
	Quantity  int64     `json:"quantity" db:"Quantity"`
	ProductId uuid.UUID `json:"productId" db:"ProductId"`
}

func NewOrder(cstId, pmtId, cstEmail, cstName, shpAddr1, shpAddr2, shpAddrCity, shpAddrState, shpAddrZip string, subTotal, tax, shipping float64) *Order {
	newId := generateUUIDv4()
	orderNbr := generateOrderNbr()
	return &Order{
		Id:                newId,
		OrderNbr:          orderNbr,
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
		CreatedAt:         time.Now().Unix(),
	}
}

func NewOrderItems(cart Cart, orderId uuid.UUID) *[]OrderItem {
	var orderItems []OrderItem

	for _, ci := range cart {
		newId := generateUUIDv4()
		oi := &OrderItem{
			Id:        newId,
			OrderId:   orderId,
			Total:     ci.Total,
			Quantity:  ci.Quantity,
			ProductId: ci.ProductId,
		}

		orderItems = append(orderItems, *oi)
	}

	return &orderItems
}

func (o *Order) Insert(userDb *sqlx.DB) error {
	query := `
		INSERT INTO "Order" (
			"Id", 
			"OrderNbr",
			"CstId", 
			"PmtIntId", 
			"CustomerEmail", 
			"CustomerName", 
			"CreatedAt", 
			"SubTotal", 
			"Tax", 
			"Shipping", 
			"ShippingAddrLine1", 
			"ShippingAddrLine2", 
			"ShippingAddrCity", 
			"ShippingAddrState", 
			"ShippingAddrZip", 
			"Shipped"
		) VALUES (
			:Id,
			:OrderNbr, 
			:CstId, 
			:PmtIntId, 
			:CustomerEmail, 
			:CustomerName, 
			:CreatedAt, 
			:SubTotal, 
			:Tax, 
			:Shipping, 
			:ShippingAddrLine1, 
			:ShippingAddrLine2, 
			:ShippingAddrCity, 
			:ShippingAddrState, 
			:ShippingAddrZip, 
			:Shipped
		)
	`

	_, err := userDb.NamedExec(query, &o)
	if err != nil {
		return err
	}

	return nil
}

func (o *Order) InsertItems(userDb *sqlx.DB) error {

	itemsMap := make([]map[string]interface{}, len(*o.OrderItems))
	for i, oi := range *o.OrderItems {
		itemsMap[i] = map[string]interface{}{
			"Id":        oi.Id,
			"OrderId":   oi.OrderId,
			"Total":     oi.Total,
			"Quantity":  oi.Quantity,
			"ProductId": oi.ProductId,
		}

	}

	query := `
		insert into "OrderItem" (
			"Id",
			"OrderId",
			"Total",
			"Quantity",
			"ProductId"
		) VALUES (
			:Id,
			:OrderId,
			:Total,
			:Quantity,
			:ProductId
		)
	`

	_, err := userDb.NamedExec(query, itemsMap)
	if err != nil {
		return err
	}

	return nil
}

// func HandleOrderSubmit(order Order) {
// 	query := `insert into Order ("Id", "CustomerEmail", "SubTotal", "Tax", "Shipped", "CreatedAt") values ($1,$2,$3,$4,$5, $6)`

// 	res, err := userDb.Exec(query, order.Id, order.CustomerEmail, order.SubTotal, order.Tax, order.Shipped, order.CreatedAt)
// }
