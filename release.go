package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
)

const usage = `Usage: go run release.go [-dry-run=false] <version>
Example: GOPRIVATE=go.pact.im go run release.go v0.0.1`

const fakeVersion = "v0.0.0-00010101000000-000000000000"

type Module struct {
	Name string
	Path string
	Deps []string
}

func main() {
	dryRun := true

	flag.BoolVar(&dryRun, "dry-run", dryRun, "Print actions without making changes")
	flag.Usage = func() { fmt.Fprintln(os.Stderr, usage) }
	flag.Parse()

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(handler)

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	version := flag.Arg(0)

	logger.Info("Discovering Go modules...")
	modules, err := discoverModules()
	if err != nil {
		logger.Error("Failed to discover modules", "error", err)
		os.Exit(1)
	}

	logger.Info("Topologically sorting modules by dependencies...")
	sorted, err := topoSortModules(modules)
	if err != nil {
		logger.Error("Failed to sort modules", "error", err)
		os.Exit(1)
	}

	reverseDeps := make(map[string][]*Module)
	for _, mod := range sorted {
		for _, dep := range mod.Deps {
			if _, ok := modules[dep]; ok {
				reverseDeps[dep] = append(reverseDeps[dep], mod)
			}
		}
	}

	logger.Info("Modules will be processed in this order:")
	for _, mod := range sorted {
		logger.Info("Module", "name", mod.Name, "path", mod.Path)
	}

	for _, mod := range sorted {
		logger.Info("Processing module", "module", mod.Name, "path", mod.Path)

		tag := buildTag(mod.Path, version)
		logger.Info("Tagging module", "module", mod.Name, "tag", tag)
		if !dryRun {
			if err := tagAndPush(tag); err != nil {
				logger.Error("Failed to tag and push", "tag", tag, "error", err)
				os.Exit(1)
			}
		}

		affectedModules := reverseDeps[mod.Name]
		for _, dependent := range affectedModules {
			logger.Info("Updating dependent module", "module", mod.Name, "dependent", dependent.Name)
			if !dryRun {
				err := updateDependency(dependent.Path, mod.Name, version)
				if err != nil {
					logger.Error("Failed to update dependency", "module", dependent.Name, "error", err)
					os.Exit(1)
				}
				if err := stageGoModChange(dependent.Path); err != nil {
					logger.Error("Failed to stage dependency update", "module", dependent.Name, "error", err)
					os.Exit(1)
				}
			}
		}
		if !dryRun && len(affectedModules) > 0 {
			if err := commitGoModChange(mod.Name, version); err != nil {
				logger.Error("Failed to commit dependency updates", "module", mod.Name, "error", err)
				os.Exit(1)
			}
		}
	}

	logger.Info("All modules processed")
	if dryRun {
		logger.Warn("Dry run enabled — no changes were made")
	}
}

func discoverModules() (map[string]*Module, error) {
	var work struct {
		Use []struct {
			DiskPath string `json:"DiskPath"`
		} `json:"Use"`
	}
	if err := parseGoToolJSON(&work, "work", "edit", "-json"); err != nil {
		return nil, fmt.Errorf("failed to parse go.work: %w", err)
	}

	result := make(map[string]*Module)
	for _, use := range work.Use {
		modDir := use.DiskPath
		if !filepath.IsLocal(modDir) {
			return nil, fmt.Errorf("non-local path %q in go.work", modDir)
		}

		var modInfo struct {
			Module struct {
				Path string `json:"Path"`
			} `json:"Module"`
			Require []struct {
				Path    string `json:"Path"`
				Version string `json:"Version"`
			} `json:"Require"`
		}
		if err := parseGoToolJSON(&modInfo, "-C", modDir, "mod", "edit", "-json"); err != nil {
			return nil, fmt.Errorf("failed to parse go.mod in %q: %w", modDir, err)
		}

		var deps []string
		for _, req := range modInfo.Require {
			if req.Version == fakeVersion {
				continue
			}
			deps = append(deps, req.Path)
		}

		modName := modInfo.Module.Path
		result[modName] = &Module{
			Name: modName,
			Path: modDir,
			Deps: deps,
		}
	}
	return result, nil
}

func topoSortModules(modules map[string]*Module) ([]*Module, error) {
	visited := make(map[string]bool)
	temp := make(map[string]bool)
	var sorted []*Module

	var visit func(string) error
	visit = func(name string) error {
		if temp[name] {
			return fmt.Errorf("circular dependency involving %s", name)
		}
		if visited[name] {
			return nil
		}
		temp[name] = true
		for _, dep := range modules[name].Deps {
			if _, ok := modules[dep]; ok {
				if err := visit(dep); err != nil {
					return err
				}
			}
		}
		temp[name] = false
		visited[name] = true
		sorted = append(sorted, modules[name])
		return nil
	}

	names := slices.Sorted(maps.Keys(modules))
	for _, name := range names {
		if err := visit(name); err != nil {
			return nil, err
		}
	}

	return sorted, nil
}

func updateDependency(modDir, dep, version string) error {
	return run("go", "-C", modDir, "get", "--", dep+"@"+version)
}

func stageGoModChange(modDir string) error {
	return run("git", "-C", modDir, "add", "go.mod", "go.sum")
}

func commitGoModChange(modName, version string) error {
	msg := fmt.Sprintf("chore: bump %s to %s", modName, version)
	return run("git", "commit", "-m", msg)
}

func tagAndPush(tag string) error {
	if err := run("git", "tag", tag); err != nil {
		return err
	}
	return run("git", "push", "origin", tag)
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func parseGoToolJSON(target any, args ...string) error {
	cmd := exec.Command("go", args...)
	return parseJSONCommand(cmd, target)
}

func parseJSONCommand(cmd *exec.Cmd, target any) error {
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("command failed: %w", err)
	}
	if err := json.Unmarshal(output, target); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	return nil
}

func buildTag(diskPath, version string) string {
	return path.Join(filepath.ToSlash(diskPath), version)
}
