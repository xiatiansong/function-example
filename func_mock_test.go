package main

import (
	"fmt"
	"github.com/aiquestion/go-simplejson"
	"sort"
	"strings"
	"time"
)

type MockMiotClient struct {
	DeviceExtraKv map[string]map[string]string
	DeviceProps   map[string]map[string]string
	DeviceEvents  map[string]map[string]string

	FaceInfos map[string]*FaceInfoMeta
	SdsData   []*SdsData
}

func (m *MockMiotClient) DeviceEventPush(token, did, attr, title, content string, push_share bool, sound, value string, extra map[string]interface{}) error {
	m._logAction("DeviceEventPush: did=%s attr=%s title=%s content=%s push_share=%v sound=%s value=%s extra=%v",
		did, attr, title, content, push_share, sound, value, extra)
	return nil
}

func (m *MockMiotClient) SetDeviceExtraData(token, did, key, value string) error {
	if m.DeviceExtraKv == nil {
		m.DeviceExtraKv = map[string]map[string]string{}
	}
	if m.DeviceExtraKv[did] == nil {
		m.DeviceExtraKv[did] = map[string]string{}
	}
	m.DeviceExtraKv[did][key] = value
	return nil
}

func (m *MockMiotClient) GetDeviceExtraData(token, did, key string) (string, error) {
	if kv, ok := m.DeviceExtraKv[did]; ok && kv != nil {
		return m.DeviceExtraKv[did][key], nil
	}
	return "", nil
}

func (m *MockMiotClient) GetDeviceProps(token, did string, propNames []string) ([]*DevicePropOrEvent, error) {
	res := make([]*DevicePropOrEvent, 0, len(propNames))
	for _, n := range propNames {
		if kv, ok := m.DeviceProps[did]; ok && kv != nil {
			res = append(res, &DevicePropOrEvent{
				Name:      n,
				Value:     m.DeviceProps[did][n],
				TimeStamp: time.Now().Unix(),
			})
		}
	}
	return res, nil
}

func (m *MockMiotClient) SetDeviceProp(token, did, prop, value string) error {
	if kv, ok := m.DeviceProps[did]; ok {
		kv[prop] = value
	} else {
		kv := make(map[string]string, 0)
		kv[prop] = value
	}
	return nil
}

func (m *MockMiotClient) GetDeviceEvents(token, did string, eventNames []string) ([]*DevicePropOrEvent, error) {
	res := make([]*DevicePropOrEvent, 0, len(eventNames))
	for _, n := range eventNames {
		if kv, ok := m.DeviceEvents[did]; ok && kv != nil {
			res = append(res, &DevicePropOrEvent{
				Name:      n,
				Value:     m.DeviceEvents[did][n],
				TimeStamp: time.Now().Unix(),
			})
		}
	}
	return res, nil
}

func (m *MockMiotClient) WriteSdsData(token string, data *SdsData) error {
	if m.SdsData == nil {
		m.SdsData = make([]*SdsData, 0)
	}
	m.SdsData = append(m.SdsData, data)
	m._sortSdsData()
	return nil
}

func (m *MockMiotClient) GetSdsData(token string, data *SdsDataKey) (*SdsData, error) {
	s := &SdsData{
		Did:  data.Did,
		Type: data.Type,
		Key:  data.Key,
		Time: data.Time,
	}
	ri := sort.Search(len(m.SdsData), func(i int) bool { // this should be a greater than or equal func
		cmp := m._compareSdsData(s, m.SdsData[i])
		return cmp >= 0
	})
	if ri > 0 {
		if m._compareSdsData(s, m.SdsData[ri]) == 0 {
			return m.SdsData[ri], nil
		}
	}
	return nil, nil
}

// time will return in desc order
func (m *MockMiotClient) ScanSdsData(token string, data *SdsDataScanRequest) ([]*SdsData, error) {
	sk := &SdsData{
		Did:  data.Did,
		Type: data.Type,
		Key:  data.Key,
		Time: data.TimeEnd,
	}
	ek := &SdsData{
		Did:  data.Did,
		Type: data.Type,
		Key:  data.Key,
		Time: data.TimeStart,
	}
	si := sort.Search(len(m.SdsData), func(i int) bool { // this should be a greater than func
		cmp := m._compareSdsData(sk, m.SdsData[i])
		return cmp <= 0
	})
	ei := sort.Search(len(m.SdsData), func(i int) bool { // this should be a greater than func
		cmp := m._compareSdsData(ek, m.SdsData[i])
		return cmp < 0
	})
	if si < 0 || si >= len(m.SdsData) {
		return nil, nil
	}
	if ei < 0 || ei >= len(m.SdsData) {
		ei = len(m.SdsData) - 1
	}
	limit := int(data.Limit)
	if limit == 0 {
		limit = 20
	}
	if ei-si >= limit {
		return m.SdsData[si : si+limit], nil
	} else {
		return m.SdsData[si : ei+1], nil
	}
}

func (m *MockMiotClient) FaceAnalyse(token string, did, file_id string, content_type int64, extra_info string) error {
	m._logAction("FaceAnalyse: did=%s file_id=%s conent_type=%s extra_info=%s",
		did, file_id, content_type, extra_info)
	return nil
}
func (m *MockMiotClient) GetFaces(token string, did string, face_ids []string) ([]*FaceInfoMeta, error) {
	res := make([]*FaceInfoMeta, 0, len(face_ids))
	for _, f := range face_ids {
		if fi, ok := m.FaceInfos[f]; ok {
			res = append(res, fi)
		}
	}
	return res, nil
}

func (m *MockMiotClient) SpecParsePropOp(token string, req *ParsePropOpReq) (*PropOp, error) {
	return nil, nil
}
func (m *MockMiotClient) SpecParseActionOp(token string, req *ParseActionOpReq) (*ActionOp, error) {
	return nil, nil
}
func (m *MockMiotClient) SpecParseEventOp(token string, req *ParseEventOpReq) (*EventOp, error) {
	return nil, nil
}

func (m *MockMiotClient) RpcDevice(token, did string, method string, params *simplejson.Json) (*RpcResult, error) {
	return &RpcResult{
		Code:    0,
		Message: "",
		Result:  simplejson.FromObject("ok"),
	}, nil
}

func (m *MockMiotClient) QueryIrThirdDb(token string, req *IrThirdDbBean) (*IrThirdDbBean, error) {
	return &IrThirdDbBean{RemoteId: req.RemoteId, ShortCmd: "11311701"}, nil
}

func (m *MockMiotClient) ReportEvent(token, did string, eventinfo []*Eventinfo) error {
	m._logAction("ReportEvent: did=%s eventinfo=%v",
		did, eventinfo)
	return nil
}

func (m *MockMiotClient) _sortSdsData() {
	sort.Slice(m.SdsData, func(i, j int) bool {
		cmp := m._compareSdsData(m.SdsData[i], m.SdsData[j])
		return cmp < 0
	})
}
func (m *MockMiotClient) _compareSdsData(iD, jD *SdsData) int {
	didComp := strings.Compare(iD.Did, jD.Did)
	if didComp != 0 {
		return didComp
	}

	typeComp := strings.Compare(iD.Type, jD.Type)
	if typeComp != 0 {
		return typeComp
	}

	keyComp := strings.Compare(iD.Key, jD.Key)
	if keyComp != 0 {
		return keyComp
	}
	// time is in desc order
	if iD.Time != jD.Time {
		if iD.Time > jD.Time {
			return -1
		}
		return 1
	}
	return 0
}

func (*MockMiotClient) _logAction(format string, args ...interface{}) {
	fmt.Printf("MockMiotClient: "+format, args...)
}
