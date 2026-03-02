package imageref

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   bool
		registry  string
		namespace string
		imgName   string
		tag       string
		digest    string
		isDigest  bool
	}{
		{
			name:      "official image no tag",
			input:     "nginx",
			registry:  "registry-1.docker.io",
			namespace: "library",
			imgName:   "nginx",
			tag:       "latest",
		},
		{
			name:      "official image with tag",
			input:     "nginx:1.25",
			registry:  "registry-1.docker.io",
			namespace: "library",
			imgName:   "nginx",
			tag:       "1.25",
		},
		{
			name:      "official image with latest tag",
			input:     "nginx:latest",
			registry:  "registry-1.docker.io",
			namespace: "library",
			imgName:   "nginx",
			tag:       "latest",
		},
		{
			name:      "docker hub user image no tag",
			input:     "myuser/myapp",
			registry:  "registry-1.docker.io",
			namespace: "myuser",
			imgName:   "myapp",
			tag:       "latest",
		},
		{
			name:      "docker hub user image with tag",
			input:     "myuser/myapp:v2",
			registry:  "registry-1.docker.io",
			namespace: "myuser",
			imgName:   "myapp",
			tag:       "v2",
		},
		{
			name:      "ghcr with tag",
			input:     "ghcr.io/user/repo:tag",
			registry:  "ghcr.io",
			namespace: "user",
			imgName:   "repo",
			tag:       "tag",
		},
		{
			name:      "ghcr no tag",
			input:     "ghcr.io/user/repo",
			registry:  "ghcr.io",
			namespace: "user",
			imgName:   "repo",
			tag:       "latest",
		},
		{
			name:      "multi-level namespace",
			input:     "mcr.microsoft.com/dotnet/sdk:8.0",
			registry:  "mcr.microsoft.com",
			namespace: "dotnet",
			imgName:   "sdk",
			tag:       "8.0",
		},
		{
			name:      "custom registry with port",
			input:     "registry.example.com:5000/myapp:v1",
			registry:  "registry.example.com:5000",
			namespace: "library",
			imgName:   "myapp",
			tag:       "v1",
		},
		{
			name:      "digest-pinned official image",
			input:     "nginx@sha256:abcdef1234567890",
			registry:  "registry-1.docker.io",
			namespace: "library",
			imgName:   "nginx",
			tag:       "",
			digest:    "sha256:abcdef1234567890",
			isDigest:  true,
		},
		{
			name:      "docker.io explicit official image",
			input:     "docker.io/library/nginx:latest",
			registry:  "registry-1.docker.io",
			namespace: "library",
			imgName:   "nginx",
			tag:       "latest",
		},
		{
			name:      "docker.io explicit user image",
			input:     "docker.io/myuser/myapp:v1",
			registry:  "registry-1.docker.io",
			namespace: "myuser",
			imgName:   "myapp",
			tag:       "v1",
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Registry != tt.registry {
				t.Errorf("Registry = %q, want %q", got.Registry, tt.registry)
			}
			if got.Namespace != tt.namespace {
				t.Errorf("Namespace = %q, want %q", got.Namespace, tt.namespace)
			}
			if got.Name != tt.imgName {
				t.Errorf("Name = %q, want %q", got.Name, tt.imgName)
			}
			if got.Tag != tt.tag {
				t.Errorf("Tag = %q, want %q", got.Tag, tt.tag)
			}
			if got.Digest != tt.digest {
				t.Errorf("Digest = %q, want %q", got.Digest, tt.digest)
			}
			if got.IsDigest != tt.isDigest {
				t.Errorf("IsDigest = %v, want %v", got.IsDigest, tt.isDigest)
			}
		})
	}
}

func TestImageRef_FullRepository(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		imgName   string
		want      string
	}{
		{
			name:      "official image",
			namespace: "library",
			imgName:   "nginx",
			want:      "library/nginx",
		},
		{
			name:      "user image",
			namespace: "myuser",
			imgName:   "myapp",
			want:      "myuser/myapp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref := &ImageRef{Namespace: tt.namespace, Name: tt.imgName}
			if got := ref.FullRepository(); got != tt.want {
				t.Errorf("FullRepository() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestImageRef_IsDockerHub(t *testing.T) {
	tests := []struct {
		name     string
		registry string
		want     bool
	}{
		{
			name:     "docker hub",
			registry: "registry-1.docker.io",
			want:     true,
		},
		{
			name:     "ghcr",
			registry: "ghcr.io",
			want:     false,
		},
		{
			name:     "custom registry",
			registry: "registry.example.com:5000",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref := &ImageRef{Registry: tt.registry}
			if got := ref.IsDockerHub(); got != tt.want {
				t.Errorf("IsDockerHub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageRef_String(t *testing.T) {
	tests := []struct {
		name string
		ref  ImageRef
		want string
	}{
		{
			name: "official image with tag",
			ref: ImageRef{
				Registry:  "registry-1.docker.io",
				Namespace: "library",
				Name:      "nginx",
				Tag:       "latest",
			},
			want: "registry-1.docker.io/library/nginx:latest",
		},
		{
			name: "user image with tag",
			ref: ImageRef{
				Registry:  "registry-1.docker.io",
				Namespace: "myuser",
				Name:      "myapp",
				Tag:       "v2",
			},
			want: "registry-1.docker.io/myuser/myapp:v2",
		},
		{
			name: "digest-pinned image",
			ref: ImageRef{
				Registry:  "registry-1.docker.io",
				Namespace: "library",
				Name:      "nginx",
				Digest:    "sha256:abcdef",
				IsDigest:  true,
			},
			want: "registry-1.docker.io/library/nginx@sha256:abcdef",
		},
		{
			name: "custom registry",
			ref: ImageRef{
				Registry:  "ghcr.io",
				Namespace: "user",
				Name:      "repo",
				Tag:       "v1",
			},
			want: "ghcr.io/user/repo:v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ref.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
