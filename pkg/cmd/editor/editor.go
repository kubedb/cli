package editor

import (
	"os"

	"k8s.io/kubernetes/pkg/kubectl/cmd/util/editor"
)

const defaultEditor = "nano"

func NewDefaultEditor() editor.Editor {
	var editorName string
	editorName = os.Getenv("EDITOR")
	if len(editorName) != 0 {
		editorName = os.ExpandEnv(editorName)
		return editor.Editor{
			Args:  []string{editorName},
			Shell: false,
		}
	}

	editorName = os.Getenv("KUBEDB_EDITOR")
	if len(editorName) != 0 {
		editorName = os.ExpandEnv(editorName)
		return editor.Editor{
			Args:  []string{editorName},
			Shell: false,
		}
	}

	return editor.Editor{
		Args:  []string{defaultEditor},
		Shell: false,
	}
}
