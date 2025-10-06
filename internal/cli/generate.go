package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/example/node-creator-cli/internal/generator"
)

func newGenerateCmd() *cobra.Command {
	opts := generator.Options{}

	cmd := &cobra.Command{
		Use:   "generate [project_name]",
		Short: "Generate a Node.js monorepo project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectName = args[0]

			if opts.OutputDir == "" {
				opts.OutputDir = opts.ProjectName
			}

			opts.Modules = normalizeModules(opts.Modules)

			return generator.New(opts).Generate(cmd.Context())
		},
	}

	cmd.Flags().StringSliceVar(&opts.Modules, "modules", []string{"routes", "controllers", "services", "repositories", "middleware", "helpers"}, "Comma-separated list of modules to include")
	cmd.Flags().StringVar(&opts.Database, "database", "mysql", "Database type to configure (mysql, postgres, mongodb)")
	cmd.Flags().StringVar(&opts.Auth, "auth", "jwt", "Authentication strategy (jwt or paseto)")
	cmd.Flags().BoolVar(&opts.IncludeAPI, "api", true, "Generate API application")
	cmd.Flags().BoolVar(&opts.IncludeWorker, "worker", true, "Generate worker application")
	cmd.Flags().StringVar(&opts.CacheNamespace, "cache-namespace", "app", "Redis cache namespace prefix")
	cmd.Flags().StringVar(&opts.CacheGroup, "cache-group", "default", "Redis cache group segment")
	cmd.Flags().IntVar(&opts.CacheTTL, "cache-ttl", 300, "Redis cache TTL in seconds")
	cmd.Flags().StringVar(&opts.OutputDir, "output", "", "Target directory for the generated project")

	return cmd
}

func normalizeModules(mods []string) []string {
	result := make([]string, 0, len(mods))
	seen := map[string]struct{}{}
	for _, m := range mods {
		m = strings.TrimSpace(strings.ToLower(m))
		if m == "" {
			continue
		}
		if _, ok := seen[m]; ok {
			continue
		}
		seen[m] = struct{}{}
		result = append(result, m)
	}
	return result
}
