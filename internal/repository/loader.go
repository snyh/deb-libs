package packages

import (
	"encoding/gob"
	"fmt"
	"os"
)

func loadPackagesaDBIndex(fpath string) (*PackageDBIndex, error) {
	index := &PackageDBIndex{}
	err := load(fpath, &index)
	return index, err
}

func loadPackageDB(fpath string) (PackageDB, error) {
	var obj = make(PackageDB)
	err := load(fpath, &obj)
	return obj, err
}

func load(fpath string, obj interface{}) error {
	f, err := os.Open(fpath)
	if err != nil {
		return fmt.Errorf("store %q failed --> %v.", fpath, err)
	}
	defer f.Close()
	return gob.NewDecoder(f).Decode(obj)
}
func store(fpath string, obj interface{}) error {
	f, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("store %q failed --> %v.", fpath, err)
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(obj)
}
