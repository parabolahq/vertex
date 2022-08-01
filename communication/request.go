package communication

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
}
