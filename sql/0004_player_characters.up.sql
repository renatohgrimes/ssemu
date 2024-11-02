CREATE TABLE player_characters (
    id INTEGER NOT NULL PRIMARY KEY ASC AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    slot INTEGER NOT NULL,
    mask INTEGER NOT NULL,
    weapon1 INTEGER,
    weapon2 INTEGER,
    weapon3 INTEGER,
    skill INTEGER,
    hair INTEGER,
    face INTEGER,
    shirt INTEGER,
    pants INTEGER,
    shoes INTEGER,
    gloves INTEGER,
    accessory INTEGER,
    is_active BOOLEAN,

    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX player_characters_user_id ON player_characters (user_id)