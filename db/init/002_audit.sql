-- 002_audit.sql (MUP audit)

CREATE TABLE IF NOT EXISTS audit_log (
  id BIGSERIAL PRIMARY KEY,
  table_name TEXT NOT NULL,
  action TEXT NOT NULL,              -- INSERT/UPDATE/DELETE
  changed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  row_pk TEXT NULL,                  -- vrednost "id" (ako postoji)
  old_data JSONB NULL,
  new_data JSONB NULL
);

CREATE OR REPLACE FUNCTION audit_if_modified() RETURNS TRIGGER AS $$
DECLARE
  pk TEXT;
BEGIN
  IF (TG_OP = 'DELETE') THEN
    pk := to_jsonb(OLD)->>'id';
    INSERT INTO audit_log(table_name, action, row_pk, old_data, new_data)
    VALUES (TG_TABLE_NAME, TG_OP, pk, to_jsonb(OLD), NULL);
    RETURN OLD;

  ELSIF (TG_OP = 'UPDATE') THEN
    pk := to_jsonb(NEW)->>'id';
    INSERT INTO audit_log(table_name, action, row_pk, old_data, new_data)
    VALUES (TG_TABLE_NAME, TG_OP, pk, to_jsonb(OLD), to_jsonb(NEW));
    RETURN NEW;

  ELSIF (TG_OP = 'INSERT') THEN
    pk := to_jsonb(NEW)->>'id';
    INSERT INTO audit_log(table_name, action, row_pk, old_data, new_data)
    VALUES (TG_TABLE_NAME, TG_OP, pk, NULL, to_jsonb(NEW));
    RETURN NEW;
  END IF;

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_audit_citizens ON citizens;
CREATE TRIGGER trg_audit_citizens
AFTER INSERT OR UPDATE OR DELETE ON citizens
FOR EACH ROW EXECUTE FUNCTION audit_if_modified();

DROP TRIGGER IF EXISTS trg_audit_service_requests ON service_requests;
CREATE TRIGGER trg_audit_service_requests
AFTER INSERT OR UPDATE OR DELETE ON service_requests
FOR EACH ROW EXECUTE FUNCTION audit_if_modified();

DROP TRIGGER IF EXISTS trg_audit_appointments ON appointments;
CREATE TRIGGER trg_audit_appointments
AFTER INSERT OR UPDATE OR DELETE ON appointments
FOR EACH ROW EXECUTE FUNCTION audit_if_modified();

DROP TRIGGER IF EXISTS trg_audit_certificates ON certificates;
CREATE TRIGGER trg_audit_certificates
AFTER INSERT OR UPDATE OR DELETE ON certificates
FOR EACH ROW EXECUTE FUNCTION audit_if_modified();

DROP TRIGGER IF EXISTS trg_audit_payments ON payments;
CREATE TRIGGER trg_audit_payments
AFTER INSERT OR UPDATE OR DELETE ON payments
FOR EACH ROW EXECUTE FUNCTION audit_if_modified();
