package cmd

import (
	"strings"

	"github.com/skyinglyh1/uniswap_v1_test/config"
	"github.com/urfave/cli"
)

var (
	LogLevelFlag = cli.UintFlag{
		Name:  "loglevel",
		Usage: "Set the log level to `<level>` (0~6). 0:Trace 1:Debug 2:Info 3:Warn 4:Error 5:Fatal 6:MaxLevel",
		Value: config.DEFAULT_LOG_LEVEL,
	}

	OntPwd = cli.StringFlag{
		Name:  "ontpwd",
		Usage: "Password for ontology wallet",
		Value: "",
	}

	AlliaPwd = cli.StringFlag{
		Name:  "alliapwd",
		Usage: "Password for poly wallet",
		Value: "",
	}

	ConfigPathFlag = cli.StringFlag{
		Name:  "cliconfig",
		Usage: "Server config file `<path>`",
		Value: config.DEFAULT_CONFIG_FILE_NAME,
	}
)

//GetFlagName deal with short flag, and return the flag name whether flag name have short name
func GetFlagName(flag cli.Flag) string {
	name := flag.GetName()
	if name == "" {
		return ""
	}
	return strings.TrimSpace(strings.Split(name, ",")[0])
}
