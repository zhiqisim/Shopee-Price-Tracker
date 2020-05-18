CREATE DATABASE itemdb;
use itemdb;

/* CREATE item TABLE */
CREATE TABLE item (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    item_id BIGINT NOT NULL UNIQUE,
    shop_id BIGINT NOT NULL,
    item_name TEXT NOT NULL
);

CREATE UNIQUE INDEX item_index ON item (item_id);

/* CREATE item_price TABLE */
CREATE TABLE item_price (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    item_id BIGINT NOT NULL,
    price_datetime TIMESTAMP NOT NULL,
    price INT NOT NULL
);

CREATE INDEX item_price_index ON item_price (item_id);

/* TEST INSERT QUERIES */
-- INSERT INTO item (item_id, shop_id, item_name) VALUES (67707640, 2245175584, "Praise-Camp-Collar-Short-Sleeves-Navy");
-- INSERT INTO item_price (item_id, price_datetime, price) VALUES (2245175584, '18-06-12 10:34:09', 1099000);