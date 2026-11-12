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

package adopt

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const objectStorageCheckTimeout = 3 * time.Minute

func objectStorageCheckPodName(name string) string { return name + "-adopt-objstore-check" }

// objectStorageCheckScript fails when the server certificate does not verify
// against the given CA, or does not name the host. Any HTTP answer is a pass:
// only the TLS handshake is checked.
const objectStorageCheckScript = `out=$(curl -sS --max-time 30 --cacert /ca/ca.crt -o /dev/null -w '%{http_code}' "https://$HOST:$PORT/" 2>&1) || { echo "$out"; exit 1; }
echo "TLS to $HOST:$PORT verified with the given CA (HTTP $out)"`

// objectStorageEndpoint is the host and port the Milvus will use.
func objectStorageEndpoint(cfg map[string]any, t *objectStorageTLS) (host, port string) {
	data := t.storageData(objectStorageData(cfg))
	host, port = data["address"], data["port"]
	if port == "" {
		port = "443"
	}
	return host, port
}

// objectStorageCheckPod is a short-lived Pod that runs the Milvus image, which has
// curl, with the CA mounted as the Milvus pods will mount it.
func objectStorageCheckPod(ns, name, release, image string, pull []core.LocalObjectReference, t *objectStorageTLS, host, port string) *core.Pod {
	vol := core.Volume{Name: "ca"}
	if t.Kind == "ConfigMap" {
		vol.ConfigMap = &core.ConfigMapVolumeSource{
			LocalObjectReference: core.LocalObjectReference{Name: t.Name},
			Items:                []core.KeyToPath{{Key: t.Key, Path: "ca.crt"}},
		}
	} else {
		vol.Secret = &core.SecretVolumeSource{SecretName: t.Name, Items: []core.KeyToPath{{Key: t.Key, Path: "ca.crt"}}}
	}
	no, yes := false, true
	uid := int64(1000)
	return &core.Pod{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{Name: objectStorageCheckPodName(name), Namespace: ns, Labels: map[string]string{LabelRelease: release}},
		Spec: core.PodSpec{
			RestartPolicy:    core.RestartPolicyNever,
			ImagePullSecrets: pull,
			SecurityContext:  &core.PodSecurityContext{RunAsNonRoot: &yes, RunAsUser: &uid, SeccompProfile: &core.SeccompProfile{Type: core.SeccompProfileTypeRuntimeDefault}},
			Volumes:          []core.Volume{vol},
			Containers: []core.Container{{
				Name:    "check",
				Image:   image,
				Command: []string{"/bin/sh", "-c", objectStorageCheckScript},
				Env:     []core.EnvVar{{Name: "HOST", Value: host}, {Name: "PORT", Value: port}},
				SecurityContext: &core.SecurityContext{
					AllowPrivilegeEscalation: &no,
					Capabilities:             &core.Capabilities{Drop: []core.Capability{"ALL"}},
				},
				VolumeMounts: []core.VolumeMount{{Name: "ca", MountPath: "/ca", ReadOnly: true}},
			}},
		},
	}
}

// checkObjectStorageTLS runs the check Pod before any object is created, so a
// wrong CA, a host the certificate does not name, or an endpoint that is not
// reachable from the cluster stops the adoption while nothing has moved.
func (c *cluster) checkObjectStorageTLS(ctx context.Context, p *adoptPlan, name, release string) error {
	host, port := objectStorageEndpoint(p.Source.Config, p.ObjStore)
	var pull []core.LocalObjectReference
	for _, u := range p.Source.Units {
		if u.Role == roleProxy {
			pull = p.Source.Templates[u.key()].Spec.ImagePullSecrets
		}
	}
	pod := objectStorageCheckPod(c.ns, name, release, p.Source.Image, pull, p.ObjStore, host, port)
	pods := c.kube.CoreV1().Pods(c.ns)
	remove := func() {
		_ = pods.Delete(context.WithoutCancel(ctx), pod.Name, metav1.DeleteOptions{})
	}
	if err := pods.Delete(ctx, pod.Name, metav1.DeleteOptions{}); err != nil && !kerr.IsNotFound(err) {
		return err
	}
	c.logf("checking TLS to the object storage %s:%s from a Pod, before anything moves", host, port)
	if _, err := pods.Create(ctx, pod, metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("failed to create the check Pod: %w", err)
	}
	defer remove()
	var phase core.PodPhase
	err := c.waitFor(ctx, objectStorageCheckTimeout, "object storage TLS check", func() (bool, string, error) {
		cur, err := pods.Get(ctx, pod.Name, metav1.GetOptions{})
		if err != nil {
			return false, "", err
		}
		phase = cur.Status.Phase
		return phase == core.PodSucceeded || phase == core.PodFailed, string(phase), nil
	})
	if err != nil {
		return fmt.Errorf("the object storage check Pod did not finish (does the cluster pull %s?): %w", p.Source.Image, err)
	}
	logs := ""
	if rc, err := pods.GetLogs(pod.Name, &core.PodLogOptions{}).Stream(ctx); err == nil {
		b, _ := io.ReadAll(io.LimitReader(rc, 4096))
		_ = rc.Close()
		logs = strings.TrimSpace(string(b))
	}
	if phase != core.PodSucceeded {
		return fmt.Errorf("TLS to the object storage %s:%s does not verify with the CA you gave: %s", host, port, logs)
	}
	c.logf("%s", logs)
	return nil
}
