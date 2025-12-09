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

package redis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"time"

	"kubedb.dev/apimachinery/apis/kubedb"
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/apimachinery/pkg/lib"

	rd "github.com/redis/go-redis/v9"
	vsecretapi "go.virtual-secrets.dev/apimachinery/apis/virtual/v1alpha1"
	core "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	DefaultDialTimeout     = 15 * time.Second
	DefaultConnMaxIdleTime = 3 * time.Second
	DefaultPoolSize        = 1
)

type KubeDBClientBuilder struct {
	kc       client.Client
	db       *dbapi.Redis
	podName  string
	url      string
	database int
}

func NewKubeDBClientBuilder(kc client.Client, db *dbapi.Redis) *KubeDBClientBuilder {
	return &KubeDBClientBuilder{
		kc: kc,
		db: db,
	}
}

func (o *KubeDBClientBuilder) WithPod(podName string) *KubeDBClientBuilder {
	o.podName = podName
	return o
}

func (o *KubeDBClientBuilder) WithURL(url string) *KubeDBClientBuilder {
	o.url = url
	return o
}

func (o *KubeDBClientBuilder) WithDatabase(database int) *KubeDBClientBuilder {
	o.database = database
	return o
}

func (o *KubeDBClientBuilder) GetRedisClient(ctx context.Context) (*Client, error) {
	var err error
	if o.podName != "" {
		o.url = o.getPodURL()
	}
	if o.url == "" {
		o.url = o.db.Address()
	}
	rdOpts := &rd.Options{
		DialTimeout:     DefaultDialTimeout,
		ConnMaxIdleTime: DefaultConnMaxIdleTime,
		PoolSize:        DefaultPoolSize,
		Addr:            o.url,
		DB:              o.database,
	}
	if !o.db.Spec.DisableAuth {
		rdOpts.Password, err = o.getClientPassword(ctx)
		if err != nil {
			return nil, err
		}
	}
	if o.db.Spec.TLS != nil {
		rdOpts.TLSConfig, err = o.getTLSConfig(ctx)
		if err != nil {
			return nil, err
		}
	}

	rdClient := rd.NewClient(rdOpts)
	_, err = rdClient.Ping(ctx).Result()
	if err != nil {
		closeErr := rdClient.Close() // nolint:errcheck
		if closeErr != nil {
			klog.Errorf("Failed to close client. error: %v", closeErr)
		}
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return &Client{
		rdClient,
	}, nil
}

func (o *KubeDBClientBuilder) GetRedisClusterClient(ctx context.Context) (*ClusterClient, error) {
	var err error
	if o.podName != "" {
		o.url = o.getPodURL()
	}
	if o.url == "" {
		o.url = o.db.Address()
	}
	rdClusterOpts := &rd.ClusterOptions{
		DialTimeout:     DefaultDialTimeout,
		ConnMaxIdleTime: DefaultConnMaxIdleTime,
		PoolSize:        DefaultPoolSize,
		Addrs:           []string{o.url},
	}
	if !o.db.Spec.DisableAuth {
		rdClusterOpts.Password, err = o.getClientPassword(ctx)
		if err != nil {
			return nil, err
		}
	}
	if o.db.Spec.TLS != nil {
		rdClusterOpts.TLSConfig, err = o.getTLSConfig(ctx)
		if err != nil {
			return nil, err
		}
		rdClusterOpts.TLSConfig.InsecureSkipVerify = true
	}
	rdClient := rd.NewClusterClient(rdClusterOpts)
	_, err = rdClient.Ping(ctx).Result()
	if err != nil {
		closeErr := rdClient.Close() // nolint:errcheck
		if closeErr != nil {
			klog.Errorf("Failed to close client. error: %v", closeErr)
		}
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return &ClusterClient{
		rdClient,
	}, nil
}

func (o *KubeDBClientBuilder) getClientPassword(ctx context.Context) (string, error) {
	if o.db.Spec.AuthSecret == nil || o.db.Spec.AuthSecret.Name == "" {
		return "", errors.New("no database secret")
	}
	var clientPass string
	if !dbapi.IsVirtualAuthSecretReferred(o.db.Spec.AuthSecret) {
		var authSecret core.Secret
		err := o.kc.Get(ctx, client.ObjectKey{Namespace: o.db.Namespace, Name: o.db.Spec.AuthSecret.Name}, &authSecret)
		if err != nil {
			return "", err
		}
		clientPass = string(authSecret.Data[core.BasicAuthPasswordKey])
		return clientPass, nil
	}
	vSecret := &vsecretapi.Secret{}
	err := o.kc.Get(context.TODO(), client.ObjectKey{Namespace: o.db.Namespace, Name: o.db.Spec.AuthSecret.Name}, vSecret)
	if err != nil {
		return "", err
	}

	if err = lib.ValidateVirtualAuthSecret(vSecret); err != nil {
		return "", err
	}
	return string(vSecret.Data[core.BasicAuthPasswordKey]), nil
}

func (o *KubeDBClientBuilder) getTLSConfig(ctx context.Context) (*tls.Config, error) {
	var sec core.Secret
	err := o.kc.Get(ctx, client.ObjectKey{Namespace: o.db.Namespace, Name: o.db.CertificateName(dbapi.RedisClientCert)}, &sec)
	if err != nil {
		klog.Error(err, "error in getting the secret")
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(sec.Data["ca.crt"])
	cert, err := tls.X509KeyPair(sec.Data["tls.crt"], sec.Data["tls.key"])
	if err != nil {
		klog.Error(err, "error in making certificate")
		return nil, err
	}
	clientTlS := &tls.Config{
		Certificates: []tls.Certificate{
			cert,
		},
		ClientCAs: pool,
		RootCAs:   pool,
	}
	return clientTlS, nil
}

func (o *KubeDBClientBuilder) getPodURL() string {
	return fmt.Sprintf("%v.%v.%v.svc:%d", o.podName, o.db.GoverningServiceName(), o.db.Namespace, kubedb.RedisDatabasePort)
}
