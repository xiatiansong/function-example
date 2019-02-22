package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync/atomic"
)

const (
	LOG_CHANNEL_MSG_TYPE_LOG = 0
	LOG_CHANNEL_MSG_TYPE_PC  = 1
)

const (
	FATAL = 0
	ERROR = 1
	WARN  = 2
	INFO  = 3
	DEBUG = 4
)

func LogFatal(format string, argv ...interface{}) {
	logToNamedPipe(FATAL, format, argv...)
}
func LogError(format string, argv ...interface{}) {
	logToNamedPipe(ERROR, format, argv...)
}
func LogWarn(format string, argv ...interface{}) {
	logToNamedPipe(WARN, format, argv...)
}
func LogInfo(format string, argv ...interface{}) {
	logToNamedPipe(INFO, format, argv...)
}
func LogDebug(format string, argv ...interface{}) {
	logToNamedPipe(DEBUG, format, argv...)
}

var traceIdVal atomic.Value
var logConn net.Conn
var inited int64

func logToNamedPipe(level int, format string, argv ...interface{}) {
	if inited <= 0 {
		return
	}
	var traceId string
	traceId, ok := traceIdVal.Load().(string)
	if !ok {
		traceId = "UNKNOWN"
	}

	var logData string
	if len(argv) == 0 {
		logData = format
	} else {
		logData = fmt.Sprintf(format, argv...)
	}

	result := map[string]interface{}{
		"trace_id": traceId,
		"l":        level,
		"message":  logData,
	}
	resultStr, _ := json.Marshal(result)

	// write a no meaning seqId here
	err := WriteMessage(logConn, append([]byte{LOG_CHANNEL_MSG_TYPE_LOG}, resultStr...), 1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "log to pipe err %v", err)
	}
}

const (
	PC_TYPE_COUNTER = 0
	PC_TYPE_COST    = 1
)

func pcIncr(key string, val uint32) {
	preBytes := make([]byte, 1 /*LOG_CHANNEL_MSG_TYPE_PC*/ +1 /*pc_type*/ +4 /*val*/, len(key)+6)
	preBytes[0] = LOG_CHANNEL_MSG_TYPE_PC
	preBytes[1] = PC_TYPE_COUNTER
	binary.BigEndian.PutUint32(preBytes[2:], val)
	err := WriteMessage(logConn, append(preBytes, []byte(key)...), 1)
	if err != nil {
		fmt.Fprintf(os.Stderr, " pc to pipe err %v", err)
	}
}

func pcCost(key string, val uint32) {
	preBytes := make([]byte, 1 /*LOG_CHANNEL_MSG_TYPE_PC*/ +1 /*pc_type*/ +4 /*val*/, len(key)+6)
	preBytes[0] = LOG_CHANNEL_MSG_TYPE_PC
	preBytes[1] = PC_TYPE_COST
	binary.BigEndian.PutUint32(preBytes[2:], val)
	err := WriteMessage(logConn, append(preBytes, []byte(key)...), 1)
	if err != nil {
		fmt.Fprintf(os.Stderr, " pc to pipe err %v", err)
	}
}
