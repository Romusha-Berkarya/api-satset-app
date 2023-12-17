package entity

type Response struct {
	Data interface{} `json:"Data"`
}

type ResponseMessage struct {
	Message string `json:"Message"`
}
