package parser_test

import (
	"testing"

	"../parser"
	"github.com/go-test/deep"
)

func TestParse(t *testing.T) {
	var expected = []parser.ServiceMetrics{
		{"2019-11-01 15:16:59.000000", "0c783358576b", "single-proxy", "0.00%", "2.723MiB", "1.953GiB", "0.14%"},
		{"2019-11-01 15:16:59.000000", "3b1681f84c07", "single-web", "0.00%", "31.96MiB", "1.953GiB", "1.60%"},
		{"2019-11-01 15:16:59.000000", "423306a17a71", "single-api", "0.01%", "47.65MiB", "1.953GiB", "2.38%"},
		{"2019-11-01 15:17:06.000000", "0c783358576b", "single-proxy", "0.01%", "2.723MiB", "1.953GiB", "0.14%"},
		{"2019-11-01 15:17:06.000000", "3b1681f84c07", "single-web", "0.01%", "31.96MiB", "1.953GiB", "1.60%"},
		{"2019-11-01 15:17:06.000000", "423306a17a71", "single-api", "0.00%", "47.65MiB", "1.953GiB", "2.38%"},
		{"2019-11-01 15:17:14.000000", "0c783358576b", "single-proxy", "0.00%", "2.723MiB", "1.953GiB", "0.14%"},
		{"2019-11-01 15:17:14.000000", "3b1681f84c07", "single-web", "0.00%", "31.96MiB", "1.953GiB", "1.60%"},
		{"2019-11-01 15:17:14.000000", "423306a17a71", "single-api", "0.01%", "47.65MiB", "1.953GiB", "2.38%"},
		{"2019-11-01 15:17:21.000000", "0c783358576b", "single-proxy", "0.00%", "2.723MiB", "1.953GiB", "0.14%"},
		{"2019-11-01 15:17:21.000000", "3b1681f84c07", "single-web", "0.00%", "31.96MiB", "1.953GiB", "1.60%"},
		{"2019-11-01 15:17:21.000000", "423306a17a71", "single-api", "0.00%", "47.65MiB", "1.953GiB", "2.38%"},
		{"2019-11-01 15:17:28.000000", "0c783358576b", "single-proxy", "0.00%", "2.723MiB", "1.953GiB", "0.14%"},
		{"2019-11-01 15:17:28.000000", "3b1681f84c07", "single-web", "0.00%", "31.96MiB", "1.953GiB", "1.60%"},
		{"2019-11-01 15:17:28.000000", "423306a17a71", "single-api", "0.00%", "47.65MiB", "1.953GiB", "2.38%"},
	}

	metrics := parser.Parse("testdata/2019-11-01.log")
	length := len(metrics)
	var expectedNumber int = 15
	if length != expectedNumber {
		t.Errorf("The number of parsed metrics should be %d. Actual %d", expectedNumber, length)
	}

	if diff := deep.Equal(expected, metrics); diff != nil {
		t.Error(diff)
	}
}
