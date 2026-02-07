-- db/init.sql
-- Initialize users table used by the example application

CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL
);
