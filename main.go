package main

import "fmt"

func main() {
	GetInstance()
	fmt.Println((*GetInstance())["2"])
	chaRetData, zipRetData := GoLumiIrZip("3068,3068,3068,4524,546,1664,546,546,546,1664,546,546,546,546,546,1664,546,1664,546,546,546,546,546,1664,546,1664,546,1664,546,546,546,1664,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,1664,546,546,546,546,546,546,546,546,546,546,546,546,546,1664,546,546,546,1664,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,1664,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,546,1664,546,546,546,1664,546,546,546,546,546,546,546,1664,546,1664,546,1664,546,1664,546,1664,546,988")
	fmt.Println(chaRetData)
	fmt.Println(zipRetData)
	fmt.Println(extractArray("22230100010000010100000101010001000000000000000000000000000000000000000100000000000001000100000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000010001000000010101010104", "38,40,46,48,156,158,164,166"))
	/**
	stdin := os.Stdin
	reader := bufio.NewReader(stdin)
	stdout := os.Stdout
	os.Stdout = os.Stderr
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "argument invalid")
	}
	logApi := os.Args[1]
	conn, err := net.Dial("unix", logApi)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect log api failed %v", err)
		return
	}
	logConn = conn
	atomic.AddInt64(&inited, 1)

	defer conn.Close()
	for {
		msg, seqId, err := ReadMessage(reader, -1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			return
		}

		req := map[string]interface{}{}
		err = json.Unmarshal(msg, &req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", "input not json")
			return
		}
		res := map[string]interface{}{}
		event, _ := req["event"].(map[string]interface{})
		context, _ := req["context"].(map[string]interface{})
		traceId, ok := event["trace_id"].(string)
		if ok && traceId != "" {
			traceIdVal.Store(traceId)
		}
		result, err := handle(event, context)
		if err != nil {
			res["code"] = -4
			res["message"] = err.Error()
		} else {
			res["code"] = 0
		}
		res["result"] = result
		os.Stdout = os.Stderr
		res["trace_id"] = traceId
		resStr, err := json.Marshal(res)
		if err != nil {
			resStr = []byte(`{"error":"cannot marshal res to json"}`)
		}
		err = WriteMessage(stdout, resStr, int64(seqId))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			return
		}
	}
	**/
}
