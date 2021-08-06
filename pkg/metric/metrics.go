package metric

import "github.com/prometheus/client_golang/prometheus"

var (
	MethodDuration = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:        "broker_call_duration",
		Help:        "calculating the latency of grpc calls",
	},[]string{"method"})

	ActiveSubscribers = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "broker_active_subscribers",
		Help: "number of active subscribers in broker",
	})
	MethodCalls = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "broker_active_subscribers",
		Help: "number of active subscribers in broker",
	},[]string{"method"})
)

func init(){
	prometheus.Register(MethodDuration)
	prometheus.Register(ActiveSubscribers)
	prometheus.Register(MethodCalls)
}
