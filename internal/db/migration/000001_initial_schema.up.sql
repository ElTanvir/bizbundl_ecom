-- #############################################################################
-- ## UP MIGRATION (optimized)
-- #############################################################################

CREATE TABLE "leads" (
    "id"          TEXT PRIMARY KEY,
    "first_name"  TEXT NOT NULL,
    "last_name"   TEXT NOT NULL,
    "email"       TEXT NOT NULL,
    "phone"       TEXT,
    "company"     TEXT,
    "message"     TEXT,
    "status"      TEXT NOT NULL DEFAULT 'new' CHECK (status IN ('new', 'contacted', 'qualified', 'converted', 'lost')),
    "source"      TEXT,
    "created_at"  TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Fast dedupe and filtering
CREATE INDEX IF NOT EXISTS idx_leads_email ON leads (email);

-- For search
CREATE INDEX IF NOT EXISTS idx_leads_first_name_lower ON leads (lower(first_name));
CREATE INDEX IF NOT EXISTS idx_leads_last_name_lower  ON leads (lower(last_name));
CREATE INDEX IF NOT EXISTS idx_leads_email_lower      ON leads (lower(email));

-- For filter + sort
CREATE INDEX IF NOT EXISTS idx_leads_status_created_at ON leads (status, created_at DESC);

-- For export/all queries
CREATE INDEX IF NOT EXISTS idx_leads_created_at ON leads (created_at DESC);
