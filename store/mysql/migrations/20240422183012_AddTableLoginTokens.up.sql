-- *****************************************************************************************
-- TABLE login_tokens
-- *****************************************************************************************
CREATE TABLE IF NOT EXISTS login_tokens(
	id SERIAL,
    user_id CHAR(36) NOT NULL,
    token TEXT NOT NULL,
    token_type VARCHAR(100) NOT NULL,
    expires_at BIGINT(20) NOT NULL,
    CONSTRAINT fk_login_tokens_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
	PRIMARY KEY (id)
);