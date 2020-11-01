package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// templateCmd creates a file using caTemplate as base
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Creates a YAML file with the data required for a Certificate creation.",
	Long: `Creates a YAML file with the data required for a Certificate creation.
	
The file will have the required fields already filled. Customize based on your needs.`,
	Run: templateFunc,
}

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.Flags().StringVarP(&global.filename, "file", "f", "", "Filename where the template will be stored in YAML format.")
	templateCmd.Flags().BoolVar(&global.bool1, "ca", false, "Template for CA")
}

func templateFunc(cmd *cobra.Command, args []string) {

	var (
		bytesFile []byte
		err       error
	)

	if global.filename == "" {
		global.filename, err = promptText("Filename", "./template.yaml", validationRequired)
		er(err)
	}

	if global.bool1 {
		bytesFile, err = yaml.Marshal(caTemplate())
		er(err)
	} else {
		bytesFile, err = yaml.Marshal(certTemplate())
		er(err)
	}

	er(ioutil.WriteFile(global.filename, bytesFile, 0640))

	fmt.Printf("\n\nTemplate Created. '%s'\n", global.filename)

}
