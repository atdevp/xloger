package model

import (
	"errors"
	"fmt"
	"time"

	_ "github.com/xloger/g"

	"github.com/xloger/db"
)

type LogFileTask struct {
	ID          int64     `json:"id" gorm:"column: id"`
	LogsetID    string    `json:"logset_id" gorm:"column:logset_id"`
	LogsetName  string    `json:"logset_name" gorm:"column:logset_name"`
	LogsetSplit string    `json:"logset_split" gorm:"column:logset_split"`
	AgentIP     string    `json:"agent_ip" gorm:"column:agent_ip"`
	AgentRole   string    `json:"agent_role" gorm:"column:agent_role"`
	SchedueTime string    `json:"schedue_time" gorm:"column:schedue_time"`
	TempPath    string    `json:"temp_path" gorm:"column:temp_path"`
	HdfsPath    string    `json:"hdfs_path" gorm:"column:hdfs_path"`
	Status      int       `json:"status" gorm:"column:status"`
	Utime       time.Time `json:"utime" gorm:"column:utime"`
}

func QueryLogFileTaskByAgentIP(ip string) ([]LogFileTask, error) {
	con := db.Con().Xlog
	sql := fmt.Sprintf("SELECT * FROM logfile_task WHERE agent_ip='%s'", ip)

	tasks := make([]LogFileTask, 0)
	if err := con.Raw(sql).Scan(&tasks).Error; err != nil {
		return tasks, err
	}
	return tasks, nil
}

type LogFileDetailTask struct {
	Task *LogFileTask
	Host []LogFileHost
}

func QueryLogFileTaskByAgentIPAndLogsetID(logsetid, ip string) (*LogFileDetailTask, error) {

	var res = new(LogFileDetailTask)

	con := db.Con().Xlog
	sql := fmt.Sprintf("SELECT * FROM logfile_task WHERE agent_ip='%s' AND logset_id='%s'", ip, logsetid)
	var task LogFileTask
	if err := con.Raw(sql).First(&task).Error; err != nil {
		return res, err
	}
	res.Task = &task

	var hosts []LogFileHost
	sql = fmt.Sprintf("SELECT * FROM logfile_host WHERE logset_id='%s' AND agent_role='%s'", logsetid, task.AgentRole)
	if err := con.Raw(sql).Scan(&hosts).Error; err != nil {
		return res, err
	}
	res.Host = hosts

	return res, nil
}

type LogFileHost struct {
	ID        int64     `json:"id" gorm:"column:id"`
	LogsetID  string    `json:"logset_id" gorm:"column:logset_id"`
	AgentRole string    `json:"agent_role" gorm:"column:agent_role"`
	Host      string    `json:"host" gorm:"column:host"`
	Passwd    string    `json:"passwd" gorm:"column:passwd"`
	FilePath  string    `json:"file_path" gorm:"column:file_path"`
	PlsTime   time.Time `json:"pl_stime" gorm:"column:pl_stime"`
	PleTime   time.Time `json:"pl_etime" gorm:"column:pl_etime"`
	PlState   string    `json:"pl_state" gorm:"column:pl_state"`
	PusTime   time.Time `json:"pu_stime" gorm:"column:pu_stime"`
	PueTime   time.Time `json:"pu_etime" gorm:"column:pu_etime"`
	PuState   string    `json:"pu_state" gorm:"column:pu_state"`
}

func QueryLogFileHostByLogsetIDAndRole(id, role string) ([]LogFileHost, error) {
	con := db.Con().Xlog

	sql := fmt.Sprintf("SELECT * FROM logfile_host WHERE logset_id='%s' AND agent_role='%s'", id, role)

	var hosts = make([]LogFileHost, 0)
	if err := con.Raw(sql).Scan(&hosts).Error; err != nil {
		return nil, err
	}
	return hosts, nil
}

func QueryAgentNodeTaskStatus(logsetId, role string, hour int) (map[string]interface{}, error) {
	var (
		err   error
		total int
		succ  []string
		fail  []string
	)

	con := db.Con().Xlog

	sql := fmt.Sprintf("SELECT * FROM logfile_host WHERE logset_id='%s' AND agent_role='%s' ", logsetId, role)
	var hosts []LogFileHost
	if err = con.Raw(sql).Scan(&hosts).Error; err != nil {
		return nil, err
	}
	total = len(hosts)
	if total == 0 {
		return nil, errors.New("no_host")
	}

	for k := range hosts {
		h := &hosts[k]
		pleh, pueh := h.PleTime.Hour(), h.PueTime.Hour()
		plst, pust := h.PlState, h.PuState

		if hour == pleh && hour == pueh && plst == "success" && pust == "success" {
			succ = append(succ, h.Host)
		} else {
			fail = append(fail, h.Host)
		}
	}

	ret := map[string]interface{}{
		"succ":  succ,
		"fail":  fail,
		"total": total,
	}
	return ret, nil

}
