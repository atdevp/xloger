package worker

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/atdevp/devlib/file"
	"github.com/pkg/sftp"
	"github.com/xloger/db"
	"github.com/xloger/g"
	"github.com/xloger/hadoop"
	"github.com/xloger/model"
	sftpClient "github.com/xloger/sftp"
	"github.com/xloger/tools"
)

type Worker struct {
	Task *model.LogFileTask
	SSH  *sftpClient.SSHConfig
}

func (w *Worker) UpdateTaskStatus(op, status string) error {
	con := db.Con().Xlog

	var (
		err      error
		logSetID = w.Task.LogsetID
		ip       = w.SSH.Host
		role     = w.Task.AgentRole
	)

	switch op {
	case "start_pull":
		res := con.Exec("UPDATE logfile_host SET pl_stime = ? , pl_state = ? WHERE logset_id = ? AND host = ? AND agent_role = ?", time.Now(), status, logSetID, ip, role)
		err = res.Error

	case "end_pull":
		res := con.Exec("UPDATE logfile_host set pl_etime = ? , pl_state = ? WHERE logset_id = ? AND host = ? AND agent_role = ?", time.Now(), status, logSetID, ip, role)
		err = res.Error

	case "start_push":
		res := con.Exec("UPDATE logfile_host set pu_stime = ? , pu_state = ? WHERE logset_id = ? AND host = ? AND agent_role = ?", time.Now(), status, logSetID, ip, role)
		err = res.Error

	case "end_push":
		res := con.Exec("UPDATE logfile_host set pu_etime = ? , pu_state = ? WHERE logset_id = ? AND host = ? AND agent_role = ?", time.Now(), status, logSetID, ip, role)
		err = res.Error
	default:
		err = errors.New("parameter op invalid")
	}

	return err

}

func (w *Worker) GetAbsLogFileName(filepath string) (map[string]string, error) {

	var ret = make(map[string]string)

	if !strings.Contains(filepath, "[") || !strings.Contains(filepath, "]") {
		errMsg := fmt.Sprintf("logfile path: %s no contain keyword [ or ]", filepath)
		return nil, errors.New(errMsg)
	}

	pre, suf := strings.Index(filepath, "["), strings.Index(filepath, "]")
	format := filepath[pre+1 : suf]
	dt, err := tools.GetLastTime(format, w.Task.LogsetSplit)
	if err != nil {
		return nil, err
	}

	absFilename := fmt.Sprintf("%s%s%s", filepath[:pre], dt, filepath[suf+1:])
	L := strings.Split(absFilename, "/")
	filename := L[len(L)-1]

	ret["absFilename"] = absFilename
	ret["filename"] = filename

	return ret, nil

}

func (w *Worker) Pull(filepath string, wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		svr *sftp.Client
		err error
	)

	if err = w.UpdateTaskStatus("start_pull", "running"); err != nil {
		log.Println(err)
		return
	}
	log.Printf("RUNNING: start pull logfile from %s", w.SSH.Host)

	ret, err := w.GetAbsLogFileName(filepath)
	if ret["absFilename"] == "" || ret["filename"] == "" {
		log.Printf("parse %s error", filepath)
		return
	}

	tempPath, err := tools.ParseDir(w.Task.TempPath, w.Task.LogsetSplit)
	if err != nil {
		log.Println(err)
		if err = w.UpdateTaskStatus("end_pull", "failed"); err != nil {
			log.Println(err)
		}
		return
	}

	if exist := file.IsExist(tempPath); !exist {
		if err = os.MkdirAll(tempPath, 0755); err != nil {
			log.Println(err)
			if err = w.UpdateTaskStatus("end_pull", "failed"); err != nil {
				log.Println(err)
			}
			return
		}
	}

	absLfilename := path.Join(tempPath, fmt.Sprintf("%s-%s", w.SSH.Host, ret["filename"]))
	absLfile, err := os.Create(absLfilename)
	if err != nil {
		log.Println(err)
		if err = w.UpdateTaskStatus("end_pull", "failed"); err != nil {
			log.Println(err)
		}
		return
	}
	defer absLfile.Close()

	if svr, err = w.SSH.NewSFTPClient(); err != nil {
		log.Println(err)
		if err = w.UpdateTaskStatus("end_pull", "failed"); err != nil {
			log.Println(err)
		}
		return
	}
	defer svr.Close()
	defer w.SSH.SSHSession.Close()

	absRfilename := ret["absFilename"]
	absRfile, err := svr.Open(absRfilename)
	if err != nil {
		log.Println(err)
		if err = w.UpdateTaskStatus("end_pull", "failed"); err != nil {
			log.Println(err)
		}
		return
	}
	defer absRfile.Close()

	if _, err = absRfile.WriteTo(absLfile); err != nil {
		log.Print(err)
		if err = w.UpdateTaskStatus("end_pull", "failed"); err != nil {
			log.Println(err)
		}
		return
	}
	log.Printf("SUCCESS: pull %s from %s", absRfile.Name(), w.SSH.Host)

	if err = w.UpdateTaskStatus("end_pull", "success"); err != nil {
		log.Println(err)
		return
	}

	// w.Push(filepath, absLfilename)

}

func (w *Worker) Push(filepath, absLfilename string) {

	var (
		err      error
		hdfsPath = w.Task.HdfsPath
		split    = w.Task.LogsetSplit
	)
	if err = w.UpdateTaskStatus("start_push", "running"); err != nil {
		log.Println(err)
		return
	}

	if hdfsPath == "" || absLfilename == "" {
		log.Println("FAILED: hdfsPath or tempPathFile is null")
		if err = w.UpdateTaskStatus("end_push", "failed"); err != nil {
			log.Println(err)
		}
		return
	}

	hd, err := tools.ParseDir(hdfsPath, split)
	if err != nil {
		log.Println(err)
		if err = w.UpdateTaskStatus("end_push", "failed"); err != nil {
			log.Println(err)
		}
		return
	}
	L := strings.Split(absLfilename, "/")
	absHdfsFilename := path.Join(hd, L[len(L)-1])

	namenodes := g.Config().HdfsConfig.NameNodes
	hdfsClient, err := hadoop.HdfsClient(namenodes)
	if err != nil {
		log.Println(err)
		if err = w.UpdateTaskStatus("end_push", "failed"); err != nil {
			log.Println(err)
		}
		return
	}
	defer hdfsClient.Close()

	_, err = hdfsClient.Stat(hd)
	if err != nil {
		if err = hdfsClient.MkdirAll(hd, 0755); err != nil {
			log.Println(err)
			if err = w.UpdateTaskStatus("end_push", "failed"); err != nil {
				log.Println(err)
			}
		}
		return
	}

	if _, err = hdfsClient.Stat(absHdfsFilename); err == nil {
		if err = w.UpdateTaskStatus("end_push", "success"); err != nil {
			log.Println(err)
		}
		return
	}

	if err = hdfsClient.CopyToRemote(absLfilename, absHdfsFilename); err != nil {
		log.Println(err)
		if err = w.UpdateTaskStatus("end_push", "failed"); err != nil {
			log.Println(err)
		}
		return
	}
	log.Printf("SUCCESS: push %s", absLfilename)

	if err = w.UpdateTaskStatus("end_push", "success"); err != nil {
		log.Printf("WARNING: push success but update mysql error %v", err)
	}
	return

}
