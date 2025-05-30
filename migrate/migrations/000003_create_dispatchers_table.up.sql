CREATE TABLE IF NOT EXISTS dispatchers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    application_id UUID NOT NULL,
    vehicle TEXT NOT NULL,
    license TEXT NOT NULL,
    approved_at TIMESTAMP DEFAULT NOW(),
    isActive BOOLEAN DEFAULT TRUE,
    rating DOUBLE PRECISION DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (application_id) REFERENCES dispatchers_apply(id) ON DELETE CASCADE
);