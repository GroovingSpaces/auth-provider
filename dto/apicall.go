package dto

type RequestAPICallResult struct {
	RequestURL         string
	Method             string
	RequestLatency     string
	RequestBody        string
	RequestQuery       string
	ResponseBody       string
	RequestHeaders     string
	ResponseHeaders    string
	ResponseStatusCode int
}
