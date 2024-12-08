package nilcheck_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/jtbonhomme/go-nilcheck"
)

func TestAnalyzer(t *testing.T) { //nolint:paralleltest
	// Find the testdata directory
	testdata := analysistest.TestData()

	// Run the Analyzer against the test files
	analysistest.Run(t, testdata, nilcheck.Analyzer, "a")
}
