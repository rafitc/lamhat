-- Users Table 

CREATE TABLE users (
    id SERIAL PRIMARY KEY, -- auto-incrementing primary key
    user_id UUID DEFAULT gen_random_uuid() UNIQUE, -- unique identifier for user
    email_id VARCHAR(255) NOT NULL UNIQUE, -- email address, unique and not null
    is_email_valid BOOLEAN DEFAULT FALSE, -- flag to check if email is valid
    first_name VARCHAR(100), -- user's first name, now nullable
    last_name VARCHAR(100), -- user's last name, now nullable
    auth_key_hash CHAR(64) NOT NULL UNIQUE, -- hashed authentication key, unique and fixed length
    otp CHAR(6), -- OTP for email verification
    otp_generated_at TIMESTAMP, -- timestamp of OTP generation
    is_user_active BOOLEAN DEFAULT TRUE, -- flag to check if user is active
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- timestamp of creation
    last_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- timestamp of last update
);

-- Indexes for better performance
CREATE INDEX idx_users_email_id ON users (email_id);
CREATE INDEX idx_users_is_user_active ON users (is_user_active);
CREATE INDEX idx_users_created_at ON users (created_at);
CREATE INDEX idx_users_otp ON users (otp);
CREATE INDEX idx_users_otp_generated_at ON users (otp_generated_at);

ALTER TABLE users OWNER TO me; -- change owner to superuser

