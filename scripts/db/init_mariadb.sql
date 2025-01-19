CREATE DATABASE tradedb;

CREATE USER 'tradebot'@'localhost' IDENTIFIED BY '{password}';

GRANT ALL PRIVILEGES ON gotradecrypto.* TO 'tradebot'@'localhost';

FLUSH PRIVILEGES;