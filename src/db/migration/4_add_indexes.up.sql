CREATE INDEX IF NOT EXISTS idx_tasks_username ON tasks(username);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at DESC);
