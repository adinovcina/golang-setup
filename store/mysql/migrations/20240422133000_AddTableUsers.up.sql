-- *****************************************************************************************
-- TABLE users
-- *****************************************************************************************
CREATE TABLE IF NOT EXISTS users (
    id CHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
    name VARCHAR(150) NOT NULL,
    email VARCHAR(250) NOT NULL,
    phone VARCHAR(50) NULL,
    language VARCHAR(2) NOT NULL DEFAULT 'en',
    password VARCHAR(500) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT FALSE,
    -- Stores the flag if user verified his account over the email
    email_verified TINYINT(1) DEFAULT 0,
    -- Stores the date when terms and conditions are accepted
    terms_accepted TINYINT(1) DEFAULT 0,
    -- Used when user is blocked to login
    login_blocked_until DATETIME NULL,
    -- Used to count how many failed login attemps there is
    failed_login_count INT NOT NULL DEFAULT 0,
    last_time_logged BIGINT(20) NULL,
    -- required for tracking purposes
    date_created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX `uq_idx_email` (`email`)
);