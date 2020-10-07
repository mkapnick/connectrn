-- Connect to our database
\connect trrip

-- Create account table
CREATE TABLE accounts (
    id varchar(255) PRIMARY KEY,
    email varchar(255) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    enabled boolean DEFAULT TRUE,
    stripe_customer_id varchar(255) NOT NULL,
    stripe_expires_at timestamp NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

-- Create account table
CREATE TABLE company_square_subscription (
    id varchar(255) PRIMARY KEY,
    company_id varchar(255) NOT NULL,
    square_plan_id varchar(255),
    square_customer_id varchar(255),
    square_card_id varchar(255),
    square_location_id varchar(255),
    square_idempotency_key varchar(255),
    square_tax_percentage varchar(255),
    square_current_subscription jsonb,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

-- Create account table
CREATE TABLE square_subscription (
    id varchar(255) PRIMARY KEY,
    square_catalog_object_id varchar(255) NOT NULL,
    square_catalog_object_version int NOT NULL,
    square_catalog_object_type varchar(255) NOT NULL,
    square_catalog_object jsonb NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);
