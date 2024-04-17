package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

func MetricsMiddleware(r *gin.Engine) {
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	// set the threshold for a slow request to 1 second
	m.SetSlowTime(1)

	m.Use(r)
}
