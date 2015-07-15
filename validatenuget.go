package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fhalim/validatenuget/dependencymanagement"
)
import linq "github.com/ahmetalpbalkan/go-linq"

func main() {
	baseDir := flag.String("dir", ".", "Base directory")
	prefixesStr := flag.String("prefixes", "", "Comma separated prefixes of interesting packages")
	flag.Parse()

	interestingPrefixes := strings.Split(*prefixesStr, ",")

	state := dependencymanagement.CreatePackagesState()

	processDir(*baseDir, state)

	deps := state.ReconcileDependencies()
	foundProblems := reportProblems(deps, interestingPrefixes)
	if foundProblems {
		os.Exit(2)
	}
}

func contains(prefixes []string, str string) bool {
	matchingPrefix := func(prefix linq.T) (bool, error) {
		return strings.HasPrefix(str, prefix.(string)), nil
	}
	result, _ := linq.From(prefixes).AnyWith(matchingPrefix)
	return result
}

func reportProblems(deps map[string]map[string][]string, interestingPrefixes []string) bool {
	foundProblems := false
	for pkgName, versions := range deps {
		if len(versions) > 1 && contains(interestingPrefixes, pkgName) {
			foundProblems = true
			fmt.Printf("Package %v has %v referenced versions\n", pkgName, len(versions))
			for version, projects := range versions {
				fmt.Printf("%v:\n", version)
				for _, project := range projects {
					fmt.Printf("\t%v\n", project)
				}
			}
		}
	}
	return foundProblems
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
