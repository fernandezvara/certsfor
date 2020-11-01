package cmd

import (
	"context"
	"fmt"

	"github.com/fernandezvara/certsfor/internal/manager"
	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/certsfor/internal/structs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCertificateCmd represents the bootstrap command
var getCertificateCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a certificate/key from the store.",
	Long:  `Get a certificate/key from the store.`,
	Run:   getCertificateFunc,
}

func init() {
	rootCmd.AddCommand(getCertificateCmd)
	getCertificateCmd.Flags().StringVarP(&global.certFile, "cert", "c", "", "Certificate file location.")
	getCertificateCmd.Flags().StringVarP(&global.keyFile, "key", "k", "", "Key file location.")
	getCertificateCmd.Flags().StringVar(&global.collection, "ca-id", "", "CA Identifier. (required). [$CFD_CA_ID]")
	getCertificateCmd.Flags().StringVar(&global.collection, "cn", "", "Common Name. (required).")
	getCertificateCmd.MarkFlagRequired("cn")
}

func getCertificateFunc(cmd *cobra.Command, args []string) {

	var (
		srv        *service.Service
		collection string
		ca         *manager.CA
		cert       structs.Certificate
		err        error
		ctx        context.Context = context.Background()
	)

	srv = buildService()
	defer srv.Close()

	collection = collectionOrExit()

	ca, err = srv.CAGet(collection)
	er(err)

	cert, err = srv.CertificateGet(ctx, collection, viper.GetString("cn"))
	er(err)

	saveFiles(ca, cert.Certificate, cert.Key)

	fmt.Println("\n\nCertificate Created.")

}
