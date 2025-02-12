package models

import (
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type CartItem struct {
	Id        uuid.UUID `json:"id" db:"Id"`
	Product   Product   `json:"product"`
	ProductId uuid.UUID `json:"productId" db:"ProductId"`
	Total     float64   `json:"total" db:"Total"`
	Quantity  int64     `json:"quantity" db:"Quantity"`
}

type Cart []CartItem

type CartRecord struct {
	Id              uuid.UUID      `db:"Id"`
	ProductId       uuid.UUID      `db:"ProductId"`
	Total           float64        `db:"Total"`
	Quantity        int32          `db:"Quantity"`
	Name            string         `db:"Name"`
	Description     sql.NullString `db:"Description"`
	Category        sql.NullString `db:"Category"`
	Image           string         `db:"Image"`
	Price           float64        `db:"Price"`
	ProductQuantity int32          `db:"ProductQuantity"`
}

func (ci *CartItem) ToJson() string {
	bytes, _ := json.Marshal(ci)

	return string(bytes)
}

func (c *Cart) ToJson() string {
	bytes, _ := json.Marshal(c)

	return string(bytes)
}

func (ci *CartItem) insert(userDb *sqlx.DB) error {
	query := `
		insert into CartItem ("Id", "ProductId", "Total", "Quantity") values (:Id, :ProductId, :Total, :Quantity)
	`

	_, err := userDb.NamedExec(query, ci)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cart) Clear(userDb *sqlx.DB) error {
	c = &Cart{}
	if err := c.Save(userDb); err != nil {
		return err
	}

	return nil
}

func (c *Cart) Save(userDb *sqlx.DB) error {
	itemMaps := make([]map[string]interface{}, len(*c))
	for i, item := range *c {
		var id string
		if item.Id == uuid.Nil {
			newId, _ := uuid.NewV4()
			id = newId.String()
		} else {
			id = item.Id.String()
		}

		var productId string
		if item.ProductId == uuid.Nil {
			productId = item.Product.Id.String()
		} else {
			productId = item.ProductId.String()
		}

		itemMaps[i] = map[string]interface{}{
			"Id":        id,
			"ProductId": productId,
			"Total":     item.Total,
			"Quantity":  item.Quantity,
		}
	}

	tx, err := userDb.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, delErr := tx.Exec(`delete from CartItem;`)
	if delErr != nil {
		return delErr
	}

	query := `
        insert into CartItem (Id, ProductId, Total, Quantity)
        values (:Id, :ProductId, :Total, :Quantity);
    `
	if len(itemMaps) > 0 {
		_, insErr := tx.NamedExec(query, itemMaps)
		if insErr != nil {
			return err
		}

	}

	cmErr := tx.Commit()
	if cmErr != nil {
		return cmErr
	}

	return nil
}

func (ci *CartItem) save(userDb *sqlx.DB) error {
	query := `
		update CartItem
		set
			"Total" = :Total,
			"Quantity" = :Quantity
		where "Id" = :Id
	`

	_, err := userDb.NamedExec(query, ci)
	if err != nil {
		return err
	}

	return nil
}

func (cart *Cart) NumOfItems(userDb *sqlx.DB) int {
	var count int

	err := userDb.QueryRow(`select sum("Quantity") from CartItem`).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	return count
}

func (cart *Cart) GetTax() float64 {
	var tax float64

	for _, ci := range *cart {
		sbt := ci.Product.Price * float64(ci.Quantity)
		itemTax := math.Round((sbt*1.08-sbt)*100) / 100
		tax += itemTax
	}

	return tax
}

func (cart *Cart) GetTotal() float64 {

	return cart.GetSubTotal() + cart.GetTax() + cart.GetShippingCost()
}

func (cart *Cart) GetSubTotal() float64 {
	var total float64
	for _, ci := range *cart {
		total += ci.Product.Price * float64(ci.Quantity)
	}

	return total
}

func (cart *Cart) GetShippingCost() float64 {
	var total float64
	for _, ci := range *cart {
		total += ci.Product.Price * float64(ci.Quantity)
	}

	if total >= 20 || total == 0 {
		return 0
	}

	return 5
}

func AddToCart(product *Product, userDb *sqlx.DB) error {
	// fmt.Println("Add to cart")
	// fmt.Print(cart)
	// stmt, err := userDb.Prepare(`
	// 	insert into CartItem ("")
	// `)
	var cartItem CartItem
	if err := GetCartItem(product.Id, &cartItem, userDb); err != nil {
		cartItem = CartItem{}
	}

	if cartItem.Id == uuid.Nil {
		//
		newId, _ := uuid.NewV4()
		cartItem := &CartItem{
			Id:        newId,
			ProductId: product.Id,
			Total:     product.Price,
			Quantity:  1,
		}

		err := cartItem.insert(userDb)
		if err != nil {
			return err
		}
	} else {
		cartItem.Quantity += 1
		cartItem.Total = cartItem.Product.Price * float64(cartItem.Quantity)

		err := cartItem.save(userDb)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetCart(w http.ResponseWriter, r *http.Request, userDb *sqlx.DB) Cart {

	query := `
	    SELECT
	        ci.Id,
	        ci.ProductId,
			ci.Total,
	        ci.Quantity,
	        p.Name,
	        p.Description,
	        p.Category,
	        p.Image,
	        p.Price,
	        p.Quantity as ProductQuantity
	    FROM
	        CartItem ci
	    JOIN
	        Product p ON ci.ProductId = p.Id;
	`

	var recordsList []CartRecord

	err := userDb.Select(&recordsList, query)
	// err := userDb.Select(&recordsList, query)
	if err != nil {
		log.Fatal(err)
	}
	var cart Cart
	for _, r := range recordsList {
		product := &Product{
			Id:          r.ProductId,
			Name:        r.Name,
			Description: r.Description,
			Category:    r.Category,
			Image:       r.Image,
			Price:       r.Price,
			Quantity:    r.ProductQuantity,
		}

		cartItem := &CartItem{
			Id:        r.Id,
			Product:   *product,
			ProductId: r.ProductId,
			Total:     r.Total,
			Quantity:  int64(r.Quantity),
		}

		cart = append(cart, *cartItem)
	}

	return cart

}

type CartItemResp struct {
	Id        uuid.UUID `db:"Id"`
	ProductId uuid.UUID `db:"ProductId"`
	Total     float64   `db:"Total"`
	Quantity  int32     `db:"Quantity"`
}

func GetCartItem(itemId uuid.UUID, ci *CartItem, userDb *sqlx.DB) error {
	var resp CartItemResp

	err := userDb.Get(&resp, `select * from CartItem where Id = ? or ProductId = ?`, itemId)
	if err != nil {
		return err
	}

	ci.Id = resp.Id
	ci.ProductId = resp.ProductId
	ci.Total = resp.Total
	ci.Quantity = int64(resp.Quantity)

	return nil
}
