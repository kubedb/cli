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
	"fmt"
	"log"
	"os"
	"time"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	"kubedb.dev/cli/pkg/common"

	cm_api "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/spf13/cobra"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	cm_util "kmodules.xyz/cert-manager-util/certmanager/v1"
	kutil "kmodules.xyz/client-go"
	kmapi "kmodules.xyz/client-go/api/v1"
	core_util "kmodules.xyz/client-go/core/v1"
	exec_util "kmodules.xyz/client-go/tools/exec"
	appApi "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	"sigs.k8s.io/yaml"
)

var (
	desLong = "generate appbinding , secrets for remote replica"
	example = "kubectl dba remote-config mysql -n <ns> -u <user_name> -p$<password> -d<dns_name>  <db_name>\n " +
		"kubectl dba remote-config mysql -n <ns> -u <user_name> -p$<password> -d<dns_name>  <db_name> \n"
)

const (
	AppcatApiVersion = "appcatalog.appscode.com/v1alpha1"
	AppcatKind       = "AppBinding"
	ApiversionV1     = "v1"
	KindSecret       = "Secret"
)

func MysqlAPP(f cmdutil.Factory) *cobra.Command {
	var userName, password, dns, ns string
	var yes bool
	cmd := cobra.Command{
		Use:     "mysql",
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

			buffer, err := generateMySQLConfig(f, userName, password, dns, ns, args[0])
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

func generateMySQLConfig(f cmdutil.Factory, userName string, password string, dns string, ns string, dbname string) ([]byte, error) {
	var buffer []byte

	opts, err := common.NewMySQLOpts(f, dbname, ns)
	if err != nil {
		return nil, fmt.Errorf("failed to get db %s, err:%v", dbname, err)
	}

	apb, err := opts.AppcatClient.AppcatalogV1alpha1().AppBindings(ns).Get(context.TODO(), dbname, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get appbinding %v", err)
	}

	authBuff, authSecretName, err := generateMySQLAuthSecret(userName, password, ns, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth secret ,%v", err)
	}
	buffer = append(buffer, authBuff...)

	// generate secret
	if apb.Spec.TLSSecret != nil {
		tlsBuff, tlsSecretName, err := generateMySQLTlsSecret(userName, apb, ns, opts)
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

func generateMySQLTlsSecret(userName string, apb *appApi.AppBinding, ns string, opts *common.MySQLOpts) ([]byte, string, error) {
	var buffer []byte
	_, err := ensureMySQLClientCert(opts, apb, opts.DB, api.MySQLClientCert, userName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to ensure client cert %v", err)
	}
	tlsSecret := &core.Secret{}

	err = wait.PollImmediate(300*time.Millisecond, 60*time.Minute, func() (done bool, err error) {
		sercretName := opts.DB.GetCertSecretName(api.MySQLClientCert) + fmt.Sprintf("-%s", userName)

		tlsSecret, err = opts.Client.CoreV1().Secrets(ns).Get(context.TODO(), sercretName, metav1.GetOptions{})
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
	tlsSecret.APIVersion = ApiversionV1
	tlsSecret.Kind = KindSecret
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

func generateMySQLAuthSecret(userName string, password string, ns string, opts *common.MySQLOpts) ([]byte, string, error) {
	var buffer []byte
	if userName != opts.Username {
		// generate user if not present
		err := generateMySQLUser(opts, userName, password)
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

func generateMySQLUser(opts *common.MySQLOpts, name string, password string) error {
	label := opts.DB.OffshootLabels()
	if *opts.DB.Spec.Replicas > 1 {
		label["kubedb.com/role"] = "primary"
	}

	pods, err := opts.Client.CoreV1().Pods(opts.DB.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.Set.String(label),
	})
	if err != nil || len(pods.Items) == 0 {
		return err
	}
	query := fmt.Sprintf("export MYSQL_PWD='%s' && mysql -uroot -e \"create user if not exists %s; alter user %s identified by '%s';"+
		"GRANT REPLICATION SLAVE,  CLONE_ADMIN, BACKUP_ADMIN ON *.* TO '%s'@'%%' WITH GRANT OPTION; \"", opts.Pass,
		name, name, password, name)
	command := exec_util.Command("bash", "-c", query)
	container := exec_util.Container("mysql")
	options := []func(options *exec_util.Options){
		command,
		container,
	}

	_, err = exec_util.ExecIntoPod(opts.Config, &pods.Items[0], options...)
	if err != nil {
		return err
	}
	return nil
}

func ensureMySQLClientCert(opts *common.MySQLOpts, apb *appApi.AppBinding, mysql *api.MySQL, alias api.MySQLCertificateAlias, username string) (kutil.VerbType, error) {
	var duration, renewBefore *metav1.Duration
	var subject *cm_api.X509Subject
	var dnsNames, ipAddresses, uriSANs, emailSANs []string
	if _, cert := kmapi.GetCertificate(mysql.Spec.TLS.Certificates, string(alias)); cert != nil {
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
			Name:      mysql.CertificateName(alias) + fmt.Sprintf("-%s", username),
			Namespace: mysql.GetNamespace(),
		},
		func(in *cm_api.Certificate) *cm_api.Certificate {
			in.Labels = mysql.OffshootLabels()
			core_util.EnsureOwnerReference(in, ref)

			in.Spec.CommonName = username
			in.Spec.Subject = subject
			in.Spec.Duration = duration
			in.Spec.RenewBefore = renewBefore
			in.Spec.DNSNames = sets.NewString(dnsNames...).List()
			in.Spec.IPAddresses = sets.NewString(ipAddresses...).List()
			in.Spec.URIs = sets.NewString(uriSANs...).List()
			in.Spec.EmailAddresses = sets.NewString(emailSANs...).List()
			in.Spec.SecretName = mysql.GetCertSecretName(alias) + fmt.Sprintf("-%s", username)
			in.Spec.IssuerRef = GetIssuerObjectRef(mysql.Spec.TLS, string(alias))
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
