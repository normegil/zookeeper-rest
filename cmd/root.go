package cmd

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/normegil/zookeeper-rest/api/node"
	"github.com/normegil/zookeeper-rest/modules/database/mongo"
	"github.com/normegil/zookeeper-rest/modules/environment"
	"github.com/normegil/zookeeper-rest/modules/log"
	"github.com/normegil/zookeeper-rest/modules/zookeeper"
	"github.com/normegil/zookeeper-rest/router"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	mgo "gopkg.in/mgo.v2"
)

var configFilePath string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "zookeeper-rest",
	Short: "Zookeeper REST server",
	Long:  `REST server connecting to a Zookeeper instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		mgoAddress := viper.GetString(MONGO_ADDRESS) + ":" + strconv.Itoa(viper.GetInt(MONGO_PORT))
		session, err := mgo.Dial(mgoAddress)
		if err != nil {
			panic(errors.Wrapf(err, "Connecting to %s", mgoAddress))
		}
		defer session.Close()

		mgoCredentials := &mgo.Credential{
			Username: viper.GetString(MONGO_USER),
			Password: viper.GetString(MONGO_PASS),
		}
		if "" != mgoCredentials.Username {
			if err := session.Login(mgoCredentials); err != nil {
				panic(err)
			}
		}

		mgoDB := viper.GetString(MONGO_DATABASE)
		connection := mongo.NewSession(session, mgoDB)
		defer connection.Close()

		logger, err := log.New(log.Options{
			Verbose: viper.GetBool(VERBOSE),
			File: log.FileOptions{
				FolderPath: viper.GetString(LOG_DIRECTORY),
				FileName:   "zk-rest",
				MaxAge:     time.Duration(viper.GetInt(LOG_RETENTION)*24) * time.Hour,
			},
			DB: log.MongoOptions{
				URL:      mgoAddress,
				Database: mgoDB,
				User:     mgoCredentials.Username,
				Password: mgoCredentials.Password,
			},
		})
		if err != nil {
			panic(errors.Wrapf(err, "Could not initialize Logger"))
		}

		zkAddress := net.TCPAddr{
			IP:   net.ParseIP(viper.GetString(ZOOKEEPER_ADDRESS)),
			Port: viper.GetInt(ZOOKEEPER_PORT),
		}
		env := environment.Env{
			Logger: logger,
			Zk: zookeeper.Zookeeper{
				Address: zkAddress,
				Logger:  logger,
			},
			Mongo: connection,
		}

		rt := router.New(env)
		if err := rt.Register(node.Controller{env}.Routes()); nil != err {
			panic(errors.Wrap(err, "Could not register Node controllers: "))
		}
		if err := rt.Listen(viper.GetInt(PORT)); nil != err {
			panic(errors.Wrapf(err, "Fatal error while server listening (port:%d)", viper.GetInt(PORT)))
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
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("ZK_REST")
	viper.AutomaticEnv()

	RootCmd.PersistentFlags().StringVar(&configFilePath, "config", "", "config file (default is $HOME/.zookeeper-rest.toml)")
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolP("verbose", "v", true, "Verbose mode")
	viper.BindPFlag(VERBOSE, RootCmd.PersistentFlags().Lookup("verbose"))

	RootCmd.PersistentFlags().IntP("port", "p", 8080, "Port on which the server will listen")
	viper.BindPFlag(PORT, RootCmd.PersistentFlags().Lookup("port"))

	RootCmd.PersistentFlags().Int("log-retention", 7, "Number of days for the log file rotation")
	viper.BindPFlag(LOG_RETENTION, RootCmd.PersistentFlags().Lookup("log-retention"))

	RootCmd.PersistentFlags().String("log-dir", "/tmp", "Path to directory where to store log files")
	viper.BindPFlag(LOG_DIRECTORY, RootCmd.PersistentFlags().Lookup("log-dir"))

	RootCmd.PersistentFlags().String("zk-address", "127.0.0.1", "Address of Zookeeper server")
	viper.BindPFlag(ZOOKEEPER_ADDRESS, RootCmd.PersistentFlags().Lookup("zk-address"))

	RootCmd.PersistentFlags().Int("zk-port", 2181, "Port of Zookeeper server")
	viper.BindPFlag(ZOOKEEPER_PORT, RootCmd.PersistentFlags().Lookup("zk-port"))

	RootCmd.PersistentFlags().String("mg-address", "localhost", "Address of the mongo DB server for logging.")
	viper.BindPFlag(MONGO_ADDRESS, RootCmd.PersistentFlags().Lookup("mg-address"))

	RootCmd.PersistentFlags().Int("mg-port", 27017, "Port of the mongo DB server.")
	viper.BindPFlag(MONGO_PORT, RootCmd.PersistentFlags().Lookup("mg-port"))

	RootCmd.PersistentFlags().String("mg-db", "zookeeper-rest", "Name of the mongo DB")
	viper.BindPFlag(MONGO_DATABASE, RootCmd.PersistentFlags().Lookup("mg-db"))

	RootCmd.PersistentFlags().String("mg-user", "", "User for Mongo DB")
	viper.BindPFlag(MONGO_USER, RootCmd.PersistentFlags().Lookup("mg-user"))

	RootCmd.PersistentFlags().String("mg-pass", "", "Password for Mongo DB")
	viper.BindPFlag(MONGO_PASS, RootCmd.PersistentFlags().Lookup("mg-pass"))
}

func initConfig() {
	if configFilePath != "" {
		viper.SetConfigFile(configFilePath)
	}

	viper.SetConfigName(".zookeeper-rest") // name of config file (without extension)
	viper.SetConfigType("toml")
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.AutomaticEnv()         // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

const (
	PORT              = "server.port"
	VERBOSE           = "logging.verbose"
	LOG_RETENTION     = "logging.retention"
	LOG_DIRECTORY     = "logging.directory"
	MONGO_ADDRESS     = "mongo.address"
	MONGO_PORT        = "mongo.port"
	MONGO_DATABASE    = "mongo.database"
	MONGO_USER        = "mongo.user"
	MONGO_PASS        = "mongo.pass"
	ZOOKEEPER_ADDRESS = "zookeeper.address"
	ZOOKEEPER_PORT    = "zookeeper.port"
)
