/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindElasticsearchNodesStats = "ElasticsearchNodesStats"
	ResourceElasticsearchNodesStats     = "elasticsearchnodesstats"
	ResourceElasticsearchNodesStatses   = "elasticsearchnodesstats"
)

// ElasticsearchNodesStatsSpec defines the desired state of ElasticsearchNodesStats
type ElasticsearchNodesStatsSpec struct {
	Nodes []ElasticsearchNodesStatSpec `json:"nodes"`
}

type ElasticsearchNodesStatSpec struct {
	// Time the node stats were collected for this response in Unix
	Timestamp *metav1.Time `json:"timestamp"`

	// Human-readable identifier for the node.
	Name string `json:"name"`

	// Transport address for the node
	TransportAddr string `json:"transportAddr"`

	// Network host for the node
	Host string `json:"host"`

	// IP address and port for the node
	IP string `json:"ip"`

	// Roles assigned to the node
	Roles []string `json:"roles"`

	// Indices returns index information.
	Indices *NodesStatsIndex `json:"indices"`

	// OS information, e.g. CPU and memory.
	OS *NodesStatsNodeOS `json:"os"`
}

type NodesStatsIndex struct {
	Docs         *NodesStatsDocsStats         `json:"docs"`
	Shards       *NodesStatsShardCountStats   `json:"shards_stats"`
	Store        *NodesStatsStoreStats        `json:"store"`
	Indexing     *NodesStatsIndexingStats     `json:"indexing"`
	Get          *NodesStatsGetStats          `json:"get"`
	Search       *NodesStatsSearchStats       `json:"search"`
	Merges       *NodesStatsMergeStats        `json:"merges"`
	Refresh      *NodesStatsRefreshStats      `json:"refresh"`
	Flush        *NodesStatsFlushStats        `json:"flush"`
	Warmer       *NodesStatsWarmerStats       `json:"warmer"`
	QueryCache   *NodesStatsQueryCacheStats   `json:"query_cache"`
	Fielddata    *NodesStatsFielddataStats    `json:"fielddata"`
	Completion   *NodesStatsCompletionStats   `json:"completion"`
	Segments     *NodesStatsSegmentsStats     `json:"segments"`
	Translog     *NodesStatsTranslogStats     `json:"translog"`
	RequestCache *NodesStatsRequestCacheStats `json:"request_cache"`
	Recovery     NodesStatsRecoveryStats      `json:"recovery"`

	IndicesLevel map[string]NodesStatsIndex `json:"indices"` // for level=indices
	ShardsLevel  map[string]NodesStatsIndex `json:"shards"`  // for level=shards
}

type NodesStatsDocsStats struct {
	Count   int64 `json:"count"`
	Deleted int64 `json:"deleted"`
}

type NodesStatsShardCountStats struct {
	TotalCount int64 `json:"total_count"`
}

type NodesStatsStoreStats struct {
	TotalSize string `json:"size"`
	SizeBytes int64  `json:"size_in_bytes"`
}

type NodesStatsIndexingStats struct {
	IndexTotal            int64  `json:"index_total"`
	IndexTime             string `json:"index_time"`
	IndexTimeInMillis     int64  `json:"index_time_in_millis"`
	IndexCurrent          int64  `json:"index_current"`
	IndexFailed           int64  `json:"index_failed"`
	DeleteTotal           int64  `json:"delete_total"`
	DeleteTime            string `json:"delete_time"`
	DeleteTimeInMillis    int64  `json:"delete_time_in_millis"`
	DeleteCurrent         int64  `json:"delete_current"`
	NoopUpdateTotal       int64  `json:"noop_update_total"`
	IsThrottled           bool   `json:"is_throttled"`
	ThrottledTime         string `json:"throttle_time"` // no typo, see https://github.com/elastic/elasticsearch/blob/ff99bc1d3f8a7ea72718872d214ec2097dfca276/server/src/main/java/org/elasticsearch/index/shard/IndexingStats.java#L244
	ThrottledTimeInMillis int64  `json:"throttle_time_in_millis"`

	Types map[string]NodesStatsIndexingStats `json:"types"` // stats for individual types
}

type NodesStatsGetStats struct {
	Total               int64  `json:"total"`
	Time                string `json:"get_time"`
	TimeInMillis        int64  `json:"time_in_millis"`
	Exists              int64  `json:"exists"`
	ExistsTime          string `json:"exists_time"`
	ExistsTimeInMillis  int64  `json:"exists_in_millis"`
	Missing             int64  `json:"missing"`
	MissingTime         string `json:"missing_time"`
	MissingTimeInMillis int64  `json:"missing_in_millis"`
	Current             int64  `json:"current"`
}

type NodesStatsSearchStats struct {
	OpenContexts       int64  `json:"open_contexts"`
	QueryTotal         int64  `json:"query_total"`
	QueryTime          string `json:"query_time"`
	QueryTimeInMillis  int64  `json:"query_time_in_millis"`
	QueryCurrent       int64  `json:"query_current"`
	FetchTotal         int64  `json:"fetch_total"`
	FetchTime          string `json:"fetch_time"`
	FetchTimeInMillis  int64  `json:"fetch_time_in_millis"`
	FetchCurrent       int64  `json:"fetch_current"`
	ScrollTotal        int64  `json:"scroll_total"`
	ScrollTime         string `json:"scroll_time"`
	ScrollTimeInMillis int64  `json:"scroll_time_in_millis"`
	ScrollCurrent      int64  `json:"scroll_current"`

	Groups map[string]NodesStatsSearchStats `json:"groups"` // stats for individual groups
}

type NodesStatsMergeStats struct {
	Current                    int64  `json:"current"`
	CurrentDocs                int64  `json:"current_docs"`
	CurrentSize                string `json:"current_size"`
	CurrentSizeBytes           int64  `json:"current_size_in_bytes"`
	Total                      int64  `json:"total"`
	TotalTime                  string `json:"total_time"`
	TotalTimeInMillis          int64  `json:"total_time_in_millis"`
	TotalDocs                  int64  `json:"total_docs"`
	TotalSize                  string `json:"total_size"`
	TotalSizeBytes             int64  `json:"total_size_in_bytes"`
	TotalStoppedTime           string `json:"total_stopped_time"`
	TotalStoppedTimeInMillis   int64  `json:"total_stopped_time_in_millis"`
	TotalThrottledTime         string `json:"total_throttled_time"`
	TotalThrottledTimeInMillis int64  `json:"total_throttled_time_in_millis"`
	TotalThrottleBytes         string `json:"total_auto_throttle"`
	TotalThrottleBytesBytes    int64  `json:"total_auto_throttle_in_bytes"`
}

type NodesStatsRefreshStats struct {
	Total             int64  `json:"total"`
	TotalTime         string `json:"total_time"`
	TotalTimeInMillis int64  `json:"total_time_in_millis"`
}

type NodesStatsFlushStats struct {
	Total             int64  `json:"total"`
	TotalTime         string `json:"total_time"`
	TotalTimeInMillis int64  `json:"total_time_in_millis"`
}

type NodesStatsWarmerStats struct {
	Current           int64  `json:"current"`
	Total             int64  `json:"total"`
	TotalTime         string `json:"total_time"`
	TotalTimeInMillis int64  `json:"total_time_in_millis"`
}

type NodesStatsQueryCacheStats struct {
	MemorySize      string `json:"memory_size"`
	MemorySizeBytes int64  `json:"memory_size_in_bytes"`
	TotalCount      int64  `json:"total_count"`
	HitCount        int64  `json:"hit_count"`
	MissCount       int64  `json:"miss_count"`
	CacheSize       int64  `json:"cache_size"`
	CacheCount      int64  `json:"cache_count"`
	Evictions       int64  `json:"evictions"`
}

type NodesStatsFielddataStats struct {
	MemorySize      string                     `json:"memory_size"`
	MemorySizeBytes int64                      `json:"memory_size_in_bytes"`
	Evictions       int64                      `json:"evictions"`
	Fields          *NodesStatsFieldDataFields `json:"fields"`
}

type NodesStatsFieldDataFields struct {
	MemorySize      string `json:"memory_size"`
	MemorySizeBytes int64  `json:"memory_size_in_bytes"`
}

type NodesStatsCompletionStats struct {
	TotalSize string                      `json:"size"`
	SizeBytes int64                       `json:"size_in_bytes"`
	Fields    *NodesStatsCompletionFields `json:"fields"`
}

type NodesStatsCompletionFields struct {
	TotalSize string `json:"size"`
	SizeBytes int64  `json:"size_in_bytes"`
}

type NodesStatsSegmentsStats struct {
	Count                     int64  `json:"count"`
	Memory                    string `json:"memory"`
	MemoryBytes               int64  `json:"memory_in_bytes"`
	TermsMemory               string `json:"terms_memory"`
	TermsMemoryBytes          int64  `json:"terms_memory_in_bytes"`
	StoredFieldsMemory        string `json:"stored_fields_memory"`
	StoredFieldsMemoryBytes   int64  `json:"stored_fields_memory_in_bytes"`
	TermVectorsMemory         string `json:"term_vectors_memory"`
	TermVectorsMemoryBytes    int64  `json:"term_vectors_memory_in_bytes"`
	NormsMemory               string `json:"norms_memory"`
	NormsMemoryBytes          int64  `json:"norms_memory_in_bytes"`
	DocValuesMemory           string `json:"doc_values_memory"`
	DocValuesMemoryBytes      int64  `json:"doc_values_memory_in_bytes"`
	IndexWriterMemory         string `json:"index_writer_memory"`
	IndexWriterMemoryBytes    int64  `json:"index_writer_memory_in_bytes"`
	IndexWriterMaxMemory      string `json:"index_writer_max_memory"`
	IndexWriterMaxMemoryBytes int64  `json:"index_writer_max_memory_in_bytes"`
	VersionMapMemory          string `json:"version_map_memory"`
	VersionMapMemoryBytes     int64  `json:"version_map_memory_in_bytes"`
	FixedBitSetMemory         string `json:"fixed_bit_set"` // not a typo
	FixedBitSetMemoryBytes    int64  `json:"fixed_bit_set_memory_in_bytes"`
}

type NodesStatsTranslogStats struct {
	Operations int64  `json:"operations"`
	TotalSize  string `json:"size"`
	SizeBytes  int64  `json:"size_in_bytes"`
}

type NodesStatsRequestCacheStats struct {
	MemorySize      string `json:"memory_size"`
	MemorySizeBytes int64  `json:"memory_size_in_bytes"`
	Evictions       int64  `json:"evictions"`
	HitCount        int64  `json:"hit_count"`
	MissCount       int64  `json:"miss_count"`
}

type NodesStatsRecoveryStats struct {
	CurrentAsSource int64 `json:"current_as_source"`
	CurrentAsTarget int64 `json:"current_as_target"`
}

type NodesStatsNodeOS struct {
	Timestamp int64                 `json:"timestamp"`
	CPU       *NodesStatsNodeOSCPU  `json:"cpu"`
	Mem       *NodesStatsNodeOSMem  `json:"mem"`
	Swap      *NodesStatsNodeOSSwap `json:"swap"`
}

type NodesStatsNodeOSCPU struct {
	Percent     int64              `json:"percent"`
	LoadAverage map[string]float64 `json:"load_average"` // keys are: 1m, 5m, and 15m
}

type NodesStatsNodeOSMem struct {
	Total       string `json:"total"`
	TotalBytes  int64  `json:"total_in_bytes"`
	Free        string `json:"free"`
	FreeBytes   int64  `json:"free_in_bytes"`
	Used        string `json:"used"`
	UsedBytes   int64  `json:"used_in_bytes"`
	FreePercent int64  `json:"free_percent"`
	UsedPercent int64  `json:"used_percent"`
}

type NodesStatsNodeOSSwap struct {
	Total      string `json:"total"`
	TotalBytes int64  `json:"total_in_bytes"`
	Free       string `json:"free"`
	FreeBytes  int64  `json:"free_in_bytes"`
	Used       string `json:"used"`
	UsedBytes  int64  `json:"used_in_bytes"`
}

// ElasticsearchNodesStatsStatus defines the observed state of ElasticsearchNodesStats
type ElasticsearchNodesStatsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// ElasticsearchNodesStats is the Schema for the ElasticsearchNodesStats API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ElasticsearchNodesStats struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ElasticsearchNodesStatsSpec   `json:"spec,omitempty"`
	Status ElasticsearchNodesStatsStatus `json:"status,omitempty"`
}

// ElasticsearchNodesStatsList contains a list of ElasticsearchNodesStats

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ElasticsearchNodesStatsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ElasticsearchNodesStats `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ElasticsearchNodesStats{}, &ElasticsearchNodesStatsList{})
}
