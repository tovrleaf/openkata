const { Pool } = require('pg');
const logger = require('../logger');

const pool = new Pool({
  connectionString: process.env.DATABASE_URL,
  max: 20,
});

// Executive dashboard summary — called on every page load for each user session.
// Aggregates 30 days of metrics across all dimensions. Takes 800–2000ms under load.
async function getExecutiveSummary(orgId, dateRange) {
  const res = await pool.query(
    `SELECT metric_type, category, SUM(value) AS total, COUNT(*) AS data_points
     FROM daily_metrics
     WHERE org_id = $1 AND recorded_at BETWEEN $2 AND $3
     GROUP BY metric_type, category
     ORDER BY metric_type, total DESC`,
    [orgId, dateRange.start, dateRange.end]
  );
  logger.debug('Executive summary query', { orgId, rows: res.rowCount });
  return res.rows;
}

// Report list — called multiple times per session as users navigate between views.
// Returns the same data until a new report is scheduled or run.
async function getReportList(orgId) {
  const res = await pool.query(
    `SELECT id, name, last_run, schedule, status
     FROM reports
     WHERE org_id = $1
     ORDER BY last_run DESC NULLS LAST`,
    [orgId]
  );
  return res.rows;
}

// KPI benchmarks — called on dashboard render and on each filter change.
// Joins across three tables; p95 is 1.4 seconds on the current dataset size.
async function getKpiBenchmarks(orgId, period) {
  const res = await pool.query(
    `SELECT k.kpi_name, k.target_value, AVG(m.value) AS actual_value
     FROM kpi_targets k
     JOIN daily_metrics m ON m.metric_type = k.metric_type AND m.org_id = k.org_id
     JOIN reporting_periods p ON m.recorded_at BETWEEN p.start_date AND p.end_date
     WHERE k.org_id = $1 AND p.label = $2
     GROUP BY k.kpi_name, k.target_value`,
    [orgId, period]
  );
  return res.rows;
}

module.exports = { getExecutiveSummary, getReportList, getKpiBenchmarks };
