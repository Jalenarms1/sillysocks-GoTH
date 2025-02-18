package db

type Product struct {
	Id          string  `db:"Id" json:"id"`
	Name        string  `db:"Name" json:"name"`
	Description *string `db:"Description" json:"description"`
	Category    *string `db:"Category" json:"category"`
	Image       *string `db:"Image" json:"image"`
	Price       int64   `db:"Price" json:"price"`
	Quantity    int     `db:"Quantity" json:"quantity"`
	Sizes       *string `db:"Sizes" json:"sizes"`
}

type CartItem struct {
	Id           string  `db:"Id" json:"id"`
	ProductId    string  `db:"ProductId" json:"productId"`
	ProductName  string  `db:"ProductName" json:"productName"`
	ProductImage string  `db:"ProductImage" json:"productImage"`
	ProductPrice string  `db:"ProductPrice" json:"productPrice"`
	Product      Product `db:"Product" json:"product,omitempty"`
	OrderId      *string `db:"OrderId" json:"orderId"`
	Quantity     int32   `db:"Quantity" json:"quantity"`
	SubTotal     int64   `db:"SubTotal" json:"subTotal"`
	Size         string  `db:"Size" json:"size"`
}

type Order struct {
	Id                 string     `db:"Id" json:"id"`
	PaymentIntentId    *string    `db:"PaymentIntentId" json:"paymentIntentId"`
	CartItems          []CartItem `db:"CartItems" json:"cartItems,omitempty"`
	SubTotal           int64      `db:"SubTotal" json:"subTotal"`
	Tax                int64      `db:"Tax" json:"tax"`
	GrandTotal         int64      `db:"GrandTotal" json:"grandTotal"`
	ShippingTotal      int64      `db:"ShippingTotal" json:"shippingTotal"`
	ShippingLine1      *string    `db:"ShippingLine1" json:"shippingLine1"`
	ShippingLine2      *string    `db:"ShippingLine2" json:"shippingLine2"`
	ShippingCity       *string    `db:"ShippingCity" json:"shippingCity"`
	ShippingState      *string    `db:"ShippingState" json:"shippingState"`
	ShippingPostalCode *string    `db:"ShippingPostalCode" json:"shippingPostalCode"`
	CustomerName       *string    `db:"CustomerName" json:"customerName"`
	CustomerEmail      *string    `db:"CustomerEmail" json:"customerEmail"`
	CreatedAt          int64      `db:"CreatedAt" json:"createdAt"`
	Status             string     `db:"Status" json:"status"`
}
