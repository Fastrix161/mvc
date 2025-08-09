CREATE TABLE IF NOT EXISTS Orders(
order_id BIGINT PRIMARY KEY AUTO_INCREMENT,
table_number BIGINT NOT NULL,
specific_instruction TEXT,
order_status VARCHAR(255) NOT NULL,
user_id BIGINT,
FOREIGN KEY(user_id) REFERENCES User(user_id)
); 