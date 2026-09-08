CREATE table todo_app.users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version BIGINT NOT NULL DEFAULT 1,
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    phone_number VARCHAR(15) NOT NULL,
    
    CHECK (char_length(name) BETWEEN 3 AND 100),
    CHECK (char_length(surname) BETWEEN 3 AND 100),
    CHECK (phone_number ~ '^\+[0-9]+$'
            AND char_length(phone_number) BETWEEN 10 AND 15));
