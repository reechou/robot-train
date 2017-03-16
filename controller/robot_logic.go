package controller

import (
	"fmt"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-train/robot_proto"
)

func (self *Logic) HandleReceiveMsg(msg *robot_proto.ReceiveMsgInfo) {
	holmes.Debug("receive robot msg: %v", msg)
	switch msg.BaseInfo.ReceiveEvent {
	case robot_proto.RECEIVE_EVENT_MSG:
		self.handleMsg(msg)
	}
}

func (self *Logic) handleMsg(msg *robot_proto.ReceiveMsgInfo) {
	user := self.user(msg.BaseInfo.WechatNick, msg.BaseInfo.FromUserName)
	if self.cfg.MainMember == msg.BaseInfo.WechatNick {
		self.mt.TrainReply(user, msg.Msg)
	} else {
		tulingReply, err := self.tulingExt.SimpleCallV1(msg.Msg, user)
		if err != nil {
			holmes.Error("tuling reply error: %v", err)
		}
		var sendReq robot_proto.SendMsgInfo
		sendReq.SendMsgs = append(sendReq.SendMsgs, robot_proto.SendBaseInfo{
			WechatNick: msg.BaseInfo.WechatNick,
			ChatType:   msg.BaseInfo.FromType,
			UserName:   msg.BaseInfo.FromUserName,
			NickName:   msg.BaseInfo.FromNickName,
			MsgType:    robot_proto.RECEIVE_MSG_TYPE_TEXT,
			Msg:        tulingReply,
		})
		err = self.robotExt.SendMsgs(self.cfg.MainMember, &sendReq)
		if err != nil {
			holmes.Error("member reply send msg[%v] error: %v", sendReq, err)
		}
	}
}

func (self *Logic) user(robot, username string) string {
	return fmt.Sprintf("%s__%s", robot, username)
}
