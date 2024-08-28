package models

import (
	"log"
	"net/http"
)

type CartItem struct {
	Product  Product `json:"product"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type Cart struct {
	SubTotal  float64    `json:"subTotal"`
	Tax       float64    `json:"tax"`
	CartItems []CartItem `json:"cartItems"`
}

func (cart *Cart) NumOfItems() int {
	count := 0

	for _, ci := range cart.CartItems {
		count += ci.Quantity
	}

	return count
}

func (cart *Cart) GetTotal() float64 {
	return cart.SubTotal + cart.Tax + cart.GetShippingCost()
}

func (cart *Cart) GetShippingCost() float64 {
	if cart.SubTotal >= 20 {
		return 0
	}

	return 5
}

func AddToCart(w http.ResponseWriter, r *http.Request, cart *Cart) error {
	// fmt.Println("Add to cart")
	// fmt.Print(cart)
	session := GetSessionValues(r)
	session.Values["cart"] = cart
	err := session.Save(r, w)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func GetCart(w http.ResponseWriter, r *http.Request) *Cart {
	session := GetSessionValues(r)
	var cart *Cart
	if cart, ok := session.Values["cart"].(*Cart); ok && cart != nil {
		return cart
	}
	if cart == nil {
		cart = &Cart{}
	}
	return cart

}
