-- *****************************************************************************************
-- TABLE password_tokens
-- *****************************************************************************************
-- This table contains all fields for user to set his password or to reset existing one.
-- *****************************************************************************************
CREATE TABLE IF NOT EXISTS password_tokens (
	id SERIAL,
    user_id CHAR(36) NOT NULL,
    token VARCHAR(100) NOT NULL,
    expires_at BIGINT(20) NOT NULL,
    date_created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
);