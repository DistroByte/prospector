package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

func MetricsMiddleware(r *gin.Engine) {
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")

	m.SetSlowTime(5)

	m.Use(r)
}
