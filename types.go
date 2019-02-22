package main

import "github.com/aiquestion/go-simplejson"

type DevicePropOrEvent struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	TimeStamp int64  `json:"timestamp"`
}

type SdsData struct {
	Did   string `json:"did"`
	Type  string `json:"type"`
	Key   string `json:"key"`
	Time  int64  `json:"time"`
	Value string `json:"value"`
}

type SdsDataKey struct {
	Did  string `json:"did"`
	Type string `json:"type"`
	Key  string `json:"key"`
	Time int64  `json:"time"`
}

type SdsDataScanRequest struct {
	Did       string `json:"did"`
	Type      string `json:"type"`
	Key       string `json:"key"`
	TimeStart int64  `json:"time_start"`
	TimeEnd   int64  `json:"time_end"`
	Limit     int64  `json:"limit"`
}

type FaceInfoMeta struct {
	FaceIdStr        string   `json:"faceIdStr"`
	Score            *float64 `json:"score,omitempty"`
	Age              *int32   `json:"age,omitempty"`
	AgeConfidence    *float64 `json:"ageConfidence,omitempty"`
	Gender           *int32   `json:"gender,omitempty"`
	GenderConfidence *float64 `json:"genderConfidence,omitempty"`
	WearGlasses      *bool    `json:"wearGlasses,omitempty"`
	FigureInfo       *string  `json:"figureInfo,omitempty"`
	OrgPhoto         []byte   `json:"orgPhoto,omitempty"`
	Found            *bool    `json:"found,omitempty"`
}

type ParsePropOpReq struct {
	QueryType    int64  `json:"queryType"`
	Pdid         int32  `json:"pdid,omitempty"`
	Model        string `json:"model,omitempty"`
	InstanceVer  int32  `json:"instanceVer,omitempty"`
	InstanceType string `json:"instanceType,omitempty"`
	ServiceType  string `json:"serviceType,omitempty"`
	ServiceDesc  string `json:"serviceDesc,omitempty"`
	PropType     string `json:"propType,omitempty"`
	PropDesc     string `json:"propDesc,omitempty"`
	Value        string `json:"value,omitempty"`
}

type ParseActionOpReq struct {
	QueryType    int64    `json:"queryType"`
	Pdid         int32    `json:"pdid,omitempty"`
	Model        string   `json:"model,omitempty"`
	InstanceVer  int32    `json:"instanceVer,omitempty"`
	InstanceType string   `json:"instanceType,omitempty"`
	ServiceType  string   `json:"serviceType,omitempty"`
	ServiceDesc  string   `json:"serviceDesc,omitempty"`
	ActionType   string   `json:"actionType,omitempty"`
	ActionDesc   string   `json:"actionDesc,omitempty"`
	Ins          []string `json:"ins,omitempty"`
	Outs         []string `json:"outs,omitempty"`
}

type ParseEventOpReq struct {
	QueryType    int64    `thrift:"queryType,1,required" json:"queryType"`
	Pdid         int32    `thrift:"pdid,2" json:"pdid,omitempty"`
	Model        string   `thrift:"model,3" json:"model,omitempty"`
	InstanceVer  int32    `thrift:"instanceVer,4" json:"instanceVer,omitempty"`
	InstanceType string   `thrift:"instanceType,5" json:"instanceType,omitempty"`
	ServiceType  string   `thrift:"serviceType,6" json:"serviceType,omitempty"`
	ServiceDesc  string   `thrift:"serviceDesc,7" json:"serviceDesc,omitempty"`
	EventType    string   `thrift:"eventType,8" json:"eventType,omitempty"`
	EventDesc    string   `thrift:"eventDesc,9" json:"eventDesc,omitempty"`
	Arguments    []string `thrift:"arguments,10" json:"arguments,omitempty"`
}

type PropOp struct {
	Did         string `json:"did"`
	Siid        int32  `json:"siid"`
	Piid        int32  `json:"piid"`
	Value       string `json:"value,omitempty"`
	Status      int32  `json:"status,omitempty"`
	Description string `json:"description,omitempty"`
}

type ActionOp struct {
	Did         string   `json:"did"`
	Siid        int32    `json:"siid"`
	Aiid        int32    `json:"aiid"`
	InArgs      []string `json:"inArgs,omitempty"`
	OutArgs     []string `json:"outArgs,omitempty"`
	Status      int32    `json:"status,omitempty"`
	Description string   `json:"description,omitempty"`
}

type EventOp struct {
	Did         string   `json:"did"`
	Siid        int32    `json:"siid"`
	Eiid        int32    `json:"eiid"`
	Arguments   []string `json:"arguments,omitempty"`
	Status      int32    `json:"status,omitempty"`
	Description string   `json:"description,omitempty"`
}

type Eventinfo struct {
	Attr      string `json:"attr"`
	Value     string `json:"value"`
	Timestamp *int64 `json:"timestamp,omitempty"`
}

type RpcResult struct {
	Code    int64            `json:"code"`
	Message string           `json:"message"`
	Result  *simplejson.Json `json:"result"`
}

type IrThirdDbBean struct {
	RemoteId        string `json:"remote_id"`
	RsId            string `json:"rs_id,omitempty"`
	Frequency       string `json:"frequency,omitempty"`
	ShortCmd        string `json:"short_cmd,omitempty"`
	OnoffArray      string `json:"onoff_array,omitempty"`
	ModelArray      string `json:"model_array,omitempty"`
	TempArray       string `json:"temp_array,omitempty"`
	WindSpeedArray  string `json:"wind_speed_array,omitempty"`
	WindDirectArray string `json:"wind_direct_array,omitempty"`
}
