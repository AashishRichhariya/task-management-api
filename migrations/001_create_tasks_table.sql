-- Create tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraint to enforce valid status values
    CONSTRAINT valid_status CHECK (status IN ('pending', 'in_progress', 'completed', 'closed'))
);

-- Index for filtering by status
CREATE INDEX idx_tasks_status ON tasks(status);

-- Index for sorting by created_at
CREATE INDEX idx_tasks_created_at ON tasks(created_at DESC);