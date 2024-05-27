-- *****************************************************************************************
-- STORED PROCEDURE UpdateUser
-- =========================================================================================
DROP PROCEDURE IF EXISTS UpdateUser;

--MYSQL_CUSTOM_STATEMENT_DELIMITER
CREATE PROCEDURE UpdateUser (
    IN inUserID CHAR(36),
    IN inName VARCHAR(150),
    IN inPhone VARCHAR(150)
) 
BEGIN

    UPDATE
        users
    SET
        name = inName,
        phone = inPhone
    WHERE
        id = inUserID;

    SELECT
        u.id,
        u.name,
        u.email,
        u.phone,
        u.language,
        r.name
    FROM users u
    JOIN user_roles ur ON ur.user_id = u.id
    JOIN roles r ON r.id = ur.role_id
    WHERE u.id = inUserID;

END;
