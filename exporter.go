package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
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
	ApiStatsDeleteLatency      *prometheus.Desc
	ApiStatsDeleteRequests     *prometheus.Desc
	ApiStatsImportLatency      *prometheus.Desc
	ApiStatsImportRequests     *prometheus.Desc
	ApiStatsLatency            *prometheus.Desc
	ApiStatsPendingWrite       *prometheus.Desc
	ApiStatsRequests           *prometheus.Desc
	ApiStatsSearchLatency      *prometheus.Desc
	ApiStatsSearchRequests     *prometheus.Desc
	ApiStatsTotalRequests      *prometheus.Desc
	ApiStatsWriteLatency       *prometheus.Desc
	ApiStatsWriteRequests      *prometheus.Desc
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
	ch <- em.ApiStatsDeleteLatency
	ch <- em.ApiStatsDeleteRequests
	ch <- em.ApiStatsImportLatency
	ch <- em.ApiStatsImportRequests
	ch <- em.ApiStatsLatency
	ch <- em.ApiStatsPendingWrite
	ch <- em.ApiStatsRequests
	ch <- em.ApiStatsSearchLatency
	ch <- em.ApiStatsSearchRequests
	ch <- em.ApiStatsTotalRequests
	ch <- em.ApiStatsWriteLatency
	ch <- em.ApiStatsWriteRequests
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

	for name, value := range v.SystemCPUxActivePercentage {
		float, err := value.Float64()
		if err != nil {
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			e.ExporterMetrics.SystemCPUxActivePercentage, prometheus.GaugeValue, percentageToRatio(float), name,
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
			e.ExporterMetrics.ApiStatsRequests, prometheus.GaugeValue, float, name,
		)
	}
	for name, value := range s.Requests {
		float, err := value.Float64()
		if err != nil {
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			e.ExporterMetrics.ApiStatsLatency, prometheus.GaugeValue, msToSeconds(float), name,
		)
	}

	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsDeleteLatency, prometheus.GaugeValue, msToSeconds(s.DeleteLatency),
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsDeleteRequests, prometheus.GaugeValue, s.DeleteRequests,
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsImportLatency, prometheus.GaugeValue, msToSeconds(s.ImportLatency),
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsImportRequests, prometheus.GaugeValue, s.ImportRequests,
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsPendingWrite, prometheus.GaugeValue, s.PendingWrite,
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsSearchLatency, prometheus.GaugeValue, msToSeconds(s.SearchLatency),
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsSearchRequests, prometheus.GaugeValue, s.SearchRequests,
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsTotalRequests, prometheus.GaugeValue, s.TotalRequests,
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsWriteLatency, prometheus.GaugeValue, msToSeconds(s.WriteLatency),
	)
	ch <- prometheus.MustNewConstMetric(
		e.ExporterMetrics.ApiStatsWriteRequests, prometheus.GaugeValue, s.WriteRequests,
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
	em.Up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last scrape of the Typesense exporter successful.",
		nil, nil,
	)
	em.SystemCPUxActivePercentage = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "system", "system_cpu_x_active_percentage"),
		"Percentage of CPU core time spent in user mode.",
		[]string{"cpu"}, nil,
	)
	em.SystemDiskTotalBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "system", "system_disk_total_bytes"),
		"Total disk space in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.SystemDiskUsedBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "system", "system_disk_used_bytes"),
		"Used disk space in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.SystemMemoryTotalBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "system", "system_memory_total_bytes"),
		"Total memory in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.SystemMemoryUsedBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "system", "system_memory_used_bytes"),
		"Used memory in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.SystemNetworkReceived = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "system", "system_network_received_bytes"),
		"Total bytes received by the network interface.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.SystemNetworkSent = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "system", "system_network_sent_bytes"),
		"Total bytes sent by the network interface.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.TypesenseMemoryActive = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "typesense", "typesense_memory_active_bytes"),
		"Total active memory in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.TypesenseMemoryAllocated = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "typesense", "typesense_memory_allocated_bytes"),
		"Total allocated memory in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.TypesenseMemoryFragment = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "typesense", "typesense_memory_fragment_bytes"),
		"Total memory fragmentation in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.TypesenseMemoryMapped = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "typesense", "typesense_memory_mapped_bytes"),
		"Total mapped memory in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.TypesenseMemoryMetadata = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "typesense", "typesense_memory_metadata_bytes"),
		"Total metadata memory in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.TypesenseMemoryResident = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "typesense", "typesense_memory_resident_bytes"),
		"Total resident memory in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.TypesenseMemoryRetained = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "typesense", "typesense_memory_retained_bytes"),
		"Total retained memory in bytes.",
		nil, prometheus.Labels{"unit": "bytes"},
	)
	em.ApiStatsRequests = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "api", "requests_per_second"),
		"Number of requests per second.",
		[]string{"endpoint"}, prometheus.Labels{"unit": "seconds"},
	)
	em.ApiStatsLatency = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "api", "endpoint_requests_latency_seconds"),
		"Endpoint request latency in seconds.",
		[]string{"endpoint"}, prometheus.Labels{"unit": "seconds"},
	)
	em.ApiStatsDeleteLatency = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "api", "delete_requests_latency_seconds"),
		"Delete request latency in seconds.",
		nil, prometheus.Labels{"unit": "seconds"},
	)
	em.ApiStatsDeleteRequests = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "api", "delete_requests_per_second"),
		"Delete requests per second.",
		nil, nil,
	)
	em.ApiStatsImportLatency = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "api", "import_requests_latency_seconds"),
		"Import request latency in seconds.",
		nil, prometheus.Labels{"unit": "seconds"},
	)
	em.ApiStatsImportRequests = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "api", "import_requests_per_second"),
		"Import requests per second.",
		nil, nil,
	)
	em.ApiStatsPendingWrite = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "api", "pending_write_batches"),
		"Number of pending write batches.",
		nil, nil,
	)
	em.ApiStatsSearchLatency = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "api", "search_requests_latency_seconds"),
		"Search request latency in seconds.",
		nil, prometheus.Labels{"unit": "seconds"},
	)
	em.ApiStatsSearchRequests = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "api", "search_requests_per_second"),
		"Search requests per second.",
		nil, nil,
	)
	em.ApiStatsTotalRequests = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "api", "total_requests_per_second"),
		"Total requests per second.",
		nil, nil,
	)
	em.ApiStatsWriteLatency = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "api", "write_requests_latency_seconds"),
		"Write request latency in seconds.",
		nil, prometheus.Labels{"unit": "seconds"},
	)
	em.ApiStatsWriteRequests = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "api", "write_requests_per_second"),
		"Write requests per second.",
		nil, nil,
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
