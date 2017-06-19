package types

import (
	"github.com/k8sdb/postgres/pkg/audit/summary/lib"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Summary struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	SummaryReport     SummaryReport `json:"summaryReport,omitempty"`
	Status            SummaryStatus `json:"status,omitempty"`
}

type SummaryReport struct {
	Postgres map[string]*lib.DBInfo `json:"postgres,omitempty"`
}

type SummaryStatus struct {
	StartTime      *metav1.Time `json:"startTime,omitempty"`
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`
}
