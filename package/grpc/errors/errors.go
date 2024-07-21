package errors

import (
	"encoding/json"
	"log"

	"bitbucket.org/revenuemonster/monster-api/kit/errcode"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Data :
type Data struct {
	StatusCode int32
	Code       string
	Message    string
}

func (d Data) String() string {
	x, _ := json.Marshal(d)
	return string(x)
}

// Parse :
func Parse(err error) Data {
	st, isOkey := status.FromError(err)
	if !isOkey || st.Code() != codes.Unknown {
		if st != nil {
			switch st.Code() {
			case codes.Unavailable, codes.DeadlineExceeded:
				return Data{
					StatusCode: 500,
					Code:       errcode.ServiceUnavailable,
					Message:    errcode.Message[errcode.ServiceUnavailable],
				}
			}
		}

		log.Println("fail to unmarshal error: ", st.Code(), st.Message(), st.Err())

		return Data{
			StatusCode: 500,
			Code:       "INTERNAL_ERROR",
			Message:    err.Error(),
		}
	}

	msg := Data{}
	err = json.Unmarshal([]byte(st.Message()), &msg)
	if err != nil {
		return Data{
			StatusCode: 500,
			Code:       "INTERNAL_ERROR",
			Message:    "Internal error",
		}
	}

	return msg
}

// Service :
func Service(statusCode int32, code string, message string) error {
	msg := Data{}
	msg.StatusCode = statusCode
	msg.Code = code
	msg.Message = message

	return status.New(codes.Unknown, msg.String()).Err()
}
