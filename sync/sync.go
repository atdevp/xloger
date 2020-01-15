package sync

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/atdevp/cronlib"
	"github.com/xloger/g"
	m "github.com/xloger/model"
	"github.com/xloger/sftp"
	w "github.com/xloger/worker"
)

func SyncLogFileTask() {

	var (
		err error
		ip  string
	)

	if ip = g.Config().XlogNodeConfig.IP; ip == "" {
		panic("ERROR agent_ip is null")
	}

	if err = syncLofileTask(ip); err != nil {
		panic(err)
	}
}

func syncLofileTask(ip string) error {
	var (
		tasks []m.LogFileTask
		err   error
	)

	for {
		tasks, err = m.QueryLogFileTaskByAgentIP(ip)
		if err != nil {
			return err
		}
		if len(tasks) == 0 {
			log.Printf("sync logfiletask from db is null,please sleep %d second", g.Config().SyncDBInterval)
			time.Sleep(g.Config().SyncDBInterval * time.Second)
		} else {
			log.Println("sync logfiletask from db success")
			break
		}
	}

	initCrontab(tasks)
	return nil

}

func initCrontab(tasks []m.LogFileTask) {
	var (
		cron = g.LogFileCrontabSet.Get()
		f    = func() {}
	)

	for k := range tasks {
		t := &tasks[k]
		log.Println(t.LogsetID, t.AgentRole, t.AgentIP)
		hosts, err := m.QueryLogFileHostByLogsetIDAndRole(t.LogsetID, t.AgentRole)
		if err != nil {
			log.Println(err)
		}
		f = GenFunc(t, hosts)
		job, err := cronlib.NewJobModel(t.SchedueTime, f)
		if err != nil {
			log.Println(err)
			continue
		}
		if err = cron.Register(t.LogsetID, job); err != nil {
			log.Println(err)
		}
		log.Printf("task %s %s register success", t.LogsetID, t.LogsetName)
	}
	cron.Start()
	log.Println("crontab is running")
	cron.Wait()
}

func GenFunc(t *m.LogFileTask, hosts []m.LogFileHost) func() {

	if t.AgentRole == "master" {
		f := func() {
			wg := &sync.WaitGroup{}
			for k := range hosts {
				h := &hosts[k]
				wg.Add(1)
				filePath := h.FilePath
				worker := &w.Worker{
					Task: t,
					SSH: &sftp.SSHConfig{
						Host:        h.Host,
						Port:        22,
						User:        "root",
						RSAFilePath: g.Config().XlogNodeConfig.RsaFilepath,
						Passwd:      h.Passwd,
					},
				}
				go worker.Pull(filePath, wg)
			}
			wg.Wait()
			// 判断任务状态，通知BI  isOrNotNeedExecSlaveTask()
		}
		return f

	}

	f := func() {
		fail, err := isOrNotNeedExecSlaveTask(t.LogsetID, "master", hosts)
		if err != nil {
			log.Println(err)
			return
		}

		if len(fail) == 0 && err == nil {
			log.Println("slave not need to exec task")
			return
		}

		wg := &sync.WaitGroup{}
		for k := range fail {
			h := &fail[k]
			wg.Add(1)
			filePath := h.FilePath
			worker := &w.Worker{
				Task: t,
				SSH: &sftp.SSHConfig{
					Host:        h.Host,
					Port:        22,
					User:        "root",
					RSAFilePath: g.Config().XlogNodeConfig.RsaFilepath,
					Passwd:      h.Passwd,
				},
			}
			go worker.Pull(filePath, wg)
		}
		wg.Wait()
		// // 判断任务状态，告警

	}
	return f

}

func isOrNotNeedExecSlaveTask(logsetID, role string, hosts []m.LogFileHost) ([]m.LogFileHost, error) {
	h := time.Now().Hour()
	ret, err := m.QueryAgentNodeTaskStatus(logsetID, role, h)

	log.Println(ret)

	if err != nil {
		if err.Error() == "no_host" {
			return []m.LogFileHost{}, nil
		}
		return nil, err
	}

	fail := ret["fail"]
	vFail, ok := fail.([]string)
	if !ok {
		return nil, errors.New("asset fail fail")
	}
	if len(vFail) == 0 {
		return []m.LogFileHost{}, nil
	}

	failHost := make([]m.LogFileHost, 0)
	for _, ip := range vFail {
		for _, h := range hosts {
			if ip == h.Host {
				failHost = append(failHost, h)
			}
		}
	}
	return failHost, nil
}
