package packages

type Architecture string

const (
	ArchI386   Architecture = "i386"
	ArchAMD64  Architecture = "amd64"
	ArchAll    Architecture = "all"
	ArchUnknow Architecture = ""
)

// ArchitectureList store the supported Architectures
// TODO: configuration this and sync it with repository.json
var AvailableArchitectures = []Architecture{ArchAMD64, ArchI386}
