package store

// CheckStatus represents the result state of an image update check.
type CheckStatus string

const (
	StatusUpToDate        CheckStatus = "up-to-date"
	StatusUpdateAvailable CheckStatus = "update-available"
	StatusCheckFailed     CheckStatus = "check-failed"
	StatusUnknown         CheckStatus = "unknown"
	StatusChecking        CheckStatus = "checking"
)

// ImageCheck records the result of comparing a container's local image
// digest against the remote registry digest.
type ImageCheck struct {
	ID            int64       `json:"id"`
	ContainerName string      `json:"containerName"`
	ContainerID   string      `json:"containerId"`
	ImageRef      string      `json:"imageRef"`
	LocalDigest   string      `json:"localDigest"`
	RemoteDigest  string      `json:"remoteDigest"`
	Status        CheckStatus `json:"status"`
	CheckedAt     string      `json:"checkedAt"`
	Registry      string      `json:"registry"`
}

// Preference stores a user-configurable key/value setting.
type Preference struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
