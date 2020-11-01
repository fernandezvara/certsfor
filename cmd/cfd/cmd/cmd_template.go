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
