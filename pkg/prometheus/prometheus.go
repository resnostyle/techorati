package prometheus

import "github.com/prometheus/client_golang/prometheus"

// RegisterPromGaugeItems registers a gauge for watched items.
func RegisterPromGaugeItems(labels prometheus.Labels) {
	watchedItems := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "notifier_items",
		Help:        "an item watched by the notifier",
		ConstLabels: labels,
	})

	prometheus.Register(watchedItems)
	watchedItems.Set(1.0)
}

// RegisterPromGaugeTotalItems registers a gauge for total items.
func RegisterPromGaugeTotalItems(total int) {
	totalItems := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "notifier",
		Subsystem: "items",
		Name:      "total",
		Help:      "total number of items to watch",
	})

	prometheus.MustRegister(totalItems)
	totalItems.Set(float64(total))
}