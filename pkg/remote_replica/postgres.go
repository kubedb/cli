/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package remote_replica

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/cli/pkg/common"

	cm_api "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	cm "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"github.com/spf13/cobra"
	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	cm_util "kmodules.xyz/cert-manager-util/certmanager/v1"
	kutil "kmodules.xyz/client-go"
	kmapi "kmodules.xyz/client-go/api/v1"
	core_util "kmodules.xyz/client-go/core/v1"
	"kmodules.xyz/client-go/meta"
	exec_util "kmodules.xyz/client-go/tools/exec"
	appApi "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	"sigs.k8s.io/yaml"
)

func PostgreSQlAPP(f cmdutil.Factory) *cobra.Command {
	var userName, password, dns, ns, authSecretName, replicaName string
	var port int32
	var yes bool

	cmd := cobra.Command{
		Use:     "postgres",
		Short:   desLong,
		Long:    desLong,
		Example: example,
		Args:    nil,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("no database name given")
			}
			if err := userPrompt(yes); err != nil {
				log.Fatal(err)
			}

			// Accept -d host:port as a convenience. An explicit --port always wins.
			if host, p, found := strings.Cut(dns, ":"); found && strings.Count(dns, ":") == 1 {
				if v, convErr := strconv.Atoi(p); convErr == nil && v > 0 && v < 65536 {
					if !cmd.Flags().Changed("port") {
						port = int32(v)
					}
					dns = host
				}
			}

			var buffer []byte
			buffer, err := generateConfig(f, userName, password, dns, ns, authSecretName, replicaName, port, args[0])
			if err != nil {
				log.Fatal(err)
			}

			directory, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			err = os.WriteFile(fmt.Sprintf(directory+"/%s-remote-config.yaml", args[0]), buffer, 0o644)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("kubectl apply -f  %s/%s-remote-config.yaml\n", directory, args[0])
		},
		DisableAutoGenTag:     false,
		DisableFlagsInUseLine: false,
	}

	cmd.PersistentFlags().StringVarP(&userName, "user", "u", "postgres", "user name for the remote replica")
	if err := cmd.MarkPersistentFlagRequired("user"); err != nil {
		log.Fatal(err)
	}
	cmd.PersistentFlags().StringVarP(&password, "pass", "p", "password", "password name for the remote replica")
	if err := cmd.MarkPersistentFlagRequired("pass"); err != nil {
		log.Fatal(err)
	}
	cmd.PersistentFlags().StringVarP(&dns, "dns", "d", "localhost", "dns name for the remote replica")
	if err := cmd.MarkPersistentFlagRequired("dns"); err != nil {
		log.Fatal(err)
	}
	cmd.PersistentFlags().StringVarP(&ns, "namespace", "n", "default", "host namespace for the remote replica")
	if err := cmd.MarkPersistentFlagRequired("namespace"); err != nil {
		log.Fatal(err)
	}
	cmd.PersistentFlags().BoolVarP(&yes, "yes", "y", false, "permission for alter password  for the remote replica")
	cmd.PersistentFlags().StringVar(&authSecretName, "auth-secret", "", "name for the auth secret on the remote cluster (default: <dbname>-remote-replica-auth)")
	cmd.PersistentFlags().Int32Var(&port, "port", 5432, "port the source is reachable on from the remote cluster; written into the generated AppBinding (also accepted as -d host:port)")
	cmd.PersistentFlags().StringVar(&replicaName, "replica-name", "", "when set, also emit a ready-to-apply remote replica Postgres manifest with this name, sized from the source spec")
	return &cmd
}

func generateConfig(f cmdutil.Factory, userName string, password string, dns string, ns string, authSecretName string, replicaName string, port int32, dbname string) ([]byte, error) {
	var buffer []byte
	opts, err := common.NewPostgresOpts(f, dbname, ns)
	if err != nil {
		return nil, fmt.Errorf("failed to get db %s, err:%v", dbname, err)
	}

	apb, err := opts.AppcatClient.AppcatalogV1alpha1().AppBindings(ns).Get(context.TODO(), dbname, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("failed to get appbinding %v", err)
	}

	authBuff, authSecretName, err := generateAuthSecret(userName, password, ns, authSecretName, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth secret ,%v", err)
	}
	buffer = append(buffer, authBuff...)

	var tlsSecretName string
	if apb.Spec.TLSSecret != nil {
		var tlsBuff []byte
		tlsBuff, tlsSecretName, err = generateTlsSecret(userName, apb, ns, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to generate tls secret %v", err)
		}
		buffer = append(buffer, tlsBuff...)
	}

	// Deep-copy the source AppBinding and replace only the ObjectMeta with a
	// clean one. This preserves all spec fields (appRef, parameters, type,
	// version, clientConfig, etc.) so nothing is silently dropped, while
	// ensuring server-managed metadata (resourceVersion, uid, generation,
	// labels, annotations) never leaks into the generated YAML. A clean
	// ObjectMeta means the last-applied annotation stays minimal, so
	// repeated kubectl apply is idempotent and the 3-way merge never tries
	// to remove metadata.resourceVersion.
	remoteApb := apb.DeepCopy()
	remoteApb.TypeMeta = metav1.TypeMeta{
		APIVersion: AppcatApiVersion,
		Kind:       AppcatKind,
	}
	remoteApb.ObjectMeta = metav1.ObjectMeta{
		Name:      apb.Name,
		Namespace: ns,
	}
	if remoteApb.Spec.ClientConfig.Service == nil {
		remoteApb.Spec.ClientConfig.Service = &appApi.ServiceReference{}
	}
	remoteApb.Spec.ClientConfig.Service.Name = dns
	// The port the source is reachable on FROM THE REMOTE CLUSTER (a load balancer
	// frontend, not necessarily 5432). The operator injects it as PRIMARY_PORT into
	// the remote replica containers.
	remoteApb.Spec.ClientConfig.Service.Port = port
	if remoteApb.Spec.Secret == nil {
		remoteApb.Spec.Secret = &appApi.TypedLocalObjectReference{}
	}
	remoteApb.Spec.Secret.Name = authSecretName
	if tlsSecretName != "" {
		if remoteApb.Spec.TLSSecret == nil {
			remoteApb.Spec.TLSSecret = &appApi.TypedLocalObjectReference{}
		}
		remoteApb.Spec.TLSSecret.Name = tlsSecretName
	}

	appbindingYaml, err := yaml.Marshal(remoteApb)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal appbind yaml %v", err)
	}

	buffer = append(buffer, appbindingYaml...)

	if replicaName != "" {
		replicaYaml, err := generateReplicaSpec(opts.DB, replicaName, ns, apb.Name, authSecretName)
		if err != nil {
			return nil, fmt.Errorf("failed to generate replica spec %v", err)
		}
		buffer = append(buffer, []byte("---\n")...)
		buffer = append(buffer, replicaYaml...)
	}
	return buffer, nil
}

// generateReplicaSpec emits a ready-to-apply remote replica Postgres manifest sized from
// the source's spec. It copies the fields that describe capacity (version, replicas,
// storage, the postgres container's resources, standby mode) and adds what a remote
// replica needs: the remoteReplica sourceRef pointing at the generated AppBinding and
// the generated auth secret. The health checker needs no tuning here — it already skips
// write checks for remote replicas on its own.
//
// Deliberately NOT copied: spec.tls (the remote cluster has its own issuer),
// spec.monitor, archiver, init, and any custom sidecars. clientAuthMode falls back from
// cert to md5, since cert auth requires the TLS stanza that is not carried over.
// deletionPolicy is set to Halt regardless of the source: a DR replica's PVCs should
// survive an accidental CR deletion. Treat the result as a starting point — a secondary
// site is often sized differently on purpose.
func generateReplicaSpec(src *dbapi.Postgres, name, ns, sourceRefName, authSecretName string) ([]byte, error) {
	replica := dbapi.Postgres{}
	replica.APIVersion = dbapi.SchemeGroupVersion.String()
	replica.Kind = dbapi.ResourceKindPostgres
	replica.Name = name
	replica.Namespace = ns

	replica.Spec.Version = src.Spec.Version
	replica.Spec.Replicas = src.Spec.Replicas
	replica.Spec.StorageType = src.Spec.StorageType
	replica.Spec.Storage = src.Spec.Storage
	replica.Spec.StandbyMode = src.Spec.StandbyMode
	replica.Spec.DeletionPolicy = dbapi.DeletionPolicyHalt

	replica.Spec.ClientAuthMode = src.Spec.ClientAuthMode
	if replica.Spec.ClientAuthMode == dbapi.ClientAuthModeCert {
		replica.Spec.ClientAuthMode = dbapi.ClientAuthModeMD5
	}

	for _, c := range src.Spec.PodTemplate.Spec.Containers {
		if c.Name == "postgres" {
			replica.Spec.PodTemplate.Spec.Containers = []core.Container{{
				Name:      c.Name,
				Resources: c.Resources,
			}}
			break
		}
	}

	replica.Spec.AuthSecret = &dbapi.SecretReference{}
	replica.Spec.AuthSecret.Name = authSecretName
	replica.Spec.RemoteReplica = &dbapi.RemoteReplicaSpec{
		SourceRef: core.ObjectReference{
			Name:      sourceRefName,
			Namespace: ns,
		},
	}

	// Marshal via a map so the empty status stanza is dropped from the manifest.
	jsonBytes, err := json.Marshal(replica)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &m); err != nil {
		return nil, err
	}
	delete(m, "status")
	// healthChecker has no omitempty and would render as an empty stanza; the
	// operator's defaulting fills it, and it already handles remote replicas.
	if spec, ok := m["spec"].(map[string]interface{}); ok {
		if hc, ok := spec["healthChecker"].(map[string]interface{}); ok && len(hc) == 0 {
			delete(spec, "healthChecker")
		}
	}
	return yaml.Marshal(m)
}

func generateTlsSecret(userName string, apb *appApi.AppBinding, ns string, opts *common.PostgresOpts) ([]byte, string, error) {
	_, err := ensureClientCert(opts, apb, opts.DB, dbapi.PostgresClientCert, userName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to ensure client cert %v", err)
	}
	tlsSecret := &core.Secret{}

	err = wait.PollUntilContextTimeout(context.Background(), 300*time.Millisecond, 60*time.Minute, true, func(ctx context.Context) (done bool, err error) {
		sercretName := opts.DB.GetCertSecretName(dbapi.PostgresClientCert) + fmt.Sprintf("-%s", userName)

		tlsSecret, err = opts.Client.CoreV1().Secrets(ns).Get(ctx, sercretName, metav1.GetOptions{})
		if kerr.IsNotFound(err) {
			return false, nil
		} else if err != nil {
			return false, err
		}

		return true, nil
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to get tls secret %v", err)
	}
	tlsSecret.APIVersion = "v1"
	tlsSecret.Kind = "Secret"
	tlsSecret.ResourceVersion = ""
	tlsSecret.UID = ""
	tlsSecret.CreationTimestamp = metav1.Time{}
	tlsSecret.Annotations = nil
	tlsSecret.Labels = nil
	tlsSecret.ManagedFields = nil
	tlsSecretYaml, err := yaml.Marshal(tlsSecret)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal tls secret yaml %v", err)
	}

	buffer := make([]byte, 0, len(tlsSecretYaml)+4)
	buffer = append(buffer, tlsSecretYaml...)
	buffer = append(buffer, []byte("---\n")...)

	return buffer, tlsSecret.Name, nil
}

func generateAuthSecret(userName string, password string, ns string, secretName string, opts *common.PostgresOpts) ([]byte, string, error) {
	if userName != opts.Username {
		// generate user if not present
		err := generateUser(opts, userName, password)
		if err != nil {
			return nil, "", fmt.Errorf("failed to generate user err:%v", err)
		}
	} else {
		password = opts.Pass
	}
	if secretName == "" {
		secretName = fmt.Sprintf("%s-remote-replica-auth", opts.DB.Name)
	}
	// generate auth secret
	AuthSecret := core.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       KindSecret,
			APIVersion: ApiversionV1,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: ns,
		},
		StringData: map[string]string{
			"username": userName,
			"password": password,
		},
		Type: core.SecretTypeBasicAuth,
	}

	authSecretYaml, err := yaml.Marshal(AuthSecret)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal authsecret yaml %v", err)
	}
	buffer := make([]byte, 0, len(authSecretYaml)+4)
	buffer = append(buffer, authSecretYaml...)
	buffer = append(buffer, []byte("---\n")...)
	return buffer, AuthSecret.Name, nil
}

func generateUser(opts *common.PostgresOpts, name string, password string) error {
	label := opts.DB.OffshootLabels()
	label["kubedb.com/role"] = "primary"
	pods, err := opts.Client.CoreV1().Pods(opts.DB.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.Set.String(label),
	})
	if err != nil || len(pods.Items) == 0 {
		return err
	}

	query := fmt.Sprintf("SELECT rolname FROM pg_roles WHERE rolname='%s'", name)

	command := exec_util.Command("psql", "-c", query)
	container := exec_util.Container("postgres")
	options := []func(options *exec_util.Options){
		command,
		container,
	}

	out, err := exec_util.ExecIntoPod(opts.Config, &pods.Items[0], options...)
	if err != nil {
		return err
	}

	query = fmt.Sprintf("create user %s with password '%s'; alter role %s with replication; GRANT execute ON function pg_read_binary_file(text) TO %s;", name, password, name, name)
	if len(out) > 30 {
		query = fmt.Sprintf("alter role %s with password '%s' replication; GRANT execute ON function pg_read_binary_file(text) TO %s;", name, password, name)
	}

	command = exec_util.Command("psql", "-c", query)
	container = exec_util.Container("postgres")
	options = []func(options *exec_util.Options){
		command,
		container,
	}

	out, err = exec_util.ExecIntoPod(opts.Config, &pods.Items[0], options...)
	if err != nil {
		return err
	}
	fmt.Println(out)
	return nil
}

func ensureClientCert(opts *common.PostgresOpts, apb *appApi.AppBinding, postgres *dbapi.Postgres, alias dbapi.PostgresCertificateAlias, username string) (kutil.VerbType, error) {
	var duration, renewBefore *metav1.Duration
	var subject *cm_api.X509Subject
	var dnsNames, ipAddresses, uriSANs, emailSANs []string
	if _, cert := kmapi.GetCertificate(postgres.Spec.TLS.Certificates, string(alias)); cert != nil {
		dnsNames = cert.DNSNames
		ipAddresses = cert.IPAddresses
		duration = cert.Duration
		renewBefore = cert.RenewBefore
		if cert.Subject != nil {
			subject = &cm_api.X509Subject{
				Organizations:       cert.Subject.Organizations,
				Countries:           cert.Subject.Countries,
				OrganizationalUnits: cert.Subject.OrganizationalUnits,
				Localities:          cert.Subject.Localities,
				Provinces:           cert.Subject.Provinces,
				StreetAddresses:     cert.Subject.StreetAddresses,
				PostalCodes:         cert.Subject.PostalCodes,
				SerialNumber:        cert.Subject.SerialNumber,
			}
		}
		uriSANs = cert.URIs
		emailSANs = cert.EmailAddresses
	}

	ref := metav1.NewControllerRef(apb, appApi.SchemeGroupVersion.WithKind(appApi.ResourceKindApp))

	_, vt, err := cm_util.CreateOrPatchCertificate(
		context.TODO(),
		opts.CertManagerClient.CertmanagerV1(),
		metav1.ObjectMeta{
			Name:      postgres.CertificateName(alias) + fmt.Sprintf("-%s", username),
			Namespace: postgres.GetNamespace(),
		},
		func(in *cm_api.Certificate) *cm_api.Certificate {
			in.Labels = postgres.OffshootLabels()
			core_util.EnsureOwnerReference(in, ref)

			in.Spec.CommonName = username
			in.Spec.Subject = subject
			in.Spec.Duration = duration
			in.Spec.RenewBefore = renewBefore
			in.Spec.DNSNames = sets.NewString(dnsNames...).List()
			in.Spec.IPAddresses = sets.NewString(ipAddresses...).List()
			in.Spec.URIs = sets.NewString(uriSANs...).List()
			in.Spec.EmailAddresses = sets.NewString(emailSANs...).List()
			in.Spec.SecretName = postgres.GetCertSecretName(alias) + fmt.Sprintf("-%s", username)
			in.Spec.IssuerRef = GetIssuerObjectRef(postgres.Spec.TLS, string(alias))
			in.Spec.Usages = []cm_api.KeyUsage{
				cm_api.UsageDigitalSignature,
				cm_api.UsageKeyEncipherment,
				cm_api.UsageClientAuth,
			}
			pemEncodeCert := isCertMangerAdditionalOutputEnabled(opts.CertManagerClient)
			if pemEncodeCert {
				in.Spec.AdditionalOutputFormats = []cm_api.CertificateAdditionalOutputFormat{
					{
						Type: cm_api.CertificateOutputFormatCombinedPEM,
					},
				}
			}

			return in
		}, metav1.PatchOptions{},
	)

	return vt, err
}

func GetIssuerObjectRef(tlsConfig *kmapi.TLSConfig, alias string) cmmeta.ObjectReference {
	if _, cert := kmapi.GetCertificate(tlsConfig.Certificates, alias); cert != nil {
		issuer := tlsConfig.IssuerRef
		if cert.IssuerRef != nil {
			issuer = cert.IssuerRef
		}

		return cmmeta.ObjectReference{
			Name:  issuer.Name,
			Kind:  issuer.Kind,
			Group: pointer.String(issuer.APIGroup),
		}
	}

	return cmmeta.ObjectReference{}
}

func isCertMangerAdditionalOutputEnabled(certManagerClient cm.Interface) bool {
	operatorNs := meta.PodNamespace()
	demoCert := cm_api.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cert",
			Namespace: operatorNs,
		},
		Spec: cm_api.CertificateSpec{
			CommonName: "example.com",
			SecretName: "test-secret",
			IssuerRef: cmmeta.ObjectReference{
				Name: "test-issuer",
			},
			AdditionalOutputFormats: []cm_api.CertificateAdditionalOutputFormat{
				{
					Type: cm_api.CertificateOutputFormatCombinedPEM,
				},
			},
		},
	}

	_, err := certManagerClient.CertmanagerV1().Certificates(operatorNs).Create(context.TODO(), &demoCert, metav1.CreateOptions{
		DryRun: []string{
			"All",
		},
	})
	if err != nil {
		return false
	}

	klog.Info("Cert-Manager feature-gate AdditionalCertificateOutputFormats is enabled, certificates will include combined PEM output")

	return true
}

func userPrompt(yes bool) error {
	fmt.Println("password will be altered with the given password if provided user  exist you want to continue/Y/N?")
	if yes {
		return nil
	}
	var inp string
	_, err := fmt.Scan(&inp)
	if err != nil {
		return err
	}
	inp = strings.ToLower(inp)
	if inp != "y" {
		return errors.New("aborting commands")
	}
	return nil
}
