// downloader.go downloaad the debian "Packages" and "Release" file
// to the directories under $DATADIR/raw/${DISTS}
package packages

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

// DownloadRepository download files from rf.FileInfos()
// it ignoring unchanged file by checking MD5 value.
// return whether changed and error if any.
func DownloadRepository(repoURL string, rf *ReleaseFile, targetDir string) (bool, error) {
	changed := false
	for _, f := range rf.FileInfos() {
		url := repoURL + "/dists/" + rf.CodeName + "/" + f.Path
		target := path.Join(targetDir, rf.CodeName, "raw", f.Path)
		if HashFile(target) == f.MD5 {
			continue
		}
		changed = true
		err := download(url, target, f.Gzip)
		if err != nil {
			return false, err
		}
	}
	return changed, nil
}

func DownloadReleaseFile(repoURL string, codeName string, targetDir string) (*ReleaseFile, error) {
	url := fmt.Sprintf("%s/dists/%s/Release", repoURL, codeName)
	reps, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("DownloadReleaseFile  http.Get(%q) failed:(%v)", url, err)
	}
	defer reps.Body.Close()
	// There is no need check reps.StatusCode, ReadReleaseFile will check the validation.

	// download Release File
	fpath := path.Join(targetDir, codeName, "Release")
	os.MkdirAll(path.Dir(fpath), 0755)
	f, err := os.Create(fpath)
	if err != nil {
		return nil, fmt.Errorf("DownloadReleaseFile(%q) can't create file %q : %v", url, fpath, err)
	}
	defer f.Close()

	// build Release File
	rf, err := NewReleaseFile(io.TeeReader(reps.Body, f))
	if err != nil {
		return nil, fmt.Errorf("DownloadReleaseFile invalid Release file(%q) : %v", url, err)
	}
	return rf, nil
}

// DownloadTee download the url content to "dest" file
// and return a io.Reader to read it.
func download(url string, dest string, gz bool) error {
	os.MkdirAll(path.Dir(dest), 0755)

	reps, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("can't download %q : %v", url, err)
	}
	defer reps.Body.Close()

	if reps.StatusCode != 200 {
		return fmt.Errorf("can't download %q : %v", url, reps.Status)
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("Can't create file %s", url)
	}
	defer f.Close()

	n, err := io.Copy(f, reps.Body)
	if err != nil {
		return fmt.Errorf("DownloadTo: write content(%d) failed:%v", n, err)
	}

	return nil
}
