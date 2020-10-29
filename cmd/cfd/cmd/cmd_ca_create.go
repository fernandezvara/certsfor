package cmd

import (
	"context"
	"fmt"

	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/certsfor/internal/structs"
	"github.com/fernandezvara/certsfor/pkg/manager"
	"github.com/spf13/cobra"
)

// caCreateCmd represents the bootstrap command
var caCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new CA certificate/key pair.",
	Long: `Creates a new CA certificate/key pair.
	
	Every CA is identified by a uuid. To operate with the CA this ID must be passed on each request.`,
	Run: caCreateFunc,
}

func init() {
	caCmd.AddCommand(caCreateCmd)
}

func caCreateFunc(cmd *cobra.Command, args []string) {

	var (
		srv       *service.Service
		request   structs.APICertificateRequest
		man       *manager.CA
		id        string
		bytesCert []byte
		bytesKey  []byte
		err       error
		ctx       context.Context = context.Background()
	)

	srv = buildService()
	defer srv.Close()
	// subject.CommonName = "ca"
	man, id, bytesCert, bytesKey, err = srv.CACreate(ctx, request)

	fmt.Println(man, id, string(bytesCert), string(bytesKey), err)

	//man, id, = srv.Ca srv.CACreate(context.Background(), subject, 1, 0,0)

	// client := service.
	// cert, err :=

}
