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
	Nodes []ElasticsearchNodesStatSpec `json:"nodes" protobuf:"bytes,1,rep,name=nodes"`
}

type ElasticsearchNodesStatSpec struct {
	// Time the node stats were collected for this response in Unix
	Timestamp int64 `json:"timestamp" protobuf:"varint,1,opt,name=timestamp"`

	// Human-readable identifier for the node.
	Name string `json:"name" protobuf:"bytes,2,opt,name=name"`

	// Transport address for the node
	TransportAddr string `json:"transport_addr" protobuf:"bytes,3,opt,name=transport_addr,json=transportAddr"`

	// Network host for the node
	Host string `json:"host" protobuf:"bytes,4,opt,name=host"`

	// IP address and port for the node
	IP string `json:"ip" protobuf:"bytes,5,opt,name=ip"`

	// Roles assigned to the node
	Roles []string `json:"roles" protobuf:"bytes,6,rep,name=roles"`

	// Indices returns index information.
	Indices *NodesStatsIndex `json:"indices" protobuf:"bytes,7,opt,name=indices"`

	// OS information, e.g. CPU and memory.
	OS *NodesStatsNodeOS `json:"os" protobuf:"bytes,8,opt,name=os"`
}

type NodesStatsIndex struct {
	Docs         *NodesStatsDocsStats         `json:"docs" protobuf:"bytes,1,opt,name=docs"`
	Shards       *NodesStatsShardCountStats   `json:"shards_stats" protobuf:"bytes,2,opt,name=shards_stats,json=shardsStats"`
	Store        *NodesStatsStoreStats        `json:"store" protobuf:"bytes,3,opt,name=store"`
	Indexing     *NodesStatsIndexingStats     `json:"indexing" protobuf:"bytes,4,opt,name=indexing"`
	Get          *NodesStatsGetStats          `json:"get" protobuf:"bytes,5,opt,name=get"`
	Search       *NodesStatsSearchStats       `json:"search" protobuf:"bytes,6,opt,name=search"`
	Merges       *NodesStatsMergeStats        `json:"merges" protobuf:"bytes,7,opt,name=merges"`
	Refresh      *NodesStatsRefreshStats      `json:"refresh" protobuf:"bytes,8,opt,name=refresh"`
	Flush        *NodesStatsFlushStats        `json:"flush" protobuf:"bytes,9,opt,name=flush"`
	Warmer       *NodesStatsWarmerStats       `json:"warmer" protobuf:"bytes,10,opt,name=warmer"`
	QueryCache   *NodesStatsQueryCacheStats   `json:"query_cache" protobuf:"bytes,11,opt,name=query_cache,json=queryCache"`
	Fielddata    *NodesStatsFielddataStats    `json:"fielddata" protobuf:"bytes,12,opt,name=fielddata"`
	Completion   *NodesStatsCompletionStats   `json:"completion" protobuf:"bytes,13,opt,name=completion"`
	Segments     *NodesStatsSegmentsStats     `json:"segments" protobuf:"bytes,14,opt,name=segments"`
	Translog     *NodesStatsTranslogStats     `json:"translog" protobuf:"bytes,15,opt,name=translog"`
	RequestCache *NodesStatsRequestCacheStats `json:"request_cache" protobuf:"bytes,16,opt,name=request_cache,json=requestCache"`
	Recovery     NodesStatsRecoveryStats      `json:"recovery" protobuf:"bytes,17,opt,name=recovery"`

	IndicesLevel map[string]NodesStatsIndex `json:"indices" protobuf:"bytes,18,rep,name=indices"` // for level=indices
	ShardsLevel  map[string]NodesStatsIndex `json:"shards" protobuf:"bytes,19,rep,name=shards"`   // for level=shards
}

type NodesStatsDocsStats struct {
	Count   int64 `json:"count" protobuf:"varint,1,opt,name=count"`
	Deleted int64 `json:"deleted" protobuf:"varint,2,opt,name=deleted"`
}

type NodesStatsShardCountStats struct {
	TotalCount int64 `json:"total_count" protobuf:"varint,1,opt,name=total_count,json=totalCount"`
}

type NodesStatsStoreStats struct {
	TotalSize   string `json:"size" protobuf:"bytes,3,opt,name=size"`
	SizeInBytes int64  `json:"size_in_bytes" protobuf:"varint,2,opt,name=size_in_bytes,json=sizeInBytes"`
}

type NodesStatsIndexingStats struct {
	IndexTotal            int64  `json:"index_total" protobuf:"varint,1,opt,name=index_total,json=indexTotal"`
	IndexTime             string `json:"index_time" protobuf:"bytes,2,opt,name=index_time,json=indexTime"`
	IndexTimeInMillis     int64  `json:"index_time_in_millis" protobuf:"varint,3,opt,name=index_time_in_millis,json=indexTimeInMillis"`
	IndexCurrent          int64  `json:"index_current" protobuf:"varint,4,opt,name=index_current,json=indexCurrent"`
	IndexFailed           int64  `json:"index_failed" protobuf:"varint,5,opt,name=index_failed,json=indexFailed"`
	DeleteTotal           int64  `json:"delete_total" protobuf:"varint,6,opt,name=delete_total,json=deleteTotal"`
	DeleteTime            string `json:"delete_time" protobuf:"bytes,7,opt,name=delete_time,json=deleteTime"`
	DeleteTimeInMillis    int64  `json:"delete_time_in_millis" protobuf:"varint,8,opt,name=delete_time_in_millis,json=deleteTimeInMillis"`
	DeleteCurrent         int64  `json:"delete_current" protobuf:"varint,9,opt,name=delete_current,json=deleteCurrent"`
	NoopUpdateTotal       int64  `json:"noop_update_total" protobuf:"varint,10,opt,name=noop_update_total,json=noopUpdateTotal"`
	IsThrottled           bool   `json:"is_throttled" protobuf:"varint,11,opt,name=is_throttled,json=isThrottled"`
	ThrottledTime         string `json:"throttle_time" protobuf:"bytes,12,opt,name=throttle_time,json=throttleTime"` // no typo, see https://github.com/elastic/elasticsearch/blob/ff99bc1d3f8a7ea72718872d214ec2097dfca276/server/src/main/java/org/elasticsearch/index/shard/IndexingStats.java#L244
	ThrottledTimeInMillis int64  `json:"throttle_time_in_millis" protobuf:"varint,13,opt,name=throttle_time_in_millis,json=throttleTimeInMillis"`

	Types map[string]NodesStatsIndexingStats `json:"types" protobuf:"bytes,14,rep,name=types"` // stats for individual types
}

type NodesStatsGetStats struct {
	Total               int64  `json:"total" protobuf:"varint,1,opt,name=total"`
	Time                string `json:"get_time" protobuf:"bytes,2,opt,name=get_time,json=getTime"`
	TimeInMillis        int64  `json:"time_in_millis" protobuf:"varint,3,opt,name=time_in_millis,json=timeInMillis"`
	Exists              int64  `json:"exists" protobuf:"varint,4,opt,name=exists"`
	ExistsTime          string `json:"exists_time" protobuf:"bytes,5,opt,name=exists_time,json=existsTime"`
	ExistsTimeInMillis  int64  `json:"exists_in_millis" protobuf:"varint,6,opt,name=exists_in_millis,json=existsInMillis"`
	Missing             int64  `json:"missing" protobuf:"varint,7,opt,name=missing"`
	MissingTime         string `json:"missing_time" protobuf:"bytes,8,opt,name=missing_time,json=missingTime"`
	MissingTimeInMillis int64  `json:"missing_in_millis" protobuf:"varint,9,opt,name=missing_in_millis,json=missingInMillis"`
	Current             int64  `json:"current" protobuf:"varint,10,opt,name=current"`
}

type NodesStatsSearchStats struct {
	OpenContexts       int64  `json:"open_contexts" protobuf:"varint,1,opt,name=open_contexts,json=openContexts"`
	QueryTotal         int64  `json:"query_total" protobuf:"varint,2,opt,name=query_total,json=queryTotal"`
	QueryTime          string `json:"query_time" protobuf:"bytes,3,opt,name=query_time,json=queryTime"`
	QueryTimeInMillis  int64  `json:"query_time_in_millis" protobuf:"varint,4,opt,name=query_time_in_millis,json=queryTimeInMillis"`
	QueryCurrent       int64  `json:"query_current" protobuf:"varint,5,opt,name=query_current,json=queryCurrent"`
	FetchTotal         int64  `json:"fetch_total" protobuf:"varint,6,opt,name=fetch_total,json=fetchTotal"`
	FetchTime          string `json:"fetch_time" protobuf:"bytes,7,opt,name=fetch_time,json=fetchTime"`
	FetchTimeInMillis  int64  `json:"fetch_time_in_millis" protobuf:"varint,8,opt,name=fetch_time_in_millis,json=fetchTimeInMillis"`
	FetchCurrent       int64  `json:"fetch_current" protobuf:"varint,9,opt,name=fetch_current,json=fetchCurrent"`
	ScrollTotal        int64  `json:"scroll_total" protobuf:"varint,10,opt,name=scroll_total,json=scrollTotal"`
	ScrollTime         string `json:"scroll_time" protobuf:"bytes,11,opt,name=scroll_time,json=scrollTime"`
	ScrollTimeInMillis int64  `json:"scroll_time_in_millis" protobuf:"varint,12,opt,name=scroll_time_in_millis,json=scrollTimeInMillis"`
	ScrollCurrent      int64  `json:"scroll_current" protobuf:"varint,13,opt,name=scroll_current,json=scrollCurrent"`

	Groups map[string]NodesStatsSearchStats `json:"groups" protobuf:"bytes,14,rep,name=groups"` // stats for individual groups
}

type NodesStatsMergeStats struct {
	Current                    int64  `json:"current" protobuf:"varint,1,opt,name=current"`
	CurrentDocs                int64  `json:"current_docs" protobuf:"varint,2,opt,name=current_docs,json=currentDocs"`
	CurrentSize                string `json:"current_size" protobuf:"bytes,3,opt,name=current_size,json=currentSize"`
	CurrentSizeInBytes         int64  `json:"current_size_in_bytes" protobuf:"varint,4,opt,name=current_size_in_bytes,json=currentSizeInBytes"`
	Total                      int64  `json:"total" protobuf:"varint,5,opt,name=total"`
	TotalTime                  string `json:"total_time" protobuf:"bytes,6,opt,name=total_time,json=totalTime"`
	TotalTimeInMillis          int64  `json:"total_time_in_millis" protobuf:"varint,7,opt,name=total_time_in_millis,json=totalTimeInMillis"`
	TotalDocs                  int64  `json:"total_docs" protobuf:"varint,8,opt,name=total_docs,json=totalDocs"`
	TotalSize                  string `json:"total_size" protobuf:"bytes,9,opt,name=total_size,json=totalSize"`
	TotalSizeInBytes           int64  `json:"total_size_in_bytes" protobuf:"varint,10,opt,name=total_size_in_bytes,json=totalSizeInBytes"`
	TotalStoppedTime           string `json:"total_stopped_time" protobuf:"bytes,11,opt,name=total_stopped_time,json=totalStoppedTime"`
	TotalStoppedTimeInMillis   int64  `json:"total_stopped_time_in_millis" protobuf:"varint,12,opt,name=total_stopped_time_in_millis,json=totalStoppedTimeInMillis"`
	TotalThrottledTime         string `json:"total_throttled_time" protobuf:"bytes,13,opt,name=total_throttled_time,json=totalThrottledTime"`
	TotalThrottledTimeInMillis int64  `json:"total_throttled_time_in_millis" protobuf:"varint,14,opt,name=total_throttled_time_in_millis,json=totalThrottledTimeInMillis"`
	TotalThrottleBytes         string `json:"total_auto_throttle" protobuf:"bytes,15,opt,name=total_auto_throttle,json=totalAutoThrottle"`
	TotalThrottleBytesInBytes  int64  `json:"total_auto_throttle_in_bytes" protobuf:"varint,16,opt,name=total_auto_throttle_in_bytes,json=totalAutoThrottleInBytes"`
}

type NodesStatsRefreshStats struct {
	Total             int64  `json:"total" protobuf:"varint,1,opt,name=total"`
	TotalTime         string `json:"total_time" protobuf:"bytes,2,opt,name=total_time,json=totalTime"`
	TotalTimeInMillis int64  `json:"total_time_in_millis" protobuf:"varint,3,opt,name=total_time_in_millis,json=totalTimeInMillis"`
}

type NodesStatsFlushStats struct {
	Total             int64  `json:"total" protobuf:"varint,1,opt,name=total"`
	TotalTime         string `json:"total_time" protobuf:"bytes,2,opt,name=total_time,json=totalTime"`
	TotalTimeInMillis int64  `json:"total_time_in_millis" protobuf:"varint,3,opt,name=total_time_in_millis,json=totalTimeInMillis"`
}

type NodesStatsWarmerStats struct {
	Current           int64  `json:"current" protobuf:"varint,1,opt,name=current"`
	Total             int64  `json:"total" protobuf:"varint,2,opt,name=total"`
	TotalTime         string `json:"total_time" protobuf:"bytes,3,opt,name=total_time,json=totalTime"`
	TotalTimeInMillis int64  `json:"total_time_in_millis" protobuf:"varint,4,opt,name=total_time_in_millis,json=totalTimeInMillis"`
}

type NodesStatsQueryCacheStats struct {
	MemorySize        string `json:"memory_size" protobuf:"bytes,1,opt,name=memory_size,json=memorySize"`
	MemorySizeInBytes int64  `json:"memory_size_in_bytes" protobuf:"varint,2,opt,name=memory_size_in_bytes,json=memorySizeInBytes"`
	TotalCount        int64  `json:"total_count" protobuf:"varint,3,opt,name=total_count,json=totalCount"`
	HitCount          int64  `json:"hit_count" protobuf:"varint,4,opt,name=hit_count,json=hitCount"`
	MissCount         int64  `json:"miss_count" protobuf:"varint,5,opt,name=miss_count,json=missCount"`
	CacheSize         int64  `json:"cache_size" protobuf:"varint,6,opt,name=cache_size,json=cacheSize"`
	CacheCount        int64  `json:"cache_count" protobuf:"varint,7,opt,name=cache_count,json=cacheCount"`
	Evictions         int64  `json:"evictions" protobuf:"varint,8,opt,name=evictions"`
}

type NodesStatsFielddataStats struct {
	MemorySize        string                     `json:"memory_size" protobuf:"bytes,1,opt,name=memory_size,json=memorySize"`
	MemorySizeInBytes int64                      `json:"memory_size_in_bytes" protobuf:"varint,2,opt,name=memory_size_in_bytes,json=memorySizeInBytes"`
	Evictions         int64                      `json:"evictions" protobuf:"varint,3,opt,name=evictions"`
	Fields            *NodesStatsFieldDataFields `json:"fields" protobuf:"bytes,4,opt,name=fields"`
}

type NodesStatsFieldDataFields struct {
	MemorySize        string `json:"memory_size" protobuf:"bytes,1,opt,name=memory_size,json=memorySize"`
	MemorySizeInBytes int64  `json:"memory_size_in_bytes" protobuf:"varint,2,opt,name=memory_size_in_bytes,json=memorySizeInBytes"`
}

type NodesStatsCompletionStats struct {
	TotalSize   string                      `json:"size" protobuf:"bytes,4,opt,name=size"`
	SizeInBytes int64                       `json:"size_in_bytes" protobuf:"varint,2,opt,name=size_in_bytes,json=sizeInBytes"`
	Fields      *NodesStatsCompletionFields `json:"fields" protobuf:"bytes,3,opt,name=fields"`
}

type NodesStatsCompletionFields struct {
	TotalSize   string `json:"size" protobuf:"bytes,3,opt,name=size"`
	SizeInBytes int64  `json:"size_in_bytes" protobuf:"varint,2,opt,name=size_in_bytes,json=sizeInBytes"`
}

type NodesStatsSegmentsStats struct {
	Count                       int64  `json:"count" protobuf:"varint,1,opt,name=count"`
	Memory                      string `json:"memory" protobuf:"bytes,2,opt,name=memory"`
	MemoryInBytes               int64  `json:"memory_in_bytes" protobuf:"varint,3,opt,name=memory_in_bytes,json=memoryInBytes"`
	TermsMemory                 string `json:"terms_memory" protobuf:"bytes,4,opt,name=terms_memory,json=termsMemory"`
	TermsMemoryInBytes          int64  `json:"terms_memory_in_bytes" protobuf:"varint,5,opt,name=terms_memory_in_bytes,json=termsMemoryInBytes"`
	StoredFieldsMemory          string `json:"stored_fields_memory" protobuf:"bytes,6,opt,name=stored_fields_memory,json=storedFieldsMemory"`
	StoredFieldsMemoryInBytes   int64  `json:"stored_fields_memory_in_bytes" protobuf:"varint,7,opt,name=stored_fields_memory_in_bytes,json=storedFieldsMemoryInBytes"`
	TermVectorsMemory           string `json:"term_vectors_memory" protobuf:"bytes,8,opt,name=term_vectors_memory,json=termVectorsMemory"`
	TermVectorsMemoryInBytes    int64  `json:"term_vectors_memory_in_bytes" protobuf:"varint,9,opt,name=term_vectors_memory_in_bytes,json=termVectorsMemoryInBytes"`
	NormsMemory                 string `json:"norms_memory" protobuf:"bytes,10,opt,name=norms_memory,json=normsMemory"`
	NormsMemoryInBytes          int64  `json:"norms_memory_in_bytes" protobuf:"varint,11,opt,name=norms_memory_in_bytes,json=normsMemoryInBytes"`
	DocValuesMemory             string `json:"doc_values_memory" protobuf:"bytes,12,opt,name=doc_values_memory,json=docValuesMemory"`
	DocValuesMemoryInBytes      int64  `json:"doc_values_memory_in_bytes" protobuf:"varint,13,opt,name=doc_values_memory_in_bytes,json=docValuesMemoryInBytes"`
	IndexWriterMemory           string `json:"index_writer_memory" protobuf:"bytes,14,opt,name=index_writer_memory,json=indexWriterMemory"`
	IndexWriterMemoryInBytes    int64  `json:"index_writer_memory_in_bytes" protobuf:"varint,15,opt,name=index_writer_memory_in_bytes,json=indexWriterMemoryInBytes"`
	IndexWriterMaxMemory        string `json:"index_writer_max_memory" protobuf:"bytes,16,opt,name=index_writer_max_memory,json=indexWriterMaxMemory"`
	IndexWriterMaxMemoryInBytes int64  `json:"index_writer_max_memory_in_bytes" protobuf:"varint,17,opt,name=index_writer_max_memory_in_bytes,json=indexWriterMaxMemoryInBytes"`
	VersionMapMemory            string `json:"version_map_memory" protobuf:"bytes,18,opt,name=version_map_memory,json=versionMapMemory"`
	VersionMapMemoryInBytes     int64  `json:"version_map_memory_in_bytes" protobuf:"varint,19,opt,name=version_map_memory_in_bytes,json=versionMapMemoryInBytes"`
	FixedBitSetMemory           string `json:"fixed_bit_set" protobuf:"bytes,20,opt,name=fixed_bit_set,json=fixedBitSet"` // not a typo
	FixedBitSetMemoryInBytes    int64  `json:"fixed_bit_set_memory_in_bytes" protobuf:"varint,21,opt,name=fixed_bit_set_memory_in_bytes,json=fixedBitSetMemoryInBytes"`
}

type NodesStatsTranslogStats struct {
	Operations  int64  `json:"operations" protobuf:"varint,1,opt,name=operations"`
	TotalSize   string `json:"size" protobuf:"bytes,4,opt,name=size"`
	SizeInBytes int64  `json:"size_in_bytes" protobuf:"varint,3,opt,name=size_in_bytes,json=sizeInBytes"`
}

type NodesStatsRequestCacheStats struct {
	MemorySize        string `json:"memory_size" protobuf:"bytes,1,opt,name=memory_size,json=memorySize"`
	MemorySizeInBytes int64  `json:"memory_size_in_bytes" protobuf:"varint,2,opt,name=memory_size_in_bytes,json=memorySizeInBytes"`
	Evictions         int64  `json:"evictions" protobuf:"varint,3,opt,name=evictions"`
	HitCount          int64  `json:"hit_count" protobuf:"varint,4,opt,name=hit_count,json=hitCount"`
	MissCount         int64  `json:"miss_count" protobuf:"varint,5,opt,name=miss_count,json=missCount"`
}

type NodesStatsRecoveryStats struct {
	CurrentAsSource int64 `json:"current_as_source" protobuf:"varint,1,opt,name=current_as_source,json=currentAsSource"`
	CurrentAsTarget int64 `json:"current_as_target" protobuf:"varint,2,opt,name=current_as_target,json=currentAsTarget"`
}

type NodesStatsNodeOS struct {
	Timestamp int64                 `json:"timestamp" protobuf:"varint,1,opt,name=timestamp"`
	CPU       *NodesStatsNodeOSCPU  `json:"cpu" protobuf:"bytes,2,opt,name=cpu"`
	Mem       *NodesStatsNodeOSMem  `json:"mem" protobuf:"bytes,3,opt,name=mem"`
	Swap      *NodesStatsNodeOSSwap `json:"swap" protobuf:"bytes,4,opt,name=swap"`
}

type NodesStatsNodeOSCPU struct {
	Percent     int64              `json:"percent" protobuf:"varint,1,opt,name=percent"`
	LoadAverage map[string]float64 `json:"load_average" protobuf:"bytes,2,rep,name=load_average,json=loadAverage"` // keys are: 1m, 5m, and 15m
}

type NodesStatsNodeOSMem struct {
	Total        string `json:"total" protobuf:"bytes,1,opt,name=total"`
	TotalInBytes int64  `json:"total_in_bytes" protobuf:"varint,2,opt,name=total_in_bytes,json=totalInBytes"`
	Free         string `json:"free" protobuf:"bytes,3,opt,name=free"`
	FreeInBytes  int64  `json:"free_in_bytes" protobuf:"varint,4,opt,name=free_in_bytes,json=freeInBytes"`
	Used         string `json:"used" protobuf:"bytes,5,opt,name=used"`
	UsedInBytes  int64  `json:"used_in_bytes" protobuf:"varint,6,opt,name=used_in_bytes,json=usedInBytes"`
	FreePercent  int64  `json:"free_percent" protobuf:"varint,7,opt,name=free_percent,json=freePercent"`
	UsedPercent  int64  `json:"used_percent" protobuf:"varint,8,opt,name=used_percent,json=usedPercent"`
}

type NodesStatsNodeOSSwap struct {
	Total        string `json:"total" protobuf:"bytes,1,opt,name=total"`
	TotalInBytes int64  `json:"total_in_bytes" protobuf:"varint,2,opt,name=total_in_bytes,json=totalInBytes"`
	Free         string `json:"free" protobuf:"bytes,3,opt,name=free"`
	FreeInBytes  int64  `json:"free_in_bytes" protobuf:"varint,4,opt,name=free_in_bytes,json=freeInBytes"`
	Used         string `json:"used" protobuf:"bytes,5,opt,name=used"`
	UsedInBytes  int64  `json:"used_in_bytes" protobuf:"varint,6,opt,name=used_in_bytes,json=usedInBytes"`
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
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   ElasticsearchNodesStatsSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status ElasticsearchNodesStatsStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// ElasticsearchNodesStatsList contains a list of ElasticsearchNodesStats

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ElasticsearchNodesStatsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []ElasticsearchNodesStats `json:"items" protobuf:"bytes,2,rep,name=items"`
}

func init() {
	SchemeBuilder.Register(&ElasticsearchNodesStats{}, &ElasticsearchNodesStatsList{})
}
