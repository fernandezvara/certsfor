/*
Copyright © 2020 @fernandezvara

Permission is hereby granted,
 free of charge, to any person obtaining a copy
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
	"context"

	_ "github.com/fernandezvara/certsfor/db/badger"    // store driver
	_ "github.com/fernandezvara/certsfor/db/firestore" // store driver
	"github.com/fernandezvara/certsfor/db/store"
	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/spf13/viper"
)

// returns a service object or fails if configuration invalid or not found
func buildService() (srv *service.Service) {

	var err error

	if viper.GetBool(configAPIEnabled) {
		var cli *client.Client
		cli, err = client.New(viper.GetString(configAPIAddr),
			viper.GetString(configAPICA),
			viper.GetString(configAPICertificate),
			viper.GetString(configAPIKey),
		)
		er(err)
		srv = service.NewAsClient(cli)
	} else {
		var sto store.Store
		sto, err = store.Open(context.Background(),
			viper.GetString(configDBType),
			viper.GetString(configDBConnectionString),
		)
		er(err)
		srv = service.NewAsServer(sto)
	}

	return

}
