-- *****************************************************************************************
-- STORED PROCEDURE ActivateUser
-- =========================================================================================
DROP PROCEDURE IF EXISTS ActivateUser;

--MYSQL_CUSTOM_STATEMENT_DELIMITER

CREATE PROCEDURE ActivateUser (
    IN inUserID	CHAR(36)
)
BEGIN

    UPDATE users
    SET active = NOT active
    WHERE id = inUserID;

END