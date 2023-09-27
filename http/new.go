package http

import "github.com/gin-gonic/gin"


func New () *gin.Engine {
	e := gin.New()
	e.Use(InjectTrace())
	e.Use(CLogger)
	e.Use(gin.Recovery())
	e.Use(cors())
	return e
}