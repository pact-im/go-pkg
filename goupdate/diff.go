package main

import (
	"golang.org/x/exp/apidiff"
	"golang.org/x/tools/go/packages"
)

type stateDiff struct {
	// Added is a list of added module dependencies.
	Added []packages.Module
	// Removed is a list of removed module dependencies.
	Removed []packages.Module
	// Updated is a list of updated module dependencies.
	Updated []moduleDiff
}

type moduleDiff struct {
	// Path is the module path of the updated module.
	Path string

	// Old is the old module.
	Old packages.Module
	// New is the new module.
	New packages.Module

	// OldBroken is a list of broken packages in the old version.
	OldBroken []*packages.Package
	// OldBroken is a list of broken packages in the new version.
	NewBroken []*packages.Package

	// Added is a list of added packages.
	Added []*packages.Package
	// Added is a list of removed packages.
	Removed []*packages.Package

	// Changed is a list of changed packages.
	Changed []packageDiff
}

// Unchanged returns true if there are no changes except for the version bump.
func (d *moduleDiff) Unchanged() bool {
	return len(d.Changed) == 0 &&
		len(d.Removed) == 0 &&
		len(d.Added) == 0 &&
		len(d.NewBroken) == 0 &&
		len(d.OldBroken) == 0
}

type packageDiff struct {
	// Path is the package import path from the updated module.
	Path string

	// Old is the old package.
	Old *packages.Package
	// New is the old package.
	New *packages.Package

	// Diff contains API changes report.
	Diff apidiff.Report
}

func (d *packageDiff) Compatible() []apidiff.Change {
	return d.changes(true)
}

func (d *packageDiff) Incompatible() []apidiff.Change {
	return d.changes(false)
}

func (d *packageDiff) changes(compat bool) []apidiff.Change {
	var changes []apidiff.Change
	for _, c := range d.Diff.Changes {
		if c.Compatible != compat {
			continue
		}
		changes = append(changes, c)
	}
	return changes
}

// diffState returns the difference in the API between the two workspace states.
func diffState(prev, final *state) stateDiff {
	added := modulesNotInSet(final.Modules, prev.ModuleByPath)
	removed := modulesNotInSet(prev.Modules, final.ModuleByPath)

	intersect := intersectModuleSets(final.ModuleByPath, prev.ModuleByPath)

	var updated []moduleDiff
	for _, modulePath := range intersect {
		oldModule := prev.ModuleByPath[modulePath]
		newModule := final.ModuleByPath[modulePath]

		if oldModule.Version == newModule.Version {
			continue
		}

		oldPackages := prev.ByModule[modulePath]
		newPackages := final.ByModule[modulePath]

		oldGood, oldBroken := splitBrokenPackages(oldPackages)
		newGood, newBroken := splitBrokenPackages(newPackages)

		oldByImportPath := packagesByImportPath(oldGood)
		newByImportPath := packagesByImportPath(newGood)

		removed := packagesNotInSet(oldGood, newByImportPath)
		added := packagesNotInSet(newGood, oldByImportPath)

		intersect := intersectPackageSets(oldByImportPath, newByImportPath)

		var changed []packageDiff
		for _, importPath := range intersect {
			oldPkg := oldByImportPath[importPath]
			newPkg := newByImportPath[importPath]
			diff := apidiff.Changes(oldPkg.Types, newPkg.Types)
			if len(diff.Changes) == 0 {
				continue
			}
			changed = append(changed, packageDiff{
				Path: importPath,
				Old:  oldPkg,
				New:  newPkg,
				Diff: diff,
			})
		}

		updated = append(updated, moduleDiff{
			Path:      modulePath,
			Old:       oldModule,
			New:       newModule,
			OldBroken: oldBroken,
			NewBroken: newBroken,
			Added:     added,
			Removed:   removed,
			Changed:   changed,
		})
	}

	return stateDiff{
		Added:   added,
		Removed: removed,
		Updated: updated,
	}
}
