package utils

type TRespDataSet struct {
	Total   int         `json:"total,omitempty"`
	Fields  []string    `json:"fields,omitempty"`
	ArrData interface{} `json:"list,omitempty"`
}
type TResponse struct {
	Code int32         `json:"code"`
	Data *TRespDataSet `json:"data,omitempty"`
	Info string        `json:"message"`
}

type TAuth struct {
	Code  int8   `json:"code"`
	User  any    `json:"user"`
	Token string `json:"token"`
}

func Failure(info string) *TResponse {
	var resp TResponse
	resp.Code = -1
	resp.Data = nil
	resp.Info = info
	return &resp
}

func ReturnID(id int32) *TResponse {
	var resp TResponse
	resp.Code = id
	resp.Data = nil
	resp.Info = "success"
	return &resp
}

func Success(data *TRespDataSet) *TResponse {
	var resp TResponse
	resp.Code = 0
	resp.Data = data
	resp.Info = "success"
	return &resp
}

func Authentication(token string, data any) *TAuth {
	var resp TAuth
	resp.Code = 0
	resp.User = data
	resp.Token = token
	return &resp
}
func RespData(total int32, fields []string, data interface{}, err error) *TResponse {
	if err != nil {
		return Failure(err.Error())
	}
	var dataSet TRespDataSet
	dataSet.Total = int(total)
	dataSet.Fields = fields
	dataSet.ArrData = data
	return Success(&dataSet)
}
