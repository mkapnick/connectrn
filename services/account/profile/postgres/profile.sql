-- Connect to our database
\connect birrdi

-- Create profiles table
CREATE TABLE profiles (
    id varchar(255) PRIMARY KEY,
    account_id varchar(255) NOT NULL,
    name varchar(255),
    profile_rel_path varchar(255),
    preferences jsonb,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);
