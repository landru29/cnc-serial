// Package main is the main application entrypoint.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	"github.com/landru29/cnc-serial/internal/application"
	"github.com/landru29/cnc-serial/internal/gcode/grbl"
	"github.com/landru29/cnc-serial/internal/lang"
	"github.com/landru29/cnc-serial/internal/stack/memory"
	"github.com/landru29/cnc-serial/internal/transport/serial"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type options struct {
	availableLanguages []lang.Language `              json:"-"             yaml:"-"`
	Language           lang.Language   `default:"EN"  json:"language"      mapstructure:"language" yaml:"language"`
	gerbil             *grbl.Gerbil    `              json:"-"             yaml:"-"`
	stacker            *memory.Stack   `              json:"-"             yaml:"-"`
	logger             *slog.Logger    `              json:"-"             yaml:"-"`
	NavigationInc      float64         `default:"1.0" json:"navigationInc" mapstructure:"navigation_inc" yaml:"navigationInc"`
	RPC                rpcOptions      `              json:"rpc"           yaml:"rpc"`
	Serial             serialOptions   `              json:"serial"        yaml:"serial"`
}

type serialOptions struct {
	PortName serial.PortName `default:"default" json:"portName" mapstructure:"port_name" yaml:"portName"`
	BitRate  int             `default:"115200"  json:"bitRate"  mapstructure:"bit_rate"  yaml:"bitRate"`
}

type rpcOptions struct {
	Clientddr string `default:"0.0.0.0:1324" json:"clientAddr" mapstructure:"client_addr" yaml:"clientAddr"`
	AgentAddr string `default:":1324"        json:"agentAddr"  mapstructure:"agent_addr"  yaml:"agentAddr"`
}

func initApp(ctx context.Context, opts *options, args []string) (*application.Client, error) {
	var program *grbl.Program

	if len(args) > 0 {
		file, err := os.Open(args[0])
		if err != nil {
			return nil, err
		}

		defer func(closer io.Closer) {
			_ = closer.Close()
		}(file)

		program, err = grbl.NewProgram(file)
		if err != nil {
			return nil, err
		}
	}

	app, err := application.NewClient(ctx, opts.stacker, opts.gerbil, program, opts.NavigationInc)
	if err != nil {
		return nil, err
	}

	app.SetLanguage(opts.Language)

	return app, nil
}

func mainCommand() (*cobra.Command, *slog.Logger, error) {
	var forceGRPC bool

	opts := options{
		stacker: memory.New(),
		logger:  slog.Default(),
	}

	gerbil, err := grbl.New()
	if err != nil {
		return nil, opts.logger, err
	}

	opts.availableLanguages = gerbil.AvailableLanguages()
	opts.gerbil = gerbil

	if err := processOptions(&opts); err != nil {
		return nil, nil, err
	}

	output := &cobra.Command{
		Use:   "cnc",
		Short: "CNC monitor",
	}

	output.Flags().BoolVarP(&forceGRPC, "grpc", "", false, "RPC connection")
	output.PersistentFlags().VarP(
		&opts.Language,
		"lang",
		"l",
		fmt.Sprintf("language (available: %s)", joinLang(opts.availableLanguages)),
	)

	output.AddCommand(
		cleanCommand(),
		agentCommand(&opts),
		clientSerialCommand(&opts),
		clientMockCommand(&opts),
		clientRPCCommand(&opts),
		configFileCommand(&opts),
	)

	return output, opts.logger, nil
}

func joinLang(languages []lang.Language) string {
	output := make([]string, len(languages))
	for key, language := range languages {
		output[key] = string(language)
	}

	return strings.Join(output, ", ")
}

func main() {
	cmd, logger, err := mainCommand()
	if err != nil {
		logger.Error("could not initialize commands", "message", err.Error())

		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalCh

		cancel()
	}()

	if err := cmd.ExecuteContext(ctx); err != nil {
		logger.Error("could not execute command", "message", err.Error())

		panic(err)
	}
}

func processOptions(opts *options) error {
	viperConfiguration := viper.New()

	viperConfiguration.SetConfigName("console")
	viperConfiguration.SetConfigType("yaml")
	viperConfiguration.AddConfigPath(".")
	viperConfiguration.AddConfigPath("$HOME/.cnc")
	viperConfiguration.AddConfigPath("/etc/cnc/")
	viperConfiguration.SetEnvPrefix("CNC")
	viperConfiguration.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viperConfiguration.AutomaticEnv()
	err := viperConfiguration.ReadInConfig()
	if err != nil {
		var notFoundErr viper.ConfigFileNotFoundError

		if !errors.As(err, &notFoundErr) {
			return err
		}
	}

	if err := envconfig.Process("cnc", opts); err != nil {
		return err
	}

	if err := viperConfiguration.Unmarshal(opts); err != nil {
		return err
	}

	return nil
}

func configFileCommand(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "save-config",
		Short: "save config gile",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := yaml.Marshal(opts)
			if err != nil {
				return err
			}

			// Sauvegarder dans un fichier
			return os.WriteFile("drill.yaml", data, 0644)
		},
	}
}
