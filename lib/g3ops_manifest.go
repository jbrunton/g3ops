package lib

// G3opsManifest - manifest with build information for github releases
type G3opsManifest struct {
	Version      string
	Environments map[string]struct {
		Host    string
		Version string
	}
}
