package generator

import "errors"

var (
	// ErrUnsupportedDatabase indicates the database option is not supported.
	ErrUnsupportedDatabase = errors.New("unsupported database option")
	// ErrUnsupportedAuth indicates the auth strategy is not supported.
	ErrUnsupportedAuth = errors.New("unsupported auth strategy option")
	// ErrMissingProjectName indicates the project name is required.
	ErrMissingProjectName = errors.New("project name is required")
	// ErrNoAppSelected indicates at least one application must be selected.
	ErrNoAppSelected = errors.New("at least one application (api or worker) must be generated")
)
