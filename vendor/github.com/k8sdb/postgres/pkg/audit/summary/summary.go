package summary

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-ini/ini"
	tcs "github.com/k8sdb/apimachinery/client/clientset"
	"github.com/k8sdb/postgres/pkg/audit/type"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
)

func GetSummaryReport(
	kubeClient clientset.Interface,
	dbClient tcs.ExtensionInterface,
	namespace string,
	kubedbName string,
	dbname string,
	w http.ResponseWriter,
) {
	postgres, err := dbClient.Postgreses(namespace).Get(kubedbName)
	if err != nil {
		if kerr.IsNotFound(err) {
			http.Error(w, fmt.Sprintf(`Postgres "%v" not found`, kubedbName), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	secret, err := kubeClient.CoreV1().Secrets(namespace).Get(postgres.Spec.DatabaseSecret.SecretName, metav1.GetOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			http.Error(w, fmt.Sprintf(`Secret "%v" not found`, postgres.Spec.DatabaseSecret.SecretName), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	cfg, err := ini.Load(secret.Data[".admin"])
	if err != nil {
		http.Error(w, fmt.Sprintf(`secret key ".admin" not found`), http.StatusNotFound)
		return
	}
	section, err := cfg.GetSection("")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	username := "postgres"
	if k, err := section.GetKey("POSTGRES_USER"); err == nil {
		username = k.Value()
	}
	var password string
	if k, err := section.GetKey("POSTGRES_PASSWORD"); err == nil {
		password = k.Value()
	} else {
		http.Error(w, fmt.Sprintf(`POSTGRES_PASSWORD not found in secret`), http.StatusNotFound)
		return
	}

	host := fmt.Sprintf("%v.%v", kubedbName, namespace)
	port := "5432"

	databases := make([]string, 0)
	if dbname == "" {
		engine, err := newXormEngine(username, password, host, port, "postgres")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		databases, err = getAllDatabase(engine)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		databases = append(databases, dbname)
	}

	dbs := make(map[string]*types.DBInfo)
	for _, db := range databases {
		engine, err := newXormEngine(username, password, host, port, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dbInfo, err := dumpDBInfo(engine)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dbs[db] = dbInfo
	}

	data, err := json.MarshalIndent(dbs, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if data != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, string(data))
	} else {
		http.Error(w, "audit data not found", http.StatusNotFound)
	}
}
