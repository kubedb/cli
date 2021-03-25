package resumer

import (
	"context"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	cs "kubedb.dev/apimachinery/client/clientset/versioned/typed/kubedb/v1alpha2"
	dbutil "kubedb.dev/apimachinery/client/clientset/versioned/typed/kubedb/v1alpha2/util"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	kmapi "kmodules.xyz/client-go/api/v1"
)

type MongoDBResumer struct {
	dbClient cs.KubedbV1alpha2Interface
}

func NewMongoDBResumer(clientConfig *rest.Config) (*MongoDBResumer, error) {
	k, err := cs.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	return &MongoDBResumer{
		dbClient: k,
	}, nil
}

func (e *MongoDBResumer) Resume(name, namespace string) error {
	db, err := e.dbClient.MongoDBs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	_, err = dbutil.UpdateMongoDBStatus(context.TODO(), e.dbClient, db.ObjectMeta, func(status *api.MongoDBStatus) (types.UID, *api.MongoDBStatus) {
		status.Conditions = kmapi.RemoveCondition(status.Conditions, api.DatabasePaused)
		return db.UID, status
	}, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
