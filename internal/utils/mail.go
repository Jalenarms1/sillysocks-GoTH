package utils

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"path/filepath"
	"slices"
	"text/template"

	"github.com/Jalenarms1/sillysocks-GoTH/internal/db"
)

func SendMail(toAddr string) error {

	from := "dev.test.jalen@gmail.com"
	to := "jalenarms@outlook.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		"This is a test"
	fmt.Println(os.Getenv("EMAIL_AP"))
	auth := smtp.PlainAuth("", "dev.test.jalen@gmail.com", os.Getenv("EMAIL_AP"), "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

type EmailCartItem struct {
	ProductPrice string
	Quantity     int32
	SubTotal     string
}

type EmailData struct {
	CartItems    []EmailCartItem
	Tax          float64
	Total        float64
	OrderId      string
	CustomerName string
}

func SendOrderPaidEmail(order *db.Order) error {
	from := "dev.test.jalen@gmail.com"

	subject := "Subject: Thank you for your order!\n"
	contentType := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	// 	body := `<html><body
	//     style="
	//       font-family: Arial, sans-serif;
	//       background-color: #f4f4f4;
	//     "
	//   >
	//     <table
	//       align="center"
	//       width="600"
	//       style="background-color: #ffffff; padding: 20px; border-radius: 8px;"
	//     >
	//       <!-- Header -->
	//       <tr>
	//         <td align="center">
	//           <h2 style="color: #333;">Thank You for Your Order!</h2>
	//           <p style="color: #666;">Your order has been confirmed.</p>
	//         </td>
	//       </tr>

	//       <!-- Product Table -->
	//       <tr>
	//         <td>
	//           <table style="border-collapse: collapse; width: 100%;">
	//             <tr>
	//               <th style="text-align: left; padding: 10px;">Product</th>
	//               <th style="text-align: left; padding: 10px;">Quantity</th>
	//               <th style="text-align: left; padding: 10px;">Price</th>
	//             </tr>

	//             <!-- Example Product (Repeat this row dynamically in your email template) -->
	//             <tr>
	//               <td style="padding: 10px; border-top: 1px solid #ddd;">
	//                 <img
	//                   src="PRODUCT_IMAGE_URL"
	//                   alt="Product Image"
	//                   style="width: 80px; height: auto; border-radius: 4px; display: block;"
	//                 />
	//                 <p style="margin: 5px 0; font-size: 14px;">Product Name</p>
	//               </td>
	//               <td style="padding: 10px; border-top: 1px solid #ddd;">2</td>
	//               <td style="padding: 10px; border-top: 1px solid #ddd;">$50.00</td>
	//             </tr>
	//             <!-- Repeat End -->

	//           </table>
	//         </td>
	//       </tr>

	//       <!-- Order Summary -->
	//       <tr>
	//         <td align="right" style="padding: 20px;">
	//           <h3 style="margin: 0;">Total: $TOTAL_AMOUNT</h3>
	//         </td>
	//       </tr>

	//       <!-- Order Link -->
	//       <tr>
	//         <td align="center" style="padding: 20px;">
	//           <a
	//             href="https://yourwebsite.com/orders/ORDER_ID"
	//             style="
	//               background-color: #28a745;
	//               color: white;
	//               padding: 12px 24px;
	//               text-decoration: none;
	//               font-size: 16px;
	//               border-radius: 5px;
	//               display: inline-block;
	//             "
	//             >View Order Details</a
	//           >
	//         </td>
	//       </tr>
	//     </table>
	//   </body></html>`

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	t, err := template.ParseFiles(filepath.Join(dir, "templates", "orderPaid.html"))
	if err != nil {
		return err
	}

	products := slices.Clone(order.CartItems)
	for i, _ := range products {
		products[i].ProductPrice = products[i].ProductPrice / 100
	}

	seq := func(yield func(cartItem EmailCartItem) bool) {
		for _, item := range order.CartItems {
			fmt.Println(item.ProductPrice / 100)
			emailCartItem := &EmailCartItem{
				ProductPrice: fmt.Sprintf("%.2f", float64(item.ProductPrice)/100),
				Quantity:     item.Quantity,
				SubTotal:     fmt.Sprintf("%.2f", float64(item.SubTotal)/100),
			}

			yield(*emailCartItem)
		}
	}

	emailData := &EmailData{
		CartItems: slices.Collect(seq),
		Total:     float64(order.GrandTotal) / 100,
	}

	var newBody bytes.Buffer
	err = t.Execute(&newBody, emailData)
	if err != nil {
		return err
	}

	msg := []byte(subject + contentType + newBody.String())

	auth := smtp.PlainAuth("", "dev.test.jalen@gmail.com", os.Getenv("EMAIL_AP"), "smtp.gmail.com")
	err = smtp.SendMail("smtp.gmail.com:587", auth, from, []string{*order.CustomerEmail}, msg)
	if err != nil {
		return err
	}

	return nil
}
