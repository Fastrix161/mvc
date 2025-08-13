DROP DATABASE IF EXISTS restaurant;
CREATE DATABASE restaurant;
USE restaurant;

CREATE TABLE User(
user_id BIGINT PRIMARY KEY AUTO_INCREMENT,
role ENUM('customer', 'chef', 'admin') NOT NULL,
name VARCHAR(255) NOT NULL,
mobile_number BIGINT,
email VARCHAR(255) NOT NULL UNIQUE,
password VARCHAR(255) NOT NULL
);

CREATE TABLE Orders(
order_id BIGINT PRIMARY KEY AUTO_INCREMENT,
table_number BIGINT NOT NULL,
specific_instruction TEXT,
order_status VARCHAR(255) NOT NULL,
user_id BIGINT,
FOREIGN KEY(user_id) REFERENCES User(user_id)
); 

CREATE TABLE Item(
item_id BIGINT PRIMARY KEY AUTO_INCREMENT,
name VARCHAR(255) NOT NULL,
price FLOAT NOT NULL, 
category ENUM('breakfast','beverages','starters','main course','dessert') NOT NULL,
img TEXT NOT NULL
);

CREATE TABLE Ordered_items(
ID BIGINT PRIMARY KEY AUTO_INCREMENT,
item_id BIGINT NOT NULL,
quantity BIGINT NOT NULL,
order_id BIGINT NOT NULL,
FOREIGN KEY(item_id) REFERENCES Item(item_id),
FOREIGN KEY(order_id) REFERENCES Orders(order_id)
);

CREATE TABLE payment(
payment_id BIGINT PRIMARY KEY AUTO_INCREMENT,
order_id BIGINT NOT NULL,
total FLOAT NOT NULL,
mode VARCHAR(255) NOT NULL,
status BOOLEAN NOT NULL,
FOREIGN KEY(order_id) REFERENCES Orders(order_id)
);
