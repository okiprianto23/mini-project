-- +migrate Up
CREATE TABLE users
(
    id             INT AUTO_INCREMENT PRIMARY KEY,
    uuid_id        CHAR(36)            NOT NULL,
    username       VARCHAR(255) UNIQUE NOT NULL,
    email          VARCHAR(255) UNIQUE NOT NULL,
    password_hash  VARCHAR(255)        NOT NULL,
    created_by     INT                 NOT NULL,
    created_client VARCHAR(256)        NOT NULL,
    created_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_by     INT                 NOT NULL,
    updated_client VARCHAR(256)        NOT NULL,
    updated_at     DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted        BOOLEAN  DEFAULT FALSE,
    is_active      BOOLEAN  DEFAULT TRUE
);

CREATE TABLE client_token
(
    id             INT AUTO_INCREMENT PRIMARY KEY,
    uuid_id        CHAR(36)     NOT NULL,
    user_id        INT,
    token          TEXT         NOT NULL,
    created_by     INT          NOT NULL,
    created_client VARCHAR(256) NOT NULL,
    created_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at     DATETIME     NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE consumer
(
    id             INT AUTO_INCREMENT PRIMARY KEY,
    uuid_id        CHAR(36)            NOT NULL,
    user_id        INT,
    nik            VARCHAR(255) UNIQUE NOT NULL,
    full_name      VARCHAR(255)        NOT NULL,
    legal_name     VARCHAR(255)        NOT NULL,
    birth_place    VARCHAR(255)        NOT NULL,
    birth_date     DATE                NOT NULL,
    salary         DECIMAL(15, 2)      NOT NULL,
    ktp_photo      VARCHAR(255)        NOT NULL,
    selfie_photo   VARCHAR(255)        NOT NULL,
    created_by     INT                 NOT NULL,
    created_client VARCHAR(256)        NOT NULL,
    created_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_by     INT                 NOT NULL,
    updated_client VARCHAR(256)        NOT NULL,
    updated_at     DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted        BOOLEAN  DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE credit_limit
(
    limit_id               INT AUTO_INCREMENT PRIMARY KEY,
    uuid_id                CHAR(36)       NOT NULL,
    consumer_id            INT,
    monthly_installments   DECIMAL(15, 2) NOT NULL,
    interest_rate          DECIMAL(5, 2)  NOT NULL,
    tenor                  INT            NOT NULL,
    limit_amount           DECIMAL(15, 2) NOT NULL,
    remaining_limit_amount DECIMAL(15, 2) NOT NULL,
    created_by             INT            NOT NULL,
    created_client         VARCHAR(256)   NOT NULL,
    created_at             DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_by             INT            NOT NULL,
    updated_client         VARCHAR(256)   NOT NULL,
    updated_at             DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted                BOOLEAN  DEFAULT FALSE
);

CREATE TABLE transaction
(
    id                 INT AUTO_INCREMENT PRIMARY KEY,
    uuid_id            CHAR(36)            NOT NULL,
    consumer_id        INT,
    limit_id           INT,
    contract_number    VARCHAR(255) UNIQUE NOT NULL,
    otr                DECIMAL(15, 2)      NOT NULL,
    admin_fee          DECIMAL(15, 2)      NOT NULL,
    installment_amount DECIMAL(15, 2)      NOT NULL,
    interest_amount    DECIMAL(15, 2)      NOT NULL,
    asset_name         VARCHAR(255)        NOT NULL,
    transaction_date   DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by         INT                 NOT NULL,
    created_client     VARCHAR(256)        NOT NULL,
    created_at         DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_by         INT                 NOT NULL,
    updated_client     VARCHAR(256)        NOT NULL,
    updated_at         DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted            BOOLEAN  DEFAULT FALSE,
    FOREIGN KEY (consumer_id) REFERENCES consumer (id),
    FOREIGN KEY (limit_id) REFERENCES credit_limit (limit_id)
);

-- +migrate Down
DROP TABLE IF EXISTS transaction;
DROP TABLE IF EXISTS credit_limit;
DROP TABLE IF EXISTS consumer;
DROP TABLE IF EXISTS client_token;
DROP TABLE IF EXISTS users;
