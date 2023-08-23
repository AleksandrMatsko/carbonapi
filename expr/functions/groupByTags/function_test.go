package groupByTags

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/go-graphite/carbonapi/expr/functions/aggregate"
	"github.com/go-graphite/carbonapi/expr/interfaces"
	"github.com/go-graphite/carbonapi/expr/metadata"
	"github.com/go-graphite/carbonapi/expr/types"
	"github.com/go-graphite/carbonapi/pkg/parser"
	th "github.com/go-graphite/carbonapi/tests"
)

var (
	md []interfaces.FunctionMetadata = New("")
	s  []interfaces.FunctionMetadata = aggregate.New("")
)

func init() {
	for _, m := range s {
		metadata.RegisterFunction(m.Name, m.F)
	}
	for _, m := range md {
		metadata.RegisterFunction(m.Name, m.F)
	}
}

func TestGroupByTags(t *testing.T) {
	now32 := int64(time.Now().Unix())

	tests := []th.MultiReturnEvalTestItem{
		{
			`groupByTags(metric1.foo.*, "avg", "dc")`,
			map[parser.MetricRequest][]*types.MetricData{
				{Metric: "metric1.foo.*", From: 0, Until: 1}: {
					types.MakeMetricData("metric1.foo;cpu=cpu1;dc=dc1", []float64{1, math.NaN(), 3, 4, math.NaN()}, 1, now32),
					types.MakeMetricData("metric1.foo;cpu=cpu2;dc=dc1", []float64{6, 7, 8, 9, math.NaN()}, 1, now32),
					types.MakeMetricData("metric1.foo;cpu=cpu3;dc=dc1", []float64{11, 12, 13, 14, math.NaN()}, 1, now32),
					types.MakeMetricData("metric1.foo;cpu=cpu4;dc=dc1", []float64{7, 8, 9, 10, math.NaN()}, 1, now32),
				},
			},
			"groupByTags",
			map[string][]*types.MetricData{
				"avg;dc=dc1": {types.MakeMetricData("avg;dc=dc1", []float64{6.25, 9, 8.25, 9.25, math.NaN()}, 1, now32)},
			},
		},
		{
			`groupByTags(metric1.foo.*, "sum", "dc")`,
			map[parser.MetricRequest][]*types.MetricData{
				{Metric: "metric1.foo.*", From: 0, Until: 1}: {
					types.MakeMetricData("metric1.foo;cpu=cpu1;dc=dc1", []float64{1, 2, 3, 4, 5}, 1, now32),
					types.MakeMetricData("metric1.foo;cpu=cpu2;dc=dc1", []float64{6, 7, 8, 9, 10}, 1, now32),
					types.MakeMetricData("metric1.foo;cpu=cpu3;dc=dc1", []float64{11, 12, 13, 14, 15}, 1, now32),
					types.MakeMetricData("metric1.foo;cpu=cpu4;dc=dc1", []float64{7, 8, 9, 10, 11}, 1, now32),
				},
			},
			"groupByTags",
			map[string][]*types.MetricData{
				"sum;dc=dc1": {types.MakeMetricData("sum;dc=dc1", []float64{25, 29, 33, 37, 41}, 1, now32)},
			},
		},
		{
			`groupByTags(metric1.foo.*, "sum", "name", "dc")`,
			map[parser.MetricRequest][]*types.MetricData{
				{Metric: "metric1.foo.*", From: 0, Until: 1}: {
					types.MakeMetricData("metric1.foo;cpu=cpu1;dc=dc1", []float64{1, 2, 3, 4, 5}, 1, now32),
					types.MakeMetricData("metric1.foo;cpu=cpu2;dc=dc1", []float64{6, 7, 8, 9, 10}, 1, now32),
					types.MakeMetricData("metric2.foo;cpu=cpu3;dc=dc1", []float64{11, 12, 13, 14, 15}, 1, now32),
					types.MakeMetricData("metric2.foo;cpu=cpu4;dc=dc1", []float64{7, 8, 9, 10, 11}, 1, now32),
				},
			},
			"groupByTags",
			map[string][]*types.MetricData{
				"metric1.foo;dc=dc1": {types.MakeMetricData("metric1.foo;dc=dc1", []float64{7, 9, 11, 13, 15}, 1, now32)},
				"metric2.foo;dc=dc1": {types.MakeMetricData("metric2.foo;dc=dc1", []float64{18, 20, 22, 24, 26}, 1, now32)},
			},
		},
		{
			`groupByTags(metric1.foo.*, "diff", "dc")`,
			map[parser.MetricRequest][]*types.MetricData{
				{Metric: "metric1.foo.*", From: 0, Until: 1}: {
					types.MakeMetricData("metric1.foo;cpu=cpu1;dc=dc1", []float64{1, 2, 3, 4, 5}, 1, now32),
					types.MakeMetricData("metric1.foo;cpu=cpu2;dc=dc1", []float64{6, 7, 8, 9, 10}, 1, now32),
				},
			},
			"groupByTags",
			map[string][]*types.MetricData{
				"diff;dc=dc1": {types.MakeMetricData("diff;dc=dc1", []float64{-5, -5, -5, -5, -5}, 1, now32)},
			},
		},
		{
			`groupByTags(metric1.foo.*, "sum", "dc", "cpu", "rack")`,
			map[parser.MetricRequest][]*types.MetricData{
				{Metric: "metric1.foo.*", From: 0, Until: 1}: {
					types.MakeMetricData("metric1.foo;cpu=cpu1;dc=dc1", []float64{1, 2, 3, 4, 5}, 1, now32),
					types.MakeMetricData("metric1.foo;cpu=cpu2;dc=dc1", []float64{6, 7, 8, 9, 10}, 1, now32),
					types.MakeMetricData("metric1.foo;cpu=cpu3;dc=dc1", []float64{11, 12, 13, 14, 15}, 1, now32),
					types.MakeMetricData("metric1.foo;cpu=cpu4;dc=dc1", []float64{7, 8, 9, 10, 11}, 1, now32),
				},
			},
			"groupByTags",
			map[string][]*types.MetricData{
				"sum;cpu=cpu1;dc=dc1;rack=": {types.MakeMetricData("sum;cpu=cpu1;dc=dc1;rack=", []float64{1, 2, 3, 4, 5}, 1, now32)},
				"sum;cpu=cpu2;dc=dc1;rack=": {types.MakeMetricData("sum;cpu=cpu2;dc=dc1;rack=", []float64{6, 7, 8, 9, 10}, 1, now32)},
				"sum;cpu=cpu3;dc=dc1;rack=": {types.MakeMetricData("sum;cpu=cpu3;dc=dc1;rack=", []float64{11, 12, 13, 14, 15}, 1, now32)},
				"sum;cpu=cpu4;dc=dc1;rack=": {types.MakeMetricData("sum;cpu=cpu4;dc=dc1;rack=", []float64{7, 8, 9, 10, 11}, 1, now32)},
			},
		},
	}

	for _, tt := range tests {
		testName := tt.Target
		t.Run(testName, func(t *testing.T) {
			eval := th.EvaluatorFromFuncWithMetadata(metadata.FunctionMD.Functions)
			th.TestMultiReturnEvalExpr(t, eval, &tt)
		})
	}

}

func BenchmarkGroupByTags(b *testing.B) {
	target := `groupByTags(metric1.foo.*, "sum", "dc", "cpu", "rack")`
	metrics := map[parser.MetricRequest][]*types.MetricData{
		{Metric: "metric1.foo.*", From: 0, Until: 1}: {
			types.MakeMetricData("metric1.foo;cpu=cpu1;dc=dc1", []float64{1, 2, 3, 4, 5}, 1, 1),
			types.MakeMetricData("metric1.foo;cpu=cpu2;dc=dc1;rack=", []float64{6, 7, 8, 9, 10}, 1, 1),
			types.MakeMetricData("metric1.foo;cpu=cpu3;dc=dc1;rack=", []float64{11, 12, 13, 14, 15}, 1, 1),
			types.MakeMetricData("metric1.foo;cpu=cpu4;dc=dc1;rack=", []float64{7, 8, 9, 10, 11}, 1, 1),
		},
	}

	eval := th.EvaluatorFromFuncWithMetadata(metadata.FunctionMD.Functions)
	exp, _, err := parser.ParseExpr(target)
	if err != nil {
		b.Fatalf("failed to parse %s: %+v", target, err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		g, err := eval.Eval(context.Background(), exp, 0, 1, metrics)
		if err != nil {
			b.Fatalf("failed to eval %s: %+v", target, err)
		}
		_ = g
	}
}
