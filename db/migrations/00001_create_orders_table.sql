-- +goose Up
-- Orders table (Main entity storing order details)
CREATE TABLE orders (
    order_uid VARCHAR(50) PRIMARY KEY,
    track_number VARCHAR(50),
    entry VARCHAR(50),
    locale VARCHAR(10),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(50),
    delivery_service VARCHAR(50),
    shardkey VARCHAR(10),
    sm_id INT,
    date_created TIMESTAMP,
    oof_shard VARCHAR(10)
);

-- Delivery details (1-to-1 relationship with orders)
CREATE TABLE deliveries (
    order_uid VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255),
    phone VARCHAR(15),
    zip VARCHAR(10),
    city VARCHAR(100),
    address VARCHAR(255),
    region VARCHAR(100),
    email VARCHAR(100),
    FOREIGN KEY (order_uid) REFERENCES orders (order_uid) ON DELETE CASCADE
);

-- Payment details (1-to-1 relationship with orders)
CREATE TABLE payments (
    order_uid VARCHAR(50) PRIMARY KEY,
    transaction VARCHAR(50),
    request_id VARCHAR(50),
    currency VARCHAR(10),
    provider VARCHAR(50),
    amount INT,
    payment_dt TIMESTAMP,
    bank VARCHAR(50),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT,
    FOREIGN KEY (order_uid) REFERENCES orders (order_uid) ON DELETE CASCADE
);

-- Items details (1-to-many relationship with orders)
CREATE TABLE items (
    item_id SERIAL PRIMARY KEY,
    order_uid VARCHAR(50),
    chrt_id INT,
    track_number VARCHAR(50),
    price INT,
    rid VARCHAR(50),
    name VARCHAR(255),
    sale INT,
    size VARCHAR(10),
    total_price INT,
    nm_id INT,
    brand VARCHAR(100),
    status INT,
    FOREIGN KEY (order_uid) REFERENCES orders (order_uid) ON DELETE CASCADE
);

CREATE INDEX idx_orders_date_created ON orders (date_created);
CREATE INDEX idx_items_order_uid ON items (order_uid);

-- +goose Down
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS deliveries;
DROP TABLE IF EXISTS orders;
DROP INDEX IF EXISTS idx_orders_date_created;
DROP INDEX IF EXISTS idx_items_order_uid;
