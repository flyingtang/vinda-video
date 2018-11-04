package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"video/constant"
	"video/handlers"
	_ "video/models"
)

type middleWareHandler struct {
	r  *httprouter.Router
	cl *handlers.ConnLimiter
}

func newMiddleWareHandler(r *httprouter.Router, cc int) middleWareHandler {
	m := middleWareHandler{}
	m.r = r
	m.cl = handlers.NewConnLimiter(cc)
	return m
}

func (mh middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !mh.cl.GetConn() {
		http.Error(w, "too many connection", http.StatusTooManyRequests)
		return
	}
	mh.r.ServeHTTP(w, r)
	defer mh.cl.ReleaseConn()
}

func registerHandlers() *httprouter.Router {

	router := httprouter.New()

	router.GET("/api/v1/video", handlers.QueryVideoList)
	router.POST("/api/v1/video", handlers.CreateVideo)
	router.PATCH("/api/v1/video/:id", handlers.PatchVideo)
	router.GET("/api/v1/video/:id", handlers.StreamHandle)
	router.DELETE("/api/v1/video/:id", handlers.DeleteVideo)
	router.DELETE("/api/v1/video", handlers.DeleteAllVideo)
	router.POST("/api/v1/video/upload", handlers.UploadHandle)

	router.GET("/testPage", handlers.TestHanle)
	return router
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("error occur errr: %s", err)
		}
	}()
	r := registerHandlers()
	mh := newMiddleWareHandler(r, constant.LimiterRate)
	http.ListenAndServe(":3001", mh)
}
