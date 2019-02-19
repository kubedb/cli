package encoder

import (
	"bytes"
	"encoding/json"
	"io"
	"k8s.io/apimachinery/pkg/runtime"
)

func Encode(obj runtime.Object) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := encode(obj, buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encode(obj runtime.Object, w io.Writer) error {
	return json.NewEncoder(w).Encode(obj)
}
