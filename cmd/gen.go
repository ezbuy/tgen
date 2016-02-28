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
	"os"

	"github.com/ezbuy/tgen/langs"
	_ "github.com/ezbuy/tgen/langs/swift"
	"github.com/ezbuy/tgen/utils"
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

		if output == "" {
			fmt.Println("-o output path must be specified")
			return
		}

		// check whether the path is existed
		if res := utils.PathExists(output); !res {
			fmt.Printf("output path [%s] is not valid\n", output)
			return
		}

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
			for key, _ := range langs.Langs {
				fmt.Printf("\t%s\n", key)
			}
		}
	},
}

var lang string
var input string
var output string

func init() {
	RootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	genCmd.PersistentFlags().StringVarP(&lang, "lang", "l", "", "language")
	genCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "input file")
	genCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output path")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
