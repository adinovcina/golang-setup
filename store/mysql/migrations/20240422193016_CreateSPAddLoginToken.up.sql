-- *****************************************************************************************
-- STORED PROCEDURE AddLoginToken
-- =========================================================================================
DROP PROCEDURE IF EXISTS AddLoginToken;

--MYSQL_CUSTOM_STATEMENT_DELIMITER

CREATE PROCEDURE AddLoginToken (
	IN inUserID CHAR(36),
    IN inToken TEXT,
    IN inTokenType VARCHAR(100),
	IN inExpirationTime BIGINT
)
BEGIN

    INSERT INTO login_tokens (user_id, token, token_type, expires_at) 
    VALUES (inUserID, inToken, inTokenType, UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL inExpirationTime MINUTE )));
    
END;
