package main

import (
	"fmt"
	"github.com/aiquestion/go-simplejson"
)

type Event struct {
	data    *IrThirdDbBean `json:"data"`
	Token   string         `json:"token"`
	TraceId string         `json:"trace_id"`
}

func handle(event map[string]interface{}, context map[string]interface{}) (interface{}, error) {
	//bytes, _ := json.Marshal(event)
	//evt := &Event{}
	//err := json.Unmarshal(bytes, evt)
	//if err != nil {
		//return "Unmarshal event err  ", err
	//}
	data := simplejson.From(event)
	zipData := data.Get("data").Get("params").Get("zip_data").MustString()
	decodeIrData := GoLumiHeatshrinkDecode(zipData)
	fmt.Println(decodeIrData)

	//提取特征值和压缩值

	//验证特征值和压缩值在范围内

	//根据配置得到提取值

	//查询数据库
	//bean, err := miot.QueryIrThirdDb(evt.Token, evt.data)
	//if err != nil {
		//return nil, err
	//}

	return fmt.Sprintf("bean:%v", ""), nil
}