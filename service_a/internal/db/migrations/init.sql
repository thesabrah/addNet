-- Active: 1751321030208@@127.0.0.1@5432@postgres

CREATE TABLE IF NOT EXISTS outbox (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP DEFAULT NOW(),
  sent_at TIMESTAMP,
  value INT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_outbox_sent_at ON outbox(sent_at);