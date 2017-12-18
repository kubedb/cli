package storage

import (
	"bytes"
	"errors"
	"net/url"
	"strconv"

	otx "github.com/appscode/osm/context"
	"github.com/ghodss/yaml"
	"github.com/graymeta/stow"
	"github.com/graymeta/stow/azure"
	gcs "github.com/graymeta/stow/google"
	"github.com/graymeta/stow/local"
	"github.com/graymeta/stow/s3"
	"github.com/graymeta/stow/swift"
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	SecretMountPath = "/etc/osm"
)

func NewOSMSecret(client kubernetes.Interface, snapshot *api.Snapshot) (*core.Secret, error) {
	osmCtx, err := NewOSMContext(client, snapshot.Spec.SnapshotStorageSpec, snapshot.Namespace)
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
	return &core.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      snapshot.OSMSecretName(),
			Namespace: snapshot.Namespace,
		},
		Data: map[string][]byte{
			"config": osmBytes,
		},
	}, nil
}

func CheckBucketAccess(client kubernetes.Interface, spec api.SnapshotStorageSpec, namespace string) error {
	cfg, err := NewOSMContext(client, spec, namespace)
	if err != nil {
		return err
	}
	loc, err := stow.Dial(cfg.Provider, cfg.Config)
	if err != nil {
		return err
	}
	c, err := spec.Container()
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

func NewOSMContext(client kubernetes.Interface, spec api.SnapshotStorageSpec, namespace string) (*otx.Context, error) {
	config := make(map[string][]byte)

	if spec.StorageSecretName != "" {
		secret, err := client.CoreV1().Secrets(namespace).Get(spec.StorageSecretName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		config = secret.Data
	}

	nc := &otx.Context{
		Name:   "kubedb",
		Config: stow.ConfigMap{},
	}

	if spec.S3 != nil {
		nc.Provider = s3.Kind
		nc.Config[s3.ConfigAccessKeyID] = string(config[api.AWS_ACCESS_KEY_ID])
		nc.Config[s3.ConfigEndpoint] = spec.S3.Endpoint
		nc.Config[s3.ConfigRegion] = "us-east-1" // only used for creating buckets
		nc.Config[s3.ConfigSecretKey] = string(config[api.AWS_SECRET_ACCESS_KEY])
		if u, err := url.Parse(spec.S3.Endpoint); err == nil {
			nc.Config[s3.ConfigDisableSSL] = strconv.FormatBool(u.Scheme == "http")
		}
		return nc, nil
	} else if spec.GCS != nil {
		nc.Provider = gcs.Kind
		nc.Config[gcs.ConfigProjectId] = string(config[api.GOOGLE_PROJECT_ID])
		nc.Config[gcs.ConfigJSON] = string(config[api.GOOGLE_SERVICE_ACCOUNT_JSON_KEY])
		return nc, nil
	} else if spec.Azure != nil {
		nc.Provider = azure.Kind
		nc.Config[azure.ConfigAccount] = string(config[api.AZURE_ACCOUNT_NAME])
		nc.Config[azure.ConfigKey] = string(config[api.AZURE_ACCOUNT_KEY])
		return nc, nil
	} else if spec.Local != nil {
		nc.Provider = local.Kind
		nc.Config[local.ConfigKeyPath] = spec.Local.Path
		return nc, nil
	} else if spec.Swift != nil {
		nc.Provider = swift.Kind
		// https://github.com/restic/restic/blob/master/src/restic/backend/swift/config.go
		for _, val := range []struct {
			stowKey   string
			secretKey string
		}{
			// v2/v3 specific
			{swift.ConfigUsername, api.OS_USERNAME},
			{swift.ConfigKey, api.OS_PASSWORD},
			{swift.ConfigRegion, api.OS_REGION_NAME},
			{swift.ConfigTenantAuthURL, api.OS_AUTH_URL},

			// v3 specific
			{swift.ConfigDomain, api.OS_USER_DOMAIN_NAME},
			{swift.ConfigTenantName, api.OS_PROJECT_NAME},
			{swift.ConfigTenantDomain, api.OS_PROJECT_DOMAIN_NAME},

			// v2 specific
			{swift.ConfigTenantId, api.OS_TENANT_ID},
			{swift.ConfigTenantName, api.OS_TENANT_NAME},

			// v1 specific
			{swift.ConfigTenantAuthURL, api.ST_AUTH},
			{swift.ConfigUsername, api.ST_USER},
			{swift.ConfigKey, api.ST_KEY},

			// Manual authentication
			{swift.ConfigStorageURL, api.OS_STORAGE_URL},
			{swift.ConfigAuthToken, api.OS_AUTH_TOKEN},
		} {
			if _, exists := nc.Config.Config(val.stowKey); !exists {
				nc.Config[val.stowKey] = string(config[val.secretKey])
			}
		}
		return nc, nil
	}
	return nil, errors.New("no storage provider is configured")
}
