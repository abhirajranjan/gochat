package handlers

import (
	"gochat/config"
	"gochat/logger"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func MockLogger(writer io.Writer) logger.ILogger {
	l := logger.NewLogger(config.LoggerConfig{
		AppName: "test",
		Level:   "debug",
		Dev:     true,
		Encoder: "json",
	})

	l.AddWriter(writer)
	return l
}

func TestGetUserMessages(t *testing.T) {
	buffer := strings.Builder{}
	h := handler{logger: MockLogger(&buffer)}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	t.Run("GetUserMessages without context", func(t *testing.T) {
		h.GetUserMessages(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
