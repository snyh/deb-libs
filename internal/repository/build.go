package packages

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"compress/bzip2"
	"fmt"
	"os"
	"path"
	"strings"
)

var _ServerCacheDate string

type PackageDBIndex struct {
	PackagePaths map[Architecture]map[string]string
	DBPaths      map[Architecture]string
}

func NewPackagesaDBIndex(indexFile string) (*PackageDBIndex, error) {
	return loadPackagesaDBIndex(indexFile)
}

func (dbi PackageDBIndex) DBPath(arch Architecture) (string, bool) {
	p, ok := dbi.DBPaths[arch]
	return p, ok
}
func (dbi PackageDBIndex) Architectures() []Architecture {
	var archs []Architecture
	for arch := range dbi.DBPaths {
		archs = append(archs, arch)
	}
	return archs
}
func (dbi PackageDBIndex) PackageArchitectures(pid string) []Architecture {
	var r []Architecture
	for arch, paths := range dbi.PackagePaths {
		_, ok := paths[pid]
		if !ok {
			continue
		}
		r = append(r, arch)
	}
	return r
}
func (dbi PackageDBIndex) PackagePath(pid string, arch Architecture) (string, bool) {
	ps, ok := dbi.PackagePaths[arch]
	if !ok {
		return "", false
	}
	p, ok := ps[pid]
	return p, ok
}

type PackageDB map[string]Type

func BuildCache(rf *ReleaseFile, rawDataDir string, targetDir string) error {
	// 1. build $arch.dat
	DBSources := make(map[Architecture][]string)
	DBIndex := make(map[Architecture]string)
	DBs := make(map[Architecture]PackageDB)
	for _, f := range rf.PackagesFileInfos() {
		source := path.Join(rawDataDir, rf.CodeName, "raw", f.Path)
		DBSources[f.Architecture] = append(DBSources[f.Architecture], source)
	}
	for arch, sources := range DBSources {
		db, err := createPackageDB(sources)
		if err != nil {
			return nil
		}
		DBs[arch] = db
		target := buildDBPath(targetDir, rf.CodeName, _DBName(arch))
		DBIndex[arch] = target
	}

    // add package files
	contentsDBSources := make(map[Architecture][]string)
	for _, f := range rf.ContentsFileInfos() {
		source := path.Join(rawDataDir, rf.CodeName, "raw", f.Path)
		contentsDBSources[f.Architecture] = append(contentsDBSources[f.Architecture], source)
	}
	for arch, sources := range contentsDBSources {
		data, err := parsePackageFiles(sources)
		if err != nil {
			return nil
		}
        for name, files := range data {
            packageType := DBs[arch][name]
            packageType.Files = append(packageType.Files, files...)
            DBs[arch][name] = packageType
        }
	}

	// 2. build index.dat
	index := createPackageIndex(DBIndex, DBs)

	// 3. store DBs
	err := store(buildDBPath(targetDir, rf.CodeName, DBIndexName), index)
	if err != nil {
		return fmt.Errorf("BuildCache: failed store index.dat --> %v", err)
	}
	for arch, fpath := range DBIndex {
		err := store(fpath, DBs[arch])
		if err != nil {
			return fmt.Errorf("BuildCache: failed store %q(%q) --> %v", fpath, arch, err)
		}
	}

	updateServerCacheDate()

	return nil
}

func createPackageIndex(dbsPath map[Architecture]string, dbs map[Architecture]PackageDB) PackageDBIndex {
	index := PackageDBIndex{
		DBPaths:      dbsPath,
		PackagePaths: make(map[Architecture]map[string]string),
	}

	for arch, db := range dbs {
		index.PackagePaths[arch] = make(map[string]string)
		for _, t := range db {
			index.PackagePaths[arch][t.Package] = t.Filename
		}
	}
	return index
}

func createPackageDB(sourcePaths []string) (PackageDB, error) {
	r := make(map[string]Type)
	for _, source := range sourcePaths {
        fmt.Printf("I: start scanning %s\n", source)
		datas, err := parsePackageDBComponent(source)
		if err != nil {
			return nil, err
		}
        fmt.Printf("I: finish scanning %s\n", source)
		for _, t := range datas {
			r[t.Package] = t
		}
	}
	return r, nil
}

func parsePackageDBComponent(path string) ([]Type, error) {
	splitFn := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		l := len(data)
		for i, c := range data {
			if c == '\n' {
				if i+1 < l && data[i+1] == '\n' {
					return i + 2, data[:i], nil
				}
				if i+1 == l && atEOF {
					return i + 1, data[:i], nil
				}
			}
		}
		if !atEOF {
			return 0, nil, nil
		}

		if atEOF && l != 0 {
			return l, data, nil
		}

		return l, data, fmt.Errorf("end of file")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parsePackageDBComponent: Can't open :%v", err)
	}
	defer f.Close()

	var s *bufio.Scanner
	if strings.HasSuffix(strings.ToLower(path), ".gz") {
		gr, err := gzip.NewReader(f)
		if err != nil {
			return nil, fmt.Errorf("can't parse gzip file %q fallback to plain text.\n", path)
		}
		defer gr.Close()
		s = bufio.NewScanner(gr)
	} else {
		s = bufio.NewScanner(f)
	}

	s.Split(splitFn)

	var ts []Type
	for s.Scan() {
		t, err := buildType(bytes.NewBufferString(s.Text()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "parsePackageDBComponent invalid Type: %v", err)
			continue
		}
		ts = append(ts, *t)
	}
	return ts, nil
}

type PackageFiles []string

func parsePackageFiles(sourcePaths []string) (map[string]PackageFiles, error) {
    r := make(map[string]PackageFiles)
    for _, source := range sourcePaths {
        fmt.Printf("I: start scanning %s\n", source)
        datas, err := parseContents(source)
		if err != nil {
			return nil, err
		}
        fmt.Printf("I: finish scanning %s\n", source)
		for name, files := range datas {
			r[name] = files
		}
    }
	return r, nil
}

func parseContents(path string) (map[string]PackageFiles, error) {
    packageFiles := make(map[string]PackageFiles)
	f, err := os.Open(path)
	if err != nil {
        return nil, fmt.Errorf("E: Can't open :%v", err)
	}
	defer f.Close()

    var s *bufio.Scanner
	if strings.HasSuffix(strings.ToLower(path), ".gz") {
		gr, err := gzip.NewReader(f)
		if err != nil {
			return nil, fmt.Errorf("can't parse gzip file %q fallback to plain text.\n", path)
		}
		defer gr.Close()
		s = bufio.NewScanner(gr)
    } else if strings.HasSuffix(strings.ToLower(path), ".bz2") {
		br := bzip2.NewReader(f)
		s = bufio.NewScanner(br)
	} else {
		s = bufio.NewScanner(f)
	}

    type Files []string
    for s.Scan() {
        l := strings.Fields(s.Text())
        file := l[0]
        pkgPart := l[len(l) - 1]
        if file == "FILE" && pkgPart == "LOCATION" {
            continue
        }
        pkgParts := strings.Split(pkgPart, "/")
        pkgName := pkgParts[len(pkgParts) - 1]
        packageFiles[pkgName] = append(packageFiles[pkgName], file)
    }

    return packageFiles, nil
}
