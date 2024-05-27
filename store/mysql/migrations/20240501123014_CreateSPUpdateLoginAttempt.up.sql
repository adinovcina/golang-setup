-- *****************************************************************************************
-- STORED PROCEDURE UpdateLoginAttempt
-- =========================================================================================
DROP PROCEDURE IF EXISTS UpdateLoginAttempt;

--MYSQL_CUSTOM_STATEMENT_DELIMITER

CREATE PROCEDURE UpdateLoginAttempt(
	IN inUserID CHAR(36),
    IN inMinutes FLOAT,
    IN inMaxLoginFailures INT
)
BEGIN
    
	DECLARE userFailedLoginCount INT;
    
    -- Get the current failed_login_count for the user
    SELECT failed_login_count INTO userFailedLoginCount FROM users WHERE id = inUserID;
    
    -- Increment the failed_login_count
    UPDATE users
    SET failed_login_count = userFailedLoginCount + 1
    WHERE id = inUserID;
    
    -- Check if failed_login_count reaches MaxLoginFailures
    IF userFailedLoginCount + 1 = inMaxLoginFailures THEN
    -- Case 1: Set login_blocked_until to current time plus specified minutes
        UPDATE users
        SET login_blocked_until = DATE_ADD(NOW(), INTERVAL InMinutes MINUTE)
        WHERE id = inUserID;
    END IF;

    -- Check if failed_login_count exceeds MaxLoginFailures
    IF userFailedLoginCount + 1 > inMaxLoginFailures THEN
        -- Case 2: Set login_blocked_until to null and failed_login_count to 1
        -- (due to missing the password for the first time after regaining access)
        UPDATE users
        SET login_blocked_until = null, 
            failed_login_count = 1
        WHERE id = inUserID;
    END IF;
    
    -- Return the updated failed_login_count
    SELECT failed_login_count FROM users WHERE id = inUserID;
    
END;
