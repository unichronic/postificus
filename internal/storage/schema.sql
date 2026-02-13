-- Users and Auth
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

-- Seed a default user for MVP flows (id=1)
INSERT INTO users (id, email, password_hash)
VALUES ('00000000-0000-0000-0000-000000000001', 'demo@postificus.local', 'demo')
ON CONFLICT DO NOTHING;

-- User Profile Details
CREATE TABLE IF NOT EXISTS user_details (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    full_name VARCHAR(255),
    username VARCHAR(100),
    headline VARCHAR(255),
    bio TEXT,
    location VARCHAR(255),
    website TEXT,
    public_email VARCHAR(255),
    skills JSONB,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- The "Drafts" table (Hot edits)
CREATE TABLE IF NOT EXISTS drafts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    title VARCHAR(255),
    content TEXT, -- Markdown
    cover_image TEXT,
    publish_targets JSONB,
    last_saved_at TIMESTAMP,
    is_published BOOLEAN DEFAULT FALSE
);

-- Ensure column exists for older databases
ALTER TABLE drafts
    ADD COLUMN IF NOT EXISTS publish_targets JSONB;

ALTER TABLE drafts
    ADD COLUMN IF NOT EXISTS cover_image TEXT;

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
    user_id UUID REFERENCES users(id),
    platform VARCHAR(50), -- 'medium', 'linkedin', 'devto'
    credentials JSONB,
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, platform)
);

-- Unified Posts (Synced Activity Cache)
CREATE TABLE IF NOT EXISTS unified_posts (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    platform VARCHAR(50) NOT NULL, -- 'medium', 'devto'
    remote_id VARCHAR(255) NOT NULL, -- URL or slug to identify uniqueness
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    status VARCHAR(50), -- 'published', 'draft'
    views INT DEFAULT 0,
    reactions INT DEFAULT 0,
    comments INT DEFAULT 0,
    published_at TIMESTAMP,
    last_synced_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, platform, remote_id)
);
