// +build go1.13

package stdlib_test

import (
	"testing"

	"github.com/influxdata/flux"
)

func reportStatistics(b *testing.B, stats flux.Statistics) {
	b.ReportMetric(float64(stats.TotalAllocated)/float64(b.N), "total_allocs/op")
	b.ReportMetric(float64(stats.MaxAllocated)/float64(b.N), "max_allocs/op")
}
