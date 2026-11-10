CREATE TABLE todo_app.tasks(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES todo_app.users(id) ON DELETE CASCADE,
    version BIGINT NOT NULL DEFAULT 1,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMPTZ,

    CONSTRAINT tasks_completed_consistency CHECK(
            (completed IS FALSE AND completed_at IS NULL)
            OR
            (completed IS TRUE AND completed_at IS NOT NULL AND completed_at > created_at)),

    CONSTRAINT tasks_title_length CHECK(char_length(title) BETWEEN 3 AND 100),
    CONSTRAINT tasks_description_length CHECK(char_length(description) BETWEEN 3 AND 1000));

CREATE INDEX tasks_user_id_idx ON todo_app.tasks(user_id);
CREATE INDEX tasks_created_at_id_idx ON todo_app.tasks(created_at, id);
