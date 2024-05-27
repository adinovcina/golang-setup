-- *****************************************************************************************
-- STORED PROCEDURE GetUserByToken
-- =========================================================================================
DROP PROCEDURE IF EXISTS GetUserByToken;

--MYSQL_CUSTOM_STATEMENT_DELIMITER

CREATE PROCEDURE GetUserByToken (
	IN inToken TEXT,
	IN inTtokenType VARCHAR(100)
)
BEGIN
	SET @Now = UNIX_TIMESTAMP(NOW());

	UPDATE users u, (SELECT (expires_at < @Now) AS exp, token FROM login_tokens) lt
	SET u.last_time_logged = @Now
	WHERE lt.exp = false AND lt.token = inToken;

	SELECT (lt.expires_at < @Now) AS expired, 
			u.id, 
			u.name, 
			u.email, 
			u.active, 
			r.name,
			r.id,
			u.language,
			u.failed_login_count,
			u.created_at
	FROM login_tokens lt 
	JOIN users u ON u.id = lt.user_id
	JOIN user_roles ur ON ur.user_id = u.id
	JOIN roles r ON r.id = ur.role_id
	WHERE lt.token = inToken AND lt.token_type = inTtokenType;
END;
