package imagecloner

import (
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func (ic *imageCloner) IsBackupImage(image string) bool {
	return strings.Contains(image, fmt.Sprint(ic.registry+"/"+ic.repository))
}

func (ic *imageCloner) Clone(image string) (string, error) {
	sourceRef, err := name.ParseReference(image)
	if err != nil {
		return "", err
	}

	sourceImg, err := remote.Image(sourceRef)
	if err != nil {
		return "", err
	}

	targetImage := ic.getTargetImage(sourceRef)
	targetTag, err := name.NewTag(targetImage)
	if err != nil {
		return "", err
	}

	opt := remote.WithAuth(ic.auth)
	if !isImagePresent(targetTag.RepositoryStr(), targetTag.TagStr(), opt) {
		if err := remote.Write(targetTag, sourceImg, opt); err != nil {
			return "", err
		}
	}
	return targetImage, nil
}

func (ic *imageCloner) getTargetImage(source name.Reference) string {
	targetImage := ic.registry + "/" + ic.repository
	targetImage += "/" + strings.Replace(source.Context().RegistryStr(), ":", "_", 1)
	targetImage += "_" + strings.ReplaceAll(source.Context().RepositoryStr(), "/", "_")
	if tag, ok := source.(name.Tag); ok {
		targetImage += ":" + tag.TagStr()
	}
	if digest, ok := source.(name.Digest); ok {
		targetImage += "@" + digest.DigestStr()
	}
	return targetImage
}

func isImagePresent(repository, tag string, opt remote.Option) bool {
	repo, _ := name.NewRepository(repository)
	list, _ := remote.List(repo, opt)
	for _, t := range list {
		if t == tag {
			return true
		}
	}
	return false
}
