package cli

import (
	"errors"
	"graph-analyzer/data-collector/repository"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	appConfig Config
	cfgFile   string
	logLevel  int
	driver    neo4j.Driver
)

type Config struct {
	DatabaseHost     string `mapstructure:"NEO4J_HOST" validate:"required"`
	DatabasePort     string `mapstructure:"NEO4J_PORT" validate:"required"`
	DatabaseUser     string `mapstructure:"NEO4J_USER" validate:"required"`
	DatabaseRealm    string `mapstructure:"NEO4J_REALM"`
	DatabasePassword string `mapstructure:"NEO4J_PASSWORD" validate:"required"`
	NetworkName      string `mapstructure:"NETWORK_NAME" validate:"ascii,max=50"`
}

var rootCmd = &cobra.Command{
	Use:   "data-collector",
	Short: "data-collector - import graph to neo4j",
	Long:  `data-collector is used to import network graphs to neo4j`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(int(log.InfoLevel) + logLevel))
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp:          true,
			DisableLevelTruncation: true,
		})

		if err := viper.Unmarshal(&appConfig); err != nil {
			log.Fatalf("Error unmarshalling config file: %s", err)
		}

		validate := validator.New()
		if err := validate.Struct(&appConfig); err != nil {
			log.Fatalf("Missing required or wrong attribute(s): %v\n", err)
		}

		driver = repository.InitDBConnection(appConfig.DatabaseHost, appConfig.DatabasePort, appConfig.DatabaseUser, appConfig.DatabasePassword, appConfig.DatabaseRealm)

	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Whoops. There was an error while executing the CLI '%s'", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().CountVarP(&logLevel, "verbose", "v", "Increase output verbosity. Example: --verbose=2 or -vv")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is .env)")
	rootCmd.PersistentFlags().StringP("neo4j-host", "", "", "Neo4j Host (including protocol neo4j/bolt)")
	rootCmd.PersistentFlags().StringP("neo4j-port", "", "", "Neo4j Port")
	rootCmd.PersistentFlags().StringP("neo4j-user", "", "", "Neo4j Username")
	rootCmd.PersistentFlags().StringP("neo4j-password", "", "", "Neo4j Password")
	rootCmd.PersistentFlags().StringP("neo4j-realm", "", "", "Neo4j Realm")
	rootCmd.PersistentFlags().StringP("network-name", "", "default", "Display name of the network (ASCII, 50 characters max)")

	err := viper.BindPFlag("NEO4J_HOST", rootCmd.PersistentFlags().Lookup("neo4j-host"))
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindPFlag("NEO4J_PORT", rootCmd.PersistentFlags().Lookup("neo4j-port"))
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindPFlag("NEO4J_USER", rootCmd.PersistentFlags().Lookup("neo4j-user"))
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindPFlag("NEO4J_PASSWORD", rootCmd.PersistentFlags().Lookup("neo4j-password"))
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindPFlag("NEO4J_REALM", rootCmd.PersistentFlags().Lookup("neo4j-realm"))
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindPFlag("NETWORK_NAME", rootCmd.PersistentFlags().Lookup("network-name"))
	if err != nil {
		log.Fatal(err)
	}
}

func initConfig() {
	// Read environment variables
	viper.AutomaticEnv()

	// If custom env file is set via --config flag - use that one
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Fallback to default .env - return if file does not exist, use it if it exists
		if _, err := os.Stat(".env"); errors.Is(err, os.ErrNotExist) {
			return
		} else {
			viper.SetConfigFile(".env")
		}
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}
