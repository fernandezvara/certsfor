package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// caTemplateCmd creates a file using caTemplate as base
var caTemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "Creates a YAML file with the data required for a CA creation.",
	Long: `Creates a YAML file with the data required for a CA creation.
	
The file will have the required fields already filled. Customize based on your needs.`,
	Run: caTemplateFunc,
}

func init() {
	caCmd.AddCommand(caTemplateCmd)
	caTemplateCmd.Flags().StringVarP(&global.filename, "file", "f", "", "Filename where the template will be stored in YAML format.")
}

func caTemplateFunc(cmd *cobra.Command, args []string) {

	var (
		bytesFile []byte
		err       error
	)

	if global.filename == "" {
		global.filename, err = promptText("Filename", "./template.yaml", validationRequired)
		er(err)
	}

	bytesFile, err = yaml.Marshal(caTemplate())
	er(err)

	er(ioutil.WriteFile(global.filename, bytesFile, 0640))

	fmt.Printf("\n\nTemplate Created. '%s'\n", global.filename)

}
