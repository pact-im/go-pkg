package main

import (
	"sort"
	"strconv"

	"golang.org/x/tools/go/packages"
)

type state struct {
	ByModule     map[string][]*packages.Package
	ModuleByPath map[string]packages.Module
	Modules      []packages.Module
}

func loadState(workspaceDir string, patterns ...string) (*state, error) {
	log("loading packages")
	pkgs, err := loadPackages(workspaceDir, patterns...)
	if err != nil {
		return nil, err
	}
	log("loaded packages: " + strconv.Itoa(len(pkgs)))
	return stateFromPacakges(pkgs), nil
}

func stateFromPacakges(pkgs []*packages.Package) *state {
	byModule, moduleByPath := packagesByModule(pkgs)
	modules := sortedModules(moduleByPath)
	return &state{
		ByModule:     byModule,
		ModuleByPath: moduleByPath,
		Modules:      modules,
	}
}

// loadPackages loads packages from the workspace that match the given pattern.
func loadPackages(workspaceDir string, patterns ...string) ([]*packages.Package, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedModule |
			packages.NeedName |
			packages.NeedTypes |
			packages.NeedImports |
			packages.NeedDeps,
		Dir:   workspaceDir,
		Tests: true,
	}, patterns...)
	if err != nil {
		return nil, err
	}
	return collectPackages(pkgs), nil
}

// collectPackages visits all packages and returns flattened list.
func collectPackages(pkgs []*packages.Package) []*packages.Package {
	var out []*packages.Package
	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		out = append(out, pkg)
	})
	return out
}

// packagesByModule groups packages by their containing Go module. It skips
// packages that do not have modules (i.e. stdlib).
func packagesByModule(pkgs []*packages.Package) (map[string][]*packages.Package, map[string]packages.Module) {
	byModule := make(map[string][]*packages.Package)
	moduleByPath := make(map[string]packages.Module)
	for _, pkg := range pkgs {
		if pkg.Module == nil {
			continue
		}
		m := *pkg.Module
		p := m.Path
		byModule[p] = append(byModule[p], pkg)
		moduleByPath[p] = m
	}
	return byModule, moduleByPath
}

// sortedModules returns modules sorted by their module path.
func sortedModules(moduleByPath map[string]packages.Module) []packages.Module {
	modules := make([]packages.Module, 0, len(moduleByPath))
	for _, m := range moduleByPath {
		modules = append(modules, m)
	}
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Path < modules[j].Path
	})
	return modules
}

// splitBrokenPackages splits packages into good and broken.
func splitBrokenPackages(pkgs []*packages.Package) (good, broken []*packages.Package) {
	for _, pkg := range pkgs {
		if len(pkg.Errors) != 0 {
			broken = append(broken, pkg)
			continue
		}
		good = append(good, pkg)
	}
	return
}

// packagesByImportPath groups packages by their import path.
func packagesByImportPath(pkgs []*packages.Package) map[string]*packages.Package {
	byImportPath := make(map[string]*packages.Package, len(pkgs))
	for _, pkg := range pkgs {
		byImportPath[pkg.PkgPath] = pkg
	}
	return byImportPath
}

// intersectPackageSets returns a list of import paths that are present in both
// sets.
func intersectPackageSets(pkgByImportPath, otherPkgByImportPath map[string]*packages.Package) []string {
	var intersect []string
	for importPath := range pkgByImportPath {
		if _, ok := otherPkgByImportPath[importPath]; !ok {
			continue
		}
		intersect = append(intersect, importPath)
	}
	sort.Strings(intersect)
	return intersect
}

// packagesNotInSet returns packages that are not present in the given set.
func packagesNotInSet(pkgs []*packages.Package, pkgByImportPath map[string]*packages.Package) []*packages.Package {
	var notInSet []*packages.Package
	for _, pkg := range pkgs {
		if _, ok := pkgByImportPath[pkg.PkgPath]; ok {
			continue
		}
		notInSet = append(notInSet, pkg)
	}
	return notInSet
}

// intersectModuleSets returns a list of modules that are present in both sets.
func intersectModuleSets(moduleByPath, otherModuleByPath map[string]packages.Module) []string {
	var intersect []string
	for modulePath := range moduleByPath {
		if _, ok := otherModuleByPath[modulePath]; !ok {
			continue
		}
		intersect = append(intersect, modulePath)
	}
	sort.Strings(intersect)
	return intersect
}

// modulesNotInSet returns a list of module versions that are not present
// in the moduleByPath set.
func modulesNotInSet(modules []packages.Module, moduleByPath map[string]packages.Module) []packages.Module {
	var notInSet []packages.Module
	for _, m := range modules {
		if _, ok := moduleByPath[m.Path]; ok {
			continue
		}
		notInSet = append(notInSet, m)
	}
	return notInSet
}
