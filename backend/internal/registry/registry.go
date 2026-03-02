package registry

import (
	"context"

	"github.com/HerbHall/DockPulse/backend/internal/imageref"
)

// Registry can fetch the remote digest for a given image reference.
type Registry interface {
	GetDigest(ctx context.Context, ref *imageref.ImageRef) (string, error)
}
