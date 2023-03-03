package grpcHandler

import (
	"context"
	"net/http"

	"github.com/abhirajranjan/gochat/internal/api-service/serviceErrors"
	"github.com/gin-gonic/gin"
)

func (h *grpcHandler) HandleLoginRequest(modelLoginRequest ILoginRequest) (ILoginResponse, error) {
	ctx := context.Background()

	if err := validateLoginRequest(modelLoginRequest); err != nil {
		handleValidationErr(err, h.logger, "handle.ValidationError")
		return nil, err
	}

	grpcLoginRequest := modelLoginReqToGrpcLoginReq(modelLoginRequest)

	grpcLoginRes, err := h.grpc.VerifyUser(ctx, grpcLoginRequest)
	if err != nil {
		handleGrpcServerError(err, h.logger, "handle.handleLoginRequest")
		return nil, serviceErrors.ErrInternalServer
	}

	modelLoginRes := grpcLoginResToModelRes(grpcLoginRes)

	switch modelLoginRes.GetErrCode() {
	case http.StatusInternalServerError:
		handleInternalGrpcEndpointError(err, h.logger, "handler.handleLoginRequest")
		return nil, serviceErrors.ErrInternalServer
	}

	return modelLoginRes, nil
}

// generate login request from gin.Context
//
// arguments :
// - *gin.Context, has a models.LoginRequest
//
// returns :
// - loginrequest, if validation matches else write StatusBadRequest and returns nil
//
// check only binding errors
func (h *grpcHandler) GenerateLoginRequest(c *gin.Context) (ILoginRequest, error) {
	var request LoginRequest
	err := c.ShouldBind(request)

	if err == nil {
		return &request, nil
	}

	bindingerr := handleBindingErr(err)
	return nil, bindingerr
}

func (h *grpcHandler) GeneratePayloadData(modelLoginResponse ILoginResponse) map[string]interface{} {
	a, _ := GenerateMap(modelLoginResponse).(map[string]interface{})
	return a
}

func (h *grpcHandler) ExtractPayloadData(claims map[string]interface{}) IPayloadData {
	payload := loginResponse{}
	return &payload
}
