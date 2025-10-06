package generator

// Options captures CLI inputs for project generation.
type Options struct {
	ProjectName    string
	OutputDir      string
	Modules        []string
	Database       string
	Auth           string
	IncludeAPI     bool
	IncludeWorker  bool
	CacheNamespace string
	CacheGroup     string
	CacheTTL       int
}

// HasModule reports whether the requested modules include the provided key.
func (o Options) HasModule(name string) bool {
	for _, m := range o.Modules {
		if m == name {
			return true
		}
	}
	return false
}

// Validate ensures the options are consistent.
func (o Options) Validate() error {
	switch o.Database {
	case "mysql", "postgres", "mongodb":
	default:
		return ErrUnsupportedDatabase
	}

	switch o.Auth {
	case "jwt", "paseto":
	default:
		return ErrUnsupportedAuth
	}

	if o.ProjectName == "" {
		return ErrMissingProjectName
	}

	if !o.IncludeAPI && !o.IncludeWorker {
		return ErrNoAppSelected
	}

	return nil
}
