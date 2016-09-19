package packages

import C "gopkg.in/check.v1"
import "os"
import "strings"
import "flag"
import "path"

var network = flag.Bool("network", false, "download test data from network")

type netSuite struct{}

func init() {
	C.Suite(&netSuite{})
}

func (*netSuite) SetUpSuite(c *C.C) {
	if !*network {
		c.Skip("no network support")
	}
}

func (*netSuite) TestBuildTestData(c *C.C) {
	repoURL := "http://10.0.4.226"
	targetDir := "testdata/packages"
	codeName := "unstable"
	rf, err := DownloadReleaseFile(repoURL, codeName, targetDir)
	c.Check(err, C.Equals, nil)
	_, err = DownloadRepository(repoURL, rf, targetDir)
	c.Check(err, C.Equals, nil)
}

func (*netSuite) TestDumpRepository(c *C.C) {
	repoURL := "http://pools.corp.deepin.com/deepin"
	targetDir := "/tmp/dump_repository"
	codeName := "unstable"

	rf, err := DownloadReleaseFile(repoURL, codeName, targetDir)
	c.Check(err, C.Equals, nil)
	_, err = DownloadRepository(repoURL, rf, targetDir)
	c.Check(err, C.Equals, nil)

	f, _ := os.Open(path.Join(targetDir, "unstable", "Release"))
	defer f.Close()
	rf, err = NewReleaseFile(f)
	c.Check(err, C.Equals, nil)

	err = BuildCache(rf, targetDir, targetDir)
	c.Check(err, C.Equals, nil)

}

func (*testWrap) TestRelease(c *C.C) {
	f, _ := os.Open("testdata/test_release")
	defer f.Close()

	rf, err := NewReleaseFile(f)

	c.Check(err, C.Equals, nil)
	c.Check(rf.CodeName, C.Equals, "experimental")

	c.Check(len(rf.Architectures), C.Equals, 1)
	c.Check(rf.Architectures[0], C.Equals, Architecture("amd64"))
	c.Check(strings.Join(rf.Components, ""), C.Equals, "non-free")

	c.Assert(len(rf.fileInfos), C.Equals, 31)
	c.Check(len(rf.FileInfos()), C.Equals, 1)
	pf := rf.fileInfos[2]
	c.Check(pf.Size, C.Equals, uint64(0x8f))
	c.Check(pf.MD5, C.Equals, "f23e539f4e40f8491b5b5512d1e7aaa9")
	c.Check(pf.Path, C.Equals, "main/binary-i386/Release")
}
