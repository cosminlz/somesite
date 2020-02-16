package config

import (
	"github.com/namsral/flag"
)

var DataDirectory = flag.String("data-directory", "", "Path to loading templates and migrations")
