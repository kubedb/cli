package editor

import (
	"bytes"
	"fmt"
	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/go/io"
	"k8s.io/kubernetes/pkg/kubectl/resource"
)

func WriteTempFile(info *resource.Info, buf *bytes.Buffer) (bool, string) {
	path := rand.WithUniqSuffix(fmt.Sprintf("/tmp/%s-%s-edit", info.Namespace, info.Name))
	data := buf.Bytes()
	if !io.WriteString(path, string(data)) {
		return false, ""
	}

	return true, path
}
