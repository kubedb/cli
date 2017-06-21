package postgres

import (
	"github.com/appscode/go/net/httpclient"
	"github.com/k8sdb/cli/pkg/audit/type"
	esaudit "github.com/k8sdb/elasticsearch/pkg/audit/type"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetReport(client *httpclient.Client, req *http.Request) (*types.Summary, error) {
	startTime := metav1.Now()
	infos := make(map[string]*esaudit.IndexInfo)
	if _, err := client.Do(req, &infos); err != nil {
		return &types.Summary{}, err
	}

	completionTime := metav1.Now()
	summary := &types.Summary{
		Report: types.SummaryReport{
			Elastic: infos,
		},
		Status: types.SummaryStatus{
			StartTime:      &startTime,
			CompletionTime: &completionTime,
		},
	}

	return summary, nil
}
