package go_redis_prometheus_collector

import (
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
)

type redisStatsCollector struct {
	s   StatsGetter
	pms []poolMetric
}

func (c *redisStatsCollector) Describe(descs chan<- *prometheus.Desc) {
	for _, pm := range c.pms {
		descs <- pm.desc
	}
}

func (c *redisStatsCollector) Collect(metrics chan<- prometheus.Metric) {
	stats := c.s.PoolStats()
	for _, pm := range c.pms {
		metrics <- prometheus.MustNewConstMetric(pm.desc, pm.tp, pm.fn(stats))
	}
}

type poolMetric struct {
	desc *prometheus.Desc
	fn   func(stats *redis.PoolStats) float64
	tp   prometheus.ValueType
}

type StatsGetter interface {
	PoolStats() *redis.PoolStats
}

// NewRedisPoolCollector returns new collector for redis pool with
// namespace `ns`.
// `ns` may be empty.
func NewRedisPoolCollector(ns string,
	redisConn StatsGetter) prometheus.Collector {

	return &redisStatsCollector{
		s: redisConn,
		pms: []poolMetric{
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(ns, "redis", "hits"),
					"Number of times free connection was found in the pool.",
					nil, nil),
				fn: func(stats *redis.PoolStats) float64 {
					return float64(stats.Hits)
				},
				tp: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(ns, "redis", "misses"),
					"Number of times free connection was NOT found in the "+
						"pool.",
					nil, nil),
				fn: func(stats *redis.PoolStats) float64 {
					return float64(stats.Misses)
				},
				tp: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(ns, "redis", "timeouts"),
					"Number of times a wait timeout occurred.",
					nil, nil),
				fn: func(stats *redis.PoolStats) float64 {
					return float64(stats.Timeouts)
				},
				tp: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(ns, "redis", "total_conns"),
					"Number of total connections in the pool.",
					nil, nil),
				fn: func(stats *redis.PoolStats) float64 {
					return float64(stats.TotalConns)
				},
				tp: prometheus.GaugeValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(ns, "redis", "idle_conns"),
					"Number of idle connections in the pool.",
					nil, nil),
				fn: func(stats *redis.PoolStats) float64 {
					return float64(stats.IdleConns)
				},
				tp: prometheus.GaugeValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(ns, "redis", "stale_conns"),
					"Number of stale connections removed from the pool.",
					nil, nil),
				fn: func(stats *redis.PoolStats) float64 {
					return float64(stats.StaleConns)
				},
				tp: prometheus.CounterValue,
			},
		},
	}
}
