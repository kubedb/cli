package pauser

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

type RedisPauser struct {
	dbClient cs.KubedbV1alpha2Interface
}

func NewRedisPauser(clientConfig *rest.Config) (*RedisPauser, error) {
	k, err := cs.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	return &RedisPauser{
		dbClient: k,
	}, nil
}

func (e *RedisPauser) Pause(name, namespace string) error {
	db, err := e.dbClient.Redises(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	_, err = dbutil.UpdateRedisStatus(context.TODO(), e.dbClient, db.ObjectMeta, func(status *api.RedisStatus) (types.UID, *api.RedisStatus) {
		status.Conditions = kmapi.SetCondition(status.Conditions, kmapi.NewCondition(
			api.DatabasePaused,
			"Paused by KubeDB CLI tool",
			db.Generation,
		))
		return db.UID, status
	}, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
