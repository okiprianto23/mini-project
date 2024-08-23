-- +migrate Up

ALTER TABLE users
    ADD COLUMN auth_user_id INT UNIQUE NOT NULL,
    ADD COLUMN client_id VARCHAR(255),
    ADD COLUMN signature_key VARCHAR(255),
    ADD COLUMN locale ENUM('en-US', 'id-ID'),
    ADD COLUMN alias_name VARCHAR(255),
    ADD COLUMN client_alias VARCHAR(255);

INSERT INTO users (uuid_id, auth_user_id, client_id, signature_key, locale, alias_name, client_alias, username, email,
                   password_hash, created_by, created_client, updated_by, updated_client)
VALUES (UUID(), 1, '2ac7f390e3b9488784869135bd9b6278', '4be386de3ed54c4ba7290d93ef4b0919', 'en-US',
        'Testing Grolog Satu',
        'client_alias_v', 'superadmin', 'superadmin1@gmail.com',
        '$2y$10$Xmp1omkWaDcPJoiiAyypfOsd5.7cxbU7mee168gwZ6tyjmtfnKrza', 1, 'SYSTEM', 1, 'SYSTEM');
