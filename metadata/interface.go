package metadata

import (
	"github.com/spf13/afero"
)

var appFS afero.Fs

// AbsPath is an advisory type that represents an absolute path. It is advisory
// in that it is not forced to be absolute, but rather, meant to indicate
// intent, and make code easier to read.
type AbsPath string

// AbsPaths is a slice of `AbsPath`.
type AbsPaths []string

// Manager abstracts over a ksonnet application's metadata, allowing users to do
// things like: create and delete environments; search for prototypes; vendor
// libraries; and other non-core-application tasks.
type Manager interface {
	Root() AbsPath
	ComponentPaths() (AbsPaths, error)
	LibPaths(envName string) (libPath, envLibPath AbsPath)
	GenerateKsonnetLibData(spec ClusterSpec) ([]byte, []byte, error)
	CreateEnvironment(name, uri string, spec ClusterSpec, extensionsLibData, k8sLibData []byte) error
	//
	// TODO: Fill in methods as we need them.
	//
	// GetPrototype(id string) Protoype
	// SearchPrototypes(query string) []Protoype
	// VendorLibrary(uri, version string) error
	// DeleteEnv(name string) error
	//
}

// Find will recursively search the current directory and its parents for a
// `.ksonnet` folder, which marks the application root. Returns error if there
// is no application root.
func Find(path AbsPath) (Manager, error) {
	return findManager(path, afero.NewOsFs())
}

// Init will retrieve a cluster API specification, generate a
// capabilities-compliant version of ksonnet-lib, and then generate the
// directory tree for an application.
func Init(rootPath AbsPath, spec ClusterSpec) (Manager, error) {
	return initManager(rootPath, spec, appFS)
}

// ClusterSpec represents the API supported by some cluster. There are several
// ways to specify a cluster, including: querying the API server, reading an
// OpenAPI spec in some file, or consulting the OpenAPI spec released in a
// specific version of Kubernetes.
type ClusterSpec interface {
	data() ([]byte, error)
	resource() string // For testing parsing logic.
}

// ParseClusterSpec will parse a cluster spec flag and output a well-formed
// ClusterSpec object. For example, if the flag is `--version:v1.7.1`, then we
// will output a ClusterSpec representing the cluster specification associated
// with the `v1.7.1` build of Kubernetes.
func ParseClusterSpec(specFlag string) (ClusterSpec, error) {
	return parseClusterSpec(specFlag, appFS)
}

func init() {
	appFS = afero.NewOsFs()
}
