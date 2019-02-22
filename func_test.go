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
