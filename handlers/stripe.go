package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Jalenarms1/sillysocks-GoTH/models"
	"github.com/go-chi/chi/v5"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
)

func RegisterStripeRouter(router *chi.Mux) {
	stripe.Key = os.Getenv("STRIPE_SECRET")
	router.Post("/api/stripe/checkout", UseHTTPHandler(handleCreateCheckout))
	router.Post("/api/stripe/webhook", UseHTTPHandler(handleStripeWebhook))
}

func handleCreateCheckout(w http.ResponseWriter, r *http.Request) error {
	domain := "http://localhost:3000"
	r.ParseForm()
	// var cart models.Cart

	// err := json.Unmarshal([]byte(r.FormValue("cart-json")), &cart)
	cart := models.GetCart(w, r)
	// fmt.Printf("\n%v\n", cart)

	var lineItems []*stripe.CheckoutSessionLineItemParams
	var total int64
	for _, i := range cart {
		priceInCents := i.Product.Price * 100
		total += int64(priceInCents)
		li := &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency:   stripe.String("usd"),
				UnitAmount: stripe.Int64(int64(priceInCents)),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name:   stripe.String(i.Product.Name),
					Images: stripe.StringSlice([]string{i.Product.Image}),
				},
			},
			Quantity: stripe.Int64(i.Quantity),
		}

		lineItems = append(lineItems, li)
	}
	taxItem := &stripe.CheckoutSessionLineItemParams{
		PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
			Currency:   stripe.String("usd"),
			UnitAmount: stripe.Int64(int64(cart.GetTax() * 100)),
			ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
				Name: stripe.String("Tax"),
			},
		},
		Quantity: stripe.Int64(1),
	}
	lineItems = append(lineItems, taxItem)

	if cart.GetSubTotal() < 20 {
		shippingItem := &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency:   stripe.String("usd"),
				UnitAmount: stripe.Int64(500),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String("Shipping"),
				},
			},
			Quantity: stripe.Int64(1),
		}

		lineItems = append(lineItems, shippingItem)

	}
	// fmt.Printf("%v", lineItems)
	params := &stripe.CheckoutSessionParams{
		LineItems:  lineItems,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(domain + "/cart?orderStatus=success"),
		CancelURL:  stripe.String(domain + "/cart?orderStatus=canceled"),
		ShippingAddressCollection: &stripe.CheckoutSessionShippingAddressCollectionParams{
			AllowedCountries: stripe.StringSlice([]string{"US"}),
		},
	}
	session, err := session.New(params)

	if err != nil {
		return err
	}
	w.Header().Set("HX-Redirect", session.URL)
	return nil
}

func handleStripeWebhook(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	event := &stripe.Event{}

	if umErr := json.Unmarshal(body, &event); umErr != nil {
		return umErr
	}
	cart := models.GetCart(w, r)
	cstDt := event.Data.Object["customer_details"].(map[string]interface{})
	shpDt := event.Data.Object["shipping_details"].(map[string]interface{})["address"].(map[string]interface{})

	cstEmail := cstDt["email"].(string)
	cstName := cstDt["name"].(string)
	cstId := event.Data.Object["id"].(string)

	shpAddr1 := shpDt["line1"].(string)
	shpAddr2 := shpDt["line2"]
	if shpAddr2 == nil {
		shpAddr2 = ""
	}
	shpAddrCity := shpDt["city"].(string)
	shpAddrState := shpDt["state"].(string)
	shpAddrZip := shpDt["postal_code"].(string)
	// pmtIntId := event.Data.Object["payment_intent"]

	order := models.NewOrder(cstId, "", cstEmail, cstName, shpAddr1, shpAddr2.(string), shpAddrCity, shpAddrState, shpAddrZip, cart.GetSubTotal(), cart.GetTax(), cart.GetShippingCost())
	fmt.Printf("%v\n", event.Data.Object["payment_intent"])
	fmt.Printf("%v", order)

	return nil
}
