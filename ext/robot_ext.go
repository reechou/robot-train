package ext

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-train/config"
	"github.com/reechou/robot-train/robot_proto"
)

const (
	ROBOT_START_WX              = "/startwx2"
	ROBOT_ALL_ROBOTS_URI        = "/allrobots"
	ROBOT_SEND_MSGS_URI         = "/sendmsgs"
	ROBOT_FIND_FRIEND_URI       = "/findfriend"
	ROBOT_REMARK_FRIEND_URI     = "/remarkfriend"
	ROBOT_GROUP_TIREN_URI       = "/grouptiren"
	ROBOT_GROUP_MEMBER_LIST_URI = "/group_member_list"
	ROBOT_ADD_FRIEND_URI        = "/addfriend"
)

type RobotExt struct {
	client *http.Client
	cfg    *config.Config
}

func NewRobotExt(cfg *config.Config) *RobotExt {
	return &RobotExt{
		client: &http.Client{},
		cfg:    cfg,
	}
}

func (self *RobotExt) LoginRobot(request *robot_proto.StartWxReq) error {
	url := "http://" + self.cfg.RobotHost.Host + ROBOT_START_WX
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return err
	}
	var response robot_proto.WxResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return err
	}
	if response.Code != 0 {
		holmes.Error("login robot result code error: %d %s", response.Code, response.Msg)
		return fmt.Errorf("login robot result error.")
	}

	return nil
}

func (self *RobotExt) AllLoginRobots() (interface{}, error) {
	url := "http://" + self.cfg.RobotHost.Host + ROBOT_ALL_ROBOTS_URI
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return nil, err
	}
	var response robot_proto.WxResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return nil, err
	}
	if response.Code != 0 {
		holmes.Error("get all login robots result code error: %d %s", response.Code, response.Msg)
		return nil, fmt.Errorf("get all login robots result error.")
	}

	return response.Data, nil
}

func (self *RobotExt) GroupTiren(request *robot_proto.RobotGroupTirenReq) (*robot_proto.GroupUserInfo, error) {
	// maybe get robot host here

	reqBytes, err := json.Marshal(request)
	if err != nil {
		holmes.Error("json encode error: %v", err)
		return nil, err
	}

	url := "http://" + self.cfg.RobotHost.Host + ROBOT_GROUP_TIREN_URI
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return nil, err
	}
	var response robot_proto.WxGroupTirenResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return nil, err
	}
	if response.Code != 0 {
		holmes.Error("group tiren result code error: %d %s", response.Code, response.Msg)
		return nil, fmt.Errorf("group tiren result error.")
	}

	return &response.Data, nil
}

func (self *RobotExt) RemarkFriend(request *robot_proto.RobotRemarkFriendReq) error {
	// maybe get robot host here

	reqBytes, err := json.Marshal(request)
	if err != nil {
		holmes.Error("json encode error: %v", err)
		return err
	}

	url := "http://" + self.cfg.RobotHost.Host + ROBOT_REMARK_FRIEND_URI
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return err
	}
	var response robot_proto.WxResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return err
	}
	if response.Code != 0 {
		holmes.Error("remark friend result code error: %d %s", response.Code, response.Msg)
		return fmt.Errorf("remark friend result error.")
	}

	return nil
}

func (self *RobotExt) FindFriend(request *robot_proto.RobotFindFriendReq) (*robot_proto.UserFriend, error) {
	// maybe get robot host here

	reqBytes, err := json.Marshal(request)
	if err != nil {
		holmes.Error("json encode error: %v", err)
		return nil, err
	}

	url := "http://" + self.cfg.RobotHost.Host + ROBOT_FIND_FRIEND_URI
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return nil, err
	}
	var response robot_proto.WxFindFriendResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return nil, err
	}
	if response.Code != 0 {
		holmes.Error("find friend result code error: %d %s", response.Code, response.Msg)
		return nil, fmt.Errorf("find friend result error.")
	}

	return &response.Data, nil
}

func (self *RobotExt) GetGroupMemberList(request *robot_proto.RobotGetGroupMemberListReq) ([]robot_proto.GroupUserInfo, error) {
	// maybe get robot host here

	reqBytes, err := json.Marshal(request)
	if err != nil {
		holmes.Error("json encode error: %v", err)
		return nil, err
	}

	url := "http://" + self.cfg.RobotHost.Host + ROBOT_GROUP_MEMBER_LIST_URI
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return nil, err
	}
	var response robot_proto.WxGroupMemberListResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return nil, err
	}
	if response.Code != 0 {
		holmes.Error("get group member list result code error: %d %s", response.Code, response.Msg)
		return nil, fmt.Errorf("get group member list result error.")
	}

	return response.Data, nil
}

func (self *RobotExt) AddFriend(request *robot_proto.RobotAddFriendReq) error {
	// maybe get robot host here
	
	reqBytes, err := json.Marshal(request)
	if err != nil {
		holmes.Error("json encode error: %v", err)
		return err
	}
	
	url := "http://" + self.cfg.RobotHost.Host + ROBOT_ADD_FRIEND_URI
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return err
	}
	var response robot_proto.WxResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return err
	}
	if response.Code != 0 {
		holmes.Error("add friend result code error: %d %s", response.Code, response.Msg)
		return fmt.Errorf("add friend result error.")
	}
	
	return nil
}

func (self *RobotExt) SendMsgs(robotWx string, msg *robot_proto.SendMsgInfo) error {
	holmes.Debug("msg: %v", msg)
	// maybe get robot host here

	reqBytes, err := json.Marshal(msg)
	if err != nil {
		holmes.Error("json encode error: %v", err)
		return err
	}

	url := "http://" + self.cfg.RobotHost.Host + ROBOT_SEND_MSGS_URI
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return err
	}
	var response robot_proto.SendMsgResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return err
	}
	if response.Code != 0 {
		holmes.Error("send msg[%v] result code error: %d %s", msg, response.Code, response.Msg)
		return fmt.Errorf("send msg result error.")
	}

	return nil
}
