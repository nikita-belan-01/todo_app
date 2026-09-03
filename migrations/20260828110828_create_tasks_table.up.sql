CREATE TABLE todo_app.tasks(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES todo_app.users(id),
    version BIGINT NOT NULL DEFAULT 1,
    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMPTZ,

    CHECK (
            (completed IS FALSE AND completed_at IS NULL)
            OR
            (completed IS TRUE AND completed_at IS NOT NULL AND completed_at > created_at)));