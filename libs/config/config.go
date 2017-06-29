package config

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	enc "github.com/hiromaily/golibs/cipher/encryption"
	u "github.com/hiromaily/golibs/utils"
	"io/ioutil"
	"os"
)

/* singleton */
var (
	conf         *Config
	tomlFileName = "./config/config.toml"
)

// Config is of root
type Config struct {
	Environment string
	Server      *ServerConfig
	Proxy       ProxyConfig
	API         *APIConfig
	MySQL       *MySQLConfig
	Redis       *RedisConfig
}

// ServerConfig is for web server
type ServerConfig struct {
	Scheme    string          `toml:"scheme"`
	Host      string          `toml:"host"`
	Port      int             `toml:"port"`
	Log       LogConfig       `toml:"log"`
	Session   SessionConfig   `toml:"session"`
	BasicAuth BasicAuthConfig `toml:"basic_auth"`
}

// LogConfig is for Log
type LogConfig struct {
	Level uint8  `toml:"level"`
	Path  string `toml:"path"`
}

// SessionConfig is for Session
type SessionConfig struct {
	Name     string `toml:"name"`
	Key      string `toml:"key"`
	MaxAge   int    `toml:"max_age"`
	Secure   bool   `toml:"secure"`
	HTTPOnly bool   `toml:"http_only"`
}

// BasicAuthConfig is for Basic Auth
type BasicAuthConfig struct {
	User string `toml:"user"`
	Pass string `toml:"pass"`
}

// ProxyConfig is for base of Reverse Proxy Server
type ProxyConfig struct {
	Mode   uint8             `toml:"mode"` //0:off, 1:go, 2,nginx
	Server ProxyServerConfig `toml:"server"`
}

// ProxyServerConfig is for Reverse Proxy Server
type ProxyServerConfig struct {
	Scheme  string    `toml:"scheme"`
	Host    string    `toml:"host"`
	Port    int       `toml:"port"`
	WebPort []int     `toml:"web_port"`
	Log     LogConfig `toml:"log"`
}

// APIConfig is for Rest API
type APIConfig struct {
	Ajax   bool          `toml:"only_ajax"`
	CORS   *CORSConfig   `toml:"cors"`
	Header *HeaderConfig `toml:"header"`
	JWT    *JWTConfig    `toml:"jwt"`
}

// CORSConfig is for CORS
type CORSConfig struct {
	Enabled     bool     `toml:"enabled"`
	Origins     []string `toml:"origins"`
	Headers     []string `toml:"headers"`
	Methods     []string `toml:"methods"`
	Credentials bool     `toml:"credentials"`
}

// HeaderConfig is added original header for authentication
type HeaderConfig struct {
	Enabled bool   `toml:"enabled"`
	Header  string `toml:"header"`
	Key     string `toml:"key"`
}

// JWTConfig is for JWT Auth
type JWTConfig struct {
	Mode       uint8  `toml:"mode"`
	Secret     string `toml:"secret_code"`
	PrivateKey string `toml:"private_key"`
	PublicKey  string `toml:"public_key"`
}

// MySQLConfig is for MySQL Server
type MySQLConfig struct {
	*MySQLContentConfig
}

// MySQLContentConfig is for MySQL Server
type MySQLContentConfig struct {
	Encrypted bool   `toml:"encrypted"`
	Host      string `toml:"host"`
	Port      uint16 `toml:"port"`
	DbName    string `toml:"dbname"`
	User      string `toml:"user"`
	Pass      string `toml:"pass"`
}

// RedisConfig is for Redis Server
type RedisConfig struct {
	Encrypted bool   `toml:"encrypted"`
	Host      string `toml:"host"`
	Port      uint16 `toml:"port"`
	Pass      string `toml:"pass"`
	Session   bool   `toml:"session"`
}

var checkTOMLKeys = [][]string{
	{"environment"},
	{"server", "scheme"},
	{"server", "host"},
	{"server", "port"},
	{"server", "log", "level"},
	{"server", "log", "path"},
	{"server", "session", "name"},
	{"server", "session", "key"},
	{"server", "session", "max_age"},
	{"server", "session", "secure"},
	{"server", "session", "http_only"},
	{"server", "basic_auth", "user"},
	{"server", "basic_auth", "pass"},
	{"api", "only_ajax"},
	{"api", "cors", "enabled"},
	{"api", "cors", "origins"},
	{"api", "cors", "headers"},
	{"api", "cors", "methods"},
	{"api", "cors", "credentials"},
	{"api", "header", "enabled"},
	{"api", "header", "header"},
	{"api", "header", "key"},
	{"api", "jwt", "mode"},
	{"api", "jwt", "secret_code"},
	{"api", "jwt", "private_key"},
	{"api", "jwt", "public_key"},
	{"mysql", "encrypted"},
	{"mysql", "host"},
	{"mysql", "port"},
	{"mysql", "dbname"},
	{"mysql", "user"},
	{"mysql", "pass"},
	{"redis", "encrypted"},
	{"redis", "host"},
	{"redis", "port"},
	{"redis", "pass"},
	{"redis", "session"},
}

func init() {
	//tomlFileName = os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-gin-wrapper/configs/settings.toml"
	pwd, _ := os.Getwd()
	tomlFileName = pwd + "/data/config.toml"
}

//check validation of config
func validateConfig(conf *Config, md *toml.MetaData) error {
	//for protection when debugging on non production environment
	var errStrings []string

	//Check added new items on toml
	// environment
	//if !md.IsDefined("environment") {
	//	errStrings = append(errStrings, "environment")
	//}

	format := "[%s]"
	inValid := false
	for _, keys := range checkTOMLKeys {
		if !md.IsDefined(keys...) {
			switch len(keys) {
			case 1:
				format = "[%s]"
			case 2:
				format = "[%s] %s"
			case 3:
				format = "[%s.%s] %s"
			default:
				//invalid check string
				inValid = true
				break
			}
			keysIfc := u.SliceStrToInterface(keys)
			errStrings = append(errStrings, fmt.Sprintf(format, keysIfc...))
		}
	}

	// Error
	if inValid {
		return errors.New("Error: Check Text has wrong number of parameter")
	}
	if len(errStrings) != 0 {
		return fmt.Errorf("Error: There are lacks of keys : %#v \n", errStrings)
	}

	return nil
}

// load configfile
func loadConfig(fileName string) (*Config, error) {
	if fileName != "" {
		tomlFileName = fileName
	}

	d, err := ioutil.ReadFile(tomlFileName)
	if err != nil {
		return nil, fmt.Errorf(
			"Error reading %s: %s", tomlFileName, err)
	}

	var config Config
	md, err := toml.Decode(string(d), &config)
	if err != nil {
		return nil, fmt.Errorf(
			"Error parsing %s: %s(%v)", tomlFileName, err, md)
	}

	//check validation of config
	err = validateConfig(&config, &md)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// New is create instance
func New(fileName string, cipherFlg bool) *Config {
	var err error
	if conf == nil {
		conf, err = loadConfig(fileName)
	}
	if err != nil {
		panic(err)
	}

	if cipherFlg {
		Cipher()
	}

	return conf
}

// GetConf is to get config instance
func GetConf() *Config {
	var err error
	if conf == nil {
		conf, err = loadConfig("")
	}
	if err != nil {
		panic(err)
	}

	return conf
}

// SetTOMLPath is to set toml file path
func SetTOMLPath(path string) {
	tomlFileName = path
}

// Cipher is to decrypt crypted string on config
func Cipher() {
	crypt := enc.GetCrypt()

	if conf.MySQL.Encrypted {
		c := conf.MySQL
		c.Host, _ = crypt.DecryptBase64(c.Host)
		c.DbName, _ = crypt.DecryptBase64(c.DbName)
		c.User, _ = crypt.DecryptBase64(c.User)
		c.Pass, _ = crypt.DecryptBase64(c.Pass)
	}

	if conf.Redis.Encrypted {
		c := conf.Redis
		c.Host, _ = crypt.DecryptBase64(c.Host)
		c.Pass, _ = crypt.DecryptBase64(c.Pass)
	}
}
