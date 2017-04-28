package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/normegil/zookeeper-rest/api/node"
	"github.com/normegil/zookeeper-rest/modules/environment"
	"github.com/normegil/zookeeper-rest/modules/log"
	"github.com/normegil/zookeeper-rest/modules/zookeeper"
	"github.com/normegil/zookeeper-rest/router"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFilePath string

var verbose bool
var serverPort int

var logPath string
var logFileAge int

var zkAddress string

var mongoAddress string
var mongoPort int
var mongoDB string
var mongoCollection string
var mongoUser string
var mongoPass string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "zookeeper-rest",
	Short: "Zookeeper REST server",
	Long:  `REST server connecting to a Zookeeper instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.New(log.Options{
			Verbose: verbose,
			File: log.FileOptions{
				FolderPath: logPath,
				FileName:   "zk-rest",
				MaxAge:     time.Duration(logFileAge*24) * time.Hour,
			},
			DB: log.MongoOptions{
				Address:    mongoAddress,
				Port:       strconv.Itoa(mongoPort),
				DB:         mongoDB,
				Collection: mongoCollection,
				User:       mongoUser,
				Password:   mongoPass,
			},
		})
		env := environment.Env{logger, zookeeper.Zookeeper{zkAddress, logger}}

		rt := router.New(env)
		if err := rt.Register(node.Controller{env}.Routes()); nil != err {
			panic(errors.Wrap(err, "Could not register Node controllers: "))
		}
		if err := rt.Listen(serverPort); nil != err {
			panic(errors.Wrapf(err, "Fatal error while server listening (port:%d)", serverPort))
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&configFilePath, "config", "", "config file (default is $HOME/.zookeeper-rest.yaml)")

	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose mode")
	RootCmd.PersistentFlags().IntVarP(&serverPort, "port", "p", 8080, "Port on which the server will listen")

	RootCmd.PersistentFlags().IntVar(&logFileAge, "log-age", 7, "Number of days for the log file rotation")
	RootCmd.PersistentFlags().StringVar(&logPath, "log-dir", "/tmp", "Path to directory where to store log files")

	RootCmd.PersistentFlags().StringVar(&zkAddress, "zk-address", "127.0.0.1", "Address of Zookeeper server")

	RootCmd.PersistentFlags().StringVar(&mongoAddress, "mg-address", "", "Address of the mongo DB server for logging. If not specified, mongoDB will not be used.")
	RootCmd.PersistentFlags().IntVar(&mongoPort, "mg-port", 27017, "Port of the mongo DB server.")
	RootCmd.PersistentFlags().StringVar(&mongoDB, "mg-db", "zookeeper-rest", "Name of the mongo DB")
	RootCmd.PersistentFlags().StringVar(&mongoCollection, "mg-log-collection", "logging", "Name of the mongo Collection where the logging message wille be stored.")
	RootCmd.PersistentFlags().StringVar(&mongoUser, "mg-user", "", "User for Mongo DB")
	RootCmd.PersistentFlags().StringVar(&mongoPass, "mg-pass", "", "Password for Mongo DB")
}

func initConfig() {
	if configFilePath != "" {
		viper.SetConfigFile(configFilePath)
	}

	viper.SetConfigName(".zookeeper-rest") // name of config file (without extension)
	viper.AddConfigPath("$HOME")           // adding home directory as first search path
	viper.AutomaticEnv()                   // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
