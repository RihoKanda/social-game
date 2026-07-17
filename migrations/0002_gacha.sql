-- キャラクターマスタ
CREATE TABLE IF NOT EXISTS characters (
    id SMALLINT UNSIGNED PRIMARY KEY,
    name VARCHAR(32) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO characters (id, name) VALUES
    (1, 'ガレン'),
    (2, 'アッシュ'),
    (3, 'マルザハール')
ON DUPLICATE KEY UPDATE name = VALUES(name);

-- ユーザーの所持キャラ(同じキャラを複数回引いたらcountが増える)
CREATE TABLE IF NOT EXISTS user_characters (
    user_id BIGINT UNSIGNED NOT NULL,
    character_id SMALLINT UNSIGNED NOT NULL,
    count INT UNSIGNED NOT NULL DEFAULT 0,
    obtained_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, character_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (character_id) REFERENCES characters(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;