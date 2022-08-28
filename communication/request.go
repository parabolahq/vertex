package communication

import (
	"errors"
	"vertex/config"
)

// UserRequest Request body, sent by user to websocket pool. Contains info about service and request parameters
type UserRequest struct {
	Service string                 `json:"service"`
	Method  string                 `json:"method"`
	Version int                    `json:"version"`
	Data    map[string]interface{} `json:"argument"`
}

// ServiceRequest Message body, sent to message queue by pool. Contains UserRequest, user and pool ids
type ServiceRequest struct {
	UserRequest UserRequest `json:"userRequest"`
	UserId      string      `json:"userId"`
	PoolId      string      `json:"poolId"`
	FromPool    bool        `json:"fromPool"`
}

func (r *UserRequest) ToServiceRequest(userId string) (ServiceRequest, error) {
	if r.Method == "" || r.Service == "" {
		return ServiceRequest{}, errors.New("method or service are empty")
	}

	return ServiceRequest{
		UserRequest: *r,
		UserId:      userId,
		PoolId:      config.Config.String("id"),
	}, nil
}

func PoolActionServiceRequest(serviceAlias, userId, method string) ServiceRequest {
	return ServiceRequest{
		UserRequest: UserRequest{
			Service: serviceAlias,
			Method:  method,
			Data:    map[string]interface{}{},
		},
		UserId:   userId,
		PoolId:   config.Config.String("id"),
		FromPool: true,
	}
}
