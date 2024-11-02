CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY ASC,
    username TEXT UNIQUE NOT NULL,
    password CHAR(60) NOT NULL,
    created_utc DATETIME NOT NULL,
    banned_utc DATETIME NOT NULL,
    is_admin INTEGER NOT NULL,
    last_login_utc DATETIME NOT NULL
);