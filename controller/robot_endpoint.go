package controller

import (
	"encoding/json"
	"net/http"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-train/proto"
	"github.com/reechou/robot-train/robot_proto"
)

func (self *Logic) RobotReceiveMsg(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	req := &robot_proto.ReceiveMsgInfo{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("RobotReceiveMsg json decode error: %v", err)
		return
	}
	self.HandleReceiveMsg(req)

	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	WriteJSON(w, http.StatusOK, rsp)
}
