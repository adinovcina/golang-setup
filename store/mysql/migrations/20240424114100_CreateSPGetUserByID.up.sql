-- *****************************************************************************************
-- STORED PROCEDURE GetUserByID
-- =========================================================================================
DROP PROCEDURE IF EXISTS GetUserByID;

--MYSQL_CUSTOM_STATEMENT_DELIMITER

CREATE PROCEDURE GetUserByID (
    IN inID CHAR(36)
)
BEGIN

    SELECT u.id, 
        u.name, 
        u.email, 
        u.phone, 
        u.language, 
        u.active, 
        r.name,
        u.created_at
    FROM users u
    JOIN user_roles ur ON ur.user_id = u.id
    JOIN roles r ON r.id = ur.role_id
    WHERE u.id = inID;

END;
