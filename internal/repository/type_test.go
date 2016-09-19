package packages

import (
	C "gopkg.in/check.v1"
)

func (*testWrap) TestGetter(c *C.C) {
	// repo.Build("testdata")
	// g := newManager("testdata", l)
	// t, ok := g.Get("deepin-artwork")
	// c.Check(ok, C.Equals, true)
	// c.Check(t.Package, C.Equals, "deepin-artwork")

	// c.Check(len(g.Search("deepin")), C.Not(C.Equals), 0)

	// p, ok := g.QueryPath("deepin-movie", "amd64")
	// c.Check(ok, C.Equals, true)
	// c.Check(p, C.Not(C.Equals), "")

	// p, ok = g.QueryPath("deepin-movie", "i386")
	// c.Check(ok, C.Equals, true)
	// c.Check(p, C.Not(C.Equals), "")

	// p, ok = g.QueryPath("deepin-movie", "whatthefuck")
	// c.Check(ok, C.Equals, false)
	// c.Check(p, C.Equals, "")
}

func (*testWrap) TestBuildPath(c *C.C) {
	// c.Check(buildPackageCachePath("abc", ArchAMD64), C.Equals, "abc/amd64_packages.dat")
	// c.Check(buildPackageCachePath("abc", ArchI386), C.Equals, "abc/i386_packages.dat")
	// c.Check(buildPackageCachePath("abc", ArchAll), C.Equals, "abc/all_packages.dat")
}

func (*testWrap) TestReduceArchitectures(c *C.C) {
	var data = []struct {
		L []Architecture
		R []Architecture
	}{
		{
			[]Architecture{ArchAll, ArchAMD64},
			[]Architecture{ArchAMD64, ArchI386},
		},
		{
			[]Architecture{ArchAMD64},
			[]Architecture{ArchAMD64},
		},
		{
			[]Architecture{ArchAll},
			AvailableArchitectures,
		},
		{
			[]Architecture{ArchI386, ArchAll, ArchAMD64, ArchAll},
			[]Architecture{ArchI386, ArchAMD64},
		},
		{
			nil,
			nil,
		},
	}
	for _, item := range data {
		for i, v := range normalizeArchitectures(item.L) {
			c.Check(item.R[i], C.Equals, v)
		}
	}
}
