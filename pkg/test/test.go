package test

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
)

const floatingPointTolerance = 1e-6

type RespPair struct {
	StatusCode int
	Response   []byte
}

type PathResp map[string]RespPair

type URIMatcher struct {
	Params       []string
	ResponsePair RespPair
}

type MockServer struct {
	Server  *httptest.Server
	Handler http.Handler
	Targets []*string
}

func SetTargetUrl(url string, targets ...*string) {
	// Add slash to ensure we are hitting the test provided url
	url = url + "/"
	for _, target := range targets {
		*target = url
	}
}

// SetHandlerRoutesFromURIMatchers This is pretty brittle...you can have different hosts with similar paths and match on incorrect matchers
// TODO(Oliver): make this more resilient, add host to URIMatcher?
func (ms *MockServer) SetHandlerRoutesFromURIMatchers(uriMatchers []URIMatcher) {
	ms.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logx.Info("Received test request: ", r.RequestURI)
		for _, uriMatcher := range uriMatchers {
			uriM := uriMatcher
			var matches bool
			if len(uriM.Params) > 0 {
				matches = true
			}
			// Check each param to see if included in URL
			// Only match if all params satisfied
			for _, param := range uriM.Params {
				if !strings.Contains(r.RequestURI, param) {
					matches = false
					break
				}
			}
			if matches {
				logx.Info("Matched mock params %s for url: %s", uriM.Params, r.RequestURI)
				w.WriteHeader(uriM.ResponsePair.StatusCode)
				if _, err := w.Write(uriM.ResponsePair.Response); err != nil {
					logx.Warn("Error writing response: %s\n", err)
				}
				return
			}
		}
		logx.Warn("No matches for: %s", r.RequestURI)
	})
	if ms.Server != nil {
		ms.Server.Close()
	}
	ms.Server = httptest.NewServer(ms.Handler)
	SetTargetUrl(ms.Server.URL, ms.Targets...)
}

func (ms *MockServer) GetMockServer() *httptest.Server {
	return ms.Server
}

func (ms *MockServer) Close() {
	if ms.Server != nil {
		ms.Server.Close()
		ms.Server = nil
	}
}

func (ms *MockServer) SetRouteOverride(urls ...*string) {
	ms.Targets = []*string{}
	ms.Targets = append(ms.Targets, urls...)
}

func GetMockServer(expRespBody []byte, httpStatus int) (targetServer *httptest.Server) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(httpStatus)
		_, err := w.Write(expRespBody)
		if err != nil {
			return
		}
	})
	targetServer = httptest.NewServer(handler)
	return
}

func StripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
}

// Takes in two values of any type and returns
// nil if they are equal (deeply), an error if they are not.
//
// @deprecated, use "github.com/go-test/deep" for deep equality in tests
func IsEqual(tested, expected interface{}) error {

	//Check underlying types are the same
	if reflect.TypeOf(expected) != reflect.TypeOf(tested) {
		return errors.New("TypeOf expected (" + reflect.TypeOf(expected).String() +
			") != TypeOf tested (" + reflect.TypeOf(tested).String() + ")")
	}

	expectedValue := reflect.ValueOf(expected)
	testedValue := reflect.ValueOf(tested)

	switch expectedValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.String, reflect.Bool, reflect.Complex64,
		reflect.Complex128:
		if expected != tested {
			return errors.New("Expected " + reflect.TypeOf(expected).String() + ":\n" + fmt.Sprint(expected) + "\nis not equal to Tested:\n" +
				fmt.Sprint(tested))
		}
	case reflect.Float32, reflect.Float64:
		if math.Abs(expectedValue.Float()-testedValue.Float()) > floatingPointTolerance {
			return errors.New("Expected " + reflect.TypeOf(expected).String() + ":\n" + fmt.Sprint(expected) + "\nis not equal to Tested:\n" +
				fmt.Sprint(tested))
		}
	case reflect.Array, reflect.Slice:
		expectedLen := expectedValue.Len()
		testedLen := testedValue.Len()
		if expectedLen != testedLen {
			return errors.New("Length of Arrays not equal\nTested Array:\n" + fmt.Sprint(tested) +
				"\nWith Length: " + strconv.Itoa(testedLen) + " is not equal to:\n" + fmt.Sprint(expected) +
				"\nWith Length: " + strconv.Itoa(expectedLen))
		}
		for i := 0; i < expectedValue.Len(); i++ {
			err := IsEqual(testedValue.Index(i).Interface(), expectedValue.Index(i).Interface())
			if err != nil {
				return fmt.Errorf("At Index %d in Expected Array:\n %#v\n it is not equal to Tested:\n %#v\n %s",
					i, expected, tested, err)
			}
		}
	case reflect.Map:
		expectedLen := len(expectedValue.MapKeys())
		testedLen := len(testedValue.MapKeys())
		if expectedLen != testedLen {
			return errors.New("Length of Arrays not equal\nTested Map:\n" + fmt.Sprint(tested) +
				"\nWith Length: " + strconv.Itoa(testedLen) + " is not equal to:\n" + fmt.Sprint(expected) +
				"\nWith Length: " + strconv.Itoa(expectedLen))
		}
		for _, k := range expectedValue.MapKeys() {
			expectedAtKey := expectedValue.MapIndex(k).Interface()
			testedAtKey := testedValue.MapIndex(k).Interface()
			err := IsEqual(testedAtKey, expectedAtKey)
			if err != nil {
				return errors.New("At Key " + fmt.Sprint(k) + "\nExpected Map has value:\n" + fmt.Sprint(expectedAtKey) +
					"\nTested Map has value:\n" + fmt.Sprint(testedAtKey) + "\nExpected Map:\n" + fmt.Sprint(expected) +
					"\n Is not equal to Tested:\n" + fmt.Sprint(tested) + "\n" + fmt.Sprint(err))
			}
		}
	case reflect.Struct:
		numFieldsExpected := expectedValue.NumField()
		numFieldsTested := testedValue.NumField()
		if numFieldsExpected != numFieldsTested {
			return fmt.Errorf("Num of Fields in expected (%d) != num of fields in Num fields of tested (%d)", numFieldsExpected, numFieldsTested)
		}

		for i := 0; i < expectedValue.NumField(); i++ {
			expectedValueField := expectedValue.Field(i)
			testedValueField := testedValue.Field(i)

			//Nil Check
			expectedFieldIsZero := isZeroValue(expectedValueField)
			testedFieldIsZero := isZeroValue(testedValueField)
			if expectedFieldIsZero && testedFieldIsZero {
				continue
			} else if !expectedFieldIsZero && testedFieldIsZero {
				return errors.New("Expected was valid, tested was nil.")
			} else if expectedFieldIsZero && !testedFieldIsZero {
				return errors.New("Expected was nil, tested was valid.")
			}

			canInterfaceExpected := expectedValueField.CanInterface()
			canInterfaceTested := testedValueField.CanInterface()
			//If we cannot interface, we cannot test.
			if !canInterfaceTested && !canInterfaceExpected {
				continue
			} else if !canInterfaceExpected && canInterfaceTested {
				return errors.New("Expected cannot be interfaced, tested can.")
			} else if canInterfaceExpected && !canInterfaceTested {
				return errors.New("Expected can be interfaced, tested cannot.")
			}

			//Check Names Match
			expectedValueFieldName := expectedValue.Type().Field(i).Name
			testedValueFieldName := testedValue.Type().Field(i).Name
			if expectedValueFieldName != testedValueFieldName {
				return errors.New("Expected name: " + expectedValueFieldName + " is not equal to " +
					"tested name: " + testedValueFieldName + "at index " + strconv.Itoa(i))
			}

			//Check Type Match
			expectedValueFieldType := reflect.TypeOf(expectedValueField.Interface())
			testedValueFieldType := reflect.TypeOf(testedValueField.Interface())
			if expectedValueFieldType != testedValueFieldType {

				return errors.New("Expected type: " + expectedValueFieldType.String() +
					" is not equal to tested type: " + testedValueFieldType.String())
			}

			//Check Value Matches
			err := IsEqual(testedValueField.Interface(), expectedValueField.Interface())
			if err != nil {
				return errors.New("Expected " + expectedValueFieldName + " of type " + expectedValueFieldType.String() +
					" is not equal to " + "Tested " + testedValueFieldName + " of type " + testedValueFieldType.String() +
					"\n" + fmt.Sprint(err))
			}
		}
	case reflect.Ptr:
		//Check Nil
		expectedIsZero := isZeroValue(expected)
		testedIsZero := isZeroValue(tested)
		if expectedIsZero && testedIsZero {
			return nil
		} else if !expectedIsZero && testedIsZero {
			return errors.New("Expected pointer was valid, tested pointer was nil.")
		} else if expectedIsZero && !testedIsZero {
			return errors.New("Expected pointer was nil, tested pointer was valid.")
		}

		//Dereference pointer and call IsEqual on the value.
		expectedValueElem := expectedValue.Elem()
		testedValueElem := testedValue.Elem()
		if expectedValueElem.CanInterface() && testedValueElem.CanInterface() {
			err := IsEqual(testedValueElem.Interface(), expectedValueElem.Interface())
			if err != nil {
				return errors.New(fmt.Sprint(err))
			}
		}
	default: //We dont know the type
		return errors.New("Unknown Error. Type passed in is not handled")
	}
	return nil
}

// Returns whether the underlying
func isZeroValue(val interface{}) bool {
	return reflect.DeepEqual(reflect.Zero(reflect.TypeOf(val)).Interface(), val)
}
