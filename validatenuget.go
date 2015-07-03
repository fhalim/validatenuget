package main

import (
	"encoding/xml"
	"flag"
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
	file, err := os.Open(*baseDir)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	processDir(file, state)
	deps := reconcileDependencies(state.dependencies)
	//log.Println(deps)
	for pkgName, versions := range deps {
		if len(versions) > 1 {
			log.Printf("Package %v has %v versions", pkgName, len(versions))
			for version, projects := range versions {
				log.Printf("%v:", version)
				for _, project := range projects {
					log.Printf("\t%v", project)
				}
			}
		}
	}
	defer file.Close()
}

func processPackages(file *os.File, state PackagesState) {
	//log.Printf("Processing %v", file.Name())
	decoder := xml.NewDecoder(file)
	var packages Packages

	err := decoder.Decode(&packages)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	state.dependencies[file.Name()] = packages.Packages
}

func processDir(dir *os.File, state PackagesState) {
	fileInfos, err := dir.Readdir(0)
	//log.Printf("Processing directory %v", dir.Name())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	for _, item := range fileInfos {
		name := item.Name()
		if item.Mode().IsRegular() && name == "packages.config" {
			file, err := os.Open(path.Join(dir.Name(), name))
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			defer file.Close()
			processPackages(file, state)
		} else if item.Mode().IsDir() {
			subDir, err := os.Open(path.Join(dir.Name(), name))
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			defer subDir.Close()
			processDir(subDir, state)
		}
	}
}
