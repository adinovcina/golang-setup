-- *****************************************************************************************
-- STORED PROCEDURE ResetFailedLoginCount
-- =========================================================================================
DROP PROCEDURE IF EXISTS ResetFailedLoginCount;

--MYSQL_CUSTOM_STATEMENT_DELIMITER

CREATE PROCEDURE ResetFailedLoginCount(
	IN inUserID CHAR(36)
)
BEGIN

	UPDATE users 
	SET failed_login_count = 0,
		login_blocked_until = null
	WHERE id = inUserID;

END;
