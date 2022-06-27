-- SETUP DATABASE
CREATE DATABASE go_gorm_sample;
USE go_gorm_sample;

-- SETUP TABLES
CREATE TABLE market (market_id INTEGER PRIMARY KEY, market_name VARCHAR(100));
CREATE TABLE fruit (fruit_id INTEGER PRIMARY KEY, fruit_name VARCHAR(100), market_id INTEGER, customer_id INTEGER);
CREATE TABLE customer (customer_id INTEGER PRIMARY KEY, customer_name VARCHAR(100));

-- SETUP SAMPLE VALUES
INSERT INTO market VALUES
(1,"Tokyo"),
(2,"Saitama"),
(3,"Yokohama"),
(4,"Sendai"),
(5,"Osaka"),
(6,"Nagoya");

INSERT INTO fruit VALUES
(1,"banana",1,8),
(2,"apple",2,9),
(3,"acerola",1,1),
(4,"akebia",3,8),
(5,"apricot",2,2),
(6,"grapes",1,8),
(7,"guava",2,1),
(8,"lemon",2,3),
(9,"quince",1,4),
(10,"sapodilla",5,7),
(11,"carambola",6,6),
(12,"prune",6,5);

INSERT INTO customer VALUES
(1,"Jon"),
(2,"Ken"),
(3,"Musashi"),
(4,"Yuka"),
(5,"Atsuya"),
(6,"Kaede"),
(7,"Takeshi"),
(8,"Mike"),
(9,"Jack");
