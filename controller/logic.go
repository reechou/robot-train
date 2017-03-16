package controller

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-train/config"
	"github.com/reechou/robot-train/ext"
)

type Logic struct {
	sync.Mutex
	
	tulingExt *ext.TulingExt
	robotExt  *ext.RobotExt
	mt *WxMainTrainMember

	cfg *config.Config
}

func NewLogic(cfg *config.Config) *Logic {
	l := &Logic{
		cfg: cfg,
	}
	l.robotExt = ext.NewRobotExt(cfg)
	if cfg.Tuling.IfEnable {
		l.tulingExt = ext.NewTulingExt(cfg)
	}
	l.mt = NewWxMainTrainMember(cfg, l.tulingExt, l.robotExt)
	l.init()

	return l
}

func (self *Logic) init() {
	http.HandleFunc("/robot/receive_msg", self.RobotReceiveMsg)
}

func (self *Logic) Run() {
	defer holmes.Start(holmes.LogFilePath("./log"),
		holmes.EveryDay,
		holmes.AlsoStdout,
		holmes.DebugLevel).Stop()

	if self.cfg.Debug {
		EnableDebug()
	}

	holmes.Info("server starting on[%s]..", self.cfg.Host)
	holmes.Infoln(http.ListenAndServe(self.cfg.Host, nil))
}

func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func WriteBytes(w http.ResponseWriter, code int, v []byte) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	w.Write(v)
}

func EnableDebug() {

}
