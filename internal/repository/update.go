package packages

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func (m *Manager) Online() bool {
	return m.index != nil
}

func (m *Manager) DataHash() string {
	rf, err := GetReleaseFile(m.dataDir, m.codeName)
	if err != nil {
		return ""
	}
	return rf.Hash()
}

func (m *Manager) loadIndex(targetDir string) error {
	index, err := loadPackagesaDBIndex(buildDBPath(targetDir, m.codeName, DBIndexName))
	if err != nil {
		return fmt.Errorf("UpdateDB: failed load db index: %v", err)
	}
	m.index = index
	m.dbs = make(map[Architecture]PackageDB)
	return nil
}

func updateServerCacheDate() {
	_ServerCacheDate = time.Now().UTC().Format(time.UnixDate)
}

func (m *Manager) UpdateDB() error {
	// always trying load index file at the end, no matter the result.
	defer m.loadIndex(m.dataDir)

	// 1. comparing ReleaseFile.Hash to check whether need update
	os.MkdirAll(m.dataDir, 0755)
	targetDir, err := ioutil.TempDir(m.dataDir, "lastore.packages.partition")
	if err != nil {
		return fmt.Errorf("UpdateDB: failed create temp directory: %v", err)
	}
	defer os.RemoveAll(targetDir)

	rf, err := DownloadReleaseFile(m.repoURL, m.codeName, targetDir)
	if err != nil {
		return fmt.Errorf("UpdateDB: failed download release file: %v", err)
	}

	h := m.DataHash()
	if h == rf.Hash() {
	//if false {
		//log.Debugf("DataHash %q is same as upstream.\n", h)
		changed, err := DownloadRepository(m.repoURL, rf, m.dataDir)
		if err != nil {
			return err
		}
		if changed {
			BuildCache(rf, m.dataDir, m.dataDir)
		}
		m.loadIndex(m.dataDir)
		return nil
	}

	// 2. download data and build DBs to tmp/${xx} directory
	//	log.Debugf("Update new datas (from %q --> %q) on %q\n", h, rf.Hash(), targetDir)
	_, err = DownloadRepository(m.repoURL, rf, targetDir)
	if err != nil {
		return fmt.Errorf("UpdateDB: failed download repository files: %v", err)
	}
	err = BuildCache(rf, targetDir, targetDir)
	if err != nil {
		return fmt.Errorf("UpdateDB: failed build dbs: %v", err)
	}

	// 3. unload index
	m.index = nil

	// 4. removing old data on filesystem.
	os.RemoveAll(buildDBPath(m.dataDir, m.codeName))

	// 5. moving tmp/${xx} to $dataDir on filesystem.
	//	log.Debugf("Renaming %q to %q\n", targetDir, buildDBPath(m.dataDir, m.codeName))
	return os.Rename(buildDBPath(targetDir, m.codeName), buildDBPath(m.dataDir, m.codeName))
}
