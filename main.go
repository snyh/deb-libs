package main

import (
    "./internal/repository"
    "fmt"
    //"encoding/json"
    "flag"
)

func main() {
    flag.Parse()
    s := flag.Args()
    if len(s) < 1 {
        fmt.Println("E: not enough arg")
        return
    }

    //m, err := packages.NewManager("test_out", "http://pools.corp.deepin.com/deepin", "unstable")
    m, err := packages.NewManager("test_out", "http://10.0.0.10", "unstable")
	//m, err := packages.NewManager("test_out", "http://packages.deepin.test", "mydist")
	if err != nil {
		fmt.Println("E:", err)
		return
	}

	// Check whether need update the cache of repository
	err = m.UpdateDB()
	fmt.Println("M:", err)

	//Current only support field below, but add new fields is trivial.

	// 	type Type struct {
	// 	Package       string         `json:"package"`
	// 	Version       string         `json:"version"`
	// 	InstalledSize int            `json:"installed_size"`
	// 	Size          int            `json:"size"`
	// 	Architectures []Architecture `json:"architectures"`
	// 	Description   string         `json:"description"`
	// 	Filename      string         `json:"filename"`
	// 	Tag           string         `json:"tag"`
	// 	Homepage      string         `json:"homepage"`
	// }

	for _, id := range m.Search(s[0]) {
		if d, ok := m.Get(id); ok {
			fmt.Println("------------------Name:", d.Package, "---------------")
			fmt.Println("Version:", d.Version)
			fmt.Println("DESC:", d.Description)
			fmt.Println("Files:")
            for _, f := range d.Files {
                fmt.Printf(" %s\n", f)
            }

			fmt.Println("\n\n")
		}
	}

    /*
	rf, err := packages.GetReleaseFile("test_out/packages", "unstable")
	if err != nil {
		fmt.Println("E:", err)
		return
	}
	// type ReleaseFile struct {
	// 	Date          string
	// 	CodeName      string
	// 	Description   string
	// 	Components    []string
	// 	Architectures []Architecture
	// 	fileInfos     []PackagesFileInfo
	// }
	fmt.Println("Release File\n")
	d, _ := json.Marshal(rf)
	fmt.Printf("%q\n", string(d))
    */
}
