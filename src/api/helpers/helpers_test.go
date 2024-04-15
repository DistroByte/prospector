package helpers

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPerformRequest(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, r := gin.CreateTestContext(w)
		r.GET("/", func(c *gin.Context) {
			ctx.JSON(200, gin.H{"message": "ok"})
		})

		req := PerformRequest(r, "GET", "/")
		if req.Code != 200 {
			t.Errorf("expected 200, got %d", req.Code)
		}
	})
}

func TestPerformAuthRequest(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, r := gin.CreateTestContext(w)
		r.GET("/", func(c *gin.Context) {
			ctx.JSON(200, gin.H{"message": "ok"})
		})

		req := PerformAuthRequest(r, "GET", "/", "username", "password")
		if req.Code != 200 {
			t.Errorf("expected 200, got %d", req.Code)
		}
	})
}

func TestCheckJobHasValidName(t *testing.T) {
	tcs := []struct {
		name     string
		id       string
		expected bool
	}{
		{
			name:     "valid id",
			id:       "1234-prospector",
			expected: true,
		},
		{
			name:     "invalid id",
			id:       "1234",
			expected: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res := CheckJobHasValidName(tc.id)
			if res != tc.expected {
				t.Errorf("expected %t, got %t", tc.expected, res)
			}
		})
	}
}
