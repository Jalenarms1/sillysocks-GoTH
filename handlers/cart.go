package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/Jalenarms1/sillysocks-GoTH/models"
	"github.com/Jalenarms1/sillysocks-GoTH/views/cartview"
	"github.com/Jalenarms1/sillysocks-GoTH/views/home"
	"github.com/Jalenarms1/sillysocks-GoTH/views/icons"
	"github.com/go-chi/chi/v5"
)

func handleAddToCart(w http.ResponseWriter, r *http.Request) error {
	// err := r.ParseForm()
	// productId := r.FormValue("product")
	// priceStr := r.FormValue("price")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// price, _ := strconv.ParseFloat(priceStr, 64)
	var product models.Product

	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nProduct: %v", product)
	cart := models.GetCart(w, r)

	itemFound := false
	for i, ci := range cart.CartItems {
		if ci.Product.Id == product.Id {
			cart.CartItems[i].Quantity += 1
			cart.SubTotal = math.Round((cart.SubTotal+ci.Price)*100) / 100
			cart.Tax = math.Round((cart.SubTotal*1.08-cart.SubTotal)*100) / 100
			itemFound = true
			break
		}
	}

	if !itemFound {
		fmt.Println("Item not found")
		cartItem := models.CartItem{
			Product:  product,
			Price:    product.Price,
			Quantity: 1,
		}
		cart.SubTotal = math.Round((cart.SubTotal+product.Price)*100) / 100
		cart.Tax = math.Round((cart.SubTotal*1.08-cart.SubTotal)*100) / 100

		cart.CartItems = append(cart.CartItems, cartItem)
	}

	addErr := models.AddToCart(w, r, cart)
	if addErr != nil {
		log.Fatal(addErr)
	}

	// fmt.Printf("%s\n", product.Id)
	// fmt.Printf("%f\n", cart.SubTotal)
	// fmt.Printf("%.2f\n", product.Price)

	w.Header().Set("HX-Trigger", "loadCartCount")

	return Render(w, r, home.AddToCartBtn(product))
}

func handleGetCartCount(w http.ResponseWriter, r *http.Request) error {
	cart := models.GetCart(w, r)

	return Render(w, r, icons.CartIcon(cart.NumOfItems()))

}

func handleGetCartTotal(w http.ResponseWriter, r *http.Request) error {
	cart := models.GetCart(w, r)
	total := fmt.Sprintf("$%.2f", cart.GetTotal())
	fmt.Fprint(w, total)

	return nil
}

func handleCartPage(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, cartview.Page())
}

func handleGetCartItems(w http.ResponseWriter, r *http.Request) error {
	cart := models.GetCart(w, r)

	return Render(w, r, cartview.CartItems(cart.CartItems))
}

func RegisterCartRouter(router *chi.Mux) {
	router.Get("/cart", UseHTTPHandler(handleCartPage))
	router.Post("/api/cart/add", UseHTTPHandler(handleAddToCart))
	router.Get("/api/cart/count", UseHTTPHandler(handleGetCartCount))
	router.Get("/api/cart/total", UseHTTPHandler(handleGetCartTotal))
	router.Get("/api/cart/items", UseHTTPHandler(handleGetCartItems))
}
