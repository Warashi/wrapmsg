package wrapmsg_test

import (
	"testing"

	"github.com/Warashi/wrapmsg"

	"github.com/gostaticanalysis/testutil"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.RunWithSuggestedFixes(t, testdata, wrapmsg.Analyzer, "a")
}
