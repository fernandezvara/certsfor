package cmd

import (
	"context"

	"github.com/fernandezvara/certsfor/internal/service"
	_ "github.com/fernandezvara/certsfor/pkg/badger" // store driver
	"github.com/fernandezvara/certsfor/pkg/client"
	_ "github.com/fernandezvara/certsfor/pkg/firestore" // store driver
	"github.com/fernandezvara/certsfor/pkg/store"
	"github.com/spf13/viper"
)

// returns a service object or fails if configuration invalid or not found
func buildService() (srv *service.Service) {

	var err error

	if viper.GetBool(configAPIEnabled) {
		var cli *client.Client
		cli, err = client.New()
		er(err)
		srv = service.NewAsClient(cli)
	} else {
		var sto store.Store
		sto, err = store.Open(context.Background(), viper.GetString(configDBType), viper.GetString(configDBConnectionString))
		er(err)
		srv = service.NewAsServer(sto)
	}

	return

}
