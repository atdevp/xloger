package apiserver

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin/binding"

	"github.com/atdevp/cronlib"
	"github.com/gin-gonic/gin"
	"github.com/xloger/g"
	"github.com/xloger/model"
	"github.com/xloger/sync"
)

// 创建一个定时任务

type Crontab struct {
	SchedueTime string `form:"schedue_time" json:"schedue_time" binding:"required"`
	ID          string `form:"id" json:"id" binding:"required"`
}

func stdout() {
	log.Println("CREATED JOB SUCCESS")
}

func CreateCrontabAPI(c *gin.Context) {
	var err error

	ip := c.Query("ip")
	if ip != g.Config().XlogNodeConfig.IP && ip == "" {
		ResponseHandler(c, 400, "invalid ip", "")
		return
	}

	var input Crontab
	ct := c.Request.Header.Get("Content-Type")
	switch ct {
	case "application/json":
		err = c.BindJSON(&input)
	case "application/x-www-form-urlencoded":
		err = c.MustBindWith(&input, binding.Form)
	default:
		err = errors.New("bind post body error")
	}

	if err != nil {
		ResponseHandler(c, 400, err.Error(), "")
		return
	}

	spec := input.SchedueTime
	f := func() {
		go stdout()
	}
	job, err := cronlib.NewJobModel(spec, f)
	if err != nil {
		ResponseHandler(c, 400, err.Error(), "")
		return
	}

	cron := g.LogFileCrontabSet.Get()
	if err = cron.UpdateJobModel(input.ID, job); err != nil {
		ResponseHandler(c, 400, err.Error(), "")
		return
	}

	ResponseHandler(c, 200, "success", "")
	return
}

func UpdateCrontabAPI(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		ResponseHandler(c, 400, "invalid id", "")
		return
	}

	ip := c.Query("ip")
	if ip != g.Config().XlogNodeConfig.IP {
		ResponseHandler(c, 400, "invalid ip", "")
		return
	}

	ret, err := model.QueryLogFileTaskByAgentIPAndLogsetID(id, ip)
	if err != nil {
		log.Println(err)
		ResponseHandler(c, 400, err.Error(), "")
		return
	}

	cron := g.LogFileCrontabSet.Get()
	hosts := ret.Host
	if len(hosts) == 0 {
		cron.StopService(id)
		ResponseHandler(c, 200, "", "")
		return
	}

	task := ret.Task

	f := sync.GenFunc(task, hosts)
	job, err := cronlib.NewJobModel(task.SchedueTime, f)
	if err != nil {
		ResponseHandler(c, 400, err.Error(), "")
		return
	}

	if err = cron.UpdateJobModel(id, job); err != nil {
		ResponseHandler(c, 400, err.Error(), "")
		return
	}
	ResponseHandler(c, 200, "success", "")
	return

}

func DeleteCrontabAPI(c *gin.Context) {
	var id = c.Query("id")
	if id == "" {
		ResponseHandler(c, 400, "invalid id", "")
		return
	}

	cron := g.LogFileCrontabSet.Get()
	cron.StopService(id)
	ResponseHandler(c, 200, "success", "")
	return
}
