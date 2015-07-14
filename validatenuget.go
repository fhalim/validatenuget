package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fhalim/validatenuget/dependencymanagement"
)

func main() {
	baseDir := flag.String("dir", ".", "Base directory")
	flag.Parse()

	state := dependencymanagement.CreatePackagesState()

	processDir(*baseDir, state)

	deps := state.ReconcileDependencies()
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

func processError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
func processDir(dirName string, state dependencymanagement.PackagesState) {
	filepath.Walk(dirName, func(path string, item os.FileInfo, err error) error {
		name := item.Name()
		if item.Mode().IsRegular() && name == "packages.config" {
			return state.ProcessPackagesFile(path)
		}
		return nil
	})
}
