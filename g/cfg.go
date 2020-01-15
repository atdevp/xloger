package g

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/atdevp/devlib/file"
)

type XlogNode struct {
	IP          string `json:"ip"`
	Port        int    `json:"port"`
	StoragePath string `json:"storage_path"`
	RsaFilepath string `json:"rsa_filepath"`
}

type Mysql struct {
	Dsn      string `json:"dsn"`
	Idle     int    `json:"idle"`
	Max      int    `json:"max"`
	LogModel bool   `json:"log_model"`
}

type Hdfs struct {
	NameNodes      []string `json:"namenodes"`
	CustomizeConf  string   `json:"customize_conf"`
	DefaultConf    string   `json:"default_conf"`
	DefaultKrbConf string   `json:"default_krb_conf"`
}

type Alarm struct {
	API      string         `json:"api"`
	Users    map[string]int `json:"users"`
	Interval int            `json:"interval"`
	Max      int            `json:"max"`
}

type Upstream struct {
	API string `json:"api"`
}

type GlobalConfig struct {
	XlogNodeConfig *XlogNode     `json:"xlog_node"`
	MysqlConfig    *Mysql        `json:"mysql"`
	HdfsConfig     *Hdfs         `json:"hdfs"`
	AlarmConfig    *Alarm        `json:"alarm"`
	BIConfig       *BI           `json:"bi"`
	SyncDBInterval time.Duration `json:"sync_db_interval"`
}

var (
	lock   = new(sync.RWMutex)
	config *GlobalConfig
)

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("please use -c specify config file")
	}

	if !file.IsExist(cfg) {
		log.Fatalf("%s is not exist", cfg)
	}

	if !file.IsFile(cfg) {
		log.Fatalf("%s is not file type", cfg)
	}

	content, err := file.ToBytes(cfg)
	if err != nil {
		log.Fatalf("read config file %s err: %v", cfg, err)
	}

	var c GlobalConfig
	err = json.Unmarshal(content, &c)
	if err != nil {
		log.Fatalln("parse config file err: ", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c
	log.Printf("load config file %s success", cfg)
}

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func GetHaoodpConfDir() (string, error) {
	dir := os.Getenv("HADOOP_CONF_DIR")
	if dir != "" {
		return dir, nil
	}

	hadoopHome := os.Getenv("HADOOP_HOME")
	if hadoopHome != "" {
		return filepath.Join(hadoopHome, "conf"), nil
	}

	customizeConf := Config().HdfsConfig.CustomizeConf

	if exist := file.IsExist(customizeConf); exist {
		return customizeConf, nil
	}

	return Config().HdfsConfig.DefaultConf, errors.New("hadoop default conf dir is not existant")

}

func GetKrb5Conf() (string, error) {
	var cfg string
	cfg = os.Getenv("KRB5_CONFIG")
	if cfg != "" {
		return cfg, nil
	}

	defaultCfg := Config().HdfsConfig.DefaultKrbConf
	if exist := file.IsExist(defaultCfg); exist {
		return defaultCfg, nil
	}

	return "", errors.New("kerbos krb5.conf is not existant")

}
