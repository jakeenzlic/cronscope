// Package exporter serialises cronscope analysis results into portable
// output formats for downstream consumption.
//
// Supported formats:
//   - JSON  — via ExportJSON, suitable for piping into other tools or
//             storing as artefacts in CI pipelines.
//
// Typical usage:
//
//	entries, _ := parser.ParseLog(rawLog)
//	stats     := aggregator.Aggregate(entries)
//	streaks   := aggregator.DetectFailureStreaks(entries, 2)
//	trend     := aggregator.ComputeTrend(entries, "")
//
//	report := exporter.BuildReport(stats, streaks, trend)
//	if err := exporter.ExportJSON(os.Stdout, report); err != nil {
//		log.Fatal(err)
//	}
package exporter
