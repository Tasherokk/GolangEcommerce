// database/database.go
package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
    connStr := "user=postgres password=Taukent2004 dbname=ecommerce sslmode=disable"

    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Println("Failed to connect to database:", err)
        os.Exit(1)
    }

    err = DB.Ping()
    if err != nil {
        log.Println("Failed to ping database:", err)
        os.Exit(1)
    }

    log.Println("Successfully connected to database!")
}



// CREATE TABLE Users (
//     user_id SERIAL PRIMARY KEY,
//     username VARCHAR(50) UNIQUE NOT NULL,
//     password_hash VARCHAR(255) NOT NULL,
//     email VARCHAR(100) UNIQUE NOT NULL,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     role INTEGER REFERENCES Roles(role_id)
// );



// CREATE TABLE Roles (
//     role_id SERIAL PRIMARY KEY,
//     role_name VARCHAR(50) UNIQUE NOT NULL
// );

// // INSERT INTO Roles (role_name) VALUES ('customer'), ('admin');


// CREATE TABLE Categories (
//     category_id SERIAL PRIMARY KEY,
//     name VARCHAR(100) NOT NULL,
//     description TEXT
// );


// CREATE TABLE Products (
//     product_id SERIAL PRIMARY KEY,
//     name VARCHAR(100) NOT NULL,
//     description TEXT,
//     price DECIMAL(10, 2) NOT NULL,
//     stock INTEGER DEFAULT 0,
//     category_id INTEGER REFERENCES Categories(category_id) ON DELETE CASCADE,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );
// ALTER TABLE Products ADD COLUMN seller_id INTEGER REFERENCES Users(user_id);
// UPDATE Products
// SET seller_id = FLOOR(RANDOM() * 100 + 1);


// CREATE TABLE Orders (
//     order_id SERIAL PRIMARY KEY,
//     user_id INTEGER REFERENCES Users(user_id) ON DELETE CASCADE,
//     order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     status VARCHAR(50) NOT NULL,
//     total_amount DECIMAL(10, 2) NOT NULL
// );


// CREATE TABLE OrderItems (
//     order_item_id SERIAL PRIMARY KEY,
//     order_id INTEGER REFERENCES Orders(order_id) ON DELETE CASCADE,
//     product_id INTEGER REFERENCES Products(product_id) ON DELETE CASCADE,
//     quantity INTEGER NOT NULL,
//     price DECIMAL(10, 2) NOT NULL
// );


// CREATE TABLE ShoppingCart (
//     cart_id SERIAL PRIMARY KEY,
//     user_id INTEGER REFERENCES Users(user_id) ON DELETE CASCADE,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );


// CREATE TABLE CartItems (
//     cart_item_id SERIAL PRIMARY KEY,
//     cart_id INTEGER REFERENCES ShoppingCart(cart_id) ON DELETE CASCADE,
//     product_id INTEGER REFERENCES Products(product_id) ON DELETE CASCADE,
//     quantity INTEGER NOT NULL
//     UNIQUE (cart_id, product_id)
// );


// CREATE TABLE Payments (
//     payment_id SERIAL PRIMARY KEY,
//     order_id INTEGER REFERENCES Orders(order_id) ON DELETE CASCADE,
//     amount DECIMAL(10, 2) NOT NULL,
//     payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     payment_method VARCHAR(50) NOT NULL
// );


// CREATE TABLE Reviews (
//     review_id SERIAL PRIMARY KEY,
//     product_id INTEGER REFERENCES Products(product_id) ON DELETE CASCADE,
//     user_id INTEGER REFERENCES Users(user_id),
//     rating INTEGER CHECK (rating BETWEEN 1 AND 5),
//     comment TEXT,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );



// CREATE TABLE Sessions (
//     session_id SERIAL PRIMARY KEY,
//     user_id INTEGER REFERENCES Users(user_id) ON DELETE CASCADE,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     expires_at TIMESTAMP NOT NULL
// );



// CREATE TABLE UserAddresses (
//     address_id SERIAL PRIMARY KEY,
//     user_id INTEGER REFERENCES Users(user_id) ON DELETE CASCADE,
//     street VARCHAR(255),
//     city VARCHAR(100),
//     state VARCHAR(100),
//     zip_code VARCHAR(20)
// );



// CREATE TABLE ProductImages (
//     image_id SERIAL PRIMARY KEY,
//     product_id INTEGER REFERENCES Products(product_id) ON DELETE CASCADE,
//     image_url TEXT NOT NULL,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );



// CREATE TABLE AuditLogs (
//     log_id SERIAL PRIMARY KEY,
//     action TEXT NOT NULL,
//     user_id INTEGER REFERENCES Users(user_id) ON DELETE CASCADE,
//     timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );



// CREATE TABLE Cache (
//     cache_key VARCHAR(255) PRIMARY KEY,
//     cache_value TEXT,
//     expiration_time TIMESTAMP NOT NULL
// );




// DELETE FROM CartItems
// WHERE cart_item_id NOT IN (
//     SELECT MIN(cart_item_id)
//     FROM CartItems
//     GROUP BY cart_id, product_id
// );



// ALTER TABLE CartItems
// ADD CONSTRAINT cart_product_unique UNIQUE (cart_id, product_id);