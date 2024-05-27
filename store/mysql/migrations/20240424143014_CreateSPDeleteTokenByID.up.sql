-- *****************************************************************************************
-- STORED PROCEDURE DeleteTokenByID
-- =========================================================================================
DROP PROCEDURE IF EXISTS DeleteTokenByID;

--MYSQL_CUSTOM_STATEMENT_DELIMITER

CREATE PROCEDURE DeleteTokenByID (
	IN inTokenID BIGINT
)
BEGIN

DELETE FROM
	login_tokens
WHERE
	id = inTokenID;

END;