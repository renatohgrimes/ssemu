DELETE FROM player_characters;
DELETE FROM player_licenses;
DELETE FROM players;
DELETE FROM users;

INSERT INTO users (id, username, password, created_utc, banned_utc, is_admin, last_login_utc) VALUES
(1, "testuser", "$2a$12$wQsoPsEz6HmQ23WJBpU6M.mGzjhv8Tb9mQd7h9qgVlx38pyeIuhiu", "2000-01-02 03:04:05+00:00", "0001-01-01 00:00:00+00:00", false, "0001-01-01 00:00:00+00:00"),
(2, "testban", "$2a$12$wQsoPsEz6HmQ23WJBpU6M.mGzjhv8Tb9mQd7h9qgVlx38pyeIuhiu", "2000-01-02 03:04:05+00:00", "2000-01-02 03:04:05+00:00", false, "0001-01-01 00:00:00+00:00"),
(3, "testnick", "$2a$12$wQsoPsEz6HmQ23WJBpU6M.mGzjhv8Tb9mQd7h9qgVlx38pyeIuhiu", "2000-01-02 03:04:05+00:00", "0001-01-01 00:00:00+00:00", false, "0001-01-01 00:00:00+00:00"),
(4, "testadmin", "$2a$12$wQsoPsEz6HmQ23WJBpU6M.mGzjhv8Tb9mQd7h9qgVlx38pyeIuhiu", "2000-01-02 03:04:05+00:00", "0001-01-01 00:00:00+00:00", true, "0001-01-01 00:00:00+00:00"),
(5, "testnotadm", "$2a$12$wQsoPsEz6HmQ23WJBpU6M.mGzjhv8Tb9mQd7h9qgVlx38pyeIuhiu", "2000-01-02 03:04:05+00:00", "0001-01-01 00:00:00+00:00", false, "0001-01-01 00:00:00+00:00"),
(6, "testplrsvc", "$2a$12$wQsoPsEz6HmQ23WJBpU6M.mGzjhv8Tb9mQd7h9qgVlx38pyeIuhiu", "2000-01-02 03:04:05+00:00", "0001-01-01 00:00:00+00:00", false, "0001-01-01 00:00:00+00:00"),
(7, "testnewplr", "$2a$12$wQsoPsEz6HmQ23WJBpU6M.mGzjhv8Tb9mQd7h9qgVlx38pyeIuhiu", "2000-01-02 03:04:05+00:00", "0001-01-01 00:00:00+00:00", false, "0001-01-01 00:00:00+00:00"),
(8, "test", "$2y$10$0eoYLyj/21L17sMkGmtWgO8u6ivNYJw6Js.rYdBcpjBUz6GqXwDBC", "2000-01-02 03:04:05+00:00", "0001-01-01 00:00:00+00:00", true, "0001-01-01 00:00:00+00:00"),
(9, "test2", "$2y$10$qWKFpkYaAVzmBrMD6Gf9w..QOPAasG3QTngXRLpe/L3yLa0.DcJ8a", "2000-01-02 03:04:05+00:00", "0001-01-01 00:00:00+00:00", false, "0001-01-01 00:00:00+00:00");

INSERT INTO players (id, user_id, nickname, created_utc, tutorial_status) VALUES
(1, 1, "TestUserPlayer", "2000-01-02 03:04:05+00:00", 0),
(2, 3, NULL, "2000-01-02 03:04:05+00:00", 0),
(3, 4, "TestAdminPlayer", "2000-01-02 03:04:05+00:00", 0),
(4, 5, "TestNonAdmin", "2000-01-02 03:04:05+00:00", 0),
(5, 6, "TestPlayerData", "2000-01-02 03:04:05+00:00", 0),
(6, 8, "test", "2000-01-02 03:04:05+00:00", 3),
(7, 9, "test2", "2000-01-02 03:04:05+00:00", 3);

INSERT INTO player_licenses(id, user_id, license_id) VALUES 
(1, 6, 101),
(2, 6, 102);

INSERT INTO player_characters (id, user_id, slot, mask, weapon1, weapon2, weapon3, skill, hair, face, shirt, pants, shoes, gloves, accessory, is_active) VALUES 
(1, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0),
(2, 6, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1),
(3, 8, 0, 0, 5, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1),
(4, 9, 0, 1, 5, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1);
