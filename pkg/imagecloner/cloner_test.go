package imagecloner

import (
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	"gotest.tools/v3/assert"
)

func TestIsBackupImage(t *testing.T) {
	tests := []struct {
		name  string
		image string
		want  bool
	}{
		{
			name:  "valid",
			image: "ttl.sh/sm43/abc",
			want:  true,
		},
		{
			name:  "invalid",
			image: "nginx:alpine",
			want:  false,
		},
		{
			name:  "invalid",
			image: "ghcr.io/abc/nginx:alpine",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ic := ImageCloner{
				registry:   "ttl.sh",
				repository: "sm43",
			}
			got := ic.IsBackupImage(tt.image)
			assert.Assert(t, got == tt.want)
		})
	}
}

func TestGetTargetImage(t *testing.T) {
	tests := []struct {
		name  string
		image string
		want  string
	}{
		{
			name:  "image-1",
			image: "nginx:alpine",
			want:  "ttl.sh/sm43/index.docker.io_library_nginx:alpine",
		},
		{
			name:  "image-2",
			image: "r.j3ss.co/clisp:1.0",
			want:  "ttl.sh/sm43/r.j3ss.co_clisp:1.0",
		},
		{
			name:  "image-3",
			image: "gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/controller:v0.37.4@sha256:fdd699f84f843a45b1e061a4eeb02f04c8234158c908af21dd2bf1b3c6c1f862",
			want:  "ttl.sh/sm43/gcr.io_tekton-releases_github.com_tektoncd_pipeline_cmd_controller@sha256:fdd699f84f843a45b1e061a4eeb02f04c8234158c908af21dd2bf1b3c6c1f862",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ic := ImageCloner{
				registry:   "ttl.sh",
				repository: "sm43",
			}
			tag, err := name.ParseReference(tt.image)
			assert.NilError(t, err)

			got := ic.getTargetImage(tag)
			assert.Assert(t, got == tt.want)
		})
	}
}
