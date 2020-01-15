package g

import (
	"github.com/atdevp/cronlib"
	"sync"
)

var (
	LogFileCrontabSet = &SafeCollectLogFileCrontab{}
)

type SafeCollectLogFileCrontab struct {
	sync.RWMutex
	crontab *cronlib.CronSchduler
}

func (this *SafeCollectLogFileCrontab) New() {
	this.RLock()
	defer this.RUnlock()
	this.crontab = cronlib.New()
}

func (this *SafeCollectLogFileCrontab) Get() *cronlib.CronSchduler {
	this.RLock()
	defer this.RUnlock()
	return this.crontab
}
