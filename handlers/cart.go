package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/Jalenarms1/sillysocks-GoTH/models"
	"github.com/Jalenarms1/sillysocks-GoTH/views/cartview"
	"github.com/Jalenarms1/sillysocks-GoTH/views/home"
	"github.com/Jalenarms1/sillysocks-GoTH/views/icons"
	"github.com/go-chi/chi/v5"
)

func handleAddToCart(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	productStr := r.FormValue("product")
	var product models.Product
	jsonErr := json.Unmarshal([]byte(productStr), &product)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Printf("\n%v\n", product)

	if err != nil {
		log.Fatal(err)
	}

	cart := models.GetCart(w, r)

	itemFound := false
	for i, ci := range cart.CartItems {
		if ci.Product.Id == product.Id {
			cart.CartItems[i].Quantity += 1
			cart.CartItems[i].Price = ci.Product.Price * float64(cart.CartItems[i].Quantity)
			cart.SubTotal = math.Round((cart.SubTotal+ci.Product.Price)*100) / 100
			cart.Tax = math.Round((cart.SubTotal*1.08-cart.SubTotal)*100) / 100
			itemFound = true

			fmt.Printf("%.2f", ci.Price)

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
	// var total float64
	// if cart == nil {
	// 	total = 0
	// } else {
	// 	total = cart.GetTotal()

	// }
	fmt.Fprint(w, fmt.Sprintf("$%.2f", cart.SubTotal))

	return nil
}

func handleCartPage(w http.ResponseWriter, r *http.Request) error {
	cart := models.GetCart(w, r)
	return Render(w, r, cartview.Page(cart))
}

func handleGetCartItems(w http.ResponseWriter, r *http.Request) error {
	cart := models.GetCart(w, r)
	return Render(w, r, cartview.CartItems(cart.CartItems))
}

func handleCartItemIncr(w http.ResponseWriter, r *http.Request) error {
	var product models.Product
	r.ParseForm()
	fmt.Printf("%s", r.FormValue("id"))
	product.Id = r.FormValue("id")
	price, cnvErr := strconv.ParseFloat(r.FormValue("price"), 64)
	if cnvErr != nil {
		log.Fatal(cnvErr)
		return cnvErr
	}
	product.Price = price
	cart := models.GetCart(w, r)

	for i, it := range cart.CartItems {
		if it.Product.Id == product.Id {
			cart.CartItems[i].Quantity += 1
			cart.CartItems[i].Price = float64(cart.CartItems[i].Quantity) * product.Price
			cart.SubTotal += product.Price
			cart.Tax = cart.SubTotal*1.08 - cart.SubTotal
		}
	}

	err := models.AddToCart(w, r, cart)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// loadCartCount
	w.Header().Set("HX-Trigger", `{"loadPriceList": "loadPriceList", "loadCartCount": "loadCartCount"}`)

	return Render(w, r, cartview.CartItems(cart.CartItems))
}
func handleCartItemDecr(w http.ResponseWriter, r *http.Request) error {
	var product models.Product

	jsonErr := json.NewDecoder(r.Body).Decode(&product)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return jsonErr
	}

	cart := models.GetCart(w, r)

	for i, it := range cart.CartItems {
		if it.Product.Id == product.Id {
			cart.CartItems[i].Quantity -= 1
			cart.CartItems[i].Price = float64(cart.CartItems[i].Quantity) * product.Price
			cart.SubTotal -= product.Price
			cart.Tax = cart.SubTotal*1.08 - cart.SubTotal
		}
	}

	err := models.AddToCart(w, r, cart)
	if err != nil {
		log.Fatal(err)
		return err
	}

	w.Header().Set("HX-Trigger", "loadPriceList")

	return Render(w, r, cartview.Page(cart))
}

func handleGetCartPriceList(w http.ResponseWriter, r *http.Request) error {
	cart := models.GetCart(w, r)
	fmt.Printf("%v", cart)
	return Render(w, r, cartview.PriceList(cart))
}

func RegisterCartRouter(router *chi.Mux) {
	router.Get("/cart", UseHTTPHandler(handleCartPage))
	router.Post("/api/cart/add", UseHTTPHandler(handleAddToCart))
	router.Get("/api/cart/count", UseHTTPHandler(handleGetCartCount))
	router.Put("/api/cart/count/increment", UseHTTPHandler(handleCartItemIncr))
	router.Put("/api/cart/count/decrement", UseHTTPHandler(handleCartItemDecr))
	router.Get("/api/cart/total", UseHTTPHandler(handleGetCartTotal))
	router.Get("/api/cart/items", UseHTTPHandler(handleGetCartItems))
	router.Get("/api/cart/price-list", UseHTTPHandler(handleGetCartPriceList))
}
