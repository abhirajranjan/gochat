package jwtHandler

import (
	"net/http"

	"github.com/abhirajranjan/gochat/internal/api-service/serviceErrors"
	"github.com/gin-gonic/gin"
)

type IUserData interface {
	GetErr() string
	GetErrCode() int64
	GetMap() (map[string]interface{}, error)
}

// extract data from gin.Context (called from auth)
func (h *jwtHandler) HandleLoginRequest(c *gin.Context) (interface {
	GetMap() (map[string]interface{}, error)
	GetErr() string
	GetErrCode() int64
}, error) {

	h.logger.Debug(c.Get("username"))
	modelLoginRequest, err := generateLoginRequest(c)
	h.logger.Debugf("extracted login request: %#v", modelLoginRequest)

	if err != nil {
		h.logger.Warnf("failed to generate login request from context %e", err)
		return nil, serviceErrors.NewBindingErr("Bad Request")
	}

	if err := validateLoginRequest(modelLoginRequest); err != nil {
		h.logger.Warnf("failed to validate login request %e", err)
		return nil, err
	}
	h.logger.Debug("validate login request succesful")

	userData, err := h.dbHandler.GetUser(modelLoginRequest.GetUsername(), modelLoginRequest.GetPassword())
	if err != nil {
		h.logger.Warnf("error while verifing user: %s", err)
		return nil, serviceErrors.ErrInternalServer
	}

	h.logger.Debugf("user verification response: %#v", userData)

	if err := validateUserData(userData); err != nil {
		return nil, err
	}

	return userData, nil
}

func (h *jwtHandler) GeneratePayloadData(userData interface {
	GetMap() (map[string]interface{}, error)
}) map[string]interface{} {

	out := make(map[string]interface{})
	sessionID := h.dbHandler.GenerateSessionID(userData)
	err := h.payloadManager.Encode(userData, sessionID, out)

	if err != nil {
		h.logger.Errorf("fail to encode data (%#v) into payload: %s", userData, err)
	}

	if !h.dbHandler.ActivateSessionByID(sessionID) {
		h.logger.Errorf("cannot activate session %s", sessionID)
		return make(map[string]interface{})
	}

	return out
}

func (h *jwtHandler) ExtractPayloadData(claims map[string]interface{}) interface {
	Version() int64
	Get(string) (interface{}, bool)
	GetSessionID() interface{}
} {
	data, err := h.payloadManager.Decode(claims)
	if err != nil {
		h.logger.Debugf("error decoding data: %s", err)
		return nil
	}
	return data
}

func (h *jwtHandler) VerifyUser(data interface {
	Version() int64
	Get(string) (interface{}, bool)
	GetSessionID() interface{}
}, reqperm []string) bool {

	if data == nil {
		h.logger.Debug("empty data passed to verify user")
		return false
	}

	sessionID, ok := data.GetSessionID().(string)
	if !ok {
		return false
	}

	if sessionID == "" {
		return false
	}

	userperm, ok := h.dbHandler.GetUserRolesFromSession(sessionID)
	if !ok {
		h.logger.Debugf("invalid session accessed")
		// TODO force session clear
		return false
	}

	if !checkIfUserHasPermission(userperm, reqperm) {
		h.logger.Debugf("user does'nt have required permissions, required %v", reqperm)
		return false
	}

	return true
}

func (h *jwtHandler) LogoutUser(claims map[string]interface{}) int {
	data, err := h.payloadManager.Decode(claims)
	if err != nil {
		h.logger.Warnf("failed to decode claims")
		return http.StatusInternalServerError
	}

	sessionID, ok := data.GetSessionID().(string)
	if !ok {
		return http.StatusAccepted
	}

	if sessionID == "" {
		return http.StatusAccepted
	}

	if err := h.dbHandler.Logout(sessionID); err != nil {
		h.logger.Warnf("error logging out data due to database bridge error: %s", err)
		return http.StatusInternalServerError
	}
	h.logger.Debug("user logout successfully")
	return http.StatusOK
}
