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
	"strings"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [APP_NAME]",
	Short: "Create new gofl project",
	Args:  cobra.MinimumNArgs(1),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if FileExists(".gofl") {
			fmt.Println("gofl project already initialized")
			return
		}

		path, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		mobileInit(path)
		apiInit(path, args[0])
		protoInit(path,  args[0])

		packageName := []byte("package: " + args[0])
		err = ioutil.WriteFile(".gofl", packageName, 0644)
		if err != nil {
			fmt.Println(err.Error())
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
		fmt.Println(err.Error())
		return
	}

	result = exec.Command("flutter", "create", "mobile")

	_, err = result.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func apiInit(path string, projectName string) {
	err := os.Mkdir(path+"/"+"api", os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
	}

	result := exec.Command("bash", "-c", "cd api; go mod init "+projectName)

	_, err = result.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	result = exec.Command("bash", "-c", "cd api; go get google.golang.org/grpc")

	_, err = result.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	result = exec.Command("bash", "-c", "cd api; go get github.com/golang/protobuf")

	_, err = result.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	generateGoMain(projectName)
}

func generateGoMain(packageName string) {
	template :=
		`package main

import (
	"` + packageName + `/services"
)

func main() {
	services.Run()
}`

	if !FileExists("./api/main.go") {
		content := []byte(template)
		err := ioutil.WriteFile("./api/main.go", content, 0644)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}

func protoInit(path string, packageName string) {
	err := os.Mkdir(path+"/"+"protos", os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
	}

	template :=
		`syntax = "proto3";

// example protocol buffer

package protos;

option objc_class_prefix = "` + strings.ToUpper(packageName) + `";
option java_package = "com.` + packageName + `.data.remote.proto";
option java_outer_classname = "` + packageName + `Proto";

enum AppPlatform {
    APP_PLATFORM_UNSPECIFIED = 0;
    ANDROID = 1;
    IOS = 2;
    WEB = 3;
}

enum BloodType {
    O_POSITIVE = 0;
    O_NEGATIVE = 1;
    B_POSITIVE = 2;
    B_NEGATIVE = 3;
    A_POSITIVE = 4;
    A_NEGATIVE = 5;
    AB_POSITIVE = 6;
    AB_NEGATIVE = 7;
}

enum AuthType {
    EMAIL = 0;
    MOBILE = 1;
    FACEBOOK = 2;
    TWITTER = 3;
    GOOGLE = 4;
    APPLE = 5;
    GITHUB = 6;
}


message User {
    string id = 1 [json_name="id"];
    string firstName = 2 [json_name="first_name"];
    string lastName = 3 [json_name="last_name"];
    string email = 4 [json_name="email"];
    string phoneNumber = 6 [json_name="phone_number"];
    int64 createdOn = 7 [json_name="created_on"];
    BloodType bloodType = 8 [json_name="blood_type"];
}

service Authentication {
    rpc Login (LoginRequest) returns (LoginReply) {}
}

message LoginRequest {
    AuthType authType = 1 [json_name="auth_type"];
    string userName = 2 [json_name="user_name"];
    string password = 3 [json_name="password"];
    string facebookId = 4 [json_name="facebook_id"];
    string mobileNumber = 5 [json_name="mobile_number"];
    string otpCode = 6 [json_name="otp_code"];

}

message LoginReply {
    string errorMessage = 1 [json_name="error_message"];
    User user = 2 [json_name="user"];
    string sessionId = 3 [json_name="session_id"];
    int64 ttl = 4 [json_name="ttl"];
}`

	if !FileExists("./protos/example.proto") {
		content := []byte(template)
		err := ioutil.WriteFile("./protos/example.proto", content, 0644)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
