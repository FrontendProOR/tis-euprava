-- 001_schema.sql (MUP-GRAĐANI) - tačno po repository upitima

CREATE TABLE IF NOT EXISTS citizens (
  id TEXT PRIMARY KEY,
  jmbg TEXT NOT NULL UNIQUE,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  date_of_birth TIMESTAMPTZ NOT NULL,
  email TEXT NOT NULL,
  phone_number TEXT NOT NULL,
  address JSONB NOT NULL,              
  created_at TIMESTAMPTZ NOT NULL      
);

CREATE TABLE IF NOT EXISTS service_requests (
  id TEXT PRIMARY KEY,
  citizen_id TEXT NOT NULL REFERENCES citizens(id) ON DELETE CASCADE,
  type TEXT NOT NULL,
  status TEXT NOT NULL,
  submitted_at TIMESTAMPTZ NOT NULL,
  processed_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_service_requests_citizen_id
ON service_requests(citizen_id);

CREATE TABLE IF NOT EXISTS appointments (
  id TEXT PRIMARY KEY,
  citizen_id TEXT NOT NULL REFERENCES citizens(id) ON DELETE CASCADE,
  date_time TIMESTAMPTZ NOT NULL,
  police_station TEXT NOT NULL,
  status TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_appointments_citizen_id
ON appointments(citizen_id);

CREATE TABLE IF NOT EXISTS certificates (
  id TEXT PRIMARY KEY,
  request_id TEXT NOT NULL REFERENCES service_requests(id) ON DELETE CASCADE,
  issued_at TIMESTAMPTZ NOT NULL,
  pdf BYTEA NOT NULL
);

CREATE TABLE IF NOT EXISTS payments (
  id TEXT PRIMARY KEY,
  request_id TEXT NOT NULL REFERENCES service_requests(id) ON DELETE CASCADE,
  amount NUMERIC(12,2) NOT NULL,
  reference TEXT NOT NULL,
  status TEXT NOT NULL,
  paid_at TIMESTAMPTZ NULL
);
