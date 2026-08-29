CREATE EXTENSION "uuid-ossp";

CREATE table todo_app.users(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    version BIGINT NOT NULL DEFAULT 1,
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    phone_number VARCHAR(15) NOT NULL,
    
    CHECK (char_length(todo_app.users.name) BETWEEN 3 AND 100),
    CHECK (char_length(todo_app.users.surname) BETWEEN 3 AND 100),
    CHECK (phone_number ~ '^\+[0-9]+$'
            AND char_length(todo_app.users.phone_number) BETWEEN 10 AND 15));
