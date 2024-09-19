package monitoring

import "github.com/prometheus/client_golang/prometheus"

var (
	TasksProducedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "task",
			Subsystem: "producer",
			Name:      "tasks_produced_total",
			Help:      "Total number of tasks produced.",
		},
	)
)
