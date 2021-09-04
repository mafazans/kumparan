# Create Database for Article :
DROP DATABASE IF EXISTS `kumparan_db`;
CREATE DATABASE IF NOT EXISTS `kumparan_db` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

DROP TABLE IF EXISTS `kumparan_db`.`article`;
    CREATE TABLE IF NOT EXISTS `kumparan_db`.`article` (
       id INT UNSIGNED NOT NULL AUTO_INCREMENT,
       author VARCHAR(32) NOT NULL COMMENT "",
       title TEXT NOT NULL,
       body TEXT NOT NULL,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       PRIMARY KEY (id)
    ) ENGINE=INNODB;

UPDATE mysql.user SET Host='%' WHERE Host='localhost' AND User='root';
UPDATE mysql.db SET Host='%' WHERE Host='localhost' AND User='root';
FLUSH PRIVILEGES;
ALTER USER 'root' IDENTIFIED WITH mysql_native_password BY 'password';
