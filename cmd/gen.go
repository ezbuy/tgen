// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ezbuy/tgen/global"
	"github.com/ezbuy/tgen/langs"
	_ "github.com/ezbuy/tgen/langs/go"
	_ "github.com/ezbuy/tgen/langs/java"
	_ "github.com/ezbuy/tgen/langs/swift"
	"github.com/samuel/go-thrift/parser"
	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate api source code",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		if lang == "" {
			fmt.Println("-l language must be specified")
			return
		}

		if input == "" {
			fmt.Println("-i input thrift file must be specified")
			return
		}

		f, err := filepath.Abs(input)
		if err != nil {
			log.Fatalf("failed to get absoulte path of input idl file: %s", err.Error())
		}

		global.InputFile = f
		global.Mode = mode
		global.NamespacePrefix = namespacePrefix
		global.GenWebApi = genWebApi
		global.GenRpcClient = genRpcCli
		global.ValidateParams = validateParams

		p := &parser.Parser{}
		parsedThrift, _, err := p.ParseFile(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(2)
		}

		if generator, ok := langs.Langs[lang]; ok {
			generator.Generate(output, parsedThrift)
		} else {
			fmt.Printf("lang %s is not supported\n", lang)
			fmt.Println("Supported language options are:")

			for key := range langs.Langs {
				fmt.Printf("\t%s\n", key)
			}
		}
	},
}

var lang string
var namespacePrefix string
var mode string
var genWebApi bool
var genRpcCli bool
var input string
var output string
var validateParams bool

func init() {
	RootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	genCmd.PersistentFlags().StringVarP(&lang, "lang", "l", "", "language")
	genCmd.PersistentFlags().StringVarP(&namespacePrefix, "prefix", "p", "", "namespace prefix")
	genCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "", "mode: rest or jsonrpc")
	genCmd.PersistentFlags().BoolVarP(&genWebApi, "webapi", "w", true, "generate webapi file(default true)")
	genCmd.PersistentFlags().BoolVarP(&genRpcCli, "rpccli", "r", false, "generate rpc client file(default false)")
	genCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "input file")
	genCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output path")
	genCmd.PersistentFlags().BoolVarP(&validateParams, "validate", "", false, "validate service method params (default false)")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
