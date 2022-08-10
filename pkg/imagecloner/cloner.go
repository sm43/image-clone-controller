package imagecloner

import (
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func (ic *ImageCloner) IsBackupImage(image string) bool {
	return strings.Contains(image, fmt.Sprint(ic.registry+"/"+ic.repository))
}

func (ic *ImageCloner) Clone(image string) (string, error) {
	sourceTag, err := name.NewTag(image)
	if err != nil {
		return "", err
	}
	sourceImg, err := remote.Image(sourceTag)
	if err != nil {
		return "", err
	}

	targetImage := ic.getTargetImage(sourceTag)
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

func (ic *ImageCloner) getTargetImage(source name.Tag) string {
	parts := strings.Split(source.RepositoryStr(), "/")
	return ic.registry + "/" + ic.repository + "/" + parts[len(parts)-1] + ":" + source.TagStr()
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
