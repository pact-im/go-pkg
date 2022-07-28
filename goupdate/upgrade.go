package main

import (
	"os/exec"
	"strconv"
)

// upgradeModules recursively updates dependencies in the workspace.
func upgradeModules(w *workspace, s *state) (*state, error) {
	cache := make(map[string]downloadedModule)
	for i := 1; ; i++ {
		log("iteration " + strconv.Itoa(i))
		var deps []string
		for _, m := range s.Modules {
			if m.Main {
				continue
			}
			deps = append(deps, m.Path)
		}
		if len(deps) == 0 {
			return s, nil
		}

		log("looking for upgrades")
		downloads, err := queryUpgradesCached(w.Root(), deps, cache)
		if err != nil {
			return nil, err
		}
		var upgrades []downloadedModule
		for _, d := range downloads {
			m, ok := s.ModuleByPath[d.Path]
			if !ok {
				continue
			}
			if m.Version == d.Version {
				continue
			}
			log(m.Path + "@" + m.Version + " => " + d.Version)
			upgrades = append(upgrades, d)
		}
		if len(upgrades) == 0 {
			return s, nil
		}
		log("found upgrades: " + strconv.Itoa(len(upgrades)))

		requires := make([]string, len(upgrades))
		for i, u := range upgrades {
			requires[i] = u.Path + "@" + u.Version
		}

		log("updating require directives")
		for _, moduleDir := range w.Paths {
			if err := addGoModRequires(moduleDir, requires); err != nil {
				return nil, err
			}
		}

		log("syncing workspace")
		if err := syncWorkspace(w); err != nil {
			return nil, err
		}

		if s, err = loadState(w.Root(), "all"); err != nil {
			return nil, err
		}

		// Run go work sync again in case entries we do not actually
		// need were added to go.work.sum.
		if err := goworksync(w); err != nil {
			return nil, err
		}
	}
}

func queryUpgradesCached(workspaceDir string, deps []string, cache map[string]downloadedModule) ([]downloadedModule, error) {
	var uncached []string

	downloads := make(map[string]downloadedModule)
	for _, modulePath := range deps {
		if m, ok := cache[modulePath]; ok {
			downloads[modulePath] = m
			continue
		}
		uncached = append(uncached, modulePath)
	}

	if len(uncached) != 0 {
		upgrades, err := queryUpgrades(workspaceDir, uncached)
		if err != nil {
			return nil, err
		}
		for _, u := range upgrades {
			downloads[u.Path] = u
			cache[u.Path] = u
		}
	}

	out := make([]downloadedModule, 0, len(deps))
	for _, modulePath := range deps {
		u, ok := downloads[modulePath]
		if !ok {
			continue
		}
		out = append(out, u)
	}
	return out, nil
}

func addGoModRequires(moduleDir string, requires []string) error {
	base := []string{"mod", "edit"}
	args := make([]string, len(requires)+len(base))
	_ = copy(args, base)
	for i, require := range requires {
		args[i+len(base)] = "-require=" + require
	}

	c := exec.Command("go", args...)
	c.Dir = moduleDir
	return c.Run()
}
