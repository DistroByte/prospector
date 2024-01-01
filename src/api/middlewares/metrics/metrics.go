package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

// CreateMiddlewares creates the middlewares for the server
func CreateMiddlewares(r *gin.Engine) {
	// Add the metrics middleware

	MetricsMiddleware(r)
}

func MetricsMiddleware(r *gin.Engine) {
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")

	m.SetSlowTime(5)

	m.Use(r)
}
