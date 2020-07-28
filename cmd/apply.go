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
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Generates code from proto files",
	Long:  ``,
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


		result := exec.Command("bash", "-c", "protoc --go_out=plugins=grpc:"+path+"/api/protos -I"+path+"/protos "+path+"/protos/*.proto")

		_, err = result.Output()

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = exec.Command("bash", "-c", "protoc --dart_out=grpc:"+path+"/mobile/lib/protos -I"+path+"/protos "+path+"/protos/*.proto")

		_, err = result.Output()

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		generateService(GetPackageName())

	},
}

func generateService(packageName string) {
	servicesDirectoryName := "api/services"
	err := os.Mkdir(servicesDirectoryName, os.ModePerm)
	if err != nil {
	}


	protoPaths, err := WalkMatch("./protos/", "*.proto")
	if err != nil {
	}

	allTypes := ""
	var services []string

	for _, protoPath := range protoPaths {
		file, err := os.Open("./" + protoPath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		template :=
			`package services

import(
	"context"
	"` + packageName + `/protos"
)`

		currentService := ""

		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "service") {
				tokens := strings.Split(scanner.Text(), " ")
				currentService = tokens[1]
				allTypes += `type ` + currentService + "Service" + ` struct{}` + "\n"
				services = append(services, currentService)
			}

			if strings.Contains(scanner.Text(), "rpc") {

				tokens := strings.Split(strings.TrimSpace(scanner.Text()), " ")
				rpc := tokens[1]
				input := strings.TrimPrefix(strings.TrimSuffix(tokens[2], ")"), "(")
				output := strings.TrimPrefix(strings.TrimSuffix(tokens[4], ")"), "(")
				serviceFile := template
				serviceFile += "\n\n"
				serviceFile +=
					`func (s *` + currentService + "Service" + `) ` + rpc + `(ctx context.Context, in *protos.` + input + `) (*protos.` + output + `, error) {
	return &protos.` + output + `{
	}, nil
}`
				content := []byte(serviceFile)
				fileName := strings.ToLower(rpc) + ".go"
				if !FileExists(servicesDirectoryName + "/" + fileName) {
					err := ioutil.WriteFile(servicesDirectoryName + "/" + fileName, content, 0644)
					if err != nil {
					}
				}

			}
		}
	}

	if allTypes != "" {
		content := []byte("package services\n" + allTypes)
		fileName := "services.go"
		err := ioutil.WriteFile(servicesDirectoryName + "/" + fileName, content, 0644)
		if err != nil {
			fmt.Println(err)
		}
	}

	generateServerGo(packageName, services)
}

func generateServerGo(packageName string, services []string) {
	template :=
	`package services

import (
	"fmt"
	"` + packageName + `/protos"
	"google.golang.org/grpc"
	"net"
)

func Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))

	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
		return
	}

	grpcServer := grpc.NewServer()`

	for _, val := range services {
		template += "\n"
		template += "  " + strings.ToLower(val) + " := " + val + "Service{}\n"
		template += "  " + "protos.Register" + val + "Server(grpcServer, &" + strings.ToLower(val) +")\n\n"
	}
	template += `  fmt.Println("LISTENING ON PORT 7777")
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %s", err)
		return
	}
}`

	content := []byte(template)
	fileName := "api/services/server.go"
	err := ioutil.WriteFile(fileName, content, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func init() {
	rootCmd.AddCommand(applyCmd)
}

func GetPackageName() string {
	file, err := os.Open(".gofl")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, "package:") {
			tokens := strings.Split(text, "package: ")
			return tokens[1]
		}

	}
	return ""
}

func WalkMatch(root string, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
