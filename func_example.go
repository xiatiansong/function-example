package main

import (
	"github.com/aiquestion/go-simplejson"
	"lumi/function-example/heatshrink"
)

func handle(event map[string]interface{}, context map[string]interface{}) (interface{}, error) {
	eventJs := simplejson.From(event)
	token := eventJs.Get("token").MustString()
	data := eventJs.Get("data")
	zipData := data.Get("params").Get("zip_data").MustString()
	decodeIrData := heatshrink.GoLumiHeatshrinkDecode(zipData)

	//提取特征值和压缩值
	_, zipRetData := heatshrink.GoLumiIrZip(decodeIrData)
	//验证特征值和压缩值在范围内
	//根据配置得到提取值
	controllerId := data.Get("params").Get("controller_id").MustString()
	brand_id := data.Get("params").Get("brand_id").MustString()
	irRemoteCfg := (*heatshrink.GetInstance())[controllerId]
	if &irRemoteCfg == nil {
		return map[string]interface{}{"code": 500,}, nil
	}
	onoffArray := heatshrink.ExtractArray(zipRetData, irRemoteCfg.OnoffArray)
	irThirdDbBean := IrThirdDbBean{RemoteId: controllerId, OnoffArray: onoffArray}
	tempArray := heatshrink.ExtractArray(zipRetData, irRemoteCfg.TempArray)
	if tempArray != "" {
		irThirdDbBean.TempArray = tempArray
	}
	modelArray := heatshrink.ExtractArray(zipRetData, irRemoteCfg.ModelArray)
	if modelArray != "" {
		irThirdDbBean.ModelArray = modelArray
	}
	windSpeedArray := heatshrink.ExtractArray(zipRetData, irRemoteCfg.WindSpeedArray)
	if windSpeedArray != "" {
		irThirdDbBean.WindSpeedArray = windSpeedArray
	}
	//查询数据库
	bean, err := miot.QueryIrThirdDb(token, &irThirdDbBean)
	if err != nil || bean == nil {
		return map[string]interface{}{"code": 500,}, err
	}
	//构建返回
	retMap := map[string]interface{}{
		"code":          200,
		"brand_id":      brand_id,
		"controller_id": controllerId,
		"ac_key":        heatshrink.ShortCmdToAcKey(bean.ShortCmd),
	}
	return retMap, nil
}
