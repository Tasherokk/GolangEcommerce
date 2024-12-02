-- DROP TABLE Cache;
-- DROP TABLE AuditLogs;
-- DROP TABLE CartItems;
-- DROP TABLE OrderItems;
-- DROP TABLE Payments;
-- DROP TABLE Reviews;
-- DROP TABLE Sessions;
-- DROP TABLE ShoppingCart;
-- DROP TABLE UserAddresses;
-- DROP TABLE ProductImages;
-- DROP TABLE Orders;
-- DROP TABLE Products;
-- DROP TABLE Users;
-- DROP TABLE Categories;
-- DROP TABLE Roles;


CREATE TABLE IF NOT EXISTS Cache (
    cache_key VARCHAR(255) PRIMARY KEY,
    cache_value TEXT,
    expiration_time TIMESTAMP NOT NULL
);
TRUNCATE TABLE "public"."cache" CASCADE;



CREATE TABLE IF NOT EXISTS Roles (
    role_id SERIAL PRIMARY KEY,
    role_name VARCHAR(50) UNIQUE NOT NULL
);
TRUNCATE TABLE "public"."roles" CASCADE;
COPY "public"."roles"("role_name")
FROM 'C:/Program Files/PostgreSQL/15/data/custom/Golang_Ecommerce//Roles-3015edee696e5be3597307f4f1ef06550e508d58503e50287f221e2b2c3095a6.csv'
WITH (FORMAT csv, DELIMITER ',', NULL 'NULL', QUOTE '"');



CREATE TABLE IF NOT EXISTS Categories (
    category_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT
);
TRUNCATE TABLE "public"."categories" CASCADE;
COPY "public"."categories"("name","description")
FROM 'C:/Program Files/PostgreSQL/15/data/custom/Golang_Ecommerce/Categories-7098fc3c09edb7febea7cfb3fa01cf6392210a9b537abff9e1338f8ae2c944d8.csv'
WITH (FORMAT csv, DELIMITER ',', NULL 'NULL', QUOTE '"');



CREATE TABLE IF NOT EXISTS Users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    role INTEGER REFERENCES Roles(role_id)
);
TRUNCATE TABLE "public"."users" CASCADE;
COPY "public"."users"("username","password_hash","email","created_at","role")
FROM 'C:/Program Files/PostgreSQL/15/data/custom/Golang_Ecommerce//Users-f444436cd877c9201a9f4780410c73d85841b506bae369e2c0d5c0a9937bb39a.csv'
WITH (FORMAT csv, DELIMITER ',', NULL 'NULL', QUOTE '"');



CREATE TABLE Orders (
    order_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES Users(user_id),
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL
);
TRUNCATE TABLE "public"."orders" CASCADE;
COPY "public"."orders"("user_id","order_date","status","total_amount")
FROM 'C:/Program Files/PostgreSQL/15/data/custom/Golang_Ecommerce/Orders-77542fe84a475bbb3a44c477a588acc6f771d2cb5785fba17debbdb25249ccc9.csv'
WITH (FORMAT csv, DELIMITER ',', NULL 'NULL', QUOTE '"');



CREATE TABLE IF NOT EXISTS Products (
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    stock INTEGER DEFAULT 0,
    category_id INTEGER REFERENCES Categories(category_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
TRUNCATE TABLE "public"."products" CASCADE;
COPY "public"."products"("name","description","price","stock","category_id","created_at")
FROM 'C:/Program Files/PostgreSQL/15/data/custom/Golang_Ecommerce//Products-31754911f0a23c1e1ad25e110a85cf04771dfe49a012347d922dc6472906999e.csv'
WITH (FORMAT csv, DELIMITER ',', NULL 'NULL', QUOTE '"');




CREATE TABLE IF NOT EXISTS OrderItems (
    order_item_id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES Orders(order_id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES Products(product_id),
    quantity INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);
TRUNCATE TABLE "public"."orderitems" CASCADE;
COPY "public"."orderitems"("order_id","product_id","quantity","price")
FROM 'C:/Program Files/PostgreSQL/15/data/custom/Golang_Ecommerce/OrderItems-91307b3f5ce3d8c2ba86639371ecad6506c8b4f4a205a1c776ff2438ee02f18d.csv'
WITH (FORMAT csv, DELIMITER ',', NULL 'NULL', QUOTE '"');




CREATE TABLE IF NOT EXISTS Payments (
    payment_id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES Orders(order_id),
    amount DECIMAL(10, 2) NOT NULL,
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    payment_method VARCHAR(50) NOT NULL
);
TRUNCATE TABLE "public"."payments" CASCADE;
COPY "public"."payments"("order_id","amount","payment_date","payment_method")
FROM 'C:/Program Files/PostgreSQL/15/data/custom/Golang_Ecommerce//Payments-7b63c58c6e1f13875a7961331c5fcc1ed1af45c9de8416a7d8380bbcb665a6a5.csv'
WITH (FORMAT csv, DELIMITER ',', NULL 'NULL', QUOTE '"');



CREATE TABLE IF NOT EXISTS ProductImages (
    image_id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES Products(product_id),
    image_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
TRUNCATE TABLE "public"."productimages" CASCADE;
COPY "public"."productimages"("product_id","image_url","created_at")
FROM 'C:/Program Files/PostgreSQL/15/data/custom/Golang_Ecommerce//ProductImages-2575fedd43162dad9c224544ace11bd560e7782ab20548d15c8aec0f4ec7c566.csv'
WITH (FORMAT csv, DELIMITER ',', NULL 'NULL', QUOTE '"');





CREATE TABLE IF NOT EXISTS Reviews (
    review_id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES Products(product_id),
    user_id INTEGER REFERENCES Users(user_id),
    rating INTEGER CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
TRUNCATE TABLE "public"."reviews" CASCADE;
COPY "public"."reviews"("product_id","user_id","rating","comment","created_at")
FROM 'C:/Program Files/PostgreSQL/15/data/custom/Golang_Ecommerce//Reviews-9de2bb82e71ead404287635d46bcbeb41cba67dff9110fb047abce8d6b2f284f.csv'
WITH (FORMAT csv, DELIMITER ',', NULL 'NULL', QUOTE '"');





CREATE TABLE IF NOT EXISTS Sessions (
    session_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES Users(user_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);
TRUNCATE TABLE "public"."sessions" CASCADE;
COPY "public"."sessions"("user_id","created_at","expires_at")
FROM 'C:/Program Files/PostgreSQL/15/data/custom/Golang_Ecommerce//Sessions-c4e174f143f841e687807ba76a819d24aced3a9f74247ba4704309ec354f4af8.csv'
WITH (FORMAT csv, DELIMITER ',', NULL 'NULL', QUOTE '"');



CREATE TABLE IF NOT EXISTS ShoppingCart (
    cart_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES Users(user_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
TRUNCATE TABLE "public"."shoppingcart" CASCADE;
COPY "public"."shoppingcart"("user_id","created_at")
FROM 'C:/Program Files/PostgreSQL/15/data/custom/Golang_Ecommerce//ShoppingCart-5bdadb3b0f117436246456daf50c0947d85e3036905b4e45f7104a047a25c5dc.csv'
WITH (FORMAT csv, DELIMITER ',', NULL 'NULL', QUOTE '"');



CREATE TABLE IF NOT EXISTS UserAddresses (
    address_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES Users(user_id),
    street VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    zip_code VARCHAR(20)
);
TRUNCATE TABLE "public"."useraddresses" CASCADE;
COPY "public"."useraddresses"("user_id","street","city","state","zip_code")
FROM 'C:/Program Files/PostgreSQL/15/data/custom/Golang_Ecommerce//UserAddresses-ca063a00335d9641c6c9f9d64b66408c6a7720ebeb525b8e3d50fcd153c586a2.csv'
WITH (FORMAT csv, DELIMITER ',', NULL 'NULL', QUOTE '"');




CREATE TABLE IF NOT EXISTS AuditLogs (
    log_id SERIAL PRIMARY KEY,
    action TEXT NOT NULL,
    user_id INTEGER REFERENCES Users(user_id),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
TRUNCATE TABLE "public"."auditlogs" CASCADE;
COPY "public"."auditlogs"("action","user_id","timestamp")
FROM 'C:/Program Files/PostgreSQL/15/data/custom/Golang_Ecommerce/AuditLogs-e0af497e7235c67b817d0b1164f07db5dc62286d3aae246a6e58f7204bd0df8c.csv'
WITH (FORMAT csv, DELIMITER ',', NULL 'NULL', QUOTE '"');



CREATE TABLE IF NOT EXISTS CartItems (
    cart_item_id SERIAL PRIMARY KEY,
    cart_id INTEGER REFERENCES ShoppingCart(cart_id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES Products(product_id),
    quantity INTEGER NOT NULL
);
TRUNCATE TABLE "public"."cartitems" CASCADE;
COPY "public"."cartitems"("cart_id","product_id","quantity")
FROM 'C:/Program Files/PostgreSQL/15/data/custom/Golang_Ecommerce/CartItems-52e1c36deb1112c24fd882733338fa2133097017d7fa4df794dfa93ebb7e38f2.csv'
WITH (FORMAT csv, DELIMITER ',', NULL 'NULL', QUOTE '"');