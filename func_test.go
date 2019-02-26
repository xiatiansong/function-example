package main

import "testing"
import (
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	// **********
	// You can add your test data here
	// **********
	miotMock := &MockMiotClient{
		DeviceExtraKv: map[string]map[string]string{
			"did1": {
				"key1": "val1",
			},
		},
		DeviceEvents: map[string]map[string]string{
			"did1": {
				"key1": "val1",
			},
		},
		DeviceProps: map[string]map[string]string{
			"did1": {
				"key1": "val1",
			},
		},
		FaceInfos: map[string]*FaceInfoMeta{
			"face_id1": &FaceInfoMeta{},
		},
		SdsData: []*SdsData{
			{"did1", "event", "motion", 1, "val1"},
			{"did1", "event", "motion", 2, "val2"},
			{"did1", "event", "motion", 3, "val3"},
			{"did1", "event", "motion", 4, "val4"},
		},
	}
	miotMock._sortSdsData()
	miot = miotMock
}

func TestProp(t *testing.T) {
	Convey("test prop intput", t, func() {
		propInput := map[string]interface{}{
			"token": "testtoken",
			"data": map[string]interface{}{
				"did":    "testdid",
				"uid":    1234567,
				"model":  "mijia.test.v1",
				"pd_id":  999,
				"attr":   "prop.aqi",
				"status": 123,
			},
		}
		res, err := handle(propInput, test_context)
		So(err, ShouldBeNil)
		resS, _ := json.Marshal(res)
		fmt.Printf("\r\n%s\r\n", resS)
	})

}

func TestEvent(t *testing.T) {
	Convey("test event intput", t, func() {
		eventInput := map[string]interface{}{
			"token": "testtoken",
			"data": map[string]interface{}{
				"did":    "testdid",
				"uid":    1234567,
				"model":  "mijia.test.v1",
				"pd_id":  999,
				"attr":   "event.motion",
				"status": []string{"url1", "pwd1"},
			},
		}
		_, err := handle(eventInput, test_context)
		So(err, ShouldBeNil)
	})
}

var test_context = map[string]interface{}{
	"functionName": "test_func",
	"timeout":      30,
}

func TestIrSync(t *testing.T) {
	Convey("test sync ir", t, func() {
		propInput := map[string]interface{}{
			"token": "testtoken",
			"data": map[string]interface{}{
				"did": "85671432",
				"model": "lumi.acpartner.mcn02",
				"method": "get_ir_data",
				"params": map[string]interface{}{
					"did": "85671432",
					"model": "lumi.acpartner.mcn02",
					"brand_id": "97",
					"controller_id": "10727",
					"param_type": "ir_code",
					"zip_data": "nE5nExlk0mswmUsms5m8smM2nE4AIKayyazibABkAWYDBzUBgpzMQEGmYEDTcBhgFEAtMBAwHCAHMAMgMSAtYCBJgAgYHTgrsBigHLgHsD+4Qbg1OC64KETIFA5lLJtMQDCmwDFAaIBYYV7gsGGo4NCgcIGyoEhABkAgoMVgMIFmgf6g5xOZyCJoFZzAC+JrOAujDmUFRgDDFHcCgwjmD7cBNAQHADMCdwD2mwQhTeZAqmDCwDkCESAmk2mMsm0yBgYGiwVPA/sE+wuUCQMC+wA3EqkD+ZlLJnOZzNprLJxOZvNyQTCjUPUgRHEOMAsxUlH+UdJwpaHGYBgpsBzZIMBbyAgQhDhQWGR44zhaOAuwkEBkOFYgKDiu0WywNECi4SqoHChRnN5gBiYCKA5uBk4A5kDIB+4YrTmcgioIS4PtgZGXo5nth3Oa8xVsAfmPXYH7neOB+5LlAT+B/4H/g3uD/g5NgL+F/JBany0AgSDkHptMAACA",
				},
			},
		}
		res, err := handle(propInput, test_context)
		So(err, ShouldBeNil)
		resS, _ := json.Marshal(res)
		fmt.Printf("\r\n%s\r\n", resS)
	})

}
