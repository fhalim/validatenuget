package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

// Package is package element in packages.config
type Package struct {
	XMLName         xml.Name `xml:"package"`
	ID              string   `xml:"id,attr"`
	Version         string   `xml:"version,attr"`
	TargetFramework string   `xml:"targetFramework,attr"`
}

// Packages is packages element in packages.config
type Packages struct {
	XMLName  xml.Name  `xml:"packages"`
	Packages []Package `xml:"package"`
}

// PackagesState contains dependencies found in package
type PackagesState struct {
	dependencies map[string][]Package
}

func createPackagesState() PackagesState {
	return PackagesState{dependencies: make(map[string][]Package)}
}

func reconcileDependencies(packagesConfigDependencies map[string][]Package) map[string]map[string][]string {
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

func main() {
	baseDir := flag.String("dir", ".", "Base directory")
	flag.Parse()

	state := createPackagesState()

	processDir(*baseDir, state)
	deps := reconcileDependencies(state.dependencies)
	reportProblems(deps)
}

func reportProblems(deps map[string]map[string][]string) {
	for pkgName, versions := range deps {
		if len(versions) > 1 {
			fmt.Printf("Package %v has %v referenced versions\n", pkgName, len(versions))
			for version, projects := range versions {
				fmt.Printf("%v:\n", version)
				for _, project := range projects {
					fmt.Printf("\t%v\n", project)
				}
			}
		}
	}
}

func processPackages(fileName string, state PackagesState) {
	file, err := os.Open(fileName)
	processError(err)
	defer file.Close()
	decoder := xml.NewDecoder(file)
	var packages Packages

	processError(decoder.Decode(&packages))
	state.dependencies[file.Name()] = packages.Packages
}

func processError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
func processDir(dirName string, state PackagesState) {
	dir, err := os.Open(dirName)
	processError(err)
	fileInfos, err := dir.Readdir(0)
	dir.Close()
	processError(err)
	for _, item := range fileInfos {
		name := item.Name()
		if item.Mode().IsRegular() && name == "packages.config" {
			fileName := path.Join(dir.Name(), name)
			processPackages(fileName, state)
		} else if item.Mode().IsDir() {
			subDirName := path.Join(dir.Name(), name)
			processDir(subDirName, state)
		}
	}
}
