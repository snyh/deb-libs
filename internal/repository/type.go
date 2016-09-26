package packages

import (
	"fmt"
	"io"
	"strconv"
)

type Type struct {
	Package       string         `json:"package"`
	Version       string         `json:"version"`
	InstalledSize int            `json:"installed_size"`
	Size          int            `json:"size"`
	Architectures []Architecture `json:"architectures"`
	Description   string         `json:"description"`
	Filename      string         `json:"filename"`
	Tag           string         `json:"tag"`
	Homepage      string         `json:"homepage"`
    Files         []string       `json:"files"`
}

func buildType(r io.Reader) (*Type, error) {
	dsc, err := NewDSCFile(r)
	if err != nil {
		return nil, err
	}

	t := &Type{}
	t.Package = dsc.GetString("package")
	t.Version = dsc.GetString("version")
	t.InstalledSize, _ = strconv.Atoi(dsc.GetString("installed-size"))
	t.Size, _ = strconv.Atoi(dsc.GetString("size"))

	for _, arch := range dsc.GetArrayString("architecture") {
		t.Architectures = append(t.Architectures, Architecture(arch))
	}
	t.Description = dsc.GetString("description")
	t.Filename = dsc.GetString("filename")
	t.Tag = dsc.GetString("tag")
	t.Homepage = dsc.GetString("homepage")

	if len(t.Package) < 2 {
		return nil, fmt.Errorf("W: pacakge name must be at least two characters long and start with an alphanumeric character: %q", t.Package)
	}
	if t.Filename == "" {
		return nil, fmt.Errorf("W: parsing DSC not enough fields: %q %q %q", t.Package, t.Tag, t.Filename)
	}
	return t, nil
}
