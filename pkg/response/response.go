package response

import (
	"context"
	"time"
)

type Response struct {
	Data   interface{} `json:"data,omitempty"`
	Status interface{} `json:"status,omitempty"`
}

// Status -
type Status struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

// Error -
func Error(ctx context.Context, err error) interface{} {
	return withStatus(ctx, nil, err)
}

// Resp -
func Success(ctx context.Context, data interface{}) interface{} {
	return withStatus(ctx, data, nil)
}

func withStatus(ctx context.Context, data interface{}, e error) interface{} {
	if e != nil {
		t := time.Now().In(time.Local)
		return Response{
			Status: Status{
				Message: e.Error(),
				Time:    t.Format(time.RFC3339),
			},
		}
	}

	if data == nil {
		return Response{Status: "ok"}
	}

	return data
}
