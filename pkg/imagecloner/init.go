package imagecloner

import (
	"fmt"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
)

type ImageCloner interface {
	IsBackupImage(string) bool
	Clone(string) (string, error)
}

type imageCloner struct {
	auth       authn.Authenticator
	registry   string
	repository string
}

func NewCloner() (ImageCloner, error) {
	registry := os.Getenv("REGISTRY_NAME")
	repository := os.Getenv("REPOSITORY_NAME")
	username := os.Getenv("REGISTRY_USERNAME")
	password := os.Getenv("REGISTRY_PASSWORD")

	if registry == "" || username == "" || password == "" {
		return nil, fmt.Errorf("failed to get registry configurations")
	}
	if repository == "" {
		repository = username
	}
	auth := authn.AuthConfig{
		Username: username,
		Password: password,
	}
	return &imageCloner{
		auth:       authn.FromConfig(auth),
		registry:   registry,
		repository: repository,
	}, nil
}
