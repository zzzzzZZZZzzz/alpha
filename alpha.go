package main

import (
	toml "github.com/BurntSushi/toml"
	bot "github.com/whyrusleeping/hellabot"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type server struct {
	Host     string
	Port     integer
	Password string
	Nick     string
}

type serverList struct {
	Servers []server
}

func getCwd() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

func getConfigPath() (configPath string) {
	currentUser, _ := user.Current()
	configPath = filepath.Join(currentUser.HomeDir, ".alpha.toml")
	return configPath
}

func main() {

	configPath = kingpin.Flag("config", "Configuration file path.").Short('c').Default(getConfigPath()).ExistingFile()
	debugMode = kingpin.Flag("debug", "Enable debugging output.").Bool()

	kingpin.Parse()

}
