package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gochat/config"
	"gochat/internal/core/domain"
	"gochat/internal/core/ports"
	"gochat/logger"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	mockserv *MockServices
	h        handler
	w        *httptest.ResponseRecorder
	c        *gin.Context
	writer   strings.Builder
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	mockserv = new(MockServices)
	writer = strings.Builder{}
	h = handler{
		logger:  MockLogger(&writer),
		service: mockserv,
	}
}

type MockServices struct {
	mock.Mock
}

var _ ports.Service = (*MockServices)(nil)

func (m *MockServices) LoginRequest(c domain.LoginRequest) (*domain.User, error) {
	args := m.Called(c)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockServices) HandleWS(w *websocket.Conn) error {
	args := m.Called(w)
	return args.Error(0)
}

func (m *MockServices) GetUserMessages(userId int64) ([]domain.ChannelBanner, error) {
	args := m.Called(userId)
	return args.Get(0).([]domain.ChannelBanner), args.Error(1)
}

func (m *MockServices) GetMessagesFromChannel(channelid int64) (*domain.ChannelMessages, error) {
	args := m.Called(channelid)
	return args.Get(0).(*domain.ChannelMessages), args.Error(1)
}

func (m *MockServices) PostMessageInChannel(channelid int64, message *domain.Message) (*domain.Message, error) {
	args := m.Called(channelid, message)
	return args.Get(0).(*domain.Message), args.Error(1)
}

func MockLogger(writer io.Writer) logger.ILogger {
	l := logger.NewLogger(config.LoggerConfig{
		AppName: "test",
		Level:   "debug",
		Dev:     true,
		Encoder: "console",
	})
	l.AddWriter(writer)
	l.InitLogger()
	return l
}

func TestGetUserMessages(t *testing.T) {
	tcs := []struct {
		name              string
		userid            int64
		isPostServiceTest bool
		serviceResp       []domain.ChannelBanner
		serviceRespErr    error
		wantStatusCode    int
	}{
		// pre service call tests
		{
			name:           "without userId context",
			userid:         0,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "with context userId of non type int",
			userid:         0,
			wantStatusCode: http.StatusInternalServerError,
		},
		// post service call tests
		{
			name:              "service failing",
			userid:            1,
			isPostServiceTest: true,
			wantStatusCode:    http.StatusInternalServerError,
			serviceRespErr:    errors.New("service failing"),
		},
		{
			name:              "no errors",
			userid:            2,
			isPostServiceTest: true,
			wantStatusCode:    http.StatusOK,
			serviceResp: []domain.ChannelBanner{
				{
					Id:   1,
					Name: "test",
				},
			},
		},
	}

	for _, tc := range tcs {
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Set("userId", tc.userid)

		t.Run(tc.name, func(t *testing.T) {

			if tc.isPostServiceTest {
				mockCall := mockserv.On("GetUserMessages", tc.userid).Return(tc.serviceResp, tc.serviceRespErr)
				h.GetUserMessages(c)
				mockserv.AssertExpectations(t)
				mockCall.Unset()
			} else {
				h.GetUserMessages(c)
			}

			if !assert.Equal(t, tc.wantStatusCode, w.Code) {
				t.Log(writer.String())
			}

			var (
				expectedbytes []byte
				actual        string
				err           error
			)

			if tc.serviceResp != nil {
				expectedbytes, err = json.Marshal(tc.serviceResp)
				if err != nil {
					t.Error(err)
				}
			}
			actual = w.Body.String()

			assert.Equal(t, string(expectedbytes), actual)
			writer.Reset()
		})
	}
}

func TestGetMessagesFromChannel(t *testing.T) {
	tcs := []struct {
		name              string
		channelid         string
		isPostServiceTest bool
		serviceResp       *domain.ChannelMessages
		serviceRespErr    error
		wantStatusCode    int
	}{
		{
			name:      "no channel id params",
			channelid: "",
		},
		{
			name:      "channel param not int",
			channelid: "not int",
		},
		{
			name:      "GetMessagesFromChannel service error",
			channelid: "1",
		},
		{
			name:      "no error",
			channelid: "1",
		},
	}

	for _, tc := range tcs {
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		if tc.channelid != "" {
			c.Params = append(c.Params, gin.Param{Key: "channelid", Value: tc.channelid})
		}

		if tc.isPostServiceTest {
			mockCall := mockserv.On("GetMessagesFromChannel", tc.channelid).Return(tc.serviceResp, tc.serviceRespErr)
			h.GetMessagesFromChannel(c)
			mockserv.AssertExpectations(t)
			mockCall.Unset()
		} else {
			h.GetMessagesFromChannel(c)
		}

		if !assert.Equal(t, tc.wantStatusCode, w.Code) {
			t.Log(writer.String())
		}

		var (
			expectedbytes []byte
			actual        string
			err           error
		)

		if tc.serviceResp != nil {
			expectedbytes, err = json.Marshal(tc.serviceResp)
			if err != nil {
				t.Error(err)
			}
		}
		actual = w.Body.String()

		assert.Equal(t, string(expectedbytes), actual)
		writer.Reset()
	}
}
