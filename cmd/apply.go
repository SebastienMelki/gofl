/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Generates code from proto files",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		path, err := os.Getwd()
		if err != nil {
			log.Println(err)
			return
		}


		err = os.Mkdir("mobile/lib/protos", os.ModePerm)
		if err != nil {
		}

		err = os.Mkdir("api/protos", os.ModePerm)
		if err != nil {
		}

		result := exec.Command("bash", "-c", "protoc --go_out=plugins=grpc:" + path + "/api/protos -I" + path + "/protos " +path + "/protos/*.proto")

		_, err = result.Output()

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = exec.Command("bash", "-c",  "protoc --dart_out=grpc:" + path + "/mobile/lib/protos -I" + path + "/protos " + path + "/protos/user.proto")

		_, err = result.Output()

		if err != nil {
			fmt.Println(err.Error())
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}