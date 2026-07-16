-- ユーザー本体
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    device_id VARCHAR(64) NOT NULL UNIQUE,
    name VARCHAR(32) NOT NULL DEFAULT 'あいうえお',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_claimed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 認証トークン
CREATE TABLE IF NOT EXISTS auth_tokens (
    token VARCHAR(64) PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    INDEX index_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ユーザーの資源
CREATE TABLE IF NOT EXISTS user_resources (
    user_id BIGINT UNSIGNED PRIMARY KEY,
    coin BIGINT UNSIGNED NOT NULL DEFAULT 0,
    production_rate INT UNSIGNED NOT NULL DEFAULT 1,
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- キャラクターマスタ
CREATE TABLE IF NOT EXISTS characters (
    id SMALLINT UNSIGNED PRIMARY KEY,
    name VARCHAR(32) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO characters (id, name) VALUES
    (1, 'ガレン'),
    (2, 'アッシュ'),
    (3, 'マルザハール'),
ON DUPLICATE KEY UPDATE name = VALUES(name);

-- ユーザーの所持キャラ(同じキャラを複数回引いたらcountが増える)
CREATE TABLE IF NOT EXISTS user_characters(
    user_id BIGINT UNSIGNED NOT NULL,
    character_id SMALLINT UNSIGNED NOT NULL,
    count INT UNSIGNED NOT NULL DEFAULT 0,
    obtained_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, character_id),
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (character_id) REFERENCES character(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;