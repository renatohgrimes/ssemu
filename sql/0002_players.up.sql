CREATE TABLE players (
    id INTEGER NOT NULL PRIMARY KEY ASC AUTOINCREMENT,
    user_id INTEGER UNIQUE NOT NULL,
    nickname TEXT UNIQUE,
    created_utc DATETIME NOT NULL,
    tutorial_status INTEGER NOT NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE UNIQUE INDEX players_user_id ON players (user_id)