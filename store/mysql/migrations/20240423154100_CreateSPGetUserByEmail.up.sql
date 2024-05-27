-- *****************************************************************************************
-- STORED PROCEDURE GetUserByEmail
-- =========================================================================================
DROP PROCEDURE IF EXISTS GetUserByEmail;

--MYSQL_CUSTOM_STATEMENT_DELIMITER

CREATE PROCEDURE GetUserByEmail (
    IN inEmail VARCHAR(250)
)
BEGIN

    SELECT u.id, 
        u.name, 
        u.email, 
        u.password, 
        u.active,
        u.failed_login_count,
        u.login_blocked_until
    FROM users u
    WHERE u.email = inEmail;

END;
