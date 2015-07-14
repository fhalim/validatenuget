package dependencymanagement

import "encoding/xml"

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
	Dependencies map[string][]Package
}
