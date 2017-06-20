package postgres

import (
	"github.com/appscode/go/net/httpclient"
	"github.com/k8sdb/cli/pkg/audit/type"
	pgaudit "github.com/k8sdb/postgres/pkg/audit/type"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetReport(client *httpclient.Client, req *http.Request) (*types.Summary, error) {
	startTime := metav1.Now()
	dbs := make(map[string]*pgaudit.DBInfo)
	if _, err := client.Do(req, &dbs); err != nil {
		return &types.Summary{}, err
	}

	completionTime := metav1.Now()
	summary := &types.Summary{
		SummaryReport: types.SummaryReport{
			Postgres: dbs,
		},
		Status: types.SummaryStatus{
			StartTime:      &startTime,
			CompletionTime: &completionTime,
		},
	}

	return summary, nil
}
