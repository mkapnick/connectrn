-- Connect to our database
\connect connectrn

-- Create accounts table
CREATE TABLE accounts (
    id varchar(255) PRIMARY KEY,
    restaurant_id varchar(255),
    email varchar(255) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

-- Create profiles table
CREATE TABLE profiles (
    id varchar(255) PRIMARY KEY,
    account_id varchar(255) NOT NULL,
    name varchar(255),
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

-- Create restaurants table
CREATE TABLE restaurants (
    id varchar(255) PRIMARY KEY,
    name varchar(255) NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

-- Create tables table
CREATE TABLE tables (
    id varchar(255) PRIMARY KEY,
    restaurant_id varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    num_seats_available int NOT NULL DEFAULT 0,
    num_seats_reserved int NOT NULL DEFAULT 0,
    start_date timestamp NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

CREATE TABLE user_reservations (
    id varchar(255) PRIMARY KEY,
    restaurant_id varchar(255) NOT NULL,
    table_id varchar(255) NOT NULL,
    profile_id varchar(255) NOT NULL,
    num_seats int NOT NULL DEFAULT 0,
    start_date timestamp NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

CREATE TABLE user_reservations_canceled (
    id varchar(255) PRIMARY KEY,
    restaurant_id varchar(255) NOT NULL,
    table_id varchar(255) NOT NULL,
    profile_id varchar(255) NOT NULL,
    num_seats int NOT NULL DEFAULT 0,
    start_date timestamp NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);
