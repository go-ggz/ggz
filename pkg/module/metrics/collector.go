package metrics

import (
	"github.com/go-ggz/ggz/pkg/model"

	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "ggz_"

// Collector implements the prometheus.Collector interface and
// exposes gitea metrics for prometheus
type Collector struct {
	Shortens *prometheus.Desc
	Users    *prometheus.Desc
}

// NewCollector returns a new Collector with all prometheus.Desc initialized
func NewCollector() Collector {
	return Collector{
		Users: prometheus.NewDesc(
			namespace+"users",
			"Number of Users",
			nil, nil,
		),
		Shortens: prometheus.NewDesc(
			namespace+"shortens",
			"Number of Shortens",
			nil, nil,
		),
	}

}

// Describe returns all possible prometheus.Desc
func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Users
	ch <- c.Shortens
}

// Collect returns the metrics with values
func (c Collector) Collect(ch chan<- prometheus.Metric) {
	stats := model.GetStatistic()

	ch <- prometheus.MustNewConstMetric(
		c.Users,
		prometheus.GaugeValue,
		float64(stats.Counter.User),
	)
	ch <- prometheus.MustNewConstMetric(
		c.Shortens,
		prometheus.GaugeValue,
		float64(stats.Counter.Shorten),
	)
}
