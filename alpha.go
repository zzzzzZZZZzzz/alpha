package main

import (
	"encoding/json"
	toml "github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	bot "github.com/therealfakemoot/hellabot"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	//log "gopkg.in/inconshreveable/log15.v2"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

var logger = log.New()

type server struct {
	Host     string
	Password string
	Nick     string
	SSL      bool
	Channels []string
}

type ircConf struct {
	Servers map[string]server
}

func getCwd() (path string) {
	path, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	return path
}

func getConfigPath() (configPath string) {
	currentUser, _ := user.Current()
	configPath = filepath.Join(currentUser.HomeDir, ".alpha.toml")
	return configPath
}

func loadConfig(path string) (config ircConf, err error) {
	var parsedConf ircConf

	logger.WithFields(log.Fields{
		"filePath": path,
	}).Info("Loading configuration.")

	rawConf, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Error(err)
	}

	_, err = toml.Decode(string(rawConf), &parsedConf)
	if err != nil {
		logger.Error(err)
	}

	return parsedConf, nil

}

func main() {

	logger.Out = os.Stderr

	configPath := kingpin.Flag("config", "Configuration file path.").Short('c').Default(getConfigPath()).ExistingFile()
	debugMode := kingpin.Flag("debug", "Enable debugging output.").Bool()

	kingpin.Parse()

	if *debugMode {
		logger.Level = log.DebugLevel
	} else {
		logger.Level = log.InfoLevel
	}

	logger.WithFields(log.Fields{
		"configPath": *configPath,
		"debugMode":  *debugMode,
	}).Info("Loading up muh bots.")

	conf, err := loadConfig(*configPath)
	if err != nil {
		logger.Error(err)
	}

	jsonConf, _ := json.Marshal(conf)

	logger.WithFields(log.Fields{
		"parsedConfig": string(jsonConf),
	}).Info("Loaded config file.")

	for _, server := range conf.Servers {
		ssl := func(bot *bot.Bot) {
			bot.SSL = server.SSL
		}

		channels := func(bot *bot.Bot) {
			bot.Channels = server.Channels
		}

		bot, err := bot.NewBot(server.Host, server.Nick, ssl, channels)
		if err != nil {
			logger.Error(err)
		}
		logger.WithFields(log.Fields{
			"server":  server.Host,
			"botNick": server.Nick,
		}).Info("Starting bot.")

		bot.Run() //This is a blocking call. This will not work if multiple bots are defined.
	}

}
