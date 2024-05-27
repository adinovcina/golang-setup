-- *****************************************************************************************
-- STORED PROCEDURE SetPassword
-- =========================================================================================
DROP PROCEDURE IF EXISTS SetPassword;

--MYSQL_CUSTOM_STATEMENT_DELIMITER

CREATE PROCEDURE SetPassword (
    IN inUserID CHAR(36),
    IN inPassword VARCHAR(500),
    IN inToken VARCHAR(100)
)
BEGIN

    UPDATE users 
    SET password = inPassword,
        active = 1,
        email_verified = 1,
        terms_accepted = 1
    WHERE id = inUserID;

    DELETE FROM password_tokens
    WHERE token = inToken;

    SELECT u.id, 
        u.name, 
        u.email, 
        u.phone, 
        u.language, 
        u.active, 
        u.phone, 
        u.password,
        r.name,
        r.id,
        u.created_at
    FROM users u
    JOIN user_roles ur ON ur.user_id = u.id
    JOIN roles r ON r.id = ur.role_id
    WHERE u.id = inUserID;

END;
