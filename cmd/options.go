package main

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/landru29/cnc-serial/internal/gcode/grbl"
	"github.com/landru29/cnc-serial/internal/lang"
	"github.com/landru29/cnc-serial/internal/stack/memory"
	"github.com/landru29/cnc-serial/internal/transport/serial"
	"github.com/spf13/viper"
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
	PortName serial.PortName `default:"default" json:"portName" mapstructure:"portname" yaml:"portName"`
	BitRate  int             `default:"115200"  json:"bitRate"  mapstructure:"bit_rate"  yaml:"bitRate"`
}

type rpcOptions struct {
	ClientAddr string `default:"0.0.0.0:1324" json:"clientAddr" mapstructure:"client_addr" yaml:"clientAddr"`
	AgentAddr  string `default:":1324"        json:"agentAddr"  mapstructure:"agent_addr"  yaml:"agentAddr"`
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
