package grpcHandler

import (
	"context"
	"net/http"

	"github.com/abhirajranjan/gochat/internal/api-service/model"
	"github.com/abhirajranjan/gochat/internal/api-service/serviceErrors"
	"github.com/gin-gonic/gin"
)

func (h *grpcHandler) HandleLoginRequest(modelLoginRequest model.ILoginRequest) (model.ILoginResponse, error) {
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
func (h *grpcHandler) GenerateLoginRequest(c *gin.Context) (model.ILoginRequest, error) {
	var request LoginRequest
	err := c.ShouldBind(request)

	if err == nil {
		return &request, nil
	}

	bindingerr := handleBindingErr(err)
	return nil, bindingerr
}

func (h *grpcHandler) GeneratePayloadData(modelLoginResponse model.ILoginResponse) (out map[string]interface{}) {
	GenerateMap(modelLoginResponse, out)
	_, err := h.payloadManager.AddPayload(out)
	if err != nil {
		h.logger.Error(serviceErrors.NewStandardErr("handler.GeneratePayloadData", "payload function failed to parse", out, err))
	}
	return
}

func (h *grpcHandler) ExtractPayloadData(claims map[string]interface{}) model.IPayloadData {
	version, ok := claims["ver"]
	if !ok {
		h.logger.Infof("[handler.ExtractPayloadData] %s %v", "jwt missing version", claims)
		return nil
	}
	ver, ok := version.(int64)
	if !ok {
		h.logger.Infof("[handler.ExtractPayloadData] %s %v", "jwt has unknown type version", version)
		return nil
	}
	data, err := h.payloadManager.Decode(claims, ver)
	if err != nil {
		h.logger.Infof("[handler.ExtractPayloadData] %s %s", "error decoding jwt data", err)
		return nil
	}
	return data
}

func (h *grpcHandler) VerifyPayloadData(data model.IPayloadData) {
	h.payloadManager.To_Proto(data)
}
