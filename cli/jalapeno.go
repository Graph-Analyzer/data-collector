package cli

import (
	"graph-analyzer/data-collector/input"
	"graph-analyzer/data-collector/repository"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	jalapenoConfig JalapenoConfig
)

type JalapenoConfig struct {
	JagwHost             string `mapstructure:"JAGW_HOST" validate:"required"`
	JagwRequestPort      int64  `mapstructure:"JAGW_REQUEST_PORT" validate:"required"`
	JagwSubscriptionPort int64  `mapstructure:"JAGW_SUBSCRIPTION_PORT" validate:"required"`
}

var jalapenoCommand = &cobra.Command{
	Use:   "jalapeno",
	Short: "Use jalapeño gateway",
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := viper.Unmarshal(&jalapenoConfig); err != nil {
			log.Fatalf("Error unmarshalling config file: %s", err)
		}

		validate := validator.New()
		if err := validate.Struct(&jalapenoConfig); err != nil {
			log.Fatalf("Missing required attributes: %v\n", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugln("Using jalapeno command")

		graphRepo := repository.NewGraphRepository(driver)
		jalapeno := &input.Jalapeno{
			Host:             jalapenoConfig.JagwHost,
			RequestPort:      jalapenoConfig.JagwRequestPort,
			SubscriptionPort: jalapenoConfig.JagwSubscriptionPort,
		}

		input := input.InitInput(jalapeno, appConfig.NetworkName)
		input.Read(graphRepo)
	},
}

func init() {
	rootCmd.AddCommand(jalapenoCommand)
	cobra.OnInitialize(initConfig)

	jalapenoCommand.Flags().StringP("jagw-host", "", "", "JAGW Host")
	jalapenoCommand.Flags().StringP("jagw-request-port", "", "9903", "Request port")
	jalapenoCommand.Flags().StringP("jagw-subscription-port", "", "9902", "Subscription port")

	err := viper.BindPFlag("JAGW_HOST", jalapenoCommand.Flags().Lookup("jagw-host"))
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindPFlag("JAGW_REQUEST_PORT", jalapenoCommand.Flags().Lookup("jagw-request-port"))
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindPFlag("JAGW_SUBSCRIPTION_PORT", jalapenoCommand.Flags().Lookup("jagw-subscription-port"))
	if err != nil {
		log.Fatal(err)
	}
}
