package main

import "github.com/normegil/zookeeper-rest/cmd"

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
