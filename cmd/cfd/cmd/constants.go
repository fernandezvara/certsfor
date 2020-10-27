package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// global variables
type globalConfig struct {
	home    string
	cfgFile string
}

// detect home folder
func init() {

	var err error
	global.home, err = homedir.Dir()
	fmt.Println("global:", global)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	configFileStringDefault = pathFromHome(".cfd", "config.yaml")
	configDBConnectionStringDefault = pathFromHome(".cfd", "cfd.db")

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
	configDBEnable                  = "db.type"
	configDBEnableEnv               = "CFD_DB_TYPE"
	configDBEnableDefault           = "badger"
	configDBConnectionString        = "db.connection"
	configDBConnectionStringEnv     = "CFD_DB_CONNECTION"
	configDBConnectionStringDefault string
	//configDBConnectionStringDefault = pathFromHome(".cfd", "db")

	// api config
	configAPIEnable        = "api.enabled"
	configAPIEnableEnv     = "CFD_API_ENABLED"
	configAPIEnableDefault = false
	configAPIAddr          = "api.addr"
	configAPIAddrEnv       = "CFD_API_ADDR"
	configAPIAddrDefault   = "127.0.0.1:8080"

	// Configuration contents
	configAPIClient        = "api.client"
	configAPIClientEnv     = "VEIL_API_CLIENT_ADDR"
	configAPIClientDefault = "http://127.0.0.1:8080"
	configAPILog           = "api.log"
	configAPILogEnv        = "VEIL_API_LOG"
)
