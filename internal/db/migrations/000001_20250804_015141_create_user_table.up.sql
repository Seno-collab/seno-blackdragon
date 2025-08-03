CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "user" (
  id  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),  
  full_name TEXT NOT NULL,
  bio  TEXT,
  email TEXT,
  password_hash TEXT,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);