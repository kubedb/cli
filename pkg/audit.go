package pkg

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"

	tapi "github.com/k8sdb/apimachinery/api"
	"github.com/k8sdb/cli/pkg/kube"
	"github.com/k8sdb/cli/pkg/util"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/kubernetes/pkg/client/unversioned/remotecommand"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

func NewCmdAuditReport(out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit_report",
		Short: "Export audit report",
		Run: func(cmd *cobra.Command, args []string) {
			f := kube.NewKubeFactory(cmd)
			cmdutil.CheckErr(exportReport(f, cmd, out, cmdErr, args))
		},
	}
	util.AddAuditReportFlags(cmd)
	return cmd
}

type auditReport struct {
	DbName string
}

const (
	valid_resources_for_report = `Valid resource types include:

    * elastic
    * postgres
    `
)

func exportReport(f cmdutil.Factory, cmd *cobra.Command, out, errOut io.Writer, args []string) error {
	namespace, _ := util.GetNamespace(cmd)
	_ = namespace

	if len(args) == 0 {
		fmt.Fprint(errOut, "You must specify the type of resource to get. ", valid_resources_for_report)
		usageString := "Required resource not specified."
		return cmdutil.UsageError(cmd, usageString)
	}

	if len(strings.Split(args[0], ",")) > 1 {
		return errors.New("audit doesn't support multiple resource")
	}

	resource := args[0]
	items := strings.Split(resource, "/")
	kubedbType, err := util.GetResourceType(items[0])
	if err != nil {
		return err
	}

	var kubedbName string
	if len(items) > 1 {
		if len(items) > 2 {
			return errors.New("audit doesn't support multiple resource")
		}
		kubedbName = items[1]
	} else {
		if len(args) > 2 {
			return errors.New("audit doesn't support multiple resource")
		}
		kubedbName = args[1]
	}
	_ = kubedbName

	switch kubedbType {
	case tapi.ResourceTypeSnapshot, tapi.ResourceTypeDormantDatabase:
		return fmt.Errorf(`resource type "%v" doesn't support audit operation`, items[0])
	}

	dbname := cmdutil.GetFlagString(cmd, "index")
	_ = dbname

	clientset, err := f.ClientSet()
	if err != nil {
		return err
	}

	operatorNamespace := cmdutil.GetFlagString(cmd, "operator-namespace")
	operatorPodList, err := clientset.Core().Pods(operatorNamespace).List(
		metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(
				operatorLabel,
			).String(),
		},
	)
	if err != nil {
		return err
	}

	if len(operatorPodList.Items) == 0 {
		return errors.New("Operator pod not found")
	}

	restClient, err := f.RESTClient()
	if err != nil {
		return err
	}

	config, err := f.ClientConfig()
	if err != nil {
		return err
	}

	tunnel := NewTunnel(restClient, config, operatorNamespace, operatorPodList.Items[0].Name, 8443)
	if err := tunnel.forwardPort(); err != nil {
		return err
	}

	fmt.Println(tunnel.Local)
	time.Sleep(time.Hour)
	return nil
}

type Tunnel struct {
	Local     int
	Remote    int
	Namespace string
	PodName   string
	Out       io.Writer
	stopChan  chan struct{}
	readyChan chan struct{}
	config    *rest.Config
	client    rest.Interface
}

func NewTunnel(client rest.Interface, config *rest.Config, namespace, podName string, remote int) *Tunnel {
	return &Tunnel{
		config:    config,
		client:    client,
		Namespace: namespace,
		PodName:   podName,
		Remote:    remote,
		stopChan:  make(chan struct{}, 1),
		readyChan: make(chan struct{}, 1),
		Out:       ioutil.Discard,
	}
}

func (t *Tunnel) forwardPort() error {

	u := t.client.Post().
		Resource("pods").
		Namespace(t.Namespace).
		Name(t.PodName).
		SubResource("portforward").URL()

	fmt.Println(u)

	dialer, err := remotecommand.NewExecutor(t.config, "GET", u)
	if err != nil {
		return err
	}

	local, err := getAvailablePort()
	if err != nil {
		return fmt.Errorf("could not find an available port: %s", err)
	}
	t.Local = local

	ports := []string{fmt.Sprintf("%d:%d", t.Local, t.Remote)}

	pf, err := portforward.New(dialer, ports, t.stopChan, t.readyChan, t.Out, t.Out)
	if err != nil {
		return err
	}

	errChan := make(chan error)
	go func() {
		errChan <- pf.ForwardPorts()
	}()

	select {
	case err = <-errChan:
		return fmt.Errorf("forwarding ports: %v", err)
	case <-pf.Ready:
		return nil
	}
}

func getAvailablePort() (int, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer l.Close()

	_, p, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		return 0, err
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		return 0, err
	}
	return port, err
}
