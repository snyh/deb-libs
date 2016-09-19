package packages

//Getter Read and Query information from repository
type Getter interface {
	Search(query string) []string
	Get(id string) (Type, bool)

	//QueryPath
	QueryPath(id string, arch Architecture) (string, bool)

	//Debug
	RepositoryInfo() (*Repository, error)

	DataHash() string
}

type Setter interface {
	Getter
	UpdateDB() error
}

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
