package docker

import (
	"net/http"

	docker "github.com/heroku/docker-registry-client/registry"
)

const (
	registryUrl = "https://registry-1.docker.io/"
)

func CheckDockerImageVersion(repository, reference string) error {
	registry := &docker.Registry{
		URL: registryUrl,
		Client: &http.Client{
			Transport: docker.WrapTransport(http.DefaultTransport, registryUrl, "", ""),
		},
		Logf: docker.Quiet,
	}

	_, err := registry.Manifest(repository, reference)
	return err
}
