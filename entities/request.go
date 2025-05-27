package entities

type StreamResponse struct {
	Data       []Stream `json:"data"`
	Pagination Pagination
}

type Pagination struct {
	Cursor string `json:"cursor"`
}

func NewStreamResponse() *StreamResponse {
	return &StreamResponse{
		Data:       make([]Stream, 0),
		Pagination: Pagination{},
	}
}
