package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Jalenarms1/sillysocks-GoTH/models"
	"github.com/Jalenarms1/sillysocks-GoTH/views/cartview"
	"github.com/Jalenarms1/sillysocks-GoTH/views/home"
	"github.com/Jalenarms1/sillysocks-GoTH/views/icons"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

func RegisterCartRouter(router *chi.Mux) {
	router.Get("/cart", UseHTTPHandler(handleCartPage))
	router.Post("/api/cart/add", UseHTTPHandler(handleAddToCart))
	router.Get("/api/cart/count", UseHTTPHandler(handleGetCartCount))
	router.Put("/api/cart/count/increment", UseHTTPHandler(handleCartItemIncr))
	router.Put("/api/cart/count/decrement", UseHTTPHandler(handleCartItemDecr))
	router.Get("/api/cart/total", UseHTTPHandler(handleGetCartTotal))
	router.Get("/api/cart/items", UseHTTPHandler(handleGetCartItems))
	router.Put("/api/cart/items/delete", UseHTTPHandler(handleDeleteCartItems))
	router.Get("/api/cart/price-list", UseHTTPHandler(handleGetCartPriceList))
}

func handleAddToCart(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	productStr := r.FormValue("product")
	var product models.Product
	jsonErr := json.Unmarshal([]byte(productStr), &product)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	if err != nil {
		log.Fatal(err)
	}

	cart := models.GetCart(w, r)

	itemFound := false
	for i, ci := range cart {
		if ci.ProductId == product.Id {
			cart[i].Quantity += 1
			cart[i].Total = ci.Product.Price * float64(cart[i].Quantity)
			itemFound = true

			break
		}
	}

	if !itemFound {
		cartItem := models.CartItem{
			Product:   product,
			ProductId: product.Id,
			Total:     product.Price,
			Quantity:  1,
		}

		cart = append(cart, cartItem)
	}

	addErr := models.AddToCart(&product)
	if addErr != nil {
		http.Error(w, addErr.Error(), http.StatusBadRequest)
	}

	w.Header().Set("HX-Trigger", "loadCartCount")

	return Render(w, r, home.AddToCartBtn(product))
}

func handleGetCartCount(w http.ResponseWriter, r *http.Request) error {
	cart := models.GetCart(w, r)
	var numOfItems int
	if cart == nil {
		numOfItems = 0
	} else {
		numOfItems = cart.NumOfItems()
	}
	return Render(w, r, icons.CartIcon(numOfItems))

}

func handleGetCartTotal(w http.ResponseWriter, r *http.Request) error {
	cart := models.GetCart(w, r)

	fmt.Fprint(w, fmt.Sprintf("$%.2f", cart.GetSubTotal()))

	return nil
}

func handleCartPage(w http.ResponseWriter, r *http.Request) error {
	cart := models.GetCart(w, r)

	orderStatus := r.URL.Query().Get("orderStatus")

	return Render(w, r, cartview.Page(cart, orderStatus))
}

func handleGetCartItems(w http.ResponseWriter, r *http.Request) error {
	cart := models.GetCart(w, r)
	return Render(w, r, cartview.CartItems(cart))
}

func handleCartItemIncr(w http.ResponseWriter, r *http.Request) error {
	var cartItem models.CartItem
	var cart models.Cart
	r.ParseForm()
	cartData := r.FormValue("data")
	fmt.Printf("\nCartFormData:\n%s\n", cartData)
	idToUpd := uuid.FromStringOrNil(r.URL.Query().Get("cartItemId"))
	if err := models.GetCartItem(idToUpd, &cartItem); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(cartData), &cart); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	for i, ci := range cart {
		if ci.Id == idToUpd {
			cart[i].Quantity++
			cart[i].Total = cart[i].Product.Price * float64(cart[i].Quantity)
		}
	}
	err := cart.Save()
	if err != nil {
		return err
	}

	w.Header().Set("HX-Trigger", `{"loadPriceList": "loadPriceList", "loadCartCount": "loadCartCount"}`)
	return Render(w, r, cartview.CartItems(cart))
}

func handleCartItemDecr(w http.ResponseWriter, r *http.Request) error {
	var cartItem models.CartItem
	var cart models.Cart
	r.ParseForm()
	cartData := r.FormValue("data")
	idToUpd := uuid.FromStringOrNil(r.URL.Query().Get("cartItemId"))
	if err := models.GetCartItem(idToUpd, &cartItem); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(cartData), &cart); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var newCart models.Cart
	for i, ci := range cart {
		if ci.Id == idToUpd {
			cart[i].Quantity--
			cart[i].Total = cart[i].Product.Price * float64(cart[i].Quantity)

			if cart[i].Quantity > 0 {
				newCart = append(newCart, cart[i])
			}
		} else {

			newCart = append(newCart, ci)
		}

	}
	cart = newCart
	err := cart.Save()
	if err != nil {
		return err
	}

	w.Header().Set("HX-Trigger", `{"loadPriceList": "loadPriceList", "loadCartCount": "loadCartCount"}`)
	return Render(w, r, cartview.CartItems(cart))
}

func handleGetCartPriceList(w http.ResponseWriter, r *http.Request) error {
	cart := models.GetCart(w, r)
	return Render(w, r, cartview.PriceList(cart))
}

func handleDeleteCartItems(w http.ResponseWriter, r *http.Request) error {
	cart := models.GetCart(w, r)
	r.ParseForm()
	cartItemIds := strings.Split(r.FormValue("cartItemIds")[1:], ",")
	var newCartItems models.Cart
	for _, ci := range cart {
		idFound := false
		for _, id := range cartItemIds {
			uid, _ := uuid.FromString(id)
			if uid == ci.Id {
				idFound = true
			}
		}

		if !idFound {
			newCartItems = append(newCartItems, ci)

		} else {
		}
	}

	cart = newCartItems

	err := cart.Save()
	if err != nil {
		return err
	}

	// models.AddToCart(w, r, &cart)

	w.Header().Set("HX-Trigger", `{"loadPriceList": "", "loadCartCount": ""}`)

	return Render(w, r, cartview.CartItems(cart))
}
