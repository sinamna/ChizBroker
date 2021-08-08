package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	MethodDuration = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:        "method_duration",
		Help:        "calculating the latency of grpc calls",
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	},[]string{"method"})

	ActiveSubscribers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "broker_active_subscribers",
		Help: "number of active subscribers in broker",
	})
	MethodCalls = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "method_count",
		Help: "number of method calls in broker",
	},[]string{"method"})

	MethodError = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "method_error_count",
		Help: "counter error of each method",
	},[]string{"method"})

)

func init(){
	prometheus.Register(MethodDuration)
	prometheus.Register(ActiveSubscribers)
	prometheus.Register(MethodCalls)
	prometheus.Register(MethodError)
}
