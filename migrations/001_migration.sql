CREATE SCHEMA IF NOT EXISTS shortener;
CREATE TABLE IF NOT EXISTS shortener.shortener(
short_url varchar(15) NOT NULL,
url varchar(255) NOT NULL,
user_id varchar(50) NOT NULL,
PRIMARY KEY (short_url)
);
