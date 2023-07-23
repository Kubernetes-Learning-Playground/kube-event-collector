package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func HttpServer(port string) {

	r := gin.New()
	// 提供给prometheus的接口
	r.GET("/metrics", PrometheusHandler())

	err := r.Run(fmt.Sprintf(":%v", port))
	fmt.Println(err)
}

func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
