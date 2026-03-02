package imageref

import (
	"errors"
	"strings"
)

// ImageRef represents a parsed container image reference.
type ImageRef struct {
	Registry  string // "registry-1.docker.io" for Docker Hub
	Namespace string // "library" for official images
	Name      string // "nginx"
	Tag       string // "latest" if omitted
	Digest    string // "sha256:abc..." if pinned
	IsDigest  bool   // true if reference uses @sha256:...
}

const (
	dockerHubRegistry = "registry-1.docker.io"
	defaultNamespace  = "library"
	defaultTag        = "latest"
)

// FullRepository returns namespace/name (e.g., "library/nginx").
func (r *ImageRef) FullRepository() string {
	return r.Namespace + "/" + r.Name
}

// IsDockerHub returns true if this image is from Docker Hub.
func (r *ImageRef) IsDockerHub() bool {
	return r.Registry == dockerHubRegistry
}

// String returns the canonical string form of the reference.
func (r *ImageRef) String() string {
	var b strings.Builder

	b.WriteString(r.Registry)
	b.WriteByte('/')
	b.WriteString(r.Namespace)
	b.WriteByte('/')
	b.WriteString(r.Name)

	if r.IsDigest {
		b.WriteByte('@')
		b.WriteString(r.Digest)
	} else if r.Tag != "" {
		b.WriteByte(':')
		b.WriteString(r.Tag)
	}

	return b.String()
}

// Parse decomposes a container image reference string into its components.
func Parse(ref string) (*ImageRef, error) {
	if ref == "" {
		return nil, errors.New("empty image reference")
	}

	result := &ImageRef{}

	// Split off digest (@sha256:...) first.
	nameTag := ref
	if before, after, ok := strings.Cut(ref, "@"); ok {
		nameTag = before
		result.Digest = after
		result.IsDigest = true
	}

	// Split off tag (:tag) if no digest.
	if !result.IsDigest {
		if idx := lastTagSeparator(nameTag); idx >= 0 {
			result.Tag = nameTag[idx+1:]
			nameTag = nameTag[:idx]
		}
	}

	// Parse the remaining name into registry/namespace/name.
	parseNameParts(nameTag, result)

	// Normalize docker.io -> registry-1.docker.io.
	if result.Registry == "docker.io" || result.Registry == "index.docker.io" {
		result.Registry = dockerHubRegistry
	}

	// Default tag when not digest-pinned.
	if !result.IsDigest && result.Tag == "" {
		result.Tag = defaultTag
	}

	return result, nil
}

// lastTagSeparator returns the index of the colon separating the tag from
// the image name. It only considers a colon after the last slash to avoid
// confusing a registry port (e.g., registry:5000) with a tag separator.
func lastTagSeparator(s string) int {
	afterSlash := s
	offset := 0
	if idx := strings.LastIndex(s, "/"); idx >= 0 {
		afterSlash = s[idx+1:]
		offset = idx + 1
	}
	if idx := strings.LastIndex(afterSlash, ":"); idx >= 0 {
		return offset + idx
	}
	return -1
}

// looksLikeRegistry returns true if the segment looks like a hostname:
// contains a dot or a colon (port).
func looksLikeRegistry(segment string) bool {
	return strings.ContainsAny(segment, ".:")
}

// parseNameParts splits the name portion into registry, namespace, and name.
func parseNameParts(name string, ref *ImageRef) {
	parts := strings.Split(name, "/")

	switch len(parts) {
	case 1:
		// "nginx" -> Docker Hub official image.
		ref.Registry = dockerHubRegistry
		ref.Namespace = defaultNamespace
		ref.Name = parts[0]

	case 2:
		if looksLikeRegistry(parts[0]) {
			// "ghcr.io/repo" -> custom registry, default namespace.
			ref.Registry = parts[0]
			ref.Namespace = defaultNamespace
			ref.Name = parts[1]
		} else {
			// "myuser/myapp" -> Docker Hub user image.
			ref.Registry = dockerHubRegistry
			ref.Namespace = parts[0]
			ref.Name = parts[1]
		}

	default:
		// 3+ parts: first is registry, second is namespace, rest is name.
		ref.Registry = parts[0]
		ref.Namespace = parts[1]
		ref.Name = strings.Join(parts[2:], "/")
	}
}
