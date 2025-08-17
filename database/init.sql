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


-- Menu
INSERT INTO Item (name, price, category, img) VALUES
    ('Paneer Tikka', '200','starters','images/placeHolder1.jpg'),
    ('Bread Omelete', '100','breakfast','images/placeHolder1.jpg'),
    ('Malai Kofte', '400','main course','images/placeHolder2.jpg'),
    ('Hot Coffee', '60','beverages','images/placeHolder2.jpg'),
    ('Chocolate Brownie', '120','dessert','images/placeHolder3.jpg'),
    ('Spring Rolls','200','starters','images/placeHolder1.jpg'),
    ('PBJ sandwitch','80','breakfast','images/placeHolder3.jpg'),
    ('Paneer Lababdaar','400','main course','images/placeHolder1.jpg'),
    ('Mojito','140','beverages','images/placeHolder2.jpg'),
    ('Ice Cream (Vanilla)','70','dessert','images/placeHolder3.jpg');
    
    -- User
    -- password: PasswordAdmin
    -- password: PrimeChef
    -- SaltRound used: 10
    INSERT INTO User (name, email, role, mobile_number, password) VALUES
    ('Admin','admin@admin.com','admin','1010101010', '$2a$12$mTbQ9lM1YJK/DjT1/PBcoullmeBpzKC2ABlf22/OUPvgY5XHrAUVS'), 
    ('Chef','chef@gmail.com','chef','1231231230', '$2a$10$jgLHoHhS9WijLymbX.direU0MMMrzYmgxpQbSW6vF8vlwbGqGHrhy'); 

    