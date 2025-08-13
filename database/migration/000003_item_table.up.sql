CREATE TABLE IF NOT EXISTS Item(
item_id BIGINT PRIMARY KEY AUTO_INCREMENT,
name VARCHAR(255) NOT NULL,
price FLOAT NOT NULL, 
category ENUM('breakfast','beverages','starters','main course','dessert') NOT NULL,
img TEXT NOT NULL
);