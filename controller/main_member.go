package controller

import (
	"fmt"
	"sync"
	"time"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-train/config"
	"github.com/reechou/robot-train/ext"
	"github.com/reechou/robot-train/robot_proto"
)

var (
	TRAIN_TOPIC = []string{
		"你好啊",
		"你觉得今天天气怎么样",
		"你现在在干嘛呢",
		"我在看书呢,是不是很厉害",
		"我在玩游戏,玩王者荣耀呢,你呢",
		"你觉得LOL这个游戏怎么样,反正我在玩",
		"我喜欢工作,热爱工作",
		"我是这个世界的神,你觉得呢.",
		"你吃饭了吗",
		"今天工作到什么时候,咱们去看电影啊",
		"咱们去逛街吧,今天据说有打折呢",
		"你知道腾讯这个公司吗,有点牛掰呢",
		"你喜欢听什么歌",
		"咱们等会去KTV吧?",
		"今天奢侈了一把,享受了一把头等舱",
		"今天我请假了,放假的心情真好",
		"明天我要去日本啦,哈哈哈",
		"我在香港买了很多好东西,要不要看看",
		"如果我是dj你会爱我吗",
		"你的眼神让我看到了你的纯洁",
		"留下太多的故事,心灵得到解脱",
		"你觉得什么时候是最美的时刻?",
	}
)

type WxMainTrainMember struct {
	cfg       *config.Config
	tulingExt *ext.TulingExt
	robotExt  *ext.RobotExt

	topicIdx        int
	trainListMutex  sync.Mutex
	trainList       []robot_proto.GroupUserInfo
	trainReplyMutex sync.Mutex
	trainReply      map[string]string
	changeTopicTime int64

	stop chan struct{}
	done chan struct{}
}

func NewWxMainTrainMember(cfg *config.Config, tulingExt *ext.TulingExt, robotExt *ext.RobotExt) *WxMainTrainMember {
	wmtm := &WxMainTrainMember{
		cfg:        cfg,
		tulingExt:  tulingExt,
		robotExt:   robotExt,
		trainReply: make(map[string]string),
		stop:       make(chan struct{}),
		done:       make(chan struct{}),
	}
	err := wmtm.getTrainList()
	if err != nil {
		holmes.Fatal("get train list error: %v", err)
	}
	go wmtm.runGetTrainList()
	go wmtm.run()

	return wmtm
}

func (self *WxMainTrainMember) Stop() {
	close(self.stop)
	<-self.done
}

func (self *WxMainTrainMember) runGetTrainList() {
	holmes.Debug("start run get train list")
	for {
		select {
		case <-time.After(5 * time.Minute):
			self.getTrainList()
		case <-self.stop:
			return
		}
	}
}

func (self *WxMainTrainMember) getTrainList() error {
	req := &robot_proto.RobotGetGroupMemberListReq{
		WechatNick:    self.cfg.MainMember,
		GroupNickName: self.cfg.TrainPool,
	}
	trainList, err := self.robotExt.GetGroupMemberList(req)
	if err != nil {
		holmes.Error("get group member list error: %v", err)
		return err
	}
	holmes.Debug("get train list: %v", trainList)
	self.trainListMutex.Lock()
	self.trainList = trainList
	self.trainListMutex.Unlock()
	// check if friend
	for _, v := range self.trainList {
		if v.NickName == self.cfg.MainMember {
			continue
		}
		findReq := &robot_proto.RobotFindFriendReq{
			WechatNick: self.cfg.MainMember,
			UserName:   v.UserName,
			NickName:   v.NickName,
		}
		uf, err := self.robotExt.FindFriend(findReq)
		if err != nil {
			holmes.Error("robot find friend error: %v", err)
			continue
		}
		if uf == nil {
			addReq := &robot_proto.RobotAddFriendReq{
				WechatNick: self.cfg.MainMember,
				UserName:   v.UserName,
			}
			err = self.robotExt.AddFriend(addReq)
			if err != nil {
				holmes.Error("robot add friend error: %v", err)
			}
		} else {
			holmes.Debug("find friend: %v", uf)
		}
	}
	return nil
}

func (self *WxMainTrainMember) run() {
	holmes.Debug("start run train")
	interval := self.cfg.TrainInterval
	if interval == 0 {
		interval = 5
	}
	for {
		select {
		case <-time.After(time.Duration(interval) * time.Minute):
			self.check()
		case <-self.stop:
			close(self.done)
			return
		}
	}
}

func (self *WxMainTrainMember) check() {
	self.trainListMutex.Lock()
	defer self.trainListMutex.Unlock()

	nowTime := time.Now()
	hour := nowTime.Hour()
	if hour >= 2 && hour <= 8 {
		holmes.Debug("train pause in night.")
		return
	}

	var ifChangeTopic bool
	now := time.Now().Unix()
	if now-self.changeTopicTime > 7200 {
		ifChangeTopic = true
		self.changeTopicTime = now
	}

	for _, v := range self.trainList {
		if v.NickName == "Mr.REE" || v.NickName == self.cfg.MainMember {
			continue
		}
		user := self.user(v.UserName)
		reply := self.getReply(user)
		newReply := ""
		if reply == "" {
			newReply = self.getTopic()
		} else {
			tulingReply, err := self.tulingExt.SimpleCallV1(reply, user)
			if err != nil {
				holmes.Error("tuling get new reply error: %v", err)
				newReply = self.getTopic()
			} else {
				if tulingReply == reply {
					newReply = self.getTopic()
				} else {
					if ifChangeTopic {
						holmes.Debug("too long time, start to change topic")
						newReply = self.getTopic()
					} else {
						newReply = tulingReply
					}
				}
			}
		}
		var sendReq robot_proto.SendMsgInfo
		sendReq.SendMsgs = append(sendReq.SendMsgs, robot_proto.SendBaseInfo{
			WechatNick: self.cfg.MainMember,
			ChatType:   robot_proto.FROM_TYPE_PEOPLE,
			UserName:   v.UserName,
			NickName:   v.NickName,
			MsgType:    robot_proto.RECEIVE_MSG_TYPE_TEXT,
			Msg:        newReply,
		})
		err := self.robotExt.SendMsgs(self.cfg.MainMember, &sendReq)
		if err != nil {
			holmes.Error("main train send msg[%v] error: %v", sendReq, err)
		}
	}
}

func (self *WxMainTrainMember) TrainReply(user, reply string) {
	self.trainReplyMutex.Lock()
	defer self.trainReplyMutex.Unlock()

	self.trainReply[user] = reply
}

func (self *WxMainTrainMember) getReply(user string) string {
	self.trainReplyMutex.Lock()
	defer self.trainReplyMutex.Unlock()

	return self.trainReply[user]
}

func (self *WxMainTrainMember) user(username string) string {
	return fmt.Sprintf("%s__%s", self.cfg.MainMember, username)
}

func (self *WxMainTrainMember) getTopic() string {
	self.topicIdx = (self.topicIdx + 1) % len(TRAIN_TOPIC)
	return TRAIN_TOPIC[self.topicIdx]
}
