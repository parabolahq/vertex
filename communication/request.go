package communication

import (
	"errors"
	"vertex/config"
)

// UserRequest Request body, sent by user to websocket pool. Contains info about service and request parameters
type UserRequest struct {
	ServiceAlias string                 `json:"serviceAlias"`
	MethodName   string                 `json:"methodName"`
	Params       map[string]interface{} `json:"params"`
}

// ServiceRequest Message body, sent to message queue by pool. Contains UserRequest, user and pool ids
type ServiceRequest struct {
	UserRequest UserRequest `json:"userRequest"`
	UserId      string      `json:"userId"`
	PoolId      string      `json:"poolId"`
	FromPool    bool        `json:"fromPool"`
}

func (r *UserRequest) ToServiceRequest(userId string) (ServiceRequest, error) {
	if r.MethodName == "" || r.ServiceAlias == "" {
		return ServiceRequest{}, errors.New("methodName or serviceAlias are empty")
	}

	return ServiceRequest{
		UserRequest: *r,
		UserId:      userId,
		PoolId:      config.Config.String("id"),
	}, nil
}

func PoolActionServiceRequest(serviceAlias, userId, actionName string) ServiceRequest {
	return ServiceRequest{
		UserRequest: UserRequest{
			ServiceAlias: serviceAlias,
			MethodName:   actionName,
			Params:       map[string]interface{}{},
		},
		UserId:   userId,
		PoolId:   config.Config.String("id"),
		FromPool: true,
	}
}
