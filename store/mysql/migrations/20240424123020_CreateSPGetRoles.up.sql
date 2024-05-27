-- *****************************************************************************************
-- STORED PROCEDURE GetRoles
-- =========================================================================================
DROP PROCEDURE IF EXISTS GetRoles;

--MYSQL_CUSTOM_STATEMENT_DELIMITER

CREATE PROCEDURE GetRoles (
	IN inUserID CHAR(36)
)
BEGIN

    SELECT r.id,
        r.name
    FROM users u
    JOIN user_roles ur ON u.id = ur.user_id
    JOIN roles r ON r.id = ur.role_id
    WHERE u.id = inUserID;

END;
