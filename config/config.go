package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var Instance *Config

type Config struct {
	MasterQQ              string  `yaml:"MasterQQ"`
	LoginQQ               int64   `yaml:"LoginQQ"`
	ImageStorePath        string  `yaml:"ImageStorePath"`
	LogFile               string  `yaml:"LogFile"`
	QQChatID              string  `yaml:"QQChatID"`
	QQChatKey             string  `yaml:"QQChatKey"`
	GuildFlagRaceQQGroup  []int   `yaml:"GuildFlagRaceQQGroup"`
	OfficialNoticeQQGroup []int64 `yaml:"OfficialNoticeQQGroup"`
	QAEditQQGroup         []int   `yaml:"QAEditQQGroup"`
	MVPMonitorGroups      []int64 `yaml:"MVPMonitorGroups"`
	MVPNoticeGroups       []int64 `yaml:"MVPNoticeGroups"`
}

func Init(filename string) *Config {
	Instance = &Config{}
	if yamlFile, err := ioutil.ReadFile(filename); err != nil {
		logrus.Error(err)
	} else if err = yaml.Unmarshal(yamlFile, Instance); err != nil {
		logrus.Error(err)
	}
	return Instance
}
