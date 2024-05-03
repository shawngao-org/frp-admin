package config

import (
	"errors"
	"fmt"
	"frp-admin/logger"
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
	"time"
)

var (
	configMutex sync.Mutex
	Conf        Config
)

type Template struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type Config struct {
	Server struct {
		Ip           string `yaml:"ip"`
		Port         uint64 `yaml:"port"`
		FrontEndAddr string `yaml:"front-end-addr"`
	} `yaml:"server"`
	Mail struct {
		Host     string     `yaml:"host"`
		Port     uint64     `yaml:"port"`
		Mail     string     `yaml:"mail"`
		NickName string     `yaml:"nick-name"`
		FromMail string     `yaml:"from-mail"`
		Password string     `yaml:"password"`
		Template []Template `yaml:"template"`
	}
	Database struct {
		Mysql struct {
			Host     string `yaml:"host"`
			Port     uint64 `yaml:"port"`
			Db       string `yaml:"db"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
		} `yaml:"mysql"`
		Redis struct {
			Host     string `yaml:"host"`
			Port     uint64 `yaml:"port"`
			Db       int    `yaml:"db"`
			Password string `yaml:"password"`
			PoolSize int    `yaml:"pool-size"`
			Timeout  uint64 `yaml:"timeout"`
		} `yaml:"redis"`
	} `yaml:"database"`
	Security struct {
		Password struct {
			Method string `yaml:"method"`
			Secret string `yaml:"secret"`
			Cost   int    `yaml:"cost"`
		} `yaml:"password"`
		Jwt struct {
			Secret  string `yaml:"secret"`
			Timeout int64  `yaml:"timeout"`
		} `yaml:"jwt"`
		Rsa struct {
			Public  string `yaml:"public"`
			Private string `yaml:"private"`
		} `yaml:"rsa"`
		Totp struct {
			Issuer string `yaml:"issuer"`
		} `yaml:"totp"`
	} `yaml:"security"`
	Nacos struct {
		Enable    bool   `yaml:"enable"`
		Ip        string `yaml:"ip"`
		Port      uint64 `yaml:"port"`
		Username  string `yaml:"username"`
		Password  string `yaml:"password"`
		Namespace string `yaml:"namespace"`
		Group     string `yaml:"group"`
		DataId    string `yaml:"dataId"`
		Timeout   uint64 `yaml:"timeout"`
		Loglevel  string `yaml:"loglevel"`
	} `yaml:"nacos"`
	Develop bool `yaml:"develop"`
}

func parseContent2Config(content []byte) Config {
	var config Config
	err := yaml.Unmarshal(content, &config)
	if err != nil {
		readFileErrLogImpl(err)
		os.Exit(-1)
	}
	return config
}

func GetConfig() {
	configFileName := "config.yml"
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		logger.LogErr("Configuration file is not exist !")
		readFileErrLogImpl(err)
		os.Exit(-1)
	}
	content, err := os.ReadFile(configFileName)
	if err != nil {
		readFileErrLogImpl(err)
		os.Exit(-1)
	}
	config := parseContent2Config(content)
	logger.LogSuccess("Configuration file '%s'", configFileName)
	if config.Develop {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	if config.Nacos.Enable {
		logger.LogInfo("Configuration file mode: Nacos unified configuration center")
		nacosMain(config)
		return
	}
	logger.LogInfo("Profile Mode: Local config file")
	configMutex.Lock()
	Conf = config
	configMutex.Unlock()
}

func readFileErrLogImpl(err error) {
	logger.LogErr("Configuration file '%s'")
	logger.LogErr("%s", err)
	os.Exit(-1)
}

func nacosMain(config Config) {
	sc := []constant.ServerConfig{{
		IpAddr: config.Nacos.Ip,
		Port:   config.Nacos.Port,
	}}
	cc := constant.ClientConfig{
		NamespaceId:         config.Nacos.Namespace,
		TimeoutMs:           config.Nacos.Timeout,
		NotLoadCacheAtStart: true,
		LogDir:              "log",
		LogLevel:            config.Nacos.Loglevel,
		Username:            config.Nacos.Username,
		Password:            config.Nacos.Password,
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	configClientErrHandle(err)
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: config.Nacos.DataId,
		Group:  config.Nacos.Group,
	})
	configClientErrHandle(err)
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: config.Nacos.DataId,
		Group:  config.Nacos.Group,
		OnChange: func(namespace, group, dataId, data string) {
			logger.LogInfo("The configuration file has changed...")
			logger.LogInfo("Group: %s, Data Id: %s", group, dataId)
			configMutex.Lock()
			Conf = parseContent2Config([]byte(data))
			configMutex.Unlock()
			err := restart()
			if err != nil {
				logger.LogErr(err.Error())
			}
		},
	})
	if err != nil {
		logger.LogErr("Failed to listen config.")
	}
	configMutex.Lock()
	Conf = parseContent2Config([]byte(content))
	configMutex.Unlock()
}

func executeRestart() error {
	cmd := exec.Command("sh", "-c", "sleep 3 && exit 0")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %s", err)
	}
	return nil
}

func restart() error {
	if runtime.GOOS == "linux" {
		return errors.New("config hot reload is not support Windows OS")
	}
	logger.LogWarn("The server is restarting, please wait a few seconds...")
	err := executeRestart()
	if err != nil {
		return fmt.Errorf("failed to execute restart comand: %s", err)
	}
	time.Sleep(1 * time.Second)
	binary, err := exec.LookPath(os.Args[0])
	if err != nil {
		return fmt.Errorf("can't get executable binary: %s", err)
	}
	err = syscall.Exec(binary, os.Args, os.Environ())
	if err != nil {
		return fmt.Errorf("failed to restart server: %s", err)
	}
	return nil
}

func configClientErrHandle(err error) {
	if err != nil {
		logger.LogErr(err.Error())
		os.Exit(-1)
	}
}
