package conf

import (
    "log"
    "flag"
    "github.com/BurntSushi/toml"
)

type Config struct {
    Port                int
    Output_cron         string
    InitDataSrc_cron    string
    OutputFile          string
    IPFileURL           string
    IncludeGroup        []string
}

func ParseConfig() *Config {
    conf := Config{}

    confFile := flag.String("f", "scanner.conf", "config file")
    flag.Parse()


    if _, err := toml.DecodeFile(*confFile, &conf); err != nil { 
        log.Fatal(err)
    }

    return &conf
}
