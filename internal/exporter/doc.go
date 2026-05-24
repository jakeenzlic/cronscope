// Package exporter provides serialisation helpers for cronscope reports.
//
// Supported formats:
//
//   - JSON  – ExportJSON writes a structured JSON document containing job
//     statistics, failure streaks, and trend data.
//
//   - CSV   – ExportCSV writes a flat comma-separated file with one row per
//     job, suitable for import into spreadsheets or further processing with
//     standard Unix tools.
//
// Both exporters accept an io.Writer so callers can direct output to a file,
// an HTTP response body, or an in-memory buffer as needed.
package exporter
