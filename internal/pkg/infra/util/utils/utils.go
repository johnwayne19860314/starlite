package utils

import (
	"encoding/json"
	"net/http"
	"regexp"

	uuid "github.com/satori/go.uuid"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetRequestId(req *http.Request) (xRequestId string) {
	return GetRequestIdByKey(req, "")
}

func GetRequestIdByKey(req *http.Request, requestIdKey string) (xRequestId string) {
	if requestIdKey == "" {
		requestIdKey = "txid"
	}
	if txid := req.Header.Get(requestIdKey); txid != "" {
		xRequestId = txid
	} else {
		xRequestId = uuid.NewV4().String()
	}
	return
}

// MatchPattern takes a string and a pattern and returns true if the string
// matches the pattern according to regular expression matching. It returns
// false and an error if the pattern is invalid.
func MatchPattern(input, pattern string) (bool, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, errorx.Errorf("invalid pattern: %v", err)
	}
	return re.MatchString(input), nil
}

// Jsonify function takes an interface and returns a JSON representation in bytes
// and any error encountered.
func Jsonify(v interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return nil, errorx.Errorf("json marshaling failed: %v", err)
	}
	return jsonData, nil
}

func DetermineFailReason(err error) int {
	// parse gRPC error and set the corresponding FailReason
	var failReason int
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			switch s.Code() {
			case grpccodes.InvalidArgument:
				failReason = 3 // invalid argument
			case grpccodes.PermissionDenied:
				failReason = 7 // permission denied
			case grpccodes.Unavailable:
				failReason = 14 // Unavailable
			case grpccodes.NotFound:
				failReason = 5 // not found
			case grpccodes.FailedPrecondition:
				failReason = 9 // failed precondition
			case grpccodes.ResourceExhausted:
				failReason = 8 // resource exhausted
			default:
				failReason = 13 // internal errors
			}
		} else {
			failReason = 13 // internal error
		}
	} else {
		failReason = 0
	}
	return failReason
}
