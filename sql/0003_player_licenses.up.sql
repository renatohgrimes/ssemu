CREATE TABLE player_licenses (
    id INTEGER NOT NULL PRIMARY KEY ASC AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    license_id INTEGER NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX player_licenses_user_id ON player_licenses (user_id)