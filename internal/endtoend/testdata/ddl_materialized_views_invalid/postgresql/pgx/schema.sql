CREATE TABLE kpis (
  ts TIMESTAMPTZ,
  event_id TEXT NOT NULL
);

CREATE MATERIALIZED VIEW IF NOT EXISTS grouped_kpis AS
SELECT date_trunc('1 day', ts) as day, COUNT(*)
FROM kpis
GROUP BY 1;
