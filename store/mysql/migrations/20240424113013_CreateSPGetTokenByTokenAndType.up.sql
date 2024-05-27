-- *****************************************************************************************
-- STORED PROCEDURE GetTokenByTokenAndType
-- =========================================================================================
DROP PROCEDURE IF EXISTS GetTokenByTokenAndType;

--MYSQL_CUSTOM_STATEMENT_DELIMITER

CREATE PROCEDURE GetTokenByTokenAndType (
	IN inToken TEXT,
	IN inTokenType VARCHAR(100)
) 
BEGIN
	SET @Now = UNIX_TIMESTAMP(NOW());

	SELECT (lt.expires_at < @Now) AS expired,
			lt.id,
			lt.user_id,
			lt.token,
			lt.token_type
	FROM login_tokens lt
	JOIN users u ON u.id = lt.user_id
	WHERE lt.token = inToken AND lt.token_type = inTokenType;

END;