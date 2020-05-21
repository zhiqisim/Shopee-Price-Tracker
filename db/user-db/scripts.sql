CREATE DATABASE userdb;
use userdb;

/* CREATE users TABLE */
CREATE TABLE user (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL UNIQUE,
    user_password CHAR(64) NOT NULL
);

CREATE UNIQUE INDEX user_index ON user (username);

/* CREATE user_item TABLE */
CREATE TABLE user_item (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL,
    item_id BIGINT NOT NULL,
    item_name TEXT NOT NULL,
    UNIQUE KEY username_item_id_key (username, item_id)
);

CREATE INDEX user_item_index ON user_item (username);


/* TEST INSERT QUERIES */
-- INSERT INTO user (username, user_password) VALUES ("zhiqisim", "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8");
-- INSERT INTO user_item (username, item_id) VALUES ("zhiqisim", "2245175584");