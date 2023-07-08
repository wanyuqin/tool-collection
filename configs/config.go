package configs

import (
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/wanyuqin/tool-collection/backend/x/xfile"
	"github.com/wanyuqin/tool-collection/logger"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var (
	DownloadPathBlankErr  = errors.New("download path is blank")
	DownloadPathNotDirErr = errors.New("download path is not dir")
)

var (
	DefaultConfigName   = "settings.yaml"
	DefaultDownloadPath = "tools_collection"
)

type Config struct {
	Download DownloadConfig `json:"download" yaml:"download"`
}

type DownloadConfig struct {
	Path string `json:"path" yaml:"path"`
}

func GetConfig() Config {
	return LoadConfig()
}

func LoadConfig() Config {
	config := Config{}
	body, err := ReadConfig()
	if err != nil {
		logger.Error(err.Error())
		return config
	}

	err = yaml.Unmarshal(body, &config)
	if err != nil {
		logger.Error(err.Error())
		return config
	}
	logger.Debug(fmt.Sprintf("%v\n", config))
	return config
}

func SaveDownloadSettings(downloadConfig DownloadConfig) error {
	cfg := LoadConfig()
	cfg.Download = downloadConfig

	body, err := yaml.Marshal(cfg)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	configPath, err := GetConfigPath()
	if err != nil {
		logger.Errorf("get config path failed: %s", err)
		return err
	}
	file, err := os.OpenFile(configPath, os.O_TRUNC|os.O_WRONLY, 0666)
	defer file.Close()

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_, err = file.Write(body)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil

}

// InitConfigFile 初始化配置文件
func InitConfigFile(rootPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	downloadPath := filepath.Join(homeDir, DefaultDownloadPath)
	err = InitDownloadPath(downloadPath)
	if err != nil {
		return err
	}
	// 创建配置文件
	cfgPath := filepath.Join(rootPath, DefaultConfigName)
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		cfgFile, err := os.Create(cfgPath)
		if err != nil {
			fmt.Printf("create config file failed: %v\n", err)
			return nil
		}
		defer cfgFile.Close()
		cfg := Config{
			Download: DownloadConfig{
				Path: downloadPath,
			},
		}

		cfgByte, err := yaml.Marshal(cfg)
		if err != nil {
			fmt.Printf("yaml marshal config failed: %v\n", err)
			return nil
		}

		_, err = cfgFile.Write(cfgByte)
		if err != nil {
			fmt.Printf("write cfg failed: %v\n", err)
			return nil

		}

	}
	return nil
}

func InitDownloadPath(downloadPath string) error {

	return xfile.CreateDirIfNotExist(downloadPath)
}

func (c *Config) CheckDownloadPath() error {
	if strutil.IsBlank(c.Download.Path) {
		return DownloadPathBlankErr
	}

	stat, err := os.Stat(c.Download.Path)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return DownloadPathNotDirErr
	}
	return nil
}

func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".tools_collection", DefaultConfigName), nil
}

func ReadConfig() ([]byte, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}
