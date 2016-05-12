package global

import (
	"fmt"
)

func Version() string {
	return fmt.Sprintf("tgen v%d.%d.%d", versionMajor, versionMinor, versionPatch)
}

const (
	versionMajor = 0
	versionMinor = 0
	versionPatch = 1
)
