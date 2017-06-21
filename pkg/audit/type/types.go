package types

import (
	esaudit "github.com/k8sdb/elasticsearch/pkg/audit/type"
	pgaudit "github.com/k8sdb/postgres/pkg/audit/type"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Summary struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Report            SummaryReport `json:"report,omitempty"`
	Status            SummaryStatus `json:"status,omitempty"`
}

type SummaryReport struct {
	Postgres map[string]*pgaudit.DBInfo    `json:"postgres,omitempty"`
	Elastic  map[string]*esaudit.IndexInfo `json:"elastic,omitempty"`
}

type SummaryStatus struct {
	StartTime      *metav1.Time `json:"startTime,omitempty"`
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`
}
