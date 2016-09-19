package packages

import (
	"fmt"
	"path"
	"sort"
	"strings"
	"sync"
)

type Manager struct {
	dataDir  string
	codeName string
	index    *PackageDBIndex
	dbs      map[Architecture]PackageDB
	dbLock   sync.Mutex
	repoURL  string
}

const DataDirectoryPrefix = "packages"

func init() {
	updateServerCacheDate()
}

func NewManager(baseDataDir string, repoURL string, codeName string) (*Manager, error) {
	if repoURL == "" || baseDataDir == "" || codeName == "" {
		return nil, fmt.Errorf("Please setup packages.newManager")
	}
	m := &Manager{
		dataDir:  path.Join(baseDataDir, DataDirectoryPrefix),
		codeName: codeName,
		dbs:      make(map[Architecture]PackageDB),
		repoURL:  repoURL,
	}

	return m, nil
}

func (m *Manager) Search(q string) []string {
	if !m.Online() {
		return make([]string, 0)
	}

	var r = make(map[string]struct{})
	for _, data := range m.index.PackagePaths {
		for id := range data {
			if strings.Contains(id, q) {
				r[id] = struct{}{}
			}
		}
	}
	return sortMapString(r)
}

func (m *Manager) QueryPath(id string, arch Architecture) (string, bool) {
	if !m.Online() {
		return "", false
	}

	data, ok := m.index.PackagePaths[arch]
	if !ok {
		return "", false
	}
	path, ok := data[id]
	return path, ok
}

// normalizeArchitectures expand "all" to Architecture and remove redundant archs
func normalizeArchitectures(archs []Architecture) []Architecture {
	var r []Architecture
	var addedFlag = make(map[Architecture]bool)
	for _, arch := range archs {
		if arch == ArchAll {
			for _, _arch := range AvailableArchitectures {
				if addedFlag[_arch] {
					continue
				}
				addedFlag[_arch] = true
				r = append(r, _arch)
			}
		} else {
			if addedFlag[arch] {
				continue
			}
			addedFlag[arch] = true
			r = append(r, arch)
		}
	}
	return r
}
func (m *Manager) Get(id string) (Type, bool) {
	if !m.Online() {
		return Type{}, false
	}

	archs := m.index.PackageArchitectures(id)
	for _, arch := range archs {
		DB, err := m.getDB(arch)
		if err != nil {
			continue
		}
		t, ok := DB[id]
		if !ok {
			continue
		}
		t.Architectures = archs
		return t, true
	}
	return Type{}, false
}

func (m *Manager) getDB(arch Architecture) (PackageDB, error) {
	// If we don't lock this, the loadPackageDB maybe invoked too many times
	// to cause memory exploded.
	m.dbLock.Lock()
	defer m.dbLock.Unlock()

	DB, ok := m.dbs[arch]

	if !ok {
		var err error
		DB, err = loadPackageDB(buildDBPath(m.dataDir, m.codeName, _DBName(arch)))
		if err != nil {
			return nil, err
		}
		m.dbs[arch] = DB
	}
	return DB, nil
}

func sortMapString(d map[string]struct{}) []string {
	var r = make([]string, 0)
	for k := range d {
		r = append(r, k)
	}
	sort.Strings(r)
	return r
}
