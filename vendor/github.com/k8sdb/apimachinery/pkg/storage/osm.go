package storage

import (
	"bytes"
	"errors"
	"net/url"
	"strconv"

	otx "github.com/appscode/osm/pkg/context"
	"github.com/ghodss/yaml"
	"github.com/graymeta/stow"
	"github.com/graymeta/stow/azure"
	gcs "github.com/graymeta/stow/google"
	"github.com/graymeta/stow/local"
	"github.com/graymeta/stow/s3"
	"github.com/graymeta/stow/swift"
	tapi "github.com/k8sdb/apimachinery/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

const (
	SecretMountPath = "/etc/osm"
)

func NewOSMSecret(client clientset.Interface, snapshot *tapi.Snapshot, namespace string) (*apiv1.Secret, error) {
	osmCtx, err := NewOSMContext(client, snapshot.Spec.SnapshotStorageSpec, namespace)
	if err != nil {
		return nil, err
	}
	osmCfg := &otx.OSMConfig{
		CurrentContext: osmCtx.Name,
		Contexts:       []*otx.Context{osmCtx},
	}
	osmBytes, err := yaml.Marshal(osmCfg)
	if err != nil {
		return nil, err
	}
	return &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      snapshot.Name,
			Namespace: snapshot.Namespace,
		},
		Data: map[string][]byte{
			"config": osmBytes,
		},
	}, nil
}

func CheckBucketAccess(client clientset.Interface, spec tapi.SnapshotStorageSpec, namespace string) error {
	cfg, err := NewOSMContext(client, spec, namespace)
	if err != nil {
		return err
	}
	loc, err := stow.Dial(cfg.Provider, cfg.Config)
	if err != nil {
		return err
	}
	c, err := GetContainer(spec)
	if err != nil {
		return err
	}
	container, err := loc.Container(c)
	if err != nil {
		return err
	}
	r := bytes.NewReader([]byte("CheckBucketAccess"))
	item, err := container.Put(".kubedb", r, r.Size(), nil)
	if err != nil {
		return err
	}
	if err := container.RemoveItem(item.ID()); err != nil {
		return err
	}
	return nil
}

func NewOSMContext(client clientset.Interface, spec tapi.SnapshotStorageSpec, namespace string) (*otx.Context, error) {
	secret, err := client.CoreV1().Secrets(namespace).Get(spec.StorageSecretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	nc := &otx.Context{
		Name:   "kubedb",
		Config: stow.ConfigMap{},
	}

	if spec.S3 != nil {
		nc.Provider = s3.Kind
		nc.Config[s3.ConfigAccessKeyID] = string(secret.Data[tapi.AWS_ACCESS_KEY_ID])
		nc.Config[s3.ConfigEndpoint] = spec.S3.Endpoint
		nc.Config[s3.ConfigRegion] = spec.S3.Region
		nc.Config[s3.ConfigSecretKey] = string(secret.Data[tapi.AWS_SECRET_ACCESS_KEY])
		if u, err := url.Parse(spec.S3.Endpoint); err == nil {
			nc.Config[s3.ConfigDisableSSL] = strconv.FormatBool(u.Scheme == "http")
		}
		return nc, nil
	} else if spec.GCS != nil {
		nc.Provider = gcs.Kind
		nc.Config[gcs.ConfigProjectId] = string(secret.Data[tapi.GOOGLE_PROJECT_ID])
		nc.Config[gcs.ConfigJSON] = string(secret.Data[tapi.GOOGLE_SERVICE_ACCOUNT_JSON_KEY])
		return nc, nil
	} else if spec.Azure != nil {
		nc.Provider = azure.Kind
		nc.Config[azure.ConfigAccount] = string(secret.Data[tapi.AZURE_ACCOUNT_NAME])
		nc.Config[azure.ConfigKey] = string(secret.Data[tapi.AZURE_ACCOUNT_KEY])
		return nc, nil
	} else if spec.Local != nil {
		nc.Provider = local.Kind
		nc.Config[local.ConfigKeyPath] = spec.Local.Path
		return nc, nil
	} else if spec.Swift != nil {
		nc.Provider = swift.Kind
		nc.Config[swift.ConfigKey] = string(secret.Data[tapi.OS_PASSWORD])
		nc.Config[swift.ConfigTenantAuthURL] = string(secret.Data[tapi.OS_AUTH_URL])
		nc.Config[swift.ConfigTenantName] = string(secret.Data[tapi.OS_TENANT_NAME])
		nc.Config[swift.ConfigUsername] = string(secret.Data[tapi.OS_USERNAME])
		return nc, nil
	}
	return nil, errors.New("No storage provider is configured.")
}

func GetContainer(spec tapi.SnapshotStorageSpec) (string, error) {
	if spec.S3 != nil {
		return spec.S3.Bucket, nil
	} else if spec.GCS != nil {
		return spec.GCS.Bucket, nil
	} else if spec.Azure != nil {
		return spec.Azure.Container, nil
	} else if spec.Local != nil {
		return "kubedb", nil
	} else if spec.Swift != nil {
		return spec.Swift.Container, nil
	}
	return "", errors.New("No storage provider is configured.")
}
