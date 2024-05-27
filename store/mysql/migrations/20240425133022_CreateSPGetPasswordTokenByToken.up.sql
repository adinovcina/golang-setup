-- *****************************************************************************************
-- STORED PROCEDURE GetPasswordTokenByToken
-- =========================================================================================
DROP PROCEDURE IF EXISTS GetPasswordTokenByToken;

--MYSQL_CUSTOM_STATEMENT_DELIMITER

CREATE PROCEDURE GetPasswordTokenByToken (
    IN inToken VARCHAR(100)
)
BEGIN

	SELECT user_id, token, expires_at
    FROM password_tokens
    WHERE token = inToken;

END;