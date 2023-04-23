package helpers

type Response struct {  
	Meta Meta        	`json:"meta"`
	Data interface{} 	`json:"data"`

}

type Meta struct {
	Message string  `json:"message"`
	Status  string  `json:"status"`
	Code    int     `json:"code"`
}

func ApiResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Status:  status,
		Code:    code,
	}

	jsonresponse := Response{
		Meta: meta,
		Data: data,
	} 

	return jsonresponse
}

	