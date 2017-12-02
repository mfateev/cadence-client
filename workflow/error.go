// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package workflow

import (
	"go.uber.org/cadence/.gen/go/shared"
	"go.uber.org/cadence/internal"
)

/*
Below are the possible errors that activity or child workflow could return:
1) *cadence.CustomError: (this should be the most common one)
	If activity or child workflow implementation returns *CustomError by using NewCustomError() API, workflow code would receive *CustomError.
	The err would contain a Reason and Details. The reason is what activity specified to NewCustomError(), which workflow
	code could check to determine what kind of error it was and take actions based on the reason. The details is encoded
	[]byte which workflow code could extract strong typed data. Workflow code needs to know what the types of the encoded
	details are before extracting them.
2) *workflow.GenericError:
	If activity or child workflow implementation returns errors other than from NewCustomError() API,
    workflow code would receive *GenericError.
	Use err.Error() to get the string representation of the actual error.
3) *cadence.CanceledError:
	If activity or child workflow was canceled, workflow code will receive instance of *CanceledError.
    When activity or child workflow finishes cleanup it can indicate it by returning error created through
    NewCancelError() and could supply optional details which could be extracted by workflow code.
4) *workflow.TimeoutError:
	If activity or child workflow was timed out (several timeout types), workflow code will receive instance of
    *TimeoutError. The err contains details about what type of timeout it was.
5) *workflow.PanicError:
	If activity code panics while executing, cadence activity worker will report it as activity failure to cadence server.
	The cadence client library will present that failure as *PanicError to workflow code. The err contains a string
	representation of the panic message and the call stack when panic was happen.
    Note that there should be no *PanicError from child workflow. When panic happen in workflow implementation code,
    cadence client library catches that panic and causing the decision timeout. That decision task will be retried at
    a later time (with exponential backoff retry intervals). Eventually either decision code is fixed to not panic or
    a workflow execution times out. In the timeout case the parent workflow receives TimeoutError.


Workflow code could handle errors based on different types of error. Below is sample code of how error handling looks like.

_, err := cadence.ExecuteActivity(ctx, MyActivity, ...).Get(nil)
if err != nil {
	switch err := err.(type) {
	case *CustomError:
		// handle activity errors (created via NewCustomError() API)
		switch err.Reason() {
		case cadence.CustomErrReasonA: // assume CustomErrReasonA is constant defined by activity implementation
			var detailMsg string // assuming activity return error by NewCustomError(CustomErrReasonA, "string details")
			err.Details(&detailMsg) // extract strong typed details (corresponding to CustomErrReasonA)
			// handle CustomErrReasonA
		case CustomErrReasonB:
			// handle CustomErrReasonB
		default:
			// newer version of activity could return new errors that workflow was not aware of.
		}
	case *cadence.GenericError:
		// handle generic error (errors created other than using NewCustomError() API)
	case *cadence.CanceledError:
		// handle cancellation
	case *cadence.TimeoutError:
		// handle timeout, could check timeout type by err.TimeoutType()
	case *cadence.PanicError:
		// handle panic
	}
}
*/

type (
	// GenericError is returned from activity or child workflow when an implementations return error
	// other than from cadence.NewCustomError() API.
	GenericError = internal.GenericError

	// TimeoutError returned when activity or child workflow timed out.
	TimeoutError = internal.TimeoutError

	// PanicError contains information about panicked workflow/activity.
	PanicError = internal.PanicError
)

// NewTimeoutError creates TimeoutError instance.
// Use NewHeartbeatTimeoutError to create heartbeat TimeoutError
// WARNING: This function is public only to support unit testing of workflows.
// It shouldn't be used by application level code.
func NewTimeoutError(timeoutType shared.TimeoutType) *TimeoutError {
	return internal.NewTimeoutError(timeoutType)
}

// NewHeartbeatTimeoutError creates TimeoutError instance
// WARNING: This function is public only to support unit testing of workflows.
// It shouldn't be used by application level code.
func NewHeartbeatTimeoutError(details ...interface{}) *TimeoutError {
	return internal.NewHeartbeatTimeoutError(details)
}
