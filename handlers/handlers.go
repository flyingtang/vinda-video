package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"video/constant"
	"video/models"
	"video/response"
)

// 播放视频
func StreamHandle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// 路径
	vid := ps.ByName("id")
	//videoPath := constant.Video_dir + vid

	v, err := models.FindVideoById(vid)
	if err != nil {
		log.Fatalf("invalid vidio id err : %s", err.Error())
		response.SendErrorResponse(w, http.StatusBadRequest, "无效的视频ID")
		return
	}
	if len(v.UrlName) == 0 {
		log.Fatal("empty v.UrlName")
		response.SendErrorResponse(w, http.StatusInternalServerError, "无效的视频")
		return
	}
	fp := filepath.Join(constant.Video_dir, v.UrlName)
	// 打开
	video, err := os.Open(fp)
	if err != nil {
		http.Error(w, "no find video", http.StatusInternalServerError)
		return
	}
	// 发送
	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, "", time.Now(), video)
	// 关闭
	defer video.Close()
}

func CreateVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var v models.Video
	if err := parseParma(r, &v); err != nil {
		log.Fatalf("parseParma(r, &v); err %s", err.Error())
		response.SendErrorResponse(w, http.StatusBadRequest, "解析参数失败")
		return
	}

	err := models.CreateVideo(&v)
	if err != nil {
		log.Fatalf("models.CreateVideo err %s", err.Error())
		response.SendErrorResponse(w, http.StatusBadRequest, "解析参数失败")
		return
	}
	response.SendOKResponse(w, http.StatusCreated, map[string]interface{}{"message": "创建成功"})
	return
}

func PatchVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	var v models.Video
	if err := parseParma(r, &v); err != nil {
		log.Fatalf("parseParma(r, &v); err %s", err.Error())
		response.SendErrorResponse(w, http.StatusBadRequest, "解析参数失败")
		return
	}

	err := models.PatchVideo(id, &v)
	if err != nil {
		log.Fatalf("parseParma(r, &v); err %s", err.Error())
		response.SendErrorResponse(w, http.StatusBadRequest, "解析参数失败")
		return
	}
	response.SendOKResponse(w, http.StatusOK, map[string]interface{}{"message": "更新成功"})
}

func parseParma(r *http.Request, m interface{}) (err error) {
	var limit int64 = constant.FormLimit // 1M 最大接受
	bdata, err := ioutil.ReadAll(io.LimitReader(r.Body, limit))
	if err != nil {
		return
	}
	if err = r.Body.Close(); err != nil {
		return
	}
	err = json.Unmarshal(bdata, m)
	return
}

func DeleteVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	id := ps.ByName("id")
	err := models.DeleteVideo(id)
	if err != nil {
		log.Fatalf("models.DeleteVideo err %s", err.Error())
		response.SendErrorResponse(w, http.StatusBadRequest, "服务错误")
		return
	}
	response.SendOKResponse(w, http.StatusOK, map[string]interface{}{"message": "删除成功"})
}

// 仅仅上传视频
func UploadHandle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	http.MaxBytesReader(w, r.Body, constant.MaxSize)
	if err := r.ParseMultipartForm(constant.MaxSize); err != nil {
		http.Error(w, "file is too big", http.StatusBadRequest)
		return
	}
	// 中间参数可用验证content
	file, h, err := r.FormFile("file")
	if err != nil {
		log.Printf("parse file err: %s", err.Error())
		response.SendErrorResponse(w, http.StatusInternalServerError, "internal server error ")
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("read file err: %s", err.Error())
		response.SendErrorResponse(w, http.StatusInternalServerError, "internal server error ")
	}
	// 生成随机唯一字符串做名字
	fn := GetRandomNumber()
	//fn := ps.ByName("id")
	err = ioutil.WriteFile(filepath.Join(constant.Video_dir, fn), data, 0666)
	if err != nil {
		log.Printf("write file err: %s", err.Error())
		response.SendErrorResponse(w, http.StatusInternalServerError, "internal server error ")
	}
	res := map[string]interface{}{
		"message":    "上传成功",
		"urlName":    fn,
		"originName": h.Filename,
	}
	response.SendOKResponse(w, http.StatusCreated, res)
	//models.CreateVideo()
	return
}

type Filter struct {
	Where map[string]interface{} `json:"where"`
	Page  int64                  `json:"page"`
	Order string                 `json:"order"`
}

type ids struct {
	Ids []int `json:"ids"`
}
func DeleteAllVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	var i ids
	err := parseParma(r, &i)
	if err != nil {
		log.Fatalf("parseParma(r, &i) err %s", err.Error())
		response.SendErrorResponse(w, http.StatusBadRequest, "服务错误")
		return
	}
	err = models.DeleteALLVideo(i.Ids)
	if err !=nil{
		log.Fatalf("models.DeleteALLVideo err %s", err.Error())
		response.SendErrorResponse(w, http.StatusBadRequest, "服务错误")
		return
	}
	response.SendOKResponse(w, http.StatusOK, map[string]interface{}{"message": "删除成功"})
}

// 查询视频列表，默认10条
func QueryVideoList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filterStr := r.FormValue("filter")
	var filter Filter
	if err := json.Unmarshal([]byte(filterStr), &filter); err != nil {
		response.SendErrorResponse(w, http.StatusInternalServerError, "服务错误")
		return
	}
	if filter.Page < 1 {
		filter.Page = 1
	}
	data, total, err := models.QueryVideos(filter.Page)
	if err != nil {
		log.Fatalf("models.QueryVideos(page) err %s", err.Error())
		response.SendErrorResponse(w, http.StatusInternalServerError, "服务错误")
		return
	}

	response.SendOKResponse(w, http.StatusOK, map[string]interface{}{"data": data, "limitRate": constant.LimiterRate ,"pageSize":constant.PageSize, "total": total})
}

// 测试上传视频前端模板

func TestHanle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	//	const tpl = `
	//		<!DOCTYPE html>
	//<html lang="en">
	//<head>
	//    <meta charset="UTF-8">
	//    <title>upload</title>
	//</head>
	//<body>
	//    <form enctype="multipart/form-data" action="http://127.0.0.1:4001/upload/123" method="post">
	//        <input type="file" name="file">
	//        <button type="submit" value="input file" >上传</button>
	//    </form>
	//</body>
	//</html>
	//	`
	//	t := template.New("new template")
	//	t , _ = t.ParseFiles("upload.html")
	t, _ := template.ParseFiles("html/upload.html")
	t.Execute(w, nil)
}

// 用md5生成一个随机字符串
func GetRandomNumber() string {
	h := md5.New()
	t := time.Now().Unix()
	i := rand.Int63()
	is := strconv.FormatInt(i, 10)

	// 三次失败尝试
	_, err := h.Write([]byte(strconv.FormatInt(t, 10)))
	if err != nil {
		_, err := h.Write([]byte(strconv.FormatInt(t, 10)))
		if err != nil {
			h.Write([]byte(strconv.FormatInt(t, 10)))
		}
	}
	return hex.EncodeToString(h.Sum([]byte(is)))
}
