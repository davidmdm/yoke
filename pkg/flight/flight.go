package flight

import (
	"os"
	"path/filepath"
)

// Release is convenience for fetching the release name within the context of an executing flight.
// This will generally be the name of release passed to "yoke takeoff"
func Release() string {
	_, release := filepath.Split(os.Args[0])
	return release
}

// Namespace is a convenience function for fetching the namespace within the context of an executing flight.
// This will generally be the -namespace flag passed to "yoke takeoff"
func Namespace() string {
	return os.Getenv("NAMESPACE")
}
