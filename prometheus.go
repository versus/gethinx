package main

import "github.com/prometheus/client_golang/prometheus"

//TODO:  добавить метрики прометеуса https://stackoverflow.com/questions/37611754/how-to-push-metrics-to-prometheus-using-client-golang

var (
	cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU.",
	})
	hdFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "hd_errors_total",
		Help: "Number of hard-disk errors.",
	})
)

func init() {
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(hdFailures)
}
