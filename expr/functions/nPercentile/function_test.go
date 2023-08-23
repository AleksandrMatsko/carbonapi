package nPercentile

import (
	"math"
	"testing"
	"time"

	"github.com/go-graphite/carbonapi/expr/interfaces"
	"github.com/go-graphite/carbonapi/expr/metadata"
	"github.com/go-graphite/carbonapi/expr/types"
	"github.com/go-graphite/carbonapi/pkg/parser"
	th "github.com/go-graphite/carbonapi/tests"
)

var (
	md []interfaces.FunctionMetadata = New("")
)

func init() {
	for _, m := range md {
		metadata.RegisterFunction(m.Name, m.F)
	}
}

func TestNPercentile(t *testing.T) {
	now32 := int64(time.Now().Unix())

	tests := []th.EvalTestItem{
		{
			`nPercentile(metric1,50)`,
			map[parser.MetricRequest][]*types.MetricData{
				{Metric: "metric1", From: 0, Until: 1}: {types.MakeMetricData("metric1", []float64{2, 4, 6, 10, 14, 20, math.NaN()}, 1, now32)},
			},
			[]*types.MetricData{types.MakeMetricData("nPercentile(metric1,50)", []float64{8, 8, 8, 8, 8, 8, 8}, 1, now32).SetTag("nPercentile", "50")},
		},
	}

	for _, tt := range tests {
		testName := tt.Target
		t.Run(testName, func(t *testing.T) {
			eval := th.EvaluatorFromFunc(md[0].F)
			th.TestEvalExpr(t, eval, &tt)
		})
	}

}
