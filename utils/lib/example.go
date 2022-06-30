package lib

import (
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"net/http"
	"os"
	"syscall"
)

const (
	URL = "https://contracts.racingtime.io"
)

const (
	Spec = "0 0/1 * * * ?"
)

type ServiceObject struct {
	c    *cron.Cron
	spec string
	cmd  func()
	eId  cron.EntryID
}

func init() {
	var _err error
	_object := ServiceObject{}
	_object.c = cron.New(cron.WithSeconds())
	_object.cmd = Run
	_object.eId, _err = _object.c.AddFunc(Spec, _object.cmd)
	if _err != nil {
		panic(_err)
	}
	_object.c.Start()

}

func Run() {
	resp, err := http.Get(URL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		println("close")
		_pId := os.Getpid()
		_ = syscall.Kill(_pId, syscall.SIGINT)
	}
}
