[![Go build](https://github.com/Netcracker/qubership-core-lib-go-error-handling/actions/workflows/go-build.yml/badge.svg)](https://github.com/Netcracker/qubership-core-lib-go-error-handling/actions/workflows/go-build.yml)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?metric=coverage&project=Netcracker_qubership-core-lib-go-error-handling)](https://sonarcloud.io/summary/overall?id=Netcracker_qubership-core-lib-go-error-handling)
[![duplicated_lines_density](https://sonarcloud.io/api/project_badges/measure?metric=duplicated_lines_density&project=Netcracker_qubership-core-lib-go-error-handling)](https://sonarcloud.io/summary/overall?id=Netcracker_qubership-core-lib-go-error-handling)
[![vulnerabilities](https://sonarcloud.io/api/project_badges/measure?metric=vulnerabilities&project=Netcracker_qubership-core-lib-go-error-handling)](https://sonarcloud.io/summary/overall?id=Netcracker_qubership-core-lib-go-error-handling)
[![bugs](https://sonarcloud.io/api/project_badges/measure?metric=bugs&project=Netcracker_qubership-core-lib-go-error-handling)](https://sonarcloud.io/summary/overall?id=Netcracker_qubership-core-lib-go-error-handling)
[![code_smells](https://sonarcloud.io/api/project_badges/measure?metric=code_smells&project=Netcracker_qubership-core-lib-go-error-handling)](https://sonarcloud.io/summary/overall?id=Netcracker_qubership-core-lib-go-error-handling)

## Error Code Exceptions library

Go library that handles exceptions with error codes

- [Requirements to error produced by the code](#requirements-to-error-produced-by-the-code)
- [ErrCodeErr interface](#errcodeerr-interface)
- [ErrCodeErr implementation](#errcodeerr-implementation)
- [How to log ErrCodeErr](#how-to-log-errcodeerr)
- [Example of log message for ErrCodeErr with cause as multi-cause error](#example-of-log-message-for-errcodeerr-with-cause-as-multi-cause-error)
- [How to implement custom ErrCodeErr](#how-to-implement-custom-errcodeerr)
- [ErrCodeErr errors and fiber http server](#errcodeerr-errors-and-fiber-http-server)
- [How to handle remote Error Code errors](#how-to-handle-remote-error-code-errors)

### Requirements to error produced by the code

1. Each error must be associated with unique id
2. Each error type must have corresponding unique code in ([A-Z]+-)+\d{4} format
3. Each error type must have title - simple description without context/secure details 
4. Each error must be able to wrap another error
5. Each error must preserve its stacktrace. Error's stacktrace method must return stacktraces chain of all wrapped errors (as long as the wrapped errors provide their stacktraces)

#### ErrCodeErr interface
| method                                 | return type | mandatory | description                                                                                                                            |
|----------------------------------------|-------------|-----------|----------------------------------------------------------------------------------------------------------------------------------------|
| Error() (from builtin interface error) | string      | true      | error's message in format '{error-type-name} [{error-code}][{error-id}] {detail/errorCode.title (if detail is empty)} \n {stacktrace}' |
| GetId()                                | string      | true      | unique id (UUID) generated on creation of the exception                                                                                |
| GetErrorCode()                         | ErrorCode   | true      | struct which provides code and title for the exception                                                                                 |
| GetDetail()                            | String      | false     | context information related to particular exception                                                                                    |
| GetStackTrace()                        | Integer     | false     | http status code, if any                                                                                                               |

### ErrCodeErr implementation
This library provides default implementation for ErrCodeErr interface via ErrCodeError struct
This struct implements error's interface and all methods from ErrCodeErr interface.
Any struct which implements ErrCodeErr interface must meet the following requirements:

* Log message must contain the following tags
    * [error_code={ErrCodeErr.GetErrorCode().getCode()}]
    * [error_id={ErrCodeErr.GetId()}]
*  Custom ErrCodeErr must implement Error() string method. It must return string in the following format:
  ```
  [{error_code}][{error_id}] {message} \n {stacktrace}
  ```

#### How to log ErrCodeErr
```
  import (
    github.com/netcracker/qubership-core-lib-go/v3/logging
	errs "github.com/netcracker/qubership-core-lib-go-error-handling/v3/errors"
  )
  //...
  logger.Errorf(errs.ToLogFormat(errT))
  // or use template directly
  logger.Errorf(errs.ErrorLogTemplate, errT.GetErrorCode().Code, errT.GetId(), errT.GetStackTrace())
```

#### Example of log message for ErrCodeErr with cause as multi-cause error
```
MultiCauseError [NC-COMMON-2100][29fa2daa-f90d-11ec-ba5c-a860b613b330] multiple independent errors have happened
Caused by (1/2): ErrCodeError [TEST-1001][29fa262a-f90d-11ec-ba5c-a860b613b330] test 1001 detail
 goroutine 21 [running]:
 runtime/debug.Stack()
 	/usr/local/go/src/runtime/debug/stack.go:24 +0x7a
 github.com/netcracker/qubership-core-lib-go-error-handling/errors.NewError({{0x12c3f89, 0x9}, {0x12c3f9b, 0x9}}, {0x12c630c, 0x10}, {0x0, 0x0})
 	/workspace/qubership-core-lib-go-error-handling/errors/errors.go:107 +0x95
 github.com/netcracker/qubership-core-lib-go-error-handling/errors.TestMultiCauseError(0xc0001169c0)
 	/workspace/qubership-core-lib-go-error-handling/errors/errors_test.go:23 +0xb2
 testing.tRunner(0xc0001169c0, 0x12d6dd0)
 	/usr/local/go/src/testing/testing.go:1439 +0x1c3
 created by testing.(*T).Run
 	/usr/local/go/src/testing/testing.go:1486 +0x67c
Caused by (2/2): ErrCodeError [TEST-1002][29fa2d64-f90d-11ec-ba5c-a860b613b330] test 1002 detail
 goroutine 21 [running]:
 runtime/debug.Stack()
 	/usr/local/go/src/runtime/debug/stack.go:24 +0x7a
 github.com/netcracker/qubership-core-lib-go-error-handling/errors.NewError({{0x12c3f92, 0x9}, {0x12c3fa4, 0x9}}, {0x12c631c, 0x10}, {0x0, 0x0})
 	/workspace/qubership-core-lib-go-error-handling/errors/errors.go:107 +0x95
 github.com/netcracker/qubership-core-lib-go-error-handling/errors.TestMultiCauseError(0xc0001169c0)
 	/workspace/qubership-core-lib-go-error-handling/errors/errors_test.go:24 +0x11f
 testing.tRunner(0xc0001169c0, 0x12d6dd0)
 	/usr/local/go/src/testing/testing.go:1439 +0x1c3
 created by testing.(*T).Run
 	/usr/local/go/src/testing/testing.go:1486 +0x67c
```

#### How to implement custom ErrCodeErr
```
  import (
    github.com/netcracker/qubership-core-lib-go/v3/logging
	errs "github.com/netcracker/qubership-core-lib-go-error-handling/errors"
  )
  
  type CustomErr struct {
	*ErrCodeError // your struct must embed *ErrCodeError
	CustomField string
  }
  
  func NewCustomErr(customField string) *CustomErr {
      	err := errs.New(CustomErr{CustomField: customField}, ErrorCode{Code: "EXAMPLE-1001", Title: "example error"}, "example detail")
    	return err
  }
  
  // use your custom error
  func example() error {
    	err := NewCustomErr("custom test error")
    	return err
  }
```

#### ErrCodeErr errors and fiber http server
To wrap ErrCodeErr into TMF format (regarding this TMF format see details [here](https://github.com/Netcracker/qubership-core-lib-go-error-handling/blob/main/tmf/errors.go) 
before sending REST response in the fiber http server use [fiber-server-utils lib](https://github.com/netcracker/qubership-core-lib-go-fiber-server-utils)
below is the snipped how to set up fiber server with DefaultErrorHandler:
```
  import (
    github.com/netcracker/qubership-core-lib-go/v3/logging
	errs "github.com/netcracker/qubership-core-lib-go-error-handling/v3/errors"
	fiberserver "github.com/netcracker/qubership-core-lib-go-fiber-server-utils/v2"
	fibererrors "github.com/netcracker/qubership-core-lib-go-fiber-server-utils/v2/errors"
  )
  
  func main() {
	app, err := fiberserver.New(fiber.Config{
		Network:      fiber.NetworkTCP,
		ErrorHandler: fibererrors.DefaultErrorHandler(errs.ErrorCode{Code: "YOUR-MS-0001", Title: "unknown error"}),
	}).Process()
  )
```
By default fibererrors.DefaultErrorHandler sends response with 500 http status code
If the custom error requires custom http status code, or response must include custom meta field then your custom error must implement 'func(ctx *fiber.Ctx) error'
See example below:
```
  import (
    github.com/netcracker/qubership-core-lib-go/v3/logging
	errs "github.com/netcracker/qubership-core-lib-go-error-handling/v3/errors"
	tmf "github.com/netcracker/qubership-core-lib-go-error-handling/v3/tmf"
	fiberserver "github.com/netcracker/qubership-core-lib-go-fiber-server-utils/v2"
	fibererrors "github.com/netcracker/qubership-core-lib-go-fiber-server-utils/v2/errors"
  )
  
  type CustomErr struct {
	*ErrCodeError // your struct must embed *ErrCodeError
	CustomField string
  }

  func (e *CustomErr) getMeta() map[string]any {
	return map[string]any{"custom": e.CustomField})
  }
  
  func (e *CustomErr) Handle(ctx *fiber.Ctx) error {
     status := http.StatusBadRequest
     response := NewResponseBuilder(e).
		Meta(e.getMeta()).
		Status(status).
		Build()
	 return ctx.Status(status).JSON(response)
    }
  }
```
Provide 'func (e *YourCustomErr) ToError() *Error {...}' conversion function for your custom ErrCodeError in case it will be used in a MultiCauseError.
NewResponseBuilder() function will use this function while building tmf.Response
from provided MultiCauseError
```
  import (
    github.com/netcracker/qubership-core-lib-go/v3/logging
	errs "github.com/netcracker/qubership-core-lib-go-error-handling/v3/errors"
	tmf "github.com/netcracker/qubership-core-lib-go-error-handling/v3/tmf"
	fiberserver "github.com/netcracker/qubership-core-lib-go-fiber-server-utils/v2"
	fibererrors "github.com/netcracker/qubership-core-lib-go-fiber-server-utils/v2/errors"
  )
  
  type CustomErr struct {
	*ErrCodeError // your struct must embed *ErrCodeError
	CustomField string
  }
 
  func (e *CustomErr) getMeta() map[string]any {
	return map[string]any{"custom": e.CustomField})
  }
  
  func (e *CustomErr) ToError() *Error {
	return NewErrorBuilder(e).Meta(e.getMeta()).Build()
  }
  
  func (e *CustomErr) Handle(ctx *fiber.Ctx) error {
     status := http.StatusBadRequest
     response := NewResponseBuilder(e).
		Meta(e.getMeta()).
		Status(status).
		Build()
	 return ctx.Status(status).JSON(response)
    }
  }
```

#### How to handle remote Error Code errors
If your library or microservice performs REST calls you need to handle error responses from these calls.
See example below how to parse and wrap remote errors
```
  import (
    github.com/netcracker/qubership-core-lib-go/v3/logging
	errs "github.com/netcracker/qubership-core-lib-go-error-handling/v3/errors"
	"github.com/netcracker/qubership-core-lib-go-error-handling/v3/tmf"
	fiberserver "github.com/netcracker/qubership-core-lib-go-fiber-server-utilsv2/"
	fibererrors "github.com/netcracker/qubership-core-lib-go-fiber-server-utils/v2/errors"
  )
  
  var tmfErrConverter = tmf.DefaultConverter{}
   
  type FailedToSendRequestError struct {
	*ErrCodeError  
  }
  
  func NewFailedToSendRequestError(detail string, cause error) *FailedToSendRequestError {
      	return errs.New(FailedToSendRequestError{}, ErrorCode{Code: "EXAMPLE-1001", Title: "failed to send request error"}, detail, cause)
  }

  func SendRequest(ctx context.Context, method, uri string, body []byte) error {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(method)
	req.SetRequestURI(uri)
	req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	req.SetBody(body)
	resp := fasthttp.AcquireResponse()
	err := forwarder.client.Do(req, resp)
	if resp.StatusCode() >= 400 {
		r := tmf.Response{}
		bytes := resp.Body()
		err = json.Unmarshal(bytes, &r)
		if err != nil {
			errCodeError := NewFailedToSendRequestError(fmt.Sprintf("response not in TMF format, response was: %s", string(bytes)), nil)
			return resp, errCodeError
		}
		// convert remote response to RemoteErrCodeError
		remoteErr := tmfErrConverter.BuildErrorCodeError(r)
		// wrap remote error to your local error code error
		errCodeError := NewFailedToSendRequestError("Received error response in NC TMF format.", remoteErr)
		return resp, errCodeError
	}
```
