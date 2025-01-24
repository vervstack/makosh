package version

import (
	_ "embed"

	"github.com/godverv/matreshka"
)

//go:embed config.yaml
var prodConfig []byte

var version string

func init() {
	cfg, err := matreshka.ParseConfig(prodConfig)
	if err != nil {
		panic(err)
	}

	version = cfg.Version
}

func Version() string {
	return version
}
