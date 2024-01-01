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
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
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
	var userName, password, dns, ns string
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

			var buffer []byte
			buffer, err := generateConfig(f, userName, password, dns, ns, args[0])
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
	return &cmd
}

func generateConfig(f cmdutil.Factory, userName string, password string, dns string, ns string, dbname string) ([]byte, error) {
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

	authBuff, authSecretName, err := generateAuthSecret(userName, password, ns, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth secret ,%v", err)
	}
	buffer = append(buffer, authBuff...)

	// generate secret
	if apb.Spec.TLSSecret != nil {
		tlsBuff, tlsSecretName, err := generateTlsSecret(userName, apb, ns, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to generate tls secret %v", err)
		}
		buffer = append(buffer, tlsBuff...)
		apb.Spec.TLSSecret.Name = tlsSecretName
	}

	apb.APIVersion = AppcatApiVersion
	apb.Kind = AppcatKind
	apb.Spec.ClientConfig.Service.Name = dns
	apb.Spec.Secret.Name = authSecretName
	apb.ObjectMeta.Annotations = nil
	apb.ObjectMeta.ManagedFields = nil
	apb.ObjectMeta.OwnerReferences = nil

	appbindingYaml, err := yaml.Marshal(apb)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal appbind yaml %v", err)
	}

	buffer = append(buffer, appbindingYaml...)
	return buffer, nil
}

func generateTlsSecret(userName string, apb *appApi.AppBinding, ns string, opts *common.PostgresOpts) ([]byte, string, error) {
	var buffer []byte
	_, err := ensureClientCert(opts, apb, opts.DB, api.PostgresClientCert, userName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to ensure client cert %v", err)
	}
	tlsSecret := &core.Secret{}

	err = wait.PollUntilContextTimeout(context.Background(), 300*time.Millisecond, 60*time.Minute, true, func(ctx context.Context) (done bool, err error) {
		sercretName := opts.DB.GetCertSecretName(api.PostgresClientCert) + fmt.Sprintf("-%s", userName)

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
	tlsSecret.ObjectMeta.Annotations = nil
	tlsSecret.ObjectMeta.ManagedFields = nil
	tlsSecretYaml, err := yaml.Marshal(tlsSecret)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal tls secret yaml %v", err)
	}

	buffer = append(buffer, tlsSecretYaml...)
	buffer = append(buffer, []byte("---\n")...)

	return buffer, tlsSecret.Name, nil
}

func generateAuthSecret(userName string, password string, ns string, opts *common.PostgresOpts) ([]byte, string, error) {
	var buffer []byte
	if userName != opts.Username {
		// generate user if not present
		err := generateUser(opts, userName, password)
		if err != nil {
			return nil, "", fmt.Errorf("failed to generate user err:%v", err)
		}
	} else {
		password = opts.Pass
	}
	// generate auth secret
	AuthSecret := core.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       KindSecret,
			APIVersion: ApiversionV1,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-remote-replica-auth", opts.DB.Name),
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

func ensureClientCert(opts *common.PostgresOpts, apb *appApi.AppBinding, postgres *api.Postgres, alias api.PostgresCertificateAlias, username string) (kutil.VerbType, error) {
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
		}, metav1.PatchOptions{})

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
	fmt.Scan(&inp)
	inp = strings.ToLower(inp)
	if inp != "y" {
		return errors.New("aborting commands")
	}
	return nil
}
