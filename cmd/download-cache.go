/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	// "io"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/docker/cli/cli/compose/loader"
	"github.com/docker/cli/cli/compose/schema"
	// interp "github.com/docker/cli/cli/compose/interpolation"
	// tpl "github.com/docker/cli/cli/compose/template"
	composetypes "github.com/docker/cli/cli/compose/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Service is a Docker Compose service entry with it's corresponding labels mapped
type Service struct {
  Labels map[string]string
}

type serviceRead struct {
  Labels interface{}
}

type composeDataRead struct {
  Services map[string]serviceRead
}

// downloadCacheCmd represents the download-cache command
var downloadCacheCmd = &cobra.Command{
	Use:   "download-cache [service_name_on_compose]",
	Short: "Downloads images from the 'cache_from' key in a compose service",
	Long: `Downloads images from the 'cache_from' key in a compose service:

docker-image-manager download-cache -f .semaphore/ci-compose.yml service_name --break-on-first`,

	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a service defined in the compose file")
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		imageList, err := getCacheImageList(composeFile, args[0])
		if err != nil {
			fmt.Println(err)
    	os.Exit(1)
		}

		fmt.Println(imageList)
	},
}

var composeFile string

func init() {
	rootCmd.AddCommand(downloadCacheCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	downloadCacheCmd.PersistentFlags().StringVarP(
		&composeFile,
		"compose-file",
		"c",
		"docker-compose.yml",
		"Specify an alternate compose file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCacheCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getCacheImageList(composeFile string, serviceName string) ([]string, error) {
	configDetails, err := getConfigDetails(composeFile)
	if err != nil {
		return nil, err
	}

	config, err := loader.Load(configDetails)
	if err != nil {
		fmt.Println(err)
		if fpe, ok := err.(*loader.ForbiddenPropertiesError); ok {
			return nil, errors.Errorf("Compose file contains unsupported options:\n\n%s\n",
				propertyWarnings(fpe.Properties))
		}

		return nil, err
	}

	var selectedService composetypes.ServiceConfig
	
	for _, service := range config.Services {
		if service.Name == serviceName {
			selectedService = service
		}
  }

	if selectedService.Name != serviceName {
		return nil, errors.New("Service '" + serviceName + "' not found")
	}
	
	return selectedService.Build.CacheFrom, nil
}

func getConfigDetails(composeFile string) (composetypes.ConfigDetails, error) {
	var details composetypes.ConfigDetails

	workingDir, err := os.Getwd()
	if err != nil {
		return details, err
	}
	
	details.WorkingDir = workingDir

	configFile, err := loadConfigFile(composeFile)
	if err != nil {
		return details, err
	}

	details.ConfigFiles = append(details.ConfigFiles, *configFile)
	details.Version = schema.Version(details.ConfigFiles[0].Config)
	
	details.Environment, err = buildEnvironment(os.Environ())
	return details, err
}

func propertyWarnings(properties map[string]string) string {
	var msgs []string
	for name, description := range properties {
		msgs = append(msgs, fmt.Sprintf("%s: %s", name, description))
	}
	sort.Strings(msgs)
	return strings.Join(msgs, "\n\n")
}

func buildEnvironment(env []string) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for _, s := range env {
		// if value is empty, s is like "K=", not "K".
		if !strings.Contains(s, "=") {
			return result, errors.Errorf("unexpected environment %q", s)
		}
		kv := strings.SplitN(s, "=", 2)
		result[kv[0]] = kv[1]
	}
	return result, nil
}

func loadConfigFile(filename string) (*composetypes.ConfigFile, error) {
	var bytes []byte
	var err error

	bytes, err = ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config, err := loader.ParseYAML(bytes)
	if err != nil {
		return nil, err
	}

	return &composetypes.ConfigFile{
		Filename: filename,
		Config:   config,
	}, nil
}