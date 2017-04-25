package configmanager

//load config.toml
import (
	"os"
	"io/ioutil"
	"github.com/naoina/toml"
)

type server struct {
	Address    string
	Password   string
	CipherType string
}

type ProxyTarget struct {
	Listen string
	Target string
	Host   string
}

type Config struct {
	Server       server
	ProxyTargets []ProxyTarget
}

// read toml and return config
func MustLoad() *Config {
	var configPath string
	if localConfigPath := "./config.toml"; Exists(localConfigPath) {
		configPath = localConfigPath
	} else {
		configPath = os.Getenv("HOME") + "/.MonkeyProxy/config.toml"
	}

	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	config := new(Config)
	err = toml.Unmarshal(configData, config)

	if err != nil {
		panic(err)
	}
	return config
}
