package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/appscode/go/runtime"
	"github.com/kubedb/cli/pkg/cmds"
	"github.com/spf13/cobra/doc"
)

const (
	version = "0.9.0-rc.1"
)

var (
	tplFrontMatter = template.Must(template.New("index").Parse(`---
title: Reference
description: KubeDB CLI Reference
menu:
  docs_{{ .Version }}:
    identifier: reference
    name: Reference
    weight: 1000
menu_name: docs_{{ .Version }}
---
`))

	_ = template.Must(tplFrontMatter.New("cmd").Parse(`---
title: {{ .Name }}
menu:
  docs_{{ .Version }}:
    identifier: {{ .ID }}
    name: {{ .Name }}
    parent: reference
{{- if .RootCmd }}
    weight: 0
{{ end }}
menu_name: docs_{{ .Version }}
section_menu_id: reference
{{- if .RootCmd }}
aliases:
  - /docs/{{ .Version }}/reference/
{{ end }}
---
`))
)

// ref: https://github.com/spf13/cobra/blob/master/doc/md_docs.md
func main() {
	rootCmd := cmds.NewKubedbCommand(os.Stdin, os.Stdout, os.Stderr)
	dir := runtime.GOPath() + "/src/github.com/kubedb/cli/docs/reference"
	fmt.Printf("Generating cli markdown tree in: %v\n", dir)
	err := os.RemoveAll(dir)
	if err != nil {
		log.Fatalln(err)
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatalln(err)
	}

	filePrepender := func(filename string) string {
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		data := struct {
			ID      string
			Name    string
			Version string
			RootCmd bool
		}{
			strings.Replace(base, "_", "-", -1),
			strings.Title(strings.Replace(base, "_", " ", -1)),
			version,
			!strings.ContainsRune(base, '_'),
		}
		var buf bytes.Buffer
		if err := tplFrontMatter.ExecuteTemplate(&buf, "cmd", data); err != nil {
			log.Fatalln(err)
		}
		return buf.String()
	}

	linkHandler := func(name string) string {
		return "/docs/reference/" + name
	}
	doc.GenMarkdownTreeCustom(rootCmd, dir, filePrepender, linkHandler)

	index := filepath.Join(dir, "_index.md")
	f, err := os.OpenFile(index, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	err = tplFrontMatter.ExecuteTemplate(f, "index", struct{ Version string }{version})
	if err != nil {
		log.Fatalln(err)
	}
	if err := f.Close(); err != nil {
		log.Fatalln(err)
	}
}
