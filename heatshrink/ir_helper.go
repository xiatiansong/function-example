package heatshrink

/*
#include "heatshrink_app.h"
#include "heatshrink_common.h"
#include "heatshrink_config.h"
#include "heatshrink_decoder.h"
#include "heatshrink_encoder.h"
#include "irzip.c"
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

//定义配置文件数据格式
type IrRemoteCfg struct {
	RemoteId       string
	Frequency      string
	CmdLen         string
	OnoffArray     string
	ModelArray     string
	TempArray      string
	WindSpeedArray string
}

var instance *map[string]IrRemoteCfg
var once sync.Once

/*
进行heatshrink热压缩编码
*/
func GoLumiHeatshrinkEncode(encode_in_buf string) string {
	cs := C.CString(encode_in_buf)
	encode_len := len(encode_in_buf)

	p := C.malloc(C.size_t(len(encode_in_buf)))
	defer C.free(p)

	ret := C.LumiHeatshrinkBase64Encode(cs, C.int(encode_len), (*C.char)(p))
	ret_len := (int)(ret)
	if ret_len > 0 && ret_len < encode_len {
		data := C.GoStringN((*C.char)(p), ret)
		return data
	}
	return ""
}

/*
进行heatshrink热压缩解码
*/
func GoLumiHeatshrinkDecode(decode_in_buf string) string {
	cs := C.CString(decode_in_buf)
	decode_len := len(decode_in_buf)

	p := C.malloc(C.size_t(decode_len * 10)) // 解码比例，不知道压缩比例，fixme
	defer C.free(p)

	ret := C.LumiHeatshrinkBase64Decode(cs, C.int(decode_len), (*C.char)(p))
	ret_len := (int)(ret)
	if ret_len > 0 {
		data := C.GoStringN((*C.char)(p), ret)
		return data
	}
	return ""
}

/*
对红外时间序列进行特征值提取
 */
func GoLumiIrZip(ir_in_buf string) (string, string) {
	cs := C.CString(ir_in_buf)
	defer C.free(unsafe.Pointer(cs))

	chararistorData := C.malloc(C.size_t(C.MAX_CHARARISTOR_UNIT*4 + 1))
	zipData := C.malloc(C.size_t(C.MAX_ZIP_UNIT*2 + 1))
	defer C.free(chararistorData)
	defer C.free(zipData)

	C.irzip(cs, (*C.char)(chararistorData), (*C.char)(zipData))
	chaRetData := C.GoString((*C.char)(chararistorData))
	zipRetData := C.GoString((*C.char)(zipData))
	return chaRetData, zipRetData
}

/*
按截取位数对压缩值进行截取
 */
func ExtractArray(cmdData string, temp string) string {
	var substring = ""
	indexArray := strings.Split(temp, ",");
	for _, value := range indexArray {
		index, _ := strconv.ParseInt(value, 10, 32)
		substring += string(cmdData[index-1])
	}
	return substring
}

/*
按截取位数对压缩值进行截取
 */
func ShortCmdToAcKey(shortCmd string) string {
	char01 := shortCmd[0:1]
	char12 := shortCmd[1:2]
	char23 := shortCmd[2:3]
	//char34 := shortCmd[3:4]
	char45 := shortCmd[4:5]
	char56 := shortCmd[5:6]
	var switchState, model, windSpeed, temperature string
	windDirect := "D0"
	//开关
	if char01 == "0" {
		switchState = "P0"
	} else if char01 == "1" {
		switchState = "P1"
	}
	//模式
	if char12 == "0" {
		//制热
		model = "M1";
	} else if char12 == "1" {
		//制冷
		model = "M0";
	} else if char12 == "2" {
		//自动
		model = "M2";
	} else if char12 == "3" {
		//除湿
		model = "M4";
	} else if char12 == "4" {
		//送风
		model = "M3";
	}
	//风速
	if char23 == "0" {
		windSpeed = "S1";
	} else if char23 == "1" {
		windSpeed = "S2";
	} else if char23 == "2" {
		windSpeed = "S3";
	} else if char23 == "3" {
		windSpeed = "S0";
	}
	//温度
	iChar45, _ := strconv.ParseInt(char45, 16, 32)
	iChar56, _ := strconv.ParseInt(char56, 16, 32)
	iTemp := iChar45*16 + iChar56
	if iTemp > -1 && iTemp < 241 {
		temperature = "T" + strconv.Itoa(int(iTemp))
	}

	return switchState + "_" + model + "_" + temperature + "_" + windSpeed + "_" + windDirect
}

/*
单例获取配置
 */
func GetInstance() *map[string]IrRemoteCfg {
	once.Do(func() {
		instance, _ = readFile("heatshrink/ir_sync_cfg.json")
	})
	return instance
}

func readFile(filename string) (*map[string]IrRemoteCfg, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return nil, err
	}
	var cfgMap = &map[string]IrRemoteCfg{}
	if err := json.Unmarshal(bytes, &cfgMap); err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		return nil, err
	}
	return cfgMap, nil
}
