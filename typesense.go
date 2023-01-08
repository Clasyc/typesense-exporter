package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const TypesenseHeaderApiKey = "X-TYPESENSE-API-KEY"

type Client struct {
	http   *http.Client
	apiKey string
	host   string
}

type ResponseMetrics struct {
	SystemCPUxActivePercentage map[string]json.Number `json:"-"`
	SystemCPUActivePercentage  float64                `json:"system_cpu_active_percentage,string"`
	SystemDiskTotalBytes       float64                `json:"system_disk_total_bytes,string"`
	SystemDiskUsedBytes        float64                `json:"system_disk_used_bytes,string"`
	SystemMemoryTotalBytes     float64                `json:"system_memory_total_bytes,string"`
	SystemMemoryUsedBytes      float64                `json:"system_memory_used_bytes,string"`
	SystemNetworkReceivedBytes float64                `json:"system_network_received_bytes,string"`
	SystemNetworkSentBytes     float64                `json:"system_network_sent_bytes,string"`
	TypesenseMemoryActiveBytes float64                `json:"typesense_memory_active_bytes,string"`
	TypesenseMemoryAllocated   float64                `json:"typesense_memory_allocated_bytes,string"`
	TypesenseMemoryFragment    float64                `json:"typesense_memory_fragmentation_ratio,string"`
	TypesenseMemoryMapped      float64                `json:"typesense_memory_mapped_bytes,string"`
	TypesenseMemoryMetadata    float64                `json:"typesense_memory_metadata_bytes,string"`
	TypesenseMemoryResident    float64                `json:"typesense_memory_resident_bytes,string"`
	TypesenseMemoryRetained    float64                `json:"typesense_memory_retained_bytes,string"`
}

type ResponseApiStats struct {
	DeleteLatency  float64                `json:"delete_latency_ms"`
	DeleteRequests float64                `json:"delete_requests_per_second"`
	ImportLatency  float64                `json:"import_latency_ms"`
	ImportRequests float64                `json:"import_requests_per_second"`
	Latency        map[string]json.Number `json:"latency_ms"`
	PendingWrite   float64                `json:"pending_write_batches"`
	Requests       map[string]json.Number `json:"requests_per_second"`
	SearchLatency  float64                `json:"search_latency_ms"`
	SearchRequests float64                `json:"search_requests_per_second"`
	TotalRequests  float64                `json:"total_requests_per_second"`
	WriteLatency   float64                `json:"write_latency_ms"`
	WriteRequests  float64                `json:"write_requests_per_second"`
}

type ResponseHealth struct {
	Ok bool `json:"ok"`
}

func NewClient(apiKey string, host string, timeout time.Duration) *Client {
	return &Client{
		http: &http.Client{
			Timeout: timeout,
		},
		apiKey: apiKey,
		host:   host,
	}
}

func (c *Client) get(path string) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/%s", c.host, path)
	client := &http.Client{}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(TypesenseHeaderApiKey, c.apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	return body, err
}

func (c *Client) GetMetrics() (*ResponseMetrics, error) {
	body, err := c.get("metrics.json")
	response := &ResponseMetrics{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &response.SystemCPUxActivePercentage)
	if err != nil {
		return nil, err
	}
	for k := range response.SystemCPUxActivePercentage {
		if !strings.HasPrefix(k, "system_cpu") {
			delete(response.SystemCPUxActivePercentage, k)
		}
	}

	return response, nil
}

func (c *Client) GetStats() (*ResponseApiStats, error) {
	body, err := c.get("stats.json")
	response := &ResponseApiStats{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) GetHealth() (bool, error) {
	body, err := c.get("health")
	if err != nil {
		return false, err
	}

	response := &ResponseHealth{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return false, err
	}

	return response.Ok, nil
}
