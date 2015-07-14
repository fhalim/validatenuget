package dependencymanagement

import (
	"encoding/xml"
	"log"
	"os"
)

func processError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// ProcessPackagesFile updates the state object with dependencies from the specified packages.config
func (state PackagesState) ProcessPackagesFile(fileName string) error {
	file, err := os.Open(fileName)
	processError(err)
	defer file.Close()
	decoder := xml.NewDecoder(file)
	var packages Packages

	err = decoder.Decode(&packages)
	processError(err)
	if err == nil {
		state.Dependencies[file.Name()] = packages.Packages
	}
	return err
}

// CreatePackagesState creates a new packages state map
func CreatePackagesState() PackagesState {
	return PackagesState{Dependencies: make(map[string][]Package)}
}

// ReconcileDependencies Collates dependencies into a map by dependency/version/project
func (state PackagesState) ReconcileDependencies() map[string]map[string][]string {
	packagesConfigDependencies := state.Dependencies
	dependencyVersionProjectMap := make(map[string]map[string][]string)
	// dependency[version] = {project1, project2}
	for projectName, dependencies := range packagesConfigDependencies {
		for _, pkg := range dependencies {
			versionsForDependency, ok := dependencyVersionProjectMap[pkg.ID]
			if !ok {
				versionsForDependency = make(map[string][]string)
				dependencyVersionProjectMap[pkg.ID] = versionsForDependency
			}
			projectsForVersion, ok := dependencyVersionProjectMap[pkg.ID][pkg.Version]
			if !ok {
				projectsForVersion = make([]string, 0, 5)
				dependencyVersionProjectMap[pkg.ID][pkg.Version] = projectsForVersion
			}
			dependencyVersionProjectMap[pkg.ID][pkg.Version] = append(dependencyVersionProjectMap[pkg.ID][pkg.Version], projectName)
		}
	}
	return dependencyVersionProjectMap
}
