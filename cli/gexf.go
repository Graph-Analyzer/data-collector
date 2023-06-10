package cli

import (
	"errors"
	"graph-analyzer/data-collector/input"
	"graph-analyzer/data-collector/repository"
	"os"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	grpcConfig GrpcConfig
)

type GrpcConfig struct {
	GrpcHost string `mapstructure:"GRPC_HOST" validate:"required"`
	GrpcPort int64  `mapstructure:"GRPC_PORT" validate:"required"`
}

var gexfCommand = &cobra.Command{
	Use:   "gexf",
	Short: "Use GEXF file format",
	PreRun: func(cmd *cobra.Command, args []string) {
		gexfFile, _ := cmd.Flags().GetString("file")
		gexfListener, _ := cmd.Flags().GetBool("listener")

		if gexfFile == "" && !gexfListener {
			log.Fatalf("one of -f [file] or -l [listener] must be used")
		}

		if gexfListener {
			if err := viper.Unmarshal(&grpcConfig); err != nil {
				log.Fatalf("Error unmarshalling config file: %s", err)
			}

			validate := validator.New()
			if err := validate.Struct(&grpcConfig); err != nil {
				log.Fatalf("Missing required attributes: %s", err)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugln("Using GEXF command")

		graphRepo := repository.NewGraphRepository(driver)

		gexfFile, gexfFileError := cmd.Flags().GetString("file")
		gexfListener, _ := cmd.Flags().GetBool("listener")

		if gexfListener {
			log.Debugln("Using GEXF listener command")

			gexf := &input.GexfListener{
				Port: grpcConfig.GrpcPort,
				Host: grpcConfig.GrpcHost,
			}

			input := input.InitInput(gexf, appConfig.NetworkName)
			input.Read(graphRepo)

			return
		}

		if gexfFile != "" {
			log.Debugln("Using GEXF file command")

			if gexfFileError != nil {
				log.Fatalf("Error getting file name: %s", gexfFileError)
			}

			if _, err := os.Stat(gexfFile); errors.Is(err, os.ErrNotExist) {
				log.Fatalf("File %s does not exists", gexfFile)
			}

			gexf := &input.GexfFile{
				Filename: gexfFile,
			}

			input := input.InitInput(gexf, appConfig.NetworkName)
			input.Read(graphRepo)

			return
		}

		log.Fatal("Nothing specified")
	},
}

func init() {
	rootCmd.AddCommand(gexfCommand)

	gexfCommand.Flags().StringP("file", "f", "", "Specify gexf file")
	gexfCommand.Flags().BoolP("listener", "l", false, "Use gexf listener (grpc)")
	gexfCommand.Flags().IntP("grpc-port", "", 8081, "gRPC listening port")
	gexfCommand.Flags().StringP("grpc-host", "", "localhost", "gRPC listening host")

	gexfCommand.MarkFlagsMutuallyExclusive("file", "listener")
	gexfCommand.MarkFlagsMutuallyExclusive("file", "grpc-port")
	gexfCommand.MarkFlagsMutuallyExclusive("file", "grpc-host")

	err := viper.BindPFlag("GRPC_HOST", gexfCommand.Flags().Lookup("grpc-host"))
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindPFlag("GRPC_PORT", gexfCommand.Flags().Lookup("grpc-port"))
	if err != nil {
		log.Fatal(err)
	}
}
