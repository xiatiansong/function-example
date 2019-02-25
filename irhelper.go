package main

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

func extractArray(cmdData string, temp string) string {
	var substring = ""
	indexArray := strings.Split(temp, ",");
	for _, value := range indexArray {
		index, _ := strconv.ParseInt(value, 10, 32)
		substring += string(cmdData[index-1])
	}
	return substring
}

func GetInstance() *map[string]IrRemoteCfg {
	once.Do(func() {
		instance, _ = readFile("irsynccfg.json")
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
