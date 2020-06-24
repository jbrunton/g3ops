package context

// G3opsContext - type of current g3ops context
type G3opsContext struct {
	Name         string
	Environments map[string]struct {
		Manifest string
	}
	Services map[string]struct {
		Manifest string
	}
}
