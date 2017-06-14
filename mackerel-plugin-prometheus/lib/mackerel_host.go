package mpprometheus

import (
	"fmt"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

type PrometheusPlugin struct {
	Prefix  string
	Metrics []MetricValue
}

func (p PrometheusPlugin) getPrefix() string {
	var prefix []string
	for i, metric := range p.Metrics {
		names := strings.Split(metric.Name, ".")
		names = names[:(len(names) - 1)]
		if i == 0 {
			prefix = names
			continue
		}
		for j, name := range names {
			if prefix[j] != name {
				prefix[j] = "#"
			}
		}
	}
	return strings.Join(prefix, ".")
}

func (p PrometheusPlugin) MetricKeyPrefix() string {
	return "prometheus"
}

func (p PrometheusPlugin) GraphDefinition() map[string]mp.Graphs {
	metrics := func() []mp.Metrics {
		var metrics []mp.Metrics
		var labels []string
		for _, metric := range p.Metrics {
			names := strings.Split(metric.Name, ".")
			name := names[len(names)-1]
			if exists, _ := InArray(name, labels); !exists {
				metrics = append(metrics, mp.Metrics{
					Name:         name,
					Label:        name,
					AbsoluteName: true,
				})
				labels = append(labels, name)
			}
		}
		return metrics
	}()

	return map[string]mp.Graphs{
		p.Prefix: mp.Graphs{
			Label:   fmt.Sprintf("Prometheus(%s)", p.Prefix),
			Unit:    "float",
			Metrics: metrics,
		},
	}
}

func (p PrometheusPlugin) FetchMetrics() (map[string]interface{}, error) {
	metrics := make(map[string]interface{})

	for _, metric := range p.Metrics {
		metrics[metric.Name] = metric.Value
	}

	return metrics, nil
}

func SendMetricsToMackerelHost(metrics []MetricValue) {
	var mpPrometheus PrometheusPlugin

	mpPrometheus.Metrics = metrics
	mpPrometheus.Prefix = mpPrometheus.getPrefix()

	helper := mp.NewMackerelPlugin(mpPrometheus)

	helper.Run()
}
