package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Jalenarms1/sillysocks-GoTH/internal/db"
	"github.com/Jalenarms1/sillysocks-GoTH/internal/utils"
	"github.com/gofrs/uuid"
	stripe "github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/webhook"
)

type CheckoutSessionReqParams struct {
	CartItems []db.CartItem
}

func HandleCreateCheckoutSession(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return errors.New("method not allowed")
	}

	// body, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(w, "Invalid data", http.StatusBadRequest)
	// 	return
	// }
	var cartItemData CheckoutSessionReqParams

	err := json.NewDecoder(r.Body).Decode(&cartItemData)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return err
	}

	var total int64
	var lineItems []*stripe.CheckoutSessionLineItemParams
	for _, item := range cartItemData.CartItems {
		total += item.Product.Price * int64(item.Quantity)
		lineItem := stripe.CheckoutSessionLineItemParams{
			Quantity: stripe.Int64(int64(item.Quantity)),
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name:   &item.Product.Name,
					Images: stripe.StringSlice([]string{*item.Product.Image})},
				UnitAmount: stripe.Int64(item.Product.Price),
				Currency:   stripe.String("usd"),
			},
		}

		lineItems = append(lineItems, &lineItem)

	}

	shippingItem := stripe.CheckoutSessionLineItemParams{
		Quantity: stripe.Int64(int64(1)),
		PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
			ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
				Name: stripe.String("Shipping"),
			},
			UnitAmount: stripe.Int64(int64(500)),
			Currency:   stripe.String("usd"),
		},
	}
	fmt.Println(total)
	fmt.Println(float64(total) / 100)
	fmt.Println(int64((float64(total/100) * 1.08) * 100))
	val, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", ((float64(total)/100)*1.08)), 64)
	fmt.Println(val)
	tax := int64(val*100) - total
	fmt.Println(tax)
	taxItem := stripe.CheckoutSessionLineItemParams{
		Quantity: stripe.Int64(int64(1)),
		PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
			ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
				Name: stripe.String("Tax"),
			},
			UnitAmount: stripe.Int64(int64(tax)),
			Currency:   stripe.String("usd"),
		},
	}

	lineItems = append(lineItems, &shippingItem)
	lineItems = append(lineItems, &taxItem)

	uid, _ := uuid.NewV4()
	order := db.Order{
		Id:            uid.String(),
		SubTotal:      total,
		Tax:           tax,
		GrandTotal:    total + tax + 500,
		ShippingTotal: 500,
		CreatedAt:     time.Now().Unix(),
		Status:        "Unpaid",
	}

	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	err = order.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = db.InsertCartItems(tx, cartItemData.CartItems, order.Id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	fmt.Printf("%s/cart", os.Getenv("CLIENT_DOMAIN"))
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String("payment"),
		LineItems:          lineItems,
		ShippingAddressCollection: &stripe.CheckoutSessionShippingAddressCollectionParams{
			AllowedCountries: stripe.StringSlice([]string{"US"}),
		},
		SuccessURL: stripe.String(fmt.Sprintf("%s/order/", os.Getenv("CLIENT_DOMAIN")) + order.Id + "?resetCart=true"),
		CancelURL:  stripe.String(fmt.Sprintf("%s/cart", os.Getenv("CLIENT_DOMAIN"))),
		Metadata: map[string]string{
			"orderId": order.Id,
		},
	}

	session, err := session.New(params)
	if err != nil {
		_ = tx.Rollback()

		return err
	}

	err = tx.Commit()
	if err != nil {
		return errors.New("Error committing the order transaction for order id, " + order.Id)
	}

	// fmt.Println(session.URL)

	resp := map[string]string{
		"sessionUrl": session.URL,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	return nil

}

func HandleCheckoutSessionWebhook(w http.ResponseWriter, r *http.Request) error {

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return err
	}

	stripeSig := r.Header.Get("Stripe-Signature")
	fmt.Println(stripeSig)
	event, err := webhook.ConstructEventWithOptions(body, stripeSig, os.Getenv("STRIPE_WHKEY"), webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true})
	if err != nil {
		return err
	}
	fmt.Print(event.Type)
	fmt.Println(stripe.EventTypeCheckoutSessionCompleted)
	if event.Type == stripe.EventTypeCheckoutSessionCompleted {
		fmt.Print(event.Data.Object["metadata"].(map[string]interface{})["orderId"])
		status := event.Data.Object["payment_status"].(string)
		fmt.Println(status)
		fmt.Println(stripe.CheckoutSessionPaymentStatusPaid)
		if status == string(stripe.CheckoutSessionPaymentStatusPaid) {
			orderId := event.Data.Object["metadata"].(map[string]interface{})["orderId"].(string)
			paymentIntent := event.Data.Object["payment_intent"].(string)
			customerDetails := event.Data.Object["customer_details"]
			address := customerDetails.(map[string]interface{})["address"].(map[string]interface{})
			line1 := address["line1"].(string)

			city := address["city"].(string)
			state := address["state"].(string)
			postalCode := address["postal_code"].(string)
			name := customerDetails.(map[string]interface{})["name"].(string)
			email := customerDetails.(map[string]interface{})["email"].(string)

			fmt.Println("Existing Order")
			existingOrder, err := db.GetOrder(orderId)
			if err != nil {
				return err
			}
			fmt.Println(existingOrder)
			if existingOrder == nil {
				return errors.New("order not found " + orderId)
			}

			existingOrder.PaymentIntentId = &paymentIntent
			existingOrder.ShippingLine1 = &line1
			line2, ok := address["line2"].(string)
			if ok {
				existingOrder.ShippingLine2 = &line2
			}
			existingOrder.ShippingCity = &city
			existingOrder.ShippingState = &state
			existingOrder.ShippingPostalCode = &postalCode
			existingOrder.Status = "Paid"
			existingOrder.CustomerName = &name
			existingOrder.CustomerEmail = &email

			err = existingOrder.Save()
			if err != nil {
				return err
			}

		}

	}

	err = utils.SendMail("")
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
