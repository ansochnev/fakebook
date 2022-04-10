CREATE DATABASE fakebook
    CHARACTER SET utf8;

USE fakebook;


CREATE TABLE accounts (
    id       INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    email    CHAR(255) NOT NULL,
    username CHAR(255) NOT NULL
);

CREATE TABLE passwords (
    account_id INT UNSIGNED PRIMARY KEY
        REFERENCES accounts (id)
            ON UPDATE CASCADE
            ON DELETE CASCADE,

    hash BINARY(32) NOT NULL,
    salt BINARY(32) NOT NULL
);

CREATE TABLE profiles (
    account_id INT UNSIGNED PRIMARY KEY
        REFERENCES accounts (id)
            ON UPDATE CASCADE
            ON DELETE CASCADE,

    first_name     CHAR(255) NOT NULL,
    last_name      CHAR(255) NOT NULL,
    date_of_birth  DATE,
    sex            ENUM ('female', 'male'),
    city           CHAR(255),
    info           TEXT
);


CREATE USER 'fakebook'@'localhost' IDENTIFIED BY 'test=Oohuluo8';
GRANT SELECT, INSERT, UPDATE, DELETE ON fakebook.* TO 'fakebook'@'localhost';
