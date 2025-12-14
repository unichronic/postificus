-- Users and Auth
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

-- The "Drafts" table (Hot edits)
CREATE TABLE IF NOT EXISTS drafts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INT REFERENCES users(id),
    title VARCHAR(255),
    content TEXT, -- Markdown
    last_saved_at TIMESTAMP,
    is_published BOOLEAN DEFAULT FALSE
);

-- The "Publish Logs" (Audit Trail)
CREATE TABLE IF NOT EXISTS publish_logs (
    id SERIAL PRIMARY KEY,
    draft_id UUID REFERENCES drafts(id),
    platform VARCHAR(50), -- 'medium', 'linkedin'
    status VARCHAR(20),   -- 'queued', 'processing', 'success', 'failed'
    external_url TEXT,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- User Credentials (Encrypted/Stored for Automation)
CREATE TABLE IF NOT EXISTS user_credentials (
    user_id INT REFERENCES users(id),
    platform VARCHAR(50), -- 'medium', 'linkedin', 'devto'
    credentials JSONB,
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, platform)
);
