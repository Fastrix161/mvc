package models

import (
	"database/sql"
	"fmt"

	"github.com/fastrix161/mvc/pkg/types"
)

func CreateUser(user types.User) (int, error) {
	query := "INSERT INTO User (role, name, mobile_number, email, password) VALUES (?, ?, ?, ?, ?)"
	result, err := DB.Exec(query, user.Role, user.Name, user.MobileNumber, user.Email, user.Password)
	if err != nil {
		return 0, fmt.Errorf("error adding user %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert id %v", err)
	}

	return int(id), nil
}

func CheckEmail(email string) (*types.User, error) {
	query := "SELECT * FROM User WHERE email = ?"
	row := DB.QueryRow(query, email)
	var u types.User
	err := row.Scan(&u.UserID, &u.Role, &u.Name, &u.MobileNumber, &u.Email, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("email not found")
		}
		return nil, fmt.Errorf("error during scanning %v", err)
	}
	return &u, nil
}

func DeleteUser(id int) error {
	query := "DELETE FROM User WHERE user_id = ?"
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting data %v", err)
	}
	return nil
}

func SetUserRole(id int, role string) error {
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM `User` WHERE user_id=?)", id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user_id %d does not exist", id)
	}

	query := "UPDATE User SET role = ? WHERE user_id = ?"
	result, err := DB.Exec(query, role, id)
	if err != nil {
		return fmt.Errorf("error changing role %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("same role")
	}
	return nil
}

func GetUser(id int) (*types.User, error) {
	query := "SELECT * FROM User WHERE user_id = ?"
	row := DB.QueryRow(query, id)
	var u types.User
	err := row.Scan(&u.UserID, &u.Role, &u.Name, &u.MobileNumber, &u.Email, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error during scanning %v", err)
	}
	return &u, nil
}

func GetAllUsers() ([]types.User, error) {
	query := "SELECT * FROM User"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error getting users %v", err)
	}
	defer rows.Close()
	var users []types.User
	for rows.Next() {
		var u types.User
		err := rows.Scan(&u.UserID, &u.Role, &u.Name, &u.MobileNumber, &u.Email, &u.Password)
		if err != nil {
			return nil, fmt.Errorf("error during scanning %v", err)
		}
		users = append(users, u)
	}
	return users, nil
}

func GetItems(search string) ([]types.Item, error) {
	query := "SELECT * FROM Item WHERE LOWER(name) LIKE CONCAT(LOWER(?),'%')"
	rows, err := DB.Query(query, search)
	if err != nil {
		return nil, fmt.Errorf("error getting items %v", err)
	}
	defer rows.Close()
	var items []types.Item
	for rows.Next() {
		var i types.Item
		err := rows.Scan(&i.ItemID, &i.Name, &i.Price, &i.Category, &i.Img)
		if err != nil {
			return nil, fmt.Errorf("error during scanning %v", err)
		}
		items = append(items, i)
	}
	return items, nil
}

func GetAllItems() ([]types.Item, error) {
	query := "SELECT * FROM Item ORDER BY item_id ASC"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error getting items %v", err)
	}
	defer rows.Close()
	var items []types.Item
	for rows.Next() {
		var i types.Item
		err := rows.Scan(&i.ItemID, &i.Name, &i.Price, &i.Category, &i.Img)
		if err != nil {
			return nil, fmt.Errorf("error during scanning %v", err)
		}
		items = append(items, i)
	}
	return items, nil
}

func GetCategoryItems(category string) ([]types.Item, error) {
	var check bool = false
	for _, item := range types.CategoryList {
		if item == category {
			check = true
		}
	}
	if !check {
		return nil, nil
	}
	query := "SELECT * FROM Item WHERE category = ?"

	rows, err := DB.Query(query, category)
	if err != nil {
		return nil, fmt.Errorf("error getting items %v", err)
	}
	defer rows.Close()
	var items []types.Item
	for rows.Next() {
		var i types.Item
		err := rows.Scan(&i.ItemID, &i.Name, &i.Price, &i.Category, &i.Img)
		if err != nil {
			return nil, fmt.Errorf("error during scanning %v", err)
		}
		items = append(items, i)
	}
	return items, nil
}

func GetOrder(id int) (*types.Order, error) {
	query := "SELECT * FROM Orders WHERE Orders.order_id = ?"
	row := DB.QueryRow(query, id)
	var o types.Order
	err := row.Scan(&o.OrderID, &o.TableNumber, &o.SpecificInstruction, &o.OrderStatus, &o.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("error during scanning order %v", err)
	}
	return &o, nil
}

func GetOrdersForUser(id int) ([]types.Order, error) {
	query := "SELECT * FROM Orders WHERE Orders.user_id = ?"
	rows, err := DB.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("error getting orders for id %v : %v", id, err)
	}
	defer rows.Close()
	var orders []types.Order
	for rows.Next() {
		var o types.Order
		err := rows.Scan(&o.OrderID, &o.TableNumber, &o.SpecificInstruction, &o.OrderStatus, &o.UserID)
		if err != nil {
			return nil, fmt.Errorf("error during scanning %v", err)
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func GetAllOrders2() ([]types.Order, error) {
	query := "SELECT O.* FROM Orders O JOIN Ordered_items Oi ON O.order_id = Oi.order_id ORDER BY Orders.order_id DESC"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error getting orders: %v", err)
	}
	defer rows.Close()
	var orders []types.Order
	for rows.Next() {
		var o types.Order
		err := rows.Scan(&o.OrderID, &o.SpecificInstruction, &o.OrderStatus, &o.TableNumber, &o.UserID)
		if err != nil {
			return nil, fmt.Errorf("error during scanning orders%v", err)
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func GetAllOrders() ([]types.Order, error) {
	query := "SELECT * FROM Orders ORDER BY Orders.order_id DESC"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error getting orders: %v", err)
	}
	defer rows.Close()
	var orders []types.Order
	for rows.Next() {
		var o types.Order
		err := rows.Scan(&o.OrderID, &o.TableNumber, &o.SpecificInstruction, &o.OrderStatus, &o.UserID)
		if err != nil {
			return nil, fmt.Errorf("error during scanning orders%v", err)
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func AddOrder(order types.Order) (int, error) {
	query := "INSERT INTO Orders (table_number, specific_instruction, order_status, user_id) VALUES (?, ?, ?, ?)"
	result, err := DB.Exec(query, order.TableNumber, order.SpecificInstruction, order.OrderStatus, order.UserID)
	if err != nil {
		return 0, fmt.Errorf("error adding order %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert order_id %v", err)
	}
	return int(id), nil
}

func DeleteOrder(id int) error {
	query := "DELETE FROM Orders WHERE order_id = ?"
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting order data %v", err)
	}
	return nil
}

func UpdateOrder(order types.Order) error {
	query := "UPDATE Orders SET specific_instruction = ?, order_status = ? WHERE order_id = ?"
	result, err := DB.Exec(query, order.SpecificInstruction, order.OrderStatus, order.OrderID)
	if err != nil {
		return fmt.Errorf("error updating order %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}

func AddOrderedItem(oi types.OrderedItem) (int, error) {
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM `Orders` WHERE order_id=?)", oi.OrderID).Scan(&exists)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, fmt.Errorf("order_id %d does not exist", oi.OrderID)
	}

	query := "INSERT INTO Ordered_items (order_id, item_id, quantity) VALUES (?, ?, ?)"
	result, err := DB.Exec(query, oi.OrderID, oi.ItemID, oi.Quantity)
	if err != nil {
		return 0, fmt.Errorf("error adding ordered item %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last inserted ordered item id %v", err)
	}
	return int(id), nil
}

func UpdateOrderedItems(oi types.OrderedItem) error {
	query := "UPDATE Ordered_items SET quantity = ? WHERE order_id = ? AND item_id = ?"
	result, err := DB.Exec(query, oi.Quantity, oi.OrderID, oi.ItemID)
	if err != nil {
		return fmt.Errorf("error updating ordered item %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no changes to be made")
	}
	return nil
}

func GetOrderedItems(id int) ([]types.OrderedItems, error) {
	query := "SELECT ID,Item.item_id,quantity,order_id,name,price,category FROM Ordered_items JOIN Item ON Ordered_items.item_id = Item.item_id WHERE order_id = ?"
	rows, err := DB.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("error getting ordered items: %v", err)
	}
	defer rows.Close()
	var ois []types.OrderedItems
	for rows.Next() {
		var oi types.OrderedItems
		err := rows.Scan(&oi.ID, &oi.ItemID, &oi.Quantity, &oi.OrderID, &oi.ItemName, &oi.Price, &oi.Category)
		if err != nil {
			return nil, fmt.Errorf("error during scanning ordered items %v", err)
		}
		ois = append(ois, oi)
	}
	return ois, nil
}

func DeleteOrderedItem(oi types.OrderedItem) error {
	query := "DELETE FROM Ordered_items WHERE order_id = ? AND item_id = ?"
	_, err := DB.Exec(query, oi.OrderID, oi.ItemID)
	if err != nil {
		return fmt.Errorf("error deleting ordered item %v", err)
	}
	return nil
}

func CreatePayment(pay types.Payment) (int, error) {
	query := "INSERT INTO payment (order_id, total, mode, status) VALUES (?,?,?,?)"
	result, err := DB.Exec(query, pay.OrderID, pay.Total, pay.Mode, pay.Status)
	if err != nil {
		return 0, fmt.Errorf("error creating payment %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last inserted payment id %v", err)
	}
	return int(id), nil
}

func GetPayment(id int) (*types.Payment, error) {
	query := "SELECT * FROM payment WHERE payment_id = ?"
	row := DB.QueryRow(query, id)
	var pay types.Payment
	err := row.Scan(&pay.PaymentID, &pay.OrderID, &pay.Total, &pay.Mode, &pay.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error during scanning %v", err)
	}
	return &pay, nil
}

func GetPaymentforOrder(id int) (*types.Payment, error) {
	query := "SELECT * FROM payment WHERE order_id = ?"
	row := DB.QueryRow(query, id)
	var pay types.Payment
	err := row.Scan(&pay.PaymentID, &pay.OrderID, &pay.Total, &pay.Mode, &pay.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error during scanning %v", err)
	}
	return &pay, nil
}

func GetAllPayments() ([]types.Payment, error) {
	query := "SELECT * FROM payment ORDER BY payment.payment_id DESC"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error getting payment datas: %v", err)
	}
	defer rows.Close()
	var pays []types.Payment
	for rows.Next() {
		var pay types.Payment
		err := rows.Scan(&pay.PaymentID, &pay.Mode, &pay.OrderID, &pay.Total, &pay.Status)
		if err != nil {
			return nil, fmt.Errorf("error during scanning orders%v", err)
		}
		pays = append(pays, pay)
	}
	return pays, nil
}

func DeletePayment(pay types.Payment) error {
	query := "DELETE FROM payment WHERE payment_id = ?"
	_, err := DB.Exec(query, pay.PaymentID)
	if err != nil {
		return fmt.Errorf("error deleting ordered item %v", err)
	}
	return nil
}
func UpdatePayment(pay types.Payment) (bool, error) {
	query := "UPDATE payment SET status = ?,mode =? WHERE payment_id = ?"
	result, err := DB.Exec(query, pay.Status, pay.Mode, pay.PaymentID)
	if err != nil {
		return false, fmt.Errorf("error updating payment %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error getting rows affected %v", err)
	}
	if rowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func ItemExistsInOrder(itemId int, orderId int) (bool, error) {
	query := "SELECT * FROM Ordered_items WHERE order_id = ? AND item_Id = ?"
	row := DB.QueryRow(query, orderId, itemId)
	var oi types.OrderedItem
	err := row.Scan(&oi.ID, &oi.OrderID, &oi.Quantity, &oi.ItemID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("error during scanning %v", err)
	}
	return true, nil
}

func GetOrderTotal(orderId int) (float32, error) {
	query := `
			SELECT COALESCE(SUM(oi.quantity * i.price), 0) AS total
			  FROM Ordered_items oi
			  JOIN Item i
			  ON oi.item_id = i.item_id
			  WHERE oi.order_id = ?
			  `
	var total float32
	err := DB.QueryRow(query, orderId).Scan(&total)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("GetOrderTotal: %v", err)
	}
	return total, nil
}
