package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aiquestion/go-simplejson"
	"github.com/patrickmn/go-cache"
)

var miot IMiotClient = nil

var httpClient *http.Client = nil

func init() {
	if miot == nil {
		miot = &MiotClient{
			cache: cache.New(time.Duration(600)*time.Second, time.Duration(20)*time.Minute),
		}
	}
	if httpClient == nil {
		httpClient = &http.Client{
			Transport: &http.Transport{},
			Timeout:   time.Duration(7) * time.Second,
		}
	}
}

type IMiotClient interface {
	DeviceEventPush(token, did, attr, title, content string, push_share bool, sound, value string, extra map[string]interface{}) error
	SetDeviceExtraData(token, did, key, value string) error
	GetDeviceExtraData(token, did, key string) (string, error)
	GetDeviceProps(token, did string, propNames []string) ([]*DevicePropOrEvent, error)
	SetDeviceProp(token, did, prop, value string) error
	GetDeviceEvents(token, did string, eventNames []string) ([]*DevicePropOrEvent, error)
	WriteSdsData(token string, data *SdsData) error
	GetSdsData(token string, data *SdsDataKey) (*SdsData, error)
	ScanSdsData(token string, data *SdsDataScanRequest) ([]*SdsData, error)
	FaceAnalyse(token string, did, file_id string, content_type int64, extra_info string) error
	GetFaces(token string, did string, face_ids []string) ([]*FaceInfoMeta, error)
	SpecParsePropOp(token string, req *ParsePropOpReq) (*PropOp, error)
	SpecParseActionOp(token string, req *ParseActionOpReq) (*ActionOp, error)
	SpecParseEventOp(token string, req *ParseEventOpReq) (*EventOp, error)
	ReportEvent(token, did string, eventinfo []*Eventinfo) error
	RpcDevice(token, did string, method string, params *simplejson.Json) (*RpcResult, error)
	QueryIrThirdDb(token string, req *IrThirdDbBean) (*IrThirdDbBean, error)
}

// this is a service domain in k8s cluster
const BASE_URL = "http://miot-proxy.miot-proxy"
const API_VERSION = "v1"
const API_GROUP = "miot"

type MiotClient struct {
	cache *cache.Cache
}

const (
	DEVICEEVENTPUSH_EXTRA_TITLE_LOCALE_MAP   = "title_locale"
	DEVICEEVENTPUSH_EXTRA_CONTENT_LOCALE_MAP = "content_locale"
)

const (
	LOCALE_ZH_CN = "zh_cn"
	LOCALE_EN    = "en"
	LOCALE_ZH_TW = "zh_tw"
	LOCALE_ZH_HK = "zh_hk"
)

func (m *MiotClient) DeviceEventPush(token, did, attr, title, content string, push_share bool, sound, value string, extra map[string]interface{}) error {
	extraStr := ""
	if extra != nil {
		// there should not be any err here
		extraBts, err := json.Marshal(extra)
		if err != nil {
			return fmt.Errorf("you should never reach here")
		}
		extraStr = string(extraBts)
	}
	data := map[string]interface{}{
		"did":        did,
		"attr":       attr,
		"title":      title,
		"content":    content,
		"push_share": push_share,
		"sound":      sound,
		"value":      value,
		"extra":      extraStr,
	}
	_, err := m._sendToMiotOpenGateway("device/event_push", token, data)
	return err
}

/*
用来存储设备配置
*/
func (m *MiotClient) SetDeviceExtraData(token, did, key, value string) error {
	data := map[string]interface{}{
		"key":   key,
		"did":   did,
		"value": value,
	}
	_, err := m._sendToMiotOpenGateway("device/extra_data/set", token, data)
	return err
}

func (m *MiotClient) GetDeviceExtraData(token, did, key string) (string, error) {
	data := map[string]interface{}{
		"did": did,
		"key": key,
	}
	res, err := m._sendToMiotOpenGateway("device/extra_data/get", token, data)
	if err != nil {
		return "", err
	}
	return res.Get("value").MustString(), nil
}

/*
用来存储设备当前状态, prop/event
*/
func (m *MiotClient) GetDeviceProps(token, did string, propNames []string) ([]*DevicePropOrEvent, error) {
	data := map[string]interface{}{
		"keys": propNames,
		"did":  did,
	}
	result := make([]*DevicePropOrEvent, 0, len(propNames))
	res, err := m._sendToMiotOpenGateway("device/data/batch_get_prop", token, data)
	if err != nil {
		return result, err
	}
	for _, v := range res.MustArray() {
		item := simplejson.FromObject(v)
		result = append(result, &DevicePropOrEvent{
			Name:      item.Get("name").MustString(),
			Value:     item.Get("value").MustString(),
			TimeStamp: item.Get("timestamp").MustInt64(),
		})
	}
	return result, nil
}

func (m *MiotClient) SetDeviceProp(token, did, prop, value string) error {
	data := map[string]string{
		"did":   did,
		"prop":  prop,
		"value": value,
	}

	_, err := m._sendToMiotOpenGateway("device/data/set_prop", token, data)
	if err != nil {
		return err
	}
	return nil
}

func (m *MiotClient) GetDeviceEvents(token, did string, eventNames []string) ([]*DevicePropOrEvent, error) {
	data := map[string]interface{}{
		"keys": eventNames,
		"did":  did,
	}
	result := make([]*DevicePropOrEvent, 0, len(eventNames))
	res, err := m._sendToMiotOpenGateway("device/data/batch_get_event", token, data)
	if err != nil {
		return result, err
	}
	for _, v := range res.MustArray() {
		item := simplejson.FromObject(v)
		result = append(result, &DevicePropOrEvent{
			Name:      item.Get("name").MustString(),
			Value:     item.Get("value").MustString(),
			TimeStamp: item.Get("timestamp").MustInt64(),
		})
	}
	return result, nil
}

/*
用来存储设备历史数据,每次上报的prop/event按秒存储。
*/
func (m *MiotClient) WriteSdsData(token string, data *SdsData) error {
	_, err := m._sendToMiotOpenGateway("device/sds_data/set", token, data)
	if err != nil {
		return err
	}
	return nil
}

func (m *MiotClient) GetSdsData(token string, data *SdsDataKey) (*SdsData, error) {
	res, err := m._sendToMiotOpenGateway("device/sds_data/get", token, data)
	if err != nil {
		return nil, err
	}
	return &SdsData{
		Did:   res.Get("did").MustString(),
		Type:  res.Get("type").MustString(),
		Key:   res.Get("key").MustString(),
		Time:  res.Get("time").MustInt64(),
		Value: res.Get("value").MustString(),
	}, nil
}

func (m *MiotClient) ScanSdsData(token string, data *SdsDataScanRequest) ([]*SdsData, error) {
	res, err := m._sendToMiotOpenGateway("device/sds_data/scan", token, data)
	if err != nil {
		return nil, err
	}
	result := make([]*SdsData, 0)
	for _, item := range res.MustArray() {
		itemWrap := simplejson.FromObject(item)
		result = append(result, &SdsData{
			Did:   itemWrap.Get("did").MustString(),
			Type:  itemWrap.Get("type").MustString(),
			Key:   itemWrap.Get("key").MustString(),
			Time:  itemWrap.Get("time").MustInt64(),
			Value: itemWrap.Get("value").MustString(),
		})
	}
	return result, nil
}

/*
触发人脸识别
content_type: 0-image 1-video
*/
func (m *MiotClient) FaceAnalyse(token string, did, file_id string, content_type int64, extra_info string) error {
	data := map[string]interface{}{
		"did":          did,
		"file_id":      file_id,
		"content_type": content_type,
		"extra_info":   extra_info,
	}
	_, err := m._sendToMiotOpenGateway("camera/face_analyse", token, data)
	if err != nil {
		return err
	}
	return nil
}

/*
获取人脸识别信息
*/
func (m *MiotClient) GetFaces(token string, did string, face_ids []string) ([]*FaceInfoMeta, error) {
	var ret []*FaceInfoMeta
	data := map[string]interface{}{
		"did":      did,
		"face_ids": face_ids,
	}
	res, err := m._sendToMiotOpenGateway("camera/get_faces", token, data)
	if err != nil {
		return nil, err
	}
	jsbt, _ := res.Encode()
	err = json.Unmarshal(jsbt, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

/*
 解析miot spec property
*/
func (m *MiotClient) SpecParsePropOp(token string, req *ParsePropOpReq) (*PropOp, error) {
	var ret *PropOp
	res, err := m._sendToMiotOpenGateway("miotspec/parse_prop_op", token, req)
	if err != nil {
		return nil, err
	}
	jsbt, _ := res.Encode()
	err = json.Unmarshal(jsbt, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

/*
 解析miot spec action
*/
func (m *MiotClient) SpecParseActionOp(token string, req *ParseActionOpReq) (*ActionOp, error) {
	var ret *ActionOp
	res, err := m._sendToMiotOpenGateway("miotspec/parse_action_op", token, req)
	if err != nil {
		return nil, err
	}
	jsbt, _ := res.Encode()
	err = json.Unmarshal(jsbt, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

/*
 解析miot spec event
*/
func (m *MiotClient) SpecParseEventOp(token string, req *ParseEventOpReq) (*EventOp, error) {
	var ret *EventOp
	res, err := m._sendToMiotOpenGateway("miotspec/parse_action_op", token, req)
	if err != nil {
		return nil, err
	}
	jsbt, _ := res.Encode()
	err = json.Unmarshal(jsbt, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (m *MiotClient) ReportEvent(token, did string, eventinfo []*Eventinfo) error {
	data := map[string]interface{}{
		"did":    did,
		"events": eventinfo,
	}
	_, err := m._sendToMiotOpenGateway("event/report_event", token, data)
	if err != nil {
		return err
	}
	return nil
}

func (m *MiotClient) RpcDevice(token, did string, method string, params *simplejson.Json) (*RpcResult, error) {
	data := map[string]interface{}{
		"did":    did,
		"method": method,
		"params": params,
	}
	res, err := m._sendToMiotOpenGateway("device/rpc", token, data)
	if err != nil {
		return nil, err
	}
	return &RpcResult{
		Code:    res.Get("code").MustInt64(-4),
		Message: res.Get("message").MustString(),
		Result:  res.Get("result"),
	}, nil
}

func (m *MiotClient) QueryIrThirdDb(token string, req *IrThirdDbBean) (*IrThirdDbBean, error) {
	var ret *IrThirdDbBean
	res, err := m._sendToMiotOpenGateway("irthirddb/query", token, req)
	if err != nil {
		return nil, err
	}

	jsbt, _ := res.Encode()
	err = json.Unmarshal(jsbt, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (m *MiotClient) _sendToMiotOpenGateway(path, token string, data interface{}) (_body *simplejson.Json, _err error) {
	st := time.Now()
	defer func() {
		pcCost(fmt.Sprintf("cost_%s", path), uint32(time.Now().Sub(st).Nanoseconds()/int64(time.Millisecond)))
	}()
	pcIncr(fmt.Sprintf("call_%s", path), 1)
	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("cannot convert data to json %+v", data)
	}
	payload := url.Values{}
	payload.Add("data", string(dataStr))
	payload.Add("token", token)
	url := fmt.Sprintf("%s/%s/%s/%s", BASE_URL, API_GROUP, API_VERSION, path)
	resp, err := httpClient.Post(url, "application/x-www-form-urlencoded", strings.NewReader(payload.Encode()))
	if err != nil {
		pcIncr(fmt.Sprintf("call_fail_%s", path), 1)
		LogError("call %s failed with %v  err=%v", url, payload, err)
		return nil, err
	}
	if resp == nil {
		pcIncr(fmt.Sprintf("call_fail_%s", path), 1)
		LogError("call %s failed with %v response nil", url, payload)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		LogError("call %s failed with %v read resp body failed %v ", url, payload, err)
	}

	if resp.StatusCode >= 300 && resp.StatusCode < 500 {
		pcIncr(fmt.Sprintf("call_cli_fail_%s", path), 1)
		LogInfo("call %s failed with %v response code %v body %s", url, payload, resp.StatusCode, body)
		return nil, fmt.Errorf("send device event failed with code %v", resp.StatusCode)
	}
	if resp.StatusCode > 500 {
		pcIncr(fmt.Sprintf("call_fail_%s", path), 1)
		LogWarn("call %s failed with %v response code %v body %s", url, payload, resp.StatusCode, body)
		return nil, fmt.Errorf("send device event failed with code %v", resp.StatusCode)
	}
	res, err := simplejson.NewJson([]byte(body))
	if err != nil {
		return nil, err
	}
	code := res.Get("code").MustInt(-4)
	if code != 0 {
		msg := res.Get("message").MustString()
		if msg == "" {
			msg = "unknown error"
		}
		if code == -4 {
			pcIncr(fmt.Sprintf("call_fail_%s", path), 1)
			LogWarn("call %s failed with %v response code %v body %s", url, payload, resp.StatusCode, body)
		} else {
			pcIncr(fmt.Sprintf("call_cli_fail_%s", path), 1)
			LogInfo("call %s failed with %v response code %v body %s", url, payload, resp.StatusCode, body)
		}
		return nil, fmt.Errorf("code=%v, message=%s", code, msg)
	}

	return res.Get("result"), nil
}
