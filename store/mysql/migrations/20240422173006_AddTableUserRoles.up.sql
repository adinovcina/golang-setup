-- *****************************************************************************************
-- TABLE user_roles
-- *****************************************************************************************
CREATE TABLE IF NOT EXISTS user_roles (
	id SERIAL,
    -- ID of the user this role belongs to
    user_id CHAR(36) NOT NULL,
    -- ID of the role this user has been assigned to
    role_id BIGINT UNSIGNED NOT NULL,
	PRIMARY KEY (id),
    CONSTRAINT fk_user_roles_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_user_roles_role_id FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE ON UPDATE CASCADE
);
