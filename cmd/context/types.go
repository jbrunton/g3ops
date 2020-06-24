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
	Ci struct {
		Defaults struct {
			Build struct {
				Env     map[string]string
				Command string
				Args    string
			}
		}
	}
}
