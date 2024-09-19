package monitoring

import "github.com/prometheus/client_golang/prometheus"

var (
	TasksProcessedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "task",
			Subsystem: "consumer",
			Name:      "tasks_processed_total",
			Help:      "Total number of tasks processed.",
		},
	)

	TasksDoneTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "task",
			Subsystem: "consumer",
			Name:      "tasks_done_total",
			Help:      "Total number of tasks done.",
		},
	)

	TasksValueTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "task",
			Subsystem: "consumer",
			Name:      "tasks_value_total",
			Help:      "Total number of the processed entries for each task type.",
		},
		[]string{"task_type"},
	)

	TasksValueSum = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "task",
			Subsystem: "consumer",
			Name:      "tasks_value_sum",
			Help:      "Total sum of the 'value' field for each task type.",
		},
		[]string{"task_type"},
	)
)
