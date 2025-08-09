package types

import(

)
type LoginUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
}

type SignupUser struct {
	Name         string `json:"name"`
	MobileNumber int    `json:"mobile_number"`
	Email        string `json:"email"`
	Password     string `json:"password"`
}

type User struct {
	UserID       int    `json:"user_id"`
	Role         string `json:"role"`
	Name         string `json:"name"`
	MobileNumber int    `json:"mobile_number"`
	Email        string `json:"email"`
	Password     string `json:"password"`
}

type Item struct {
	ItemID   int     `json:"item_id"`
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Category string  `json:"category"`
	Img      string  `json:"img"`
}

type Order struct {
	OrderID             int    `json:"order_id"`
	TableNumber         int    `json:"table_number"`
	SpecificInstruction string `json:"specific_instruction"`
	OrderStatus         string `json:"order_status"`
	UserID              int    `json:"user_id"`
}

type OrderedItem struct {
	ID       int `json:"ID"`
	ItemID   int `json:"item_id"`
	Quantity int `json:"quantity"`
	OrderID  int `json:"order_id"`
}

type Payment struct {
	PaymentID int     `json:"payment_id"`
	OrderID   int     `json:"order_id"`
	Total     float32 `json:"total"`
	Mode      string  `json:"mode"`
	Status    bool    `json:"status"`
}

var CategoryList = []string{
	"breakfast",
	"beverages",
	"starters",
	"main course",
	"dessert"}

var OrderStatus = []string{
	"In Queue",
	"In Progress",
	"Completed"}

var Mode=[]string{
	"unselected",
	"Cash",
	"Card",
	"UPI",
	"Net Banking"}