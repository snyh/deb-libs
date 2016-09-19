package packages

import (
	"fmt"
)

type Repository struct {
	Source             string               `json:"source"`
	CodeName           string               `json:"code_name"`
	Components         []string             `json:"components"`
	Architectures      []Architecture       `json:"architectures"`
	PackageNums        map[Architecture]int `json:"package_nums"`
	UpstreamUpdateDate string               `json:"upstream_update_date"`
	ServerCacheDate    string               `json:"server_cache_date"`
}

func (m *manager) RepositoryInfo() (*Repository, error) {
	if !m.Online() {
		return nil, fmt.Errorf("system offline")
	}
	rf, err := GetReleaseFile(m.dataDir, m.codeName)
	if err != nil {
		return nil, err
	}
	r := &Repository{
		Source:             m.repoURL,
		CodeName:           m.codeName,
		Components:         rf.Components,
		Architectures:      rf.Architectures,
		UpstreamUpdateDate: rf.Date,
		ServerCacheDate:    _ServerCacheDate,
		PackageNums:        make(map[Architecture]int),
	}

	for arch, data := range m.index.PackagePaths {
		for range data {
			r.PackageNums[arch] = r.PackageNums[arch] + 1
		}
	}
	return r, nil
}
