package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"strconv"
	"sync"
)

type Exporter struct {
	ExporterMetrics
	TypeSense *Client
}

type ExporterMetrics struct {
	Up                         *prometheus.Desc
	SystemCPUxActivePercentage *prometheus.Desc
	SystemDiskTotalBytes       *prometheus.Desc
	SystemDiskUsedBytes        *prometheus.Desc
	SystemMemoryTotalBytes     *prometheus.Desc
	SystemMemoryUsedBytes      *prometheus.Desc
	SystemNetworkReceived      *prometheus.Desc
	SystemNetworkSent          *prometheus.Desc
	TypesenseMemoryActive      *prometheus.Desc
	TypesenseMemoryAllocated   *prometheus.Desc
	TypesenseMemoryFragment    *prometheus.Desc
	TypesenseMemoryMapped      *prometheus.Desc
	TypesenseMemoryMetadata    *prometheus.Desc
	TypesenseMemoryResident    *prometheus.Desc
	TypesenseMemoryRetained    *prometheus.Desc
	ApiStatsOperationLatency   *prometheus.Desc
	ApiStatsOperationRequests  *prometheus.Desc
	ApiStatsEndpointLatency    *prometheus.Desc
	ApiStatsEndpointRequests   *prometheus.Desc
	ApiStatsPendingWrite       *prometheus.Desc
}

func (em *ExporterMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- em.Up
	ch <- em.SystemCPUxActivePercentage
	ch <- em.SystemDiskTotalBytes
	ch <- em.SystemDiskUsedBytes
	ch <- em.SystemMemoryTotalBytes
	ch <- em.SystemMemoryUsedBytes
	ch <- em.SystemNetworkReceived
	ch <- em.SystemNetworkSent
	ch <- em.TypesenseMemoryActive
	ch <- em.TypesenseMemoryAllocated
	ch <- em.TypesenseMemoryFragment
	ch <- em.TypesenseMemoryMapped
	ch <- em.TypesenseMemoryMetadata
	ch <- em.TypesenseMemoryResident
	ch <- em.TypesenseMemoryRetained
	ch <- em.ApiStatsOperationLatency
	ch <- em.ApiStatsOperationRequests
	ch <- em.ApiStatsEndpointLatency
	ch <- em.ApiStatsEndpointRequests
	ch <- em.ApiStatsPendingWrite
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.collectUp(ch)
	var wg sync.WaitGroup
	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		e.collectMetrics(ch)
		wg.Done()
	}(&wg)
	go func(wg *sync.WaitGroup) {
		e.collectStats(ch)
		wg.Done()
	}(&wg)
	wg.Wait()
}

func (e *Exporter) collectMetrics(ch chan<- prometheus.Metric) {
	v, err := e.TypeSense.GetMetrics()
	if err != nil {
		log.Println(err)
		return
	}

	i := 0
	for _, value := range v.SystemCPUxActivePercentage {
		i++
		float, err := value.Float64()
		if err != nil {
			fmt.Println(err)
			continue
		}
		label := strconv.Itoa(i)
		if i == len(v.SystemCPUxActivePercentage) {
			label = "all"
		}
		ch <- prometheus.MustNewConstMetric(
			e.ExporterMetrics.SystemCPUxActivePercentage, prometheus.GaugeValue, percentageToRatio(float),
			label,
		)
	}

	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.SystemDiskTotalBytes, prometheus.GaugeValue, v.SystemDiskTotalBytes,
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.SystemDiskUsedBytes, prometheus.GaugeValue, v.SystemDiskUsedBytes,
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.SystemMemoryTotalBytes, prometheus.GaugeValue, v.SystemMemoryTotalBytes,
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.SystemMemoryUsedBytes, prometheus.GaugeValue, v.SystemMemoryUsedBytes,
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.SystemNetworkReceived, prometheus.GaugeValue, v.SystemNetworkReceivedBytes,
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.SystemNetworkSent, prometheus.GaugeValue, v.SystemNetworkSentBytes,
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.TypesenseMemoryActive, prometheus.GaugeValue, v.TypesenseMemoryActiveBytes,
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.TypesenseMemoryAllocated, prometheus.GaugeValue, v.TypesenseMemoryAllocated,
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.TypesenseMemoryFragment, prometheus.GaugeValue, v.TypesenseMemoryFragment,
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.TypesenseMemoryMapped, prometheus.GaugeValue, v.TypesenseMemoryMapped,
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.TypesenseMemoryMetadata, prometheus.GaugeValue, v.TypesenseMemoryMetadata,
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.TypesenseMemoryResident, prometheus.GaugeValue, v.TypesenseMemoryResident,
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.TypesenseMemoryRetained, prometheus.GaugeValue, v.TypesenseMemoryRetained,
	)
}

func (e *Exporter) collectStats(ch chan<- prometheus.Metric) {
	s, err := e.TypeSense.GetStats()
	if err != nil {
		log.Println(err)
		return
	}

	for name, value := range s.Latency {
		float, err := value.Float64()
		if err != nil {
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			e.ExporterMetrics.ApiStatsEndpointRequests, prometheus.GaugeValue, float, name,
		)
	}
	for name, value := range s.Requests {
		float, err := value.Float64()
		if err != nil {
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			e.ExporterMetrics.ApiStatsEndpointLatency, prometheus.GaugeValue, msToSeconds(float), name,
		)
	}

	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsOperationLatency, prometheus.GaugeValue, msToSeconds(s.DeleteLatency), "delete",
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsOperationRequests, prometheus.GaugeValue, s.DeleteRequests, "delete",
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsOperationLatency, prometheus.GaugeValue, msToSeconds(s.ImportLatency), "import",
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsOperationRequests, prometheus.GaugeValue, s.ImportRequests, "import",
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsOperationLatency, prometheus.GaugeValue, msToSeconds(s.SearchLatency), "search",
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsOperationRequests, prometheus.GaugeValue, s.SearchRequests, "search",
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsOperationLatency, prometheus.GaugeValue, msToSeconds(s.WriteLatency), "write",
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsOperationRequests, prometheus.GaugeValue, s.WriteRequests, "write",
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsOperationRequests, prometheus.GaugeValue, s.TotalRequests, "all",
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsPendingWrite, prometheus.GaugeValue, s.PendingWrite,
	)
}

func (e *Exporter) collectUp(ch chan<- prometheus.Metric) {
	healthy, err := e.TypeSense.GetHealth()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(e.ExporterMetrics.Up, prometheus.GaugeValue, 0)
		return
	}
	if healthy {
		ch <- prometheus.MustNewConstMetric(e.ExporterMetrics.Up, prometheus.GaugeValue, 1)
		return
	}
	ch <- prometheus.MustNewConstMetric(e.ExporterMetrics.Up, prometheus.GaugeValue, 0)
}

func (em *ExporterMetrics) initializeDescriptors() {
	const namespace = "typesense"
	const subsystemSystem = "system"
	const subsystemApplication = "application"
	const subsystemApi = "api"

	em.Up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last scrape of the Typesense exporter successful.",
		nil, nil,
	)
	em.SystemCPUxActivePercentage = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemSystem, "cpu_x_active"),
		"Ratio of CPU core time spent in user mode.",
		[]string{"cpu"}, nil,
	)
	em.SystemDiskTotalBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemSystem, "disk_total_bytes"),
		"Total disk space in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.SystemDiskUsedBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemSystem, "disk_used_bytes"),
		"Used disk space in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.SystemMemoryTotalBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemSystem, "memory_total_bytes"),
		"Total memory in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.SystemMemoryUsedBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemSystem, "memory_used_bytes"),
		"Used memory in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.SystemNetworkReceived = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemSystem, "network_received_bytes"),
		"Total bytes received by the network interface.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.SystemNetworkSent = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemSystem, "network_sent_bytes"),
		"Total bytes sent by the network interface.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.TypesenseMemoryActive = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemApplication, "memory_active_bytes"),
		"Total active memory in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.TypesenseMemoryAllocated = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemApplication, "memory_allocated_bytes"),
		"Total allocated memory in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.TypesenseMemoryFragment = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemApplication, "memory_fragment_bytes"),
		"Total memory fragmentation in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.TypesenseMemoryMapped = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemApplication, "memory_mapped_bytes"),
		"Total mapped memory in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.TypesenseMemoryMetadata = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemApplication, "memory_metadata_bytes"),
		"Total metadata memory in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.TypesenseMemoryResident = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemApplication, "memory_resident_bytes"),
		"Total resident memory in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.TypesenseMemoryRetained = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemApplication, "memory_retained_bytes"),
		"Total retained memory in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.ApiStatsPendingWrite = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemApi, "pending_write_batches"),
		"Number of pending write batches.",
		nil, nil,
	)
	em.ApiStatsOperationLatency = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemApi, "operation_latency_seconds"),
		"Latency of the operation in seconds.",
		[]string{"operation"}, nil,
	)
	em.ApiStatsOperationRequests = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemApi, "operation_requests_per_second"),
		"Number of requests per second for the operation.",
		[]string{"operation"}, nil,
	)
	em.ApiStatsEndpointLatency = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemApi, "endpoint_latency_seconds"),
		"Latency of the endpoint in seconds.",
		[]string{"endpoint"}, nil,
	)
	em.ApiStatsEndpointRequests = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemApi, "endpoint_requests_per_second"),
		"Number of requests per second for the endpoint.",
		[]string{"endpoint"}, nil,
	)
}

func NewExporter(typesenseClient *Client) *Exporter {
	em := ExporterMetrics{}
	em.initializeDescriptors()
	return &Exporter{
		ExporterMetrics: em,
		TypeSense:       typesenseClient,
	}
}
