package main

import (
	"os"

	"fmt"
	"github.com/hnode/cmd/utils"
	"github.com/hnode/common/log"
	"github.com/hnode/config"
	"github.com/naoina/toml"
	"gopkg.in/urfave/cli.v1"
	"io"
	"reflect"
	"unicode"
)

var (
	dumpConfigCommand = cli.Command{
		Action:      utils.MigrateFlags(dumpConfig),
		Name:        "dumpconfig",
		Usage:       "Show configuration values",
		ArgsUsage:   "",
		Flags:       append(append(nodeFlags, rpcFlags...)),
		Category:    "MISCELLANEOUS COMMANDS",
		Description: `The dumpconfig command shows configuration values.`,
	}

	configFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TOML configuration file",
	}
)

// These settings ensure that TOML keys use the same names as Go struct fields.
var tomlSettings = toml.Config{
	NormFieldName: func(rt reflect.Type, key string) string {
		return key
	},
	FieldToKey: func(rt reflect.Type, field string) string {
		return field
	},
	MissingField: func(rt reflect.Type, field string) error {
		link := ""
		if unicode.IsUpper(rune(rt.Name()[0])) && rt.PkgPath() != "main" {
			link = fmt.Sprintf(", see https://godoc.org/%s#%s for available fields", rt.PkgPath(), rt.Name())
		}
		return fmt.Errorf("field '%s' is not defined in %s%s", field, rt.String(), link)
	},
}

// dumpConfig is the dumpconfig command.
func dumpConfig(ctx *cli.Context) error {
	cfg := MakeConfigNode(ctx)
	comment := ""

	out, err := tomlSettings.Marshal(&cfg)
	if err != nil {
		return err
	}
	io.WriteString(os.Stdout, comment)
	os.Stdout.Write(out)
	return nil
}

func MakeConfigNode(ctx *cli.Context) *config.HpbConfig {
	// Load defaults config
	cfg := config.New()
	if cfg == nil {
		log.Error("Get Hpb config fail, so exit")
		os.Exit(1)
	}
	// Apply flags.
	utils.SetConfig(ctx, cfg)
	//set cfg version
	cfg.Node.Version = config.VersionWithCommit(GitCommit)

	utils.SetTxPool(ctx, &cfg.TxPool)

	if ctx.GlobalIsSet(utils.HpbStatsURLFlag.Name) {
		cfg.HpbStats.URL = ctx.GlobalString(utils.HpbStatsURLFlag.Name)
	}

	return cfg
}
