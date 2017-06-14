package mpprometheus

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"

	client "github.com/prometheus/client_golang/api"
	api_v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type MetricValue struct {
	Name  string  `json:"name"`
	Time  int64   `json:"time"`
	Value float64 `json:"value"`
}

func Do() {
	var (
		optMackerelServiceName    string
		optPrometheusApiQuery     string
		optPrometheusMetricLabels string
		optDebug                  bool
		optDryRun                 bool
		optHelp                   bool
	)

	flag.StringVar(&optMackerelServiceName, "s", "", "Service Name where to send metrics, or Host")
	flag.StringVar(&optPrometheusApiQuery, "q", "", "Query string to retrieve metrics")
	flag.StringVar(&optPrometheusMetricLabels, "l", "", "Labels to map on metric names")
	flag.BoolVar(&optDebug, "d", false, "Debug mode")
	flag.BoolVar(&optDryRun, "n", false, "Not send metrics to Mackerel, just check metrics with -d")
	flag.BoolVar(&optHelp, "h", false, "Show this help message")
	flag.Parse()

	if optHelp {
		flag.PrintDefaults()
		return
	}

	if os.Getenv("DEBUG") != "" {
		optDebug = true
	}

	optPrometheusApiServer := os.Getenv("PROMETHEUS_API_SERVER")

	if optPrometheusApiServer == "" {
		log.Fatal("Please set PROMETHEUS_API_SERVER environment variable to get metrics.")
	}
	if optPrometheusApiQuery == "" {
		log.Fatal("Please specify a query string to retrieve metrics.")
	}
	if optPrometheusMetricLabels == "" {
		log.Fatal("Please specify labels to map on metric names.")
	}

	c, err := client.NewClient(client.Config{Address: optPrometheusApiServer})
	if err != nil {
		log.Fatal(err)
	}
	api := api_v1.NewAPI(c)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	val, err := api.Query(ctx, optPrometheusApiQuery, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	vectors := val.(model.Vector)

	if optDebug {
		PrintInJSON(os.Stdout, vectors)
	}

	var metrics []MetricValue
	labels := strings.Split(optPrometheusMetricLabels, ".")

	for _, vector := range vectors {
		var names []string

		for _, label := range labels {
			if value, ok := vector.Metric[model.LabelName(label)]; ok {
				names = append(names, string(value))
			} else {
				names = append(names, label)
			}
		}

		metrics = append(metrics, MetricValue{
			Name:  strings.Join(names, "."),
			Value: float64(vector.Value),
			Time:  vector.Timestamp.Unix(),
		})
	}

	if optDebug {
		PrintInJSON(os.Stdout, metrics)
	}

	if optDryRun {
		return
	}

	if len(metrics) > 0 {
		log.Printf("[Prometheus]: %s on %s\n", optPrometheusApiQuery,
			model.TimeFromUnix(metrics[0].Time).Time().Format(time.UnixDate))

		if optMackerelServiceName != "" {
			optMackerelApiKey := os.Getenv("MACKEREL_API_KEY")

			if optMackerelApiKey == "" {
				log.Fatal("Please set MACKEREL_API_KEY environment variable.")
			}

			SendMetricsToMackerelService(optMackerelApiKey, optMackerelServiceName, metrics)
		} else {
			SendMetricsToMackerelHost(metrics)
		}
	}
}
