package mpprometheus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func SendMetricsToMackerelService(optApiKey string, optServiceName string, metrics []MetricValue) {
	data, _ := json.Marshal(metrics)

	client := &http.Client{}

	url := fmt.Sprintf("https://mackerel.io/api/v0/services/%s/tsdb", optServiceName)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(string(data))))
	req.Header.Add("X-Api-Key", optApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal(fmt.Sprintf("Mackerel API: %s: %s", url, resp.Status))
	}

	log.Printf("[Prometheus]: Sent to \"%s\" Service\n", optServiceName)
}
