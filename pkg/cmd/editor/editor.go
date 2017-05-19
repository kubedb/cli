package editor

import (
	"os"

	//"github.com/appscode/go-term"
	"os/exec"
)


func GetEditor() string {
	//preferred
	var editor string
	editor = os.Getenv("EDITOR")
	if len(editor) != 0 {
		editor = os.ExpandEnv(editor)
		return editor
	}

	editor = os.Getenv("KUBEDB_EDITOR")
	if len(editor) != 0 {
		editor = os.ExpandEnv(editor)
		return editor
	}

	//errMessage := "Unable to launch an interactive text editor. Set the EDITOR or KUBEDB_EDITOR environment variable to an appropriate editor."
	//term.Fatalln(errMessage)
	return ""
}

func OpenEditor(editorName, tempFilePath string) error {
	cmd := exec.Command(editorName, tempFilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
