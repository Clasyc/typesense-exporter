package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"time"
)

const DefaultExporterPort = "9101"
const DefaultInsecureSkipVerify = false

func main() {
	if _, ok := os.LookupEnv("TYPESENSE_API_KEY"); !ok {
		log.Fatal("TYPESENSE_API_KEY is not set")
	}
	if _, ok := os.LookupEnv("TYPESENSE_URL"); !ok {
		log.Fatal("TYPESENSE_URL is not set")
	}
	typesenseApiKey := os.Getenv("TYPESENSE_API_KEY")
	typesenseUrl := os.Getenv("TYPESENSE_URL")
	exporterPort := Getenv("EXPORTER_PORT", DefaultExporterPort)
	insecure := GetBoolEnv("INSECURE_SKIP_VERIFY", DefaultInsecureSkipVerify)

	fmt.Println(insecure)

	client := NewClient(typesenseApiKey, typesenseUrl, insecure, 5*time.Second)
	log.Printf("Using: %s", typesenseUrl)
	_, err := client.GetHealth()
	if err != nil {
		log.Printf("WARNING: can't connect to typesense: %s", err)
	}
	exporter := NewExporter(client)
	registry := prometheus.NewRegistry()
	registry.MustRegister(exporter)

	prometheus.Unregister(prometheus.NewGoCollector())
	prometheus.Unregister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.Handle("/", http.RedirectHandler("/metrics", http.StatusMovedPermanently))
	http.Handle(
		"/health", http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		),
	)
	log.Printf("Starting typesense-exporter on %s", exporterPort)
	log.Fatal(http.ListenAndServe(":"+exporterPort, nil))
}
