-- *****************************************************************************************
-- STORED PROCEDURE AddPasswordToken
-- =========================================================================================
DROP PROCEDURE IF EXISTS AddPasswordToken;

--MYSQL_CUSTOM_STATEMENT_DELIMITER

CREATE PROCEDURE AddPasswordToken (
    IN inUserID CHAR(36),
    IN inToken VARCHAR(100),
    IN inExpiresAt BIGINT
)
BEGIN

	INSERT INTO password_tokens(user_id, token, expires_at)
	VALUES(inUserID, inToken, inExpiresAt);

    -- Select last inserted ID from the PasswordTokens
    SET @passwordTokenID = LAST_INSERT_ID();

    SELECT user_id, token, expires_at
    FROM password_tokens
    WHERE id = @passwordTokenID;
    
END;
