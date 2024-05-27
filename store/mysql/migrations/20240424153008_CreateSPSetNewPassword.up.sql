-- *****************************************************************************************
-- STORED PROCEDURE SetNewPassword
-- =========================================================================================
DROP PROCEDURE IF EXISTS SetNewPassword;

--MYSQL_CUSTOM_STATEMENT_DELIMITER

CREATE PROCEDURE SetNewPassword (
    IN inUserID	CHAR(36),
    IN inNewPassword VARCHAR(500)
)
BEGIN

    UPDATE users
    SET password = inNewPassword
    WHERE id = inUserID;

END;
