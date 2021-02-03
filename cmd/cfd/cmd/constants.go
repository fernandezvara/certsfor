/*
Copyright © 2020 @fernandezvara

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// Version variable to be updated at build time
var Version = "0.1"

// global variables
type globalConfig struct {
	home        string // $HOME
	cfgFile     string // configuration file location
	filename    string // file to parse/store/execute the action
	certFile    string // certificate file location
	caCertFile  string // ca certificate file location
	bundleFile  string // bundle file location
	keyFile     string // key file location
	collection  string // ca id <-> colelction
	cn          string // common name as argument
	bool1       bool   // common bool option
	bool2       bool   // common bool option
	remaining   int    // remaining percert (integer) for expiration
	listen      string // webserver ip:port to listen
	root        string // webserver root
	quiet       bool   // run in quiet mode
	pfxFile     string // pfx file to save
	pfxPassword string // password for pfx file
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

	configQuiet             = "quiet"
	configQuietShort        = "q"
	configQuietEnv          = "CFD_QUIET"
	configQuietDefault bool = false

	// database config
	configDBType                    = "db.type"
	configDBTypeEnv                 = "CFD_DB_TYPE"
	configDBTypeDefault             = "badger"
	configDBConnectionString        = "db.connection"
	configDBConnectionStringEnv     = "CFD_DB_CONNECTION"
	configDBConnectionStringDefault string

	// api config
	configAPIEnabled          = "api.enabled"
	configAPIEnabledEnv       = "CFD_API_ENABLED"
	configAPIEnabledDefault   = false
	configAPIAddr             = "api.addr"
	configAPIAddrEnv          = "CFD_API_ADDR"
	configAPIAddrDefault      = "127.0.0.1:8080"
	configAPIAccessLog        = "api.log.access"
	configAPIAccessLogEnv     = "CFD_LOG_ACCESS"
	configAPIAccessLogDefault = []string{"stdout"}
	configAPIErrorLog         = "api.log.error"
	configAPIErrorLogEnv      = "CFD_LOG_ERROR"
	configAPIErrorLogDefault  = []string{"stderr"}
	configAPIDebugLog         = "api.log.debug"
	configAPIDebugLogEnv      = "CFD_LOG_DEBUG"
	configAPIDebugLogDefault  = false
	// configWebEnabled            = "api.web"
	// configWebEnabledEnv         = "CFD_WEB_ENABLED"
	// configWebEnabledDefault     = false

	// tls
	configTLSCA                              = "tls.ca"
	configTLSCAEnv                           = "CFD_TLS_CA"
	configTLSCADefault                       = ""
	configTLSCertificate                     = "tls.certificate"
	configTLSCertificateEnv                  = "CFD_TLS_CERT"
	configTLSCertificateDefault              = ""
	configTLSKey                             = "tls.key"
	configTLSKeyEnv                          = "CFD_TLS_KEY"
	configTLSKeyDefault                      = ""
	configTLSRequireClientCertificate        = "tls.require_client_certificate"
	configTLSRequireClientCertificateEnv     = "CFD_TLS_CLIENT_CERT"
	configTLSRequireClientCertificateDefault = false
	configTLSUseForce                        = "tls.force"
	configTLSUseForceEnv                     = "CFD_TLS_FORCE"
	configTLSUseForceDefault                 = false

	// ca id
	configCAID        = "ca-id"
	configCAIDEnv     = "CFD_CA_ID"
	configCAIDDefault = ""
)
