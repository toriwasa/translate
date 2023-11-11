package completion

// 以下のJSONを元に構造体を定義する
// {
//     "error": {
//         "message": "Incorrect API key provided: sk-apZBs**************************************1cj4. You can find your API key at https://platform.openai.com/account/api-keys.",
//         "type": "invalid_request_error",
//         "param": null,
//         "code": "invalid_api_key"
//     }
// }

type ErrorDetail struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

type Error struct {
	Error ErrorDetail `json:"error"`
}
