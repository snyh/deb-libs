package packages

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

type FileInfo struct {
	Size         uint64
	Path         string
	Gzip         bool
	Bzip2        bool
	MD5          string
	Architecture Architecture
}

const DBIndexName = "index.dat"
const ReleaseFileName = "Release"

func _DBName(arch Architecture) string { return string(arch) + ".dat" }

type ReleaseFile struct {
	Date          string
	CodeName      string
	Description   string
	Components    []string
	Architectures []Architecture
	fileInfos     []FileInfo
}

// GetReleaseFile load ReleaseFile from dataDir with codeName
func GetReleaseFile(dataDir string, codeName string) (*ReleaseFile, error) {
	f, err := os.Open(buildDBPath(dataDir, codeName, ReleaseFileName))
	if err != nil {
		return nil, fmt.Errorf("GetReleaseFile open file error: %v", err)
	}
	defer f.Close()
	return NewReleaseFile(f)
}

// NewReleaseFile build a new ReleaseFile by reading contents from r
func NewReleaseFile(r io.Reader) (*ReleaseFile, error) {
	dsc, err := NewDSCFile(r)
	if err != nil {
		return nil, err
	}
	if len(dsc) == 0 {
		return nil, fmt.Errorf("empty dsc file")
	}

	rf := &ReleaseFile{}

	for _, arch := range dsc.GetArrayString("architectures") {
		rf.Architectures = append(rf.Architectures, Architecture(arch))
	}
	rf.Date = dsc.GetString("date")

	rf.Date = dsc.GetString("date")
	rf.CodeName = dsc.GetString("codename")
	rf.Description = dsc.GetString("description")
	rf.Date = dsc.GetString("date")
	rf.Components = dsc.GetArrayString("components")

	var ps []FileInfo
	for _, v := range dsc.GetMultiline("md5sum") {
		fs := strings.Split(strings.TrimSpace(v), " ")
		if len(fs) != 3 {
			continue
		}
		size, err := strconv.Atoi(fs[1])
		if err != nil {
			continue
		}

		ps = append(ps, FileInfo{
			Size: uint64(size),
			Path: fs[2],
			Gzip: strings.HasSuffix(fs[2], ".gz"),
			Bzip2: strings.HasSuffix(fs[2], ".bz2"),
			MD5:  fs[0],
		})
	}
	rf.fileInfos = ps
	if rf.CodeName == "" || len(rf.FileInfos()) == 0 || len(rf.Components) == 0 {
		return nil, fmt.Errorf("NewReleaseFile input data is invalid. %v", dsc)
	}
	return rf, nil
}

func buildDBPath(dataDir string, codeName string, name ...string) string {
	return path.Join(append([]string{dataDir, codeName}, name...)...)
}

func (rf ReleaseFile) Hash() string {
	var data []byte
	for _, finfo := range rf.FileInfos() {
		data = append(data, ([]byte)(finfo.MD5)...)
	}
	return fmt.Sprintf("%x", md5.Sum(data))
}

type FileInfos []FileInfo

func (infos FileInfos) Len() int {
	return len(infos)
}
func (infos FileInfos) Less(i, j int) bool {
	return infos[i].Path < infos[j].Path
}
func (infos FileInfos) Swap(i, j int) {
	infos[i], infos[j] = infos[j], infos[i]
}

func (rf ReleaseFile) FileInfos() []FileInfo {
    return append(rf.PackagesFileInfos(), rf.ContentsFileInfos()...)
}

func (rf ReleaseFile) ContentsFileInfos() []FileInfo {
	var set = make(map[string]FileInfo)
	for _, arch := range rf.Architectures {
		for _, component := range rf.Components {
			raw := component + "/Contents-" + string(arch)
			zip := raw + ".bz2"
			for _, f := range rf.fileInfos {
				if f.Path != raw && f.Path != zip {
					continue
				}
				_, ok := set[raw]
				if !ok {
					//store it if there hasn't content
					f.Architecture = arch
					set[raw] = f
				}
				if f.Bzip2 {
					//overwrite if it support bzip2
					f.Architecture = arch
					set[raw] = f
				}
			}
		}
	}

	var r = make(FileInfos, 0)
	for _, f := range set {
		r = append(r, f)
	}
	sort.Sort(r)
	return r
}


func (rf ReleaseFile) PackagesFileInfos() []FileInfo {
	var set = make(map[string]FileInfo)
	for _, arch := range rf.Architectures {
		for _, component := range rf.Components {
			raw := component + "/binary-" + string(arch) + "/Packages"
			zip := raw + ".gz"
			for _, f := range rf.fileInfos {
				if f.Path != raw && f.Path != zip {
					continue
				}
				_, ok := set[raw]
				if !ok {
					//store it if there hasn't content
					f.Architecture = arch
					set[raw] = f
				}
				if f.Gzip {
					//overwrite if it support gzip
					f.Architecture = arch
					set[raw] = f
				}
			}
		}
	}

	var r = make(FileInfos, 0)
	for _, f := range set {
		r = append(r, f)
	}
	sort.Sort(r)
	return r
}

func hashFile(fpath string) string {
	f, err := os.Open(fpath)
	if err != nil {
		return ""
	}
	defer f.Close()

	hash := md5.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		return ""
	}
	var r [16]byte
	copy(r[:], hash.Sum(nil))
	return fmt.Sprintf("%x", r)
}
