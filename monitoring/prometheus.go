package monitoring

import "github.com/prometheus/client_golang/prometheus"

//TODO:  добавить метрики прометеуса https://stackoverflow.com/questions/37611754/how-to-push-metrics-to-prometheus-using-client-golang

var (
	PromLastBlock = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "eth_lastblock_average",
		Help: "Current LastBlock from gethinx.",
	})
	PromRequest = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "eth_request_total",
		Help: "Number of request to proxy.",
	})
	PromResponse = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "eth_response_total",
		Help: "Number of response from proxy.",
	})
)
