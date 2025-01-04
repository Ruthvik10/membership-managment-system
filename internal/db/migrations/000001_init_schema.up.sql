CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    phone TEXT UNIQUE NOT NULL,
    address TEXT,
    status INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE sports (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT
);

CREATE TABLE memberships (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    member_id UUID REFERENCES members(id) ON DELETE CASCADE,
    sport_id SERIAL REFERENCES sports(id) ON DELETE CASCADE,
    type TEXT NOT NULL CHECK(type IN ('membership', 'training')),
    start_date TIMESTAMPTZ NOT NULL,
    due_date TIMESTAMPTZ NOT NULL,
    status INTEGER NOT NULL DEFAULT 1,
    fees NUMERIC(10, 2) NOT NULL
);

ALTER TABLE memberships ADD UNIQUE (member_id, sport_id, type);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    membership_id UUID REFERENCES memberships(id),
    amount NUMERIC(10, 2) NOT NULL,
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status INTEGER NOT NULL DEFAULT 1,
    payment_link TEXT
);