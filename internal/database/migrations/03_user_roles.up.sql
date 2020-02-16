CREATE TYPE user_role AS ENUM (
    'admin'
);

CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users,
    role user_role NOT NULL,
    PRIMARY KEY(user_id, role)
);

CREATE INDEX user_roles_user
    ON user_roles (user_id);
