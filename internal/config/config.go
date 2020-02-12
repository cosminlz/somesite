package config

import (
	"github.com/namsral/flag"
)

var DataDirectory = flag.String("data-directory", "/go/src/cabhelp.ro/backend", "Path to loading templates and migrations")
