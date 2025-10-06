package generator

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateFullProject(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()

	opts := Options{
		ProjectName:    "sample-app",
		OutputDir:      filepath.Join(tmp, "sample-app"),
		Modules:        []string{"routes", "controllers", "services", "repositories", "middleware", "helpers"},
		Database:       "mysql",
		Auth:           "jwt",
		IncludeAPI:     true,
		IncludeWorker:  true,
		CacheNamespace: "core",
		CacheGroup:     "v1",
		CacheTTL:       600,
	}

	if err := New(opts).Generate(context.Background()); err != nil {
		t.Fatalf("generate project: %v", err)
	}

	root := opts.OutputDir

	assertFileContains(t, filepath.Join(root, "config/env/.env.local"), "REDIS_NAMESPACE=core")
	assertFileContains(t, filepath.Join(root, "config/env/.env.local"), "REDIS_TTL=600")
	assertFileContains(t, filepath.Join(root, "packages/shared/database/index.ts"), "driver: 'mysql'")
	assertFileContains(t, filepath.Join(root, "packages/shared/package.json"), "\"mysql2\"")
	assertFileExists(t, filepath.Join(root, "packages/shared/auth/jwt.ts"))
	assertFileNotExists(t, filepath.Join(root, "packages/shared/auth/paseto.ts"))
	assertFileExists(t, filepath.Join(root, "apps/api/src/routes/index.ts"))
	assertFileExists(t, filepath.Join(root, "apps/worker/src/index.ts"))
}

func TestGenerateWorkerOnlyWithMongoPaseto(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()

	opts := Options{
		ProjectName:    "worker-only",
		OutputDir:      filepath.Join(tmp, "worker-only"),
		Modules:        []string{},
		Database:       "mongodb",
		Auth:           "paseto",
		IncludeAPI:     false,
		IncludeWorker:  true,
		CacheNamespace: "ns",
		CacheGroup:     "group",
		CacheTTL:       120,
	}

	if err := New(opts).Generate(context.Background()); err != nil {
		t.Fatalf("generate project: %v", err)
	}

	root := opts.OutputDir

	assertFileContains(t, filepath.Join(root, "packages/shared/database/index.ts"), "MongoClient")
	assertFileContains(t, filepath.Join(root, "packages/shared/package.json"), "\"mongodb\"")
	assertFileExists(t, filepath.Join(root, "packages/shared/auth/paseto.ts"))
	assertFileNotExists(t, filepath.Join(root, "packages/shared/auth/jwt.ts"))

	if _, err := os.Stat(filepath.Join(root, "apps/api")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected API app to be skipped, got err=%v", err)
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file %s to exist: %v", path, err)
	}
}

func assertFileNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected file %s to be absent", path)
	}
}

func assertFileContains(t *testing.T, path string, snippet string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file %s: %v", path, err)
	}
	if !strings.Contains(string(data), snippet) {
		t.Fatalf("expected file %s to contain %q", path, snippet)
	}
}
