package config

import "flag"

var CliConf CliConfig

type CliConfig struct {
	ConfigFile        string
	OutputPath        string //optional
	LastProcessedPath string //optional
}

func ParseCliFlags() {
	configFile := flag.String("config", "", "the path to the configuration yaml")
	lastProcessedPath := flag.String("last-processed-path", ".", "the path where the last processed date file should be saved")
	output := flag.String("output", "out", "the path where output files should be saved")

	flag.Parse()

	CliConf.ConfigFile = *configFile
	CliConf.LastProcessedPath = *lastProcessedPath
	CliConf.OutputPath = *output
}
