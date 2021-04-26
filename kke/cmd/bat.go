/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"kke/pkg/bat"
	"kke/pkg/utils"
	"log"

	"github.com/spf13/cobra"
)

var hostStr string

// batCmd represents the bat command
var batCmd = &cobra.Command{
	Use:   "bat",
	Short: "execute tasks in batches",
	Args:  cobra.MinimumNArgs(1),
	// Run: func(cmd *cobra.Command, args []string) {
	// 	mode := strings.ToLower(args[0])
	// 	hostList := strings.Split(batHosts, " ")
	// 	if len(hostList) == 1 {
	// 		hostList = strings.Split(batHosts, ",")
	// 	}

	// 	var hosts []bat.Host

	// 	for _, hoststr := range hostList {
	// 		user, cred, addr, port, err := utils.ParseHost(hoststr)
	// 		if err != nil {
	// 			log.Panic(err)
	// 		}
	// 		host := bat.Host{
	// 			Addr: addr,
	// 			Port: port,
	// 			User: user,
	// 			Cred: cred,
	// 		}

	// 		hosts = append(hosts, host)
	// 	}

	// 	task := bat.Task{
	// 		Hosts: hosts,
	// 		Mode:  mode,
	// 		Args:  args[1:],
	// 	}

	// 	fmt.Println(task.Do().Result)
	// },
}

var batPingCmd = &cobra.Command{
	Use:   "ping",
	Short: "execute ping task in batches",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		task := bat.Task{
			Hosts: initHost(hostStr),
			Mode:  "ping",
			Args:  args,
		}

		fmt.Println(task.Do().Result)
	},
}

var batExecuteCmd = &cobra.Command{
	Use:   "execute",
	Short: "execute command task in batches",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		task := bat.Task{
			Hosts: initHost(hostStr),
			Mode:  "execute",
			Args:  args,
		}

		fmt.Println(task.Do().Result)
	},
}

var batScriptCmd = &cobra.Command{
	Use:   "script",
	Short: "Bexecute script task in batches",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		task := bat.Task{
			Hosts: initHost(hostStr),
			Mode:  "script",
			Args:  args,
		}

		fmt.Println(task.Do().Result)
	},
}

var batTemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "execute file template task in batches",
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		task := bat.Task{
			Hosts: initHost(hostStr),
			Mode:  "template",
			Args:  args,
		}

		fmt.Println(task.Do().Result)
	},
}

var batShellCmd = &cobra.Command{
	Use:   "shell",
	Short: "execute intractive shell in batches",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		task := bat.Task{
			Hosts: initHost(hostStr),
			Mode:  "shell",
			Args:  args,
		}

		task.Do()
	},
}

var batPushCmd = &cobra.Command{
	Use:   "push",
	Short: "execute push task in batches",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		task := bat.Task{
			Hosts: initHost(hostStr),
			Mode:  "push",
			Args:  args,
		}

		fmt.Println(task.Do().Result)
	},
}

var batPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "execute pull task in batches",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		task := bat.Task{
			Hosts: initHost(hostStr),
			Mode:  "pull",
			Args:  args,
		}

		fmt.Println(task.Do().Result)
	},
}

func initHost(hoststr string) []bat.Host {
	hostList := []bat.Host{}

	hosts := utils.StringSplit(hoststr)

	for _, h := range hosts {
		user, cred, addr, port, err := utils.ParseHost(h)
		if err != nil {
			log.Panic(err)
		}
		host := bat.Host{
			Addr: addr,
			Port: port,
			User: user,
			Cred: cred,
		}

		hostList = append(hostList, host)
	}

	return hostList
}

func init() {

	batCmd.PersistentFlags().StringVarP(&hostStr, "host", "H", "", "host list to execute batch task. full format:<user>/<credential>@<address>:<port>")
	batCmd.AddCommand(batPingCmd)
	batCmd.AddCommand(batExecuteCmd)
	batCmd.AddCommand(batScriptCmd)
	batCmd.AddCommand(batTemplateCmd)
	batCmd.AddCommand(batShellCmd)
	batCmd.AddCommand(batPushCmd)
	batCmd.AddCommand(batPullCmd)
	rootCmd.AddCommand(batCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// batCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// batCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
