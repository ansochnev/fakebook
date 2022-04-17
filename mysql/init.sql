CREATE DATABASE fakebook
    CHARACTER SET utf8;

USE fakebook;


CREATE TABLE accounts (
    id       INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    username CHAR(255) UNIQUE NOT NULL,
    email    CHAR(255) UNIQUE NOT NULL
);

CREATE TABLE sequence (id INT UNSIGNED PRIMARY KEY);
INSERT INTO sequence VALUES (1000);

DELIMITER ##
CREATE TRIGGER accounts_set_default_username_from_id BEFORE INSERT ON accounts
    FOR EACH ROW
    BEGIN
        UPDATE sequence SET id = LAST_INSERT_ID(id + 1);
        SET NEW.id = LAST_INSERT_ID();
        IF NEW.username = '' THEN
            SET NEW.username = NEW.id;
        END IF;
    END##
DELIMITER ;

CREATE TABLE passwords (
    account_id INT UNSIGNED PRIMARY KEY,
    hash BINARY(32) NOT NULL,
    salt BINARY(32) NOT NULL,

    CONSTRAINT FOREIGN KEY (account_id)
        REFERENCES accounts (id)
            ON UPDATE CASCADE
            ON DELETE CASCADE
);

CREATE TABLE profiles (
    account_id INT UNSIGNED PRIMARY KEY,
    first_name     CHAR(255) NOT NULL,
    last_name      CHAR(255) NOT NULL,
    date_of_birth  DATE,
    sex            ENUM ('female', 'male'),
    city           CHAR(255),
    info           TEXT,

    CONSTRAINT FOREIGN KEY (account_id)
        REFERENCES accounts (id)
            ON UPDATE CASCADE
            ON DELETE CASCADE
);


CREATE USER 'fakebook'@'localhost' IDENTIFIED BY 'test=Oohuluo8';
GRANT SELECT, INSERT, UPDATE, DELETE ON fakebook.* TO 'fakebook'@'localhost';
