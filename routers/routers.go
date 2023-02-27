package routers

import (
	"github.com/gin-gonic/gin"
	"utils/job"
)

func IndexInit() *gin.Engine {
	r := gin.New()

	r.GET("/jobs", job.GetJobs)
	r.POST("/jobs", job.AddJob)
	r.DELETE("/jobs", job.DeleteJob)

	return r
}
