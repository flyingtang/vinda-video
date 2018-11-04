package response

import (
	"encoding/json"
	"net/http"
)

func SendErrorResponse(w http.ResponseWriter, sc int, msg string) {
	w.WriteHeader(sc)
	//io.WriteString(w, msg)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//http.Error(w, msg, sc)
	var data = map[string]interface{}{
		"message": msg,
	}
	// TODO 没有错误处理怎么办
	json.NewEncoder(w).Encode(data)
	return
}

func SendOKResponse(w http.ResponseWriter, sc int, d interface{}) {
	w.WriteHeader(sc)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//io.WriteString(w, msg)
	//http.Error(w, msg, sc)
	//var data = map[string]interface{}{
	//	"message": "执行成功",
	//	"data":    d,
	//}
	// TODO 没有错误处理怎么办
	json.NewEncoder(w).Encode(d)
	return
}
