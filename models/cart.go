package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Jalenarms1/sillysocks-GoTH/db"
	"github.com/gofrs/uuid"
)

type CartItem struct {
	Id        uuid.UUID `json:"id" db:"Id"`
	Product   Product   `json:"product"`
	ProductId uuid.UUID `json:"productId" db:"ProductId"`
	Total     float64   `json:"total" db:"Total"`
	Quantity  int32     `json:"quantity" db:"Quantity"`
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
	fmt.Printf("\n\n\n")
	fmt.Printf("%v\n", c)
	bytes, _ := json.Marshal(c)

	fmt.Printf("\n%s", string(bytes))
	return string(bytes)
}

func (ci *CartItem) insert() error {
	query := `
		insert into "CartItem" ("Id", "ProductId", "Total", "Quantity") values (:Id, :ProductId, :Total, :Quantity)
	`

	_, err := db.DB.NamedExec(query, ci)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cart) Save() error {
	query := `
		delete from "CartItem";
		insert into "CartItem" ("Id", "ProductId", "Total", "Quantity") values (:Id, :ProductId, :Total, :Quantity);
	`

	_, err := db.DB.NamedExec(query, c)
	if err != nil {
		return err
	}

	return nil
}

func (ci *CartItem) save() error {
	query := `
		update "CartItem"
		set
			"Total" = :Total,
			"Quantity" = :Quantity
		where "Id" = :Id
	`

	_, err := db.DB.NamedExec(query, ci)
	if err != nil {
		return err
	}

	return nil
}

func (cart *Cart) NumOfItems() int {
	var count int

	err := db.DB.QueryRow(`select sum("Quantity") from "CartItem"`).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nNumber of items: %d\n", count)

	return count
}

func (cart *Cart) GetTax() float64 {
	sbt := cart.GetSubTotal()
	return sbt*1.08 - sbt
}

func (cart *Cart) GetTotal() float64 {
	var total float64
	for _, ci := range *cart {
		total += ci.Product.Price * float64(ci.Quantity)
	}
	return total + (total*1.08 - total) + cart.GetShippingCost()
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

func AddToCart(product *Product) error {
	// fmt.Println("Add to cart")
	// fmt.Print(cart)
	// stmt, err := db.DB.Prepare(`
	// 	insert into "CartItem" ("")
	// `)
	var cartItem CartItem
	if err := GetCartItem(product.Id, &cartItem); err != nil {
		return err
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

		err := cartItem.insert()
		if err != nil {
			return err
		}
	} else {
		cartItem.Quantity += 1
		cartItem.Total = cartItem.Product.Price * float64(cartItem.Quantity)

		err := cartItem.save()
		if err != nil {
			return err
		}
	}

	return nil
}

func GetCart(w http.ResponseWriter, r *http.Request) Cart {

	query := `
	    SELECT
	        ci."Id",
	        ci."ProductId",
			ci."Total",
	        ci."Quantity",
	        p."Name",
	        p."Description",
	        p."Category",
	        p."Image",
	        p."Price",
	        p."Quantity" as "ProductQuantity"
	    FROM
	        "CartItem" ci
	    JOIN
	        "Product" p ON ci."ProductId" = p."Id";
	`

	var recordsList []CartRecord

	err := db.DB.Select(&recordsList, query)
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
			Quantity:  r.Quantity,
		}

		cart = append(cart, *cartItem)
	}

	fmt.Printf("Cart: %v", cart)

	return cart

}

func GetCartItem(productId uuid.UUID, ci *CartItem) error {

	err := db.DB.Get(&ci, `select * from "CartItem" where "ProductId" = $1`, productId.String())
	if err != nil {
		return err
	}

	return nil
}
