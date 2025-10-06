package generator

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/example/node-creator-cli/templates"
)

type templateDefinition struct {
	TemplatePath string
	TargetPath   string
	Condition    func(Options) bool
}

type Generator struct {
	opts Options
}

// New creates a generator instance with the provided options.
func New(opts Options) *Generator {
	return &Generator{opts: opts}
}

// Generate renders the Node.js project scaffolding to disk.
func (g *Generator) Generate(ctx context.Context) error { // ctx currently unused but reserved for future IO cancellation.
	if err := g.opts.Validate(); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	targetDir, err := filepath.Abs(g.opts.OutputDir)
	if err != nil {
		return fmt.Errorf("resolve output directory: %w", err)
	}

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	data := buildTemplateData(g.opts)

	for _, def := range templateCatalog {
		if def.Condition != nil && !def.Condition(g.opts) {
			continue
		}
		if err := g.renderTemplate(targetDir, def, data); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) renderTemplate(baseDir string, def templateDefinition, data templateData) error {
	content, err := templates.FS.ReadFile(def.TemplatePath)
	if err != nil {
		return fmt.Errorf("load template %s: %w", def.TemplatePath, err)
	}

	tmpl, err := template.New(filepath.Base(def.TemplatePath)).Funcs(template.FuncMap{
		"ToLower": strings.ToLower,
		"ToUpper": strings.ToUpper,
	}).Parse(string(content))
	if err != nil {
		return fmt.Errorf("parse template %s: %w", def.TemplatePath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("render template %s: %w", def.TemplatePath, err)
	}

	targetPath := filepath.Join(baseDir, def.TargetPath)
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return fmt.Errorf("create directory for %s: %w", targetPath, err)
	}

	if err := os.WriteFile(targetPath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write file %s: %w", targetPath, err)
	}

	return nil
}

type templateData struct {
	Options
	ProjectSlug     string
	ProjectRoot     string
	HasRoutes       bool
	HasControllers  bool
	HasServices     bool
	HasRepositories bool
	HasMiddleware   bool
	HasHelpers      bool
}

func buildTemplateData(opts Options) templateData {
	slug := slugify(opts.ProjectName)
	return templateData{
		Options:         opts,
		ProjectSlug:     slug,
		ProjectRoot:     opts.ProjectName,
		HasRoutes:       opts.HasModule("routes"),
		HasControllers:  opts.HasModule("controllers"),
		HasServices:     opts.HasModule("services"),
		HasRepositories: opts.HasModule("repositories"),
		HasMiddleware:   opts.HasModule("middleware"),
		HasHelpers:      opts.HasModule("helpers"),
	}
}

func slugify(name string) string {
	lower := strings.ToLower(name)
	sanitized := strings.Builder{}
	for _, r := range lower {
		switch {
		case r >= 'a' && r <= 'z':
			sanitized.WriteRune(r)
		case r >= '0' && r <= '9':
			sanitized.WriteRune(r)
		default:
			sanitized.WriteRune('-')
		}
	}
	return strings.Trim(sanitized.String(), "-")
}

var templateCatalog = []templateDefinition{
	{TemplatePath: "root/package.json.tmpl", TargetPath: "package.json"},
	{TemplatePath: "root/README.md.tmpl", TargetPath: "README.md"},
	{TemplatePath: "root/tsconfig.json.tmpl", TargetPath: "tsconfig.json"},
	{TemplatePath: "root/eslintrc.cjs.tmpl", TargetPath: ".eslintrc.cjs"},
	{TemplatePath: "root/prettier.config.cjs.tmpl", TargetPath: "prettier.config.cjs"},
	{TemplatePath: "root/gitignore.tmpl", TargetPath: ".gitignore"},
	{TemplatePath: "env/.env.local.tmpl", TargetPath: "config/env/.env.local"},
	{TemplatePath: "env/.env.dev.tmpl", TargetPath: "config/env/.env.dev"},
	{TemplatePath: "env/.env.stg.tmpl", TargetPath: "config/env/.env.stg"},
	{TemplatePath: "env/.env.prod.tmpl", TargetPath: "config/env/.env.prod"},
	{TemplatePath: "packages/shared/package.json.tmpl", TargetPath: "packages/shared/package.json"},
	{TemplatePath: "packages/shared/index.ts.tmpl", TargetPath: "packages/shared/index.ts"},
	{TemplatePath: "packages/shared/database/index.ts.tmpl", TargetPath: "packages/shared/database/index.ts"},
	{TemplatePath: "packages/shared/cache/redis.ts.tmpl", TargetPath: "packages/shared/cache/redis.ts"},
	{TemplatePath: "packages/shared/libs/logger.ts.tmpl", TargetPath: "packages/shared/libs/logger.ts"},
	{TemplatePath: "packages/shared/libs/httpLogger.ts.tmpl", TargetPath: "packages/shared/libs/httpLogger.ts"},
	{TemplatePath: "packages/shared/helpers/index.ts.tmpl", TargetPath: "packages/shared/helpers/index.ts"},
	{TemplatePath: "packages/shared/auth/index.ts.tmpl", TargetPath: "packages/shared/auth/index.ts"},
	{TemplatePath: "packages/shared/auth/jwt.ts.tmpl", TargetPath: "packages/shared/auth/jwt.ts", Condition: func(o Options) bool { return o.Auth == "jwt" }},
	{TemplatePath: "packages/shared/auth/paseto.ts.tmpl", TargetPath: "packages/shared/auth/paseto.ts", Condition: func(o Options) bool { return o.Auth == "paseto" }},
	{TemplatePath: "apps/api/package.json.tmpl", TargetPath: "apps/api/package.json", Condition: func(o Options) bool { return o.IncludeAPI }},
	{TemplatePath: "apps/api/tsconfig.json.tmpl", TargetPath: "apps/api/tsconfig.json", Condition: func(o Options) bool { return o.IncludeAPI }},
	{TemplatePath: "apps/api/src/app.ts.tmpl", TargetPath: "apps/api/src/app.ts", Condition: func(o Options) bool { return o.IncludeAPI }},
	{TemplatePath: "apps/api/src/server.ts.tmpl", TargetPath: "apps/api/src/server.ts", Condition: func(o Options) bool { return o.IncludeAPI }},
	{TemplatePath: "apps/api/src/routes/index.ts.tmpl", TargetPath: "apps/api/src/routes/index.ts", Condition: func(o Options) bool { return o.IncludeAPI && o.HasModule("routes") }},
	{TemplatePath: "apps/api/src/controllers/index.ts.tmpl", TargetPath: "apps/api/src/controllers/index.ts", Condition: func(o Options) bool { return o.IncludeAPI && o.HasModule("controllers") }},
	{TemplatePath: "apps/api/src/services/index.ts.tmpl", TargetPath: "apps/api/src/services/index.ts", Condition: func(o Options) bool { return o.IncludeAPI && o.HasModule("services") }},
	{TemplatePath: "apps/api/src/repositories/index.ts.tmpl", TargetPath: "apps/api/src/repositories/index.ts", Condition: func(o Options) bool { return o.IncludeAPI && o.HasModule("repositories") }},
	{TemplatePath: "apps/api/src/middleware/index.ts.tmpl", TargetPath: "apps/api/src/middleware/index.ts", Condition: func(o Options) bool { return o.IncludeAPI && o.HasModule("middleware") }},
	{TemplatePath: "apps/api/src/helpers/index.ts.tmpl", TargetPath: "apps/api/src/helpers/index.ts", Condition: func(o Options) bool { return o.IncludeAPI && o.HasModule("helpers") }},
	{TemplatePath: "apps/worker/package.json.tmpl", TargetPath: "apps/worker/package.json", Condition: func(o Options) bool { return o.IncludeWorker }},
	{TemplatePath: "apps/worker/tsconfig.json.tmpl", TargetPath: "apps/worker/tsconfig.json", Condition: func(o Options) bool { return o.IncludeWorker }},
	{TemplatePath: "apps/worker/src/index.ts.tmpl", TargetPath: "apps/worker/src/index.ts", Condition: func(o Options) bool { return o.IncludeWorker }},
}
