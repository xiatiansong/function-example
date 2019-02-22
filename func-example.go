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

	//bean, err := miot.QueryIrThirdDb(evt.Token, evt.data)
	//if err != nil {
		//return nil, err
	//}

	// 业务逻辑

	return fmt.Sprintf("bean:%v", ""), nil
}