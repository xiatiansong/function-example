package main

import (
	"encoding/json"
	"fmt"
	"github.com/aiquestion/go-simplejson"
)

type Event struct {
	data    *IrThirdDbBean `json:"data"`
	Token   string         `json:"token"`
	TraceId string         `json:"trace_id"`
}

func handle(event map[string]interface{}, context map[string]interface{}) (interface{}, error) {
	bytes, _ := json.Marshal(event)
	evt := &Event{}
	err := json.Unmarshal(bytes, evt)
	if err != nil {
		return "Unmarshal event err  ", err
	}
	data := simplejson.From(event)
	zipData := data.Get("data").Get("params").Get("zip_data").MustString()
	decodeIrData := GoLumiHeatshrinkDecode(zipData)

	//提取特征值和压缩值
	_, zipRetData := GoLumiIrZip(decodeIrData)
	//验证特征值和压缩值在范围内
	//根据配置得到提取值
	controllerId := data.Get("data").Get("params").Get("controller_id").MustString()
	irRemoteCfg := (*GetInstance())[controllerId]
	tempArray := extractArray(zipRetData, irRemoteCfg.TempArray)
	onoffArray := extractArray(zipRetData, irRemoteCfg.OnoffArray)
	modelArray := extractArray(zipRetData, irRemoteCfg.ModelArray)
	windSpeedArray := extractArray(zipRetData, irRemoteCfg.WindSpeedArray)
	irThirdDbBean := IrThirdDbBean{RemoteId: controllerId, OnoffArray: onoffArray, ModelArray: modelArray, TempArray: tempArray, WindSpeedArray: windSpeedArray}
	//查询数据库
	bean, err := miot.QueryIrThirdDb(evt.Token, &irThirdDbBean)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("bean:%v", bean), nil
}
