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
	"io/ioutil"
	"os"
	"os/exec"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [APP_NAME]",
	Short: "Create new gofl project",
	Args:    cobra.MinimumNArgs(1),
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if fileExists(".gofl") {
			fmt.Println("gofl project already initialized")
			return
		}

		path, err := os.Getwd()
		if err != nil {
			return
		}

		mobileInit(path)
		apiInit(path)
		protoInit(path)

		packageName := []byte("package: " + args[0])
		err = ioutil.WriteFile(".gofl", packageName, 0644)
		if err != nil {
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func mobileInit(path string) {
	result := exec.Command("cd", path)

	_, err := result.Output()

	if err != nil {
		return
	}

	result = exec.Command("flutter", "create", "mobile")

	_, err = result.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func apiInit(path string) {
	err := os.Mkdir(path + "/" + "api", os.ModePerm)
	if err != nil {
	}


}

func protoInit(path string) {
	err := os.Mkdir(path + "/" + "protos", os.ModePerm)
	if err != nil {
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
