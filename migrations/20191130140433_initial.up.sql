-- superuser

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "citext";

CREATE DOMAIN ts_default as timestamptz NOT NULL DEFAULT NOW();
