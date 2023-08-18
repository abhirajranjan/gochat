package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"gochat/internal/core/domain"
	"gochat/internal/core/ports"

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
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	mockserv = new(MockServices)
	h = handler{
		service: mockserv,
	}
}

type MockServices struct {
	mock.Mock
}

var _ ports.Service = (*MockServices)(nil)

func (m *MockServices) HandleWS(w *websocket.Conn) error {
	args := m.Called(w)
	return args.Error(0)
}

func (m *MockServices) LoginRequest(c domain.LoginRequest) (*domain.User, error) {
	args := m.Called(c)
	return args.Get(0).(*domain.User), args.Error(1)
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
				mockCall := mockserv.On("GetUserMessages", tc.userid).
					Return(tc.serviceResp, tc.serviceRespErr)
				h.GetUserMessages(c)
				mockserv.AssertExpectations(t)
				mockCall.Unset()
			} else {
				h.GetUserMessages(c)
			}

			if !assert.Equal(t, tc.wantStatusCode, w.Code) {
				t.FailNow()
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
		})
	}
}

func TestGetMessagesFromChannel(t *testing.T) {
	tcs := []struct {
		name      string
		channelid string

		// if post service test then channelid must be convertable to int64
		isPostServiceTest bool
		serviceResp       *domain.ChannelMessages
		serviceRespErr    error
		wantStatusCode    int
	}{
		{
			name:           "no channel id params",
			channelid:      "",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "channel param not int",
			channelid:      "not int",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:              "service error",
			channelid:         "1",
			isPostServiceTest: true,
			serviceRespErr:    errors.New("service error"),
			wantStatusCode:    http.StatusInternalServerError,
		},
		{
			name:              "no error",
			channelid:         "1",
			isPostServiceTest: true,
			serviceResp: &domain.ChannelMessages{
				Id: 2,
				Messages: []domain.Message{
					{
						Id:      1,
						UserId:  1023,
						At:      time.Now(),
						Type:    domain.MessageTypeText,
						Content: []byte("test message"),
					},
				},
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tc := range tcs {
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		t.Run(tc.name, func(t *testing.T) {
			if tc.channelid != "" {
				c.Params = append(c.Params, gin.Param{Key: "channelid", Value: tc.channelid})
			}

			if tc.isPostServiceTest {
				channelid, err := strconv.Atoi(tc.channelid)
				if err != nil {
					t.Error(errors.Errorf("channelid cannot be converted to type int %+v", tc.channelid))
				}

				mockCall := mockserv.On("GetMessagesFromChannel", int64(channelid)).
					Return(tc.serviceResp, tc.serviceRespErr)
				h.GetMessagesFromChannel(c)
				mockserv.AssertExpectations(t)
				mockCall.Unset()
			} else {
				h.GetMessagesFromChannel(c)
			}

			if !assert.Equal(t, tc.wantStatusCode, w.Code) {
				t.FailNow()
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

			if !assert.Equal(t, string(expectedbytes), actual) {
				t.FailNow()
			}
		})
	}
}

func TestPostMessageInChannel(t *testing.T) {
	tcs := []struct {
		name              string
		userid            int64
		channelid         string
		isPostServiceTest bool
		serviceResp       *domain.Message
		serviceRespErr    error
		wantStatusCode    int
		postmessage       *domain.Message
	}{
		// pre service call tests
		{
			name:           "without userId context",
			userid:         0,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "no channel id params",
			userid:         1,
			channelid:      "",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "channel param not int",
			userid:         1,
			channelid:      "not int",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "message userid not equals auth userid",
			userid:         1,
			channelid:      "2",
			wantStatusCode: http.StatusBadRequest,
			postmessage: &domain.Message{
				UserId:  2,
				At:      time.Now(),
				Type:    domain.MessageTypeText,
				Content: []byte("test message"),
			},
		},
		{
			name:              "service error",
			userid:            1,
			channelid:         "1",
			isPostServiceTest: true,
			wantStatusCode:    http.StatusInternalServerError,
			serviceRespErr:    errors.New("service error"),
			postmessage: &domain.Message{
				UserId:  1,
				At:      time.Now(),
				Type:    domain.MessageTypeText,
				Content: []byte("test message"),
			},
		},
		{
			name:              "no error",
			userid:            1,
			channelid:         "1",
			isPostServiceTest: true,
			wantStatusCode:    http.StatusOK,
			postmessage: &domain.Message{
				UserId:  1,
				At:      time.Now(),
				Type:    domain.MessageTypeText,
				Content: []byte("test message"),
			},
			serviceResp: &domain.Message{
				Id:      2,
				UserId:  1,
				At:      time.Now(),
				Type:    domain.MessageTypeText,
				Content: []byte("test message"),
			},
		},
	}

	for _, tc := range tcs {
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		if tc.postmessage != nil {
			b, err := json.Marshal(tc.postmessage)
			if err != nil {
				t.Error(err)
			}

			byteReaderCloser := io.NopCloser(bytes.NewReader(b))
			c.Request = &http.Request{
				Body: byteReaderCloser,
			}
			c.Header("Content-Type", "text/json")
		}

		if tc.userid != 0 {
			c.Set("userId", tc.userid)
		}

		if tc.channelid != "" {
			c.Params = append(c.Params, gin.Param{Key: "channelid", Value: tc.channelid})
		}

		t.Run(tc.name, func(t *testing.T) {
			if tc.isPostServiceTest {
				a, _ := strconv.Atoi(tc.channelid)
				mockcall := mockserv.On("PostMessageInChannel", int64(a), mock.Anything).
					Return(tc.serviceResp, tc.serviceRespErr)

				h.PostMessageInChannel(c)

				if !mockserv.AssertExpectations(t) {
					t.FailNow()
				}

				mockcall.Unset()
			} else {
				h.PostMessageInChannel(c)
			}

			if !assert.Equal(t, tc.wantStatusCode, w.Code) {
				t.FailNow()
			}

			var (
				expectedbytes []byte
				actual        string
				err           error
			)

			if tc.serviceResp != nil {
				expectedbytes, err = json.Marshal(tc.serviceResp.Id)
				if err != nil {
					t.Error(err)
				}
			}

			actual = w.Body.String()

			if !assert.Equal(t, string(expectedbytes), actual) {
				t.Logf("%s %s", string(expectedbytes), actual)
				t.FailNow()
			}
		})
	}
}
