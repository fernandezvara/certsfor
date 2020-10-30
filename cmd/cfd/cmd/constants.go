package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// global variables
type globalConfig struct {
	home     string // $HOME
	cfgFile  string // configuration file location
	filename string // file to parse/store/execute the action
	certFile string // certificate file location
	keyFile  string // key file location
}

// detect home folder
func init() {

	var err error
	global.home, err = homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	configFileStringDefault = pathFromHome(".cfd", "config.yaml")
	configDBConnectionStringDefault = pathFromHome(".cfd", "db")

}

func pathFromHome(paths ...string) (fullpath string) {

	fullpath = global.home

	for _, part := range paths {
		fullpath = filepath.Join(fullpath, part)
	}

	return

}

var (
	global globalConfig

	configFileString        = "config"
	configFileStringEnv     = "CFD_CONFIG"
	configFileStringDefault string
	// configFileStringDefault = pathFromHome(".cfd", "config.yaml")

	// database config
	configDBType                    = "db.type"
	configDBTypeEnv                 = "CFD_DB_TYPE"
	configDBTypeDefault             = "badger"
	configDBConnectionString        = "db.connection"
	configDBConnectionStringEnv     = "CFD_DB_CONNECTION"
	configDBConnectionStringDefault string
	//configDBConnectionStringDefault = pathFromHome(".cfd", "db")

	// api config
	configAPIEnabled            = "api.enabled"
	configAPIEnabledEnv         = "CFD_API_ENABLED"
	configAPIEnabledDefault     = false
	configAPIAddr               = "api.addr"
	configAPIAddrEnv            = "CFD_API_ADDR"
	configAPIAddrDefault        = "127.0.0.1:8080"
	configAPICA                 = "api.tls.ca"
	configAPICAEnv              = "CFD_API_CA"
	configAPICADefault          = ""
	configAPICertificate        = "api.tls.certificate"
	configAPICertificateEnv     = "CFD_API_ADDR"
	configAPICertificateDefault = ""
	configAPIKey                = "api.tls.key"
	configAPIKeyEnv             = "CFD_API_ADDR"
	configAPIKeyDefault         = ""

	// Configuration contents
	configAPIClient        = "api.client"
	configAPIClientEnv     = "VEIL_API_CLIENT_ADDR"
	configAPIClientDefault = "http://127.0.0.1:8080"
	configAPILog           = "api.log"
	configAPILogEnv        = "VEIL_API_LOG"
)
