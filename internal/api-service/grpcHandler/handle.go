package grpcHandler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/abhirajranjan/gochat/internal/api-service/model"
	"github.com/abhirajranjan/gochat/internal/api-service/serviceErrors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (h *grpcHandler) HandleLoginRequest(c *gin.Context) (model.ILoginResponse, error) {
	modelLoginRequest, err := generateLoginRequest(c)
	if err != nil {
		h.logger.Error(err)
		return nil, serviceErrors.NewBindingErr("Bad Request")
	}
	ctx := context.Background()

	if err := validateLoginRequest(modelLoginRequest); err != nil {
		h.logger.Error(err)
		return nil, err
	}

	grpcLoginRequest := modelLoginReqToGrpcLoginReq(modelLoginRequest)

	// TODO: add retry in case of service failure
	grpcLoginRes, err := h.grpc.VerifyUser(ctx, grpcLoginRequest)
	if err != nil {
		h.logger.Error(err)
		return nil, serviceErrors.ErrInternalServer
	}

	modelLoginRes := grpcLoginResToModelRes(grpcLoginRes)

	switch modelLoginRes.GetErrCode() {
	case http.StatusInternalServerError:
		handleInternalGrpcEndpointError(err, h.logger, "handler.handleLoginRequest")
		return nil, serviceErrors.ErrInternalServer
	case http.StatusOK:
		return modelLoginRes, nil
	default:
		h.logger.Warnf("unknown error from grpc Endpoint: %s", err)
		return nil, serviceErrors.ErrInternalServer
	}
}

func (h *grpcHandler) GeneratePayloadData(modelLoginResponse model.ILoginResponse) map[string]interface{} {
	out := map[string]interface{}{}
	if err := GenerateMap(modelLoginResponse, out); err != nil {
		h.logger.Error(errors.Wrap(err, "generate map failed, using classic marshal method"))
		out = map[string]interface{}{}
		b, err := json.Marshal(modelLoginResponse)
		if err != nil {
			// ? what Now ?
			h.logger.Error(errors.Wrap(err, "classic method failed"))
		}
		if err := json.Unmarshal(b, &out); err != nil {
			// ? what Now ?
			h.logger.Error(err)
			h.logger.Error(errors.Wrap(err, "classic method failed"))
		}
	}

	out, err := h.payloadManager.AddPayload(out)
	if err != nil {
		h.logger.Error(serviceErrors.NewStandardErr("handler.GeneratePayloadData", "payload function failed to parse", out, err))
		out, err := h.payloadManager.Encode(out, h.payloadManager.GetMinimumVersion())
		if err != nil {
			h.logger.Error(serviceErrors.NewStandardErr("handler.GeneratePayloadData", "payload function minimum version failed to parse", out, err))
		}
	}
	return out
}

func (h *grpcHandler) ExtractPayloadData(claims map[string]interface{}) model.IPayloadData {
	ver, err := GetVersionFromClaims(claims)
	if err != nil {
		h.logger.Info(errors.Wrap(err, "handler.ExtractPayloadData"))
		return nil
	}

	data, err := h.payloadManager.Decode(claims, ver)

	if err != nil {
		h.logger.Info(serviceErrors.NewStandardErr("handler.ExtractPayloadData", "error decoding jwt data", err))
		return nil
	}
	return data
}

func (h *grpcHandler) VerifyUser(data model.IPayloadData) bool {
	return h.payloadManager.VerifyUser(data)
}

func (h *grpcHandler) LogoutUser(claims map[string]interface{}) int {
	ver, err := GetVersionFromClaims(claims)
	if err != nil {
		h.logger.Info(errors.Wrap(err, "handler.LogoutUser"))
		return http.StatusUnprocessableEntity
	}
	if !h.payloadManager.LogoutUser(claims, ver) {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}
