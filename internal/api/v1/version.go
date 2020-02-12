package v1

import (
	"encoding/json"
	"net/http"

	"cabhelp.ro/backend/internal/config"
	"github.com/sirupsen/logrus"
)

// ServerVersion is the type representing the binary version
type ServerVersion struct {
	Version string `json:"version"`
}

// Marshalled json
var versionJSON []byte

func init() {
	var err error
	versionJSON, err = json.Marshal(ServerVersion{
		Version: config.Version,
	})

	if err != nil {
		panic(err)
	}
}

// VersionHandler is the handler for the /version route
func VersionHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
	if _, err := w.Write(versionJSON); err != nil {
		logrus.WithError(err).Error("Error writing version")
	}
}
