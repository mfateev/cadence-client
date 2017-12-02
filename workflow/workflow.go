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
	"github.com/uber-go/tally"
	"go.uber.org/cadence/encoded"
	"go.uber.org/cadence/internal"
	"go.uber.org/zap"
)

type (

	// ChildWorkflowFuture represents the result of a child workflow execution
	ChildWorkflowFuture = internal.ChildWorkflowFuture

	// Type identifies a workflow type.
	Type = internal.WorkflowType

	// Execution Details.
	Execution = internal.WorkflowExecution

	// Version represents a change version. See GetVersion call.
	Version = internal.Version

	// ChildWorkflowOptions stores all child workflow specific parameters that will be stored inside of a Context.
	ChildWorkflowOptions = internal.ChildWorkflowOptions

	// ChildWorkflowPolicy defines child workflow behavior when parent workflow is terminated.
	ChildWorkflowPolicy = internal.ChildWorkflowPolicy

	// RegisterWorkflowOptions consists of options for registering a workflow
	RegisterWorkflowOptions = internal.RegisterWorkflowOptions

	// Info information about currently executing workflow
	Info = internal.WorkflowInfo

	// ContinueAsNewError contains information about how to continue the workflow as new.
	ContinueAsNewError = internal.ContinueAsNewError
)

const (
	// ChildWorkflowPolicyTerminate is policy that will terminate all child workflows when parent workflow is terminated.
	ChildWorkflowPolicyTerminate ChildWorkflowPolicy = internal.ChildWorkflowPolicyTerminate
	// ChildWorkflowPolicyRequestCancel is policy that will send cancel request to all open child workflows when parent
	// workflow is terminated.
	ChildWorkflowPolicyRequestCancel ChildWorkflowPolicy = internal.ChildWorkflowPolicyRequestCancel
	// ChildWorkflowPolicyAbandon is policy that will have no impact to child workflow execution when parent workflow is
	// terminated.
	ChildWorkflowPolicyAbandon ChildWorkflowPolicy = internal.ChildWorkflowPolicyAbandon
)

// Register - registers a workflow function with the framework.
// A workflow takes a cadence context and input and returns a (result, error) or just error.
// Examples:
//	func sampleWorkflow(ctx cadence.Context, input []byte) (result []byte, err error)
//	func sampleWorkflow(ctx cadence.Context, arg1 int, arg2 string) (result []byte, err error)
//	func sampleWorkflow(ctx cadence.Context) (result []byte, err error)
//	func sampleWorkflow(ctx cadence.Context, arg1 int) (result string, err error)
// Serialization of all primitive types, structures is supported ... except channels, functions, variadic, unsafe pointer.
// This method calls panic if workflowFunc doesn't comply with the expected format.
func Register(workflowFunc interface{}) {
	internal.RegisterWorkflow(workflowFunc)
}

// RegisterWithOptions registers the workflow function with options
// The user can use options to provide an external name for the workflow or leave it empty if no
// external name is required. This can be used as
//  client.RegisterWorkflow(sampleWorkflow, RegisterWorkflowOptions{})
//  client.RegisterWorkflow(sampleWorkflow, RegisterWorkflowOptions{Name: "foo"})
// A workflow takes a cadence context and input and returns a (result, error) or just error.
// Examples:
//	func sampleWorkflow(ctx cadence.Context, input []byte) (result []byte, err error)
//	func sampleWorkflow(ctx cadence.Context, arg1 int, arg2 string) (result []byte, err error)
//	func sampleWorkflow(ctx cadence.Context) (result []byte, err error)
//	func sampleWorkflow(ctx cadence.Context, arg1 int) (result string, err error)
// Serialization of all primitive types, structures is supported ... except channels, functions, variadic, unsafe pointer.
// This method calls panic if workflowFunc doesn't comply with the expected format.
func RegisterWithOptions(workflowFunc interface{}, opts RegisterWorkflowOptions) {
	internal.RegisterWorkflowWithOptions(workflowFunc, opts)
}

// ExecuteActivity requests activity execution in the context of a workflow.
// Context can be used to pass the settings for this activity.
// For example: task list that this need to be routed, timeouts that need to be configured.
// Use ActivityOptions to pass down the options.
//  ao := ActivityOptions{
// 	    TaskList: "exampleTaskList",
// 	    ScheduleToStartTimeout: 10 * time.Second,
// 	    StartToCloseTimeout: 5 * time.Second,
// 	    ScheduleToCloseTimeout: 10 * time.Second,
// 	    HeartbeatTimeout: 0,
// 	}
//	ctx := WithActivityOptions(ctx, ao)
// Or to override a single option
//  ctx := WithTaskList(ctx, "exampleTaskList")
// Input activity is either an activity name (string) or a function representing an activity that is getting scheduled.
// Input args are the arguments that need to be passed to the scheduled activity.
//
// If the activity failed to complete then the future get error would indicate the failure, and it can be one of
// CustomError, TimeoutError, CanceledError, PanicError, GenericError.
// You can cancel the pending activity using context(cadence.WithCancel(ctx)) and that will fail the activity with
// error CanceledError.
//
// ExecuteActivity returns Future with activity result or failure.
func ExecuteActivity(ctx Context, activity interface{}, args ...interface{}) Future {
	return internal.ExecuteActivity(ctx, activity, args)
}

// ExecuteChildWorkflow requests child workflow execution in the context of a workflow.
// Context can be used to pass the settings for the child workflow.
// For example: task list that this child workflow should be routed, timeouts that need to be configured.
// Use ChildWorkflowOptions to pass down the options.
//  cwo := ChildWorkflowOptions{
// 	    ExecutionStartToCloseTimeout: 10 * time.Minute,
// 	    TaskStartToCloseTimeout: time.Minute,
// 	}
//  ctx := WithChildWorkflowOptions(ctx, cwo)
// Input childWorkflow is either a workflow name or a workflow function that is getting scheduled.
// Input args are the arguments that need to be passed to the child workflow function represented by childWorkflow.
// If the child workflow failed to complete then the future get error would indicate the failure and it can be one of
// CustomError, TimeoutError, CanceledError, GenericError.
// You can cancel the pending child workflow using context(cadence.WithCancel(ctx)) and that will fail the workflow with
// error CanceledError.
// ExecuteChildWorkflow returns ChildWorkflowFuture.
func ExecuteChildWorkflow(ctx Context, childWorkflow interface{}, args ...interface{}) ChildWorkflowFuture {
	return internal.ExecuteChildWorkflow(ctx, childWorkflow, args)
}

// GetWorkflowInfo extracts info of a current workflow from a context.
func GetWorkflowInfo(ctx Context) *Info {
	return internal.GetWorkflowInfo(ctx)
}

// GetLogger returns a logger to be used in workflow's context
func GetLogger(ctx Context) *zap.Logger {
	return internal.GetLogger(ctx)
}

// GetMetricsScope returns a metrics scope to be used in workflow's context
func GetMetricsScope(ctx Context) tally.Scope {
	return GetMetricsScope(ctx)
}

// RequestCancelWorkflow can be used to request cancellation of an external workflow.
// Input workflowID is the workflow ID of target workflow.
// Input runID indicates the instance of a workflow. Input runID is optional (default is ""). When runID is not specified,
// then the currently running instance of that workflowID will be used.
// By default, the current workflow's domain will be used as target domain. However, you can specify a different domain
// of the target workflow using the context like:
//	ctx := WithWorkflowDomain(ctx, "domain-name")
func RequestCancelWorkflow(ctx Context, workflowID, runID string) error {
	return internal.RequestCancelWorkflow(ctx, workflowID, runID)
}

// GetSignalChannel returns channel corresponding to the signal name.
func GetSignalChannel(ctx Context, signalName string) Channel {
	return internal.GetSignalChannel(ctx, signalName)
}

// SideEffect executes provided function once, records its result into the workflow history. The recorded result on
// history will be returned without executing the provided function during replay. This guarantees the deterministic
// requirement for workflow as the exact same result will be returned in replay.
// Common use case is to run some short non-deterministic code in workflow, like getting random number or new UUID.
// The only way to fail SideEffect is to panic which causes decision task failure. The decision task after timeout is
// rescheduled and re-executed giving SideEffect another chance to succeed.
//
// Caution: do not use SideEffect to modify closures, always retrieve result from SideEffect's encoded return value.
// For example this code is BROKEN:
//  // Bad example:
//  var random int
//  cadence.SideEffect(func(ctx cadence.Context) interface{} {
//         random = rand.Intn(100)
//         return nil
//  })
//  // random will always be 0 in replay, thus this code is non-deterministic
//  if random < 50 {
//         ....
//  } else {
//         ....
//  }
// On replay the provided function is not executed, the random will always be 0, and the workflow could takes a
// different path breaking the determinism.
//
// Here is the correct way to use SideEffect:
//  // Good example:
//  encodedRandom := SideEffect(func(ctx cadence.Context) interface{} {
//        return rand.Intn(100)
//  })
//  var random int
//  encodedRandom.Get(&random)
//  if random < 50 {
//         ....
//  } else {
//         ....
//  }
func SideEffect(ctx Context, f func(ctx Context) interface{}) encoded.Value {
	return internal.SideEffect(ctx, f)
}

// DefaultVersion is a version returned by GetVersion for code that wasn't versioned before
const DefaultVersion Version = internal.DefaultVersion

// GetVersion is used to safely perform backwards incompatible changes to workflow definitions.
// It is not allowed to update workflow code while there are workflows running as it is going to break
// determinism. The solution is to have both old code that is used to replay existing workflows
// as well as the new one that is used when it is executed for the first time.
// GetVersion returns maxSupported version when is executed for the first time. This version is recorded into the
// workflow history as a marker event. Even if maxSupported version is changed the version that was recorded is
// returned on replay. DefaultVersion constant contains version of code that wasn't versioned before.
// For example initially workflow has the following code:
//  err = cadence.ExecuteActivity(ctx, foo).Get(ctx, nil)
// it should be updated to
//  err = cadence.ExecuteActivity(ctx, bar).Get(ctx, nil)
// The backwards compatible way to execute the update is
//  v :=  GetVersion(ctx, "fooChange", DefaultVersion, 1)
//  if v  == DefaultVersion {
//      err = cadence.ExecuteActivity(ctx, foo).Get(ctx, nil)
//  } else {
//      err = cadence.ExecuteActivity(ctx, bar).Get(ctx, nil)
//  }
//
// Then bar has to be changed to baz:
//  v :=  GetVersion(ctx, "fooChange", DefaultVersion, 2)
//  if v  == DefaultVersion {
//      err = cadence.ExecuteActivity(ctx, foo).Get(ctx, nil)
//  } else if v == 1 {
//      err = cadence.ExecuteActivity(ctx, bar).Get(ctx, nil)
//  } else {
//      err = cadence.ExecuteActivity(ctx, baz).Get(ctx, nil)
//  }
//
// Later when there are no workflow executions running DefaultVersion the correspondent branch can be removed:
//  v :=  GetVersion(ctx, "fooChange", 1, 2)
//  if v == 1 {
//      err = cadence.ExecuteActivity(ctx, bar).Get(ctx, nil)
//  } else {
//      err = cadence.ExecuteActivity(ctx, baz).Get(ctx, nil)
//  }
//
// Currently there is no supported way to completely remove GetVersion call after it was introduced.
// Keep it even if single branch is left:
//  GetVersion(ctx, "fooChange", 2, 2)
//  err = cadence.ExecuteActivity(ctx, baz).Get(ctx, nil)
//
// It is necessary as GetVersion performs validation of a version against a workflow history and fails decisions if
// a workflow code is not compatible with it.
func GetVersion(ctx Context, changeID string, minSupported, maxSupported Version) Version {
	return internal.GetVersion(ctx, changeID, minSupported, maxSupported)
}

// SetQueryHandler sets the query handler to handle workflow query. The queryType specify which query type this handler
// should handle. The handler must be a function that returns 2 values. The first return value must be a serializable
// result. The second return value must be an error. The handler function could receive any number of input parameters.
// All the input parameter must be serializable. You should call cadence.SetQueryHandler() at the beginning of the workflow
// code. When client calls Client.QueryWorkflow() to cadence server, a task will be generated on server that will be dispatched
// to a workflow worker, which will replay the history events and then execute a query handler based on the query type.
// The query handler will be invoked out of the context of the workflow, meaning that the handler code must not use cadence
// context to do things like cadence.NewChannel(), cadence.Go() or to call any workflow blocking functions like
// Channel.Get() or Future.Get(). Trying to do so in query handler code will fail the query and client will receive
// QueryFailedError.
// Example of workflow code that support query type "current_state":
//  func MyWorkflow(ctx cadence.Context, input string) error {
//    currentState := "started" // this could be any serializable struct
//    err := cadence.SetQueryHandler(ctx, "current_state", func() (string, error) {
//      return currentState, nil
//    })
//    if err != nil {
//      currentState = "failed to register query handler"
//      return err
//    }
//    // your normal workflow code begins here, and you update the currentState as the code makes progress.
//    currentState = "waiting timer"
//    err = NewTimer(ctx, time.Hour).Get(ctx, nil)
//    if err != nil {
//      currentState = "timer failed"
//      return err
//    }
//
//    currentState = "waiting activity"
//    ctx = WithActivityOptions(ctx, myActivityOptions)
//    err = ExecuteActivity(ctx, MyActivity, "my_input").Get(ctx, nil)
//    if err != nil {
//      currentState = "activity failed"
//      return err
//    }
//    currentState = "done"
//    return nil
//  }
func SetQueryHandler(ctx Context, queryType string, handler interface{}) error {
	return internal.SetQueryHandler(ctx, queryType, handler)
}

// NewContinueAsNewError creates ContinueAsNewError instance
// If the workflow main function returns this error then the current execution is ended and
// the new execution with same workflow ID is started automatically with options
// provided to this function.
//  ctx - use context to override any options for the new workflow like execution time out, decision task time out, task list.
//	  if not mentioned it would use the defaults that the current workflow is using.
//        ctx := WithExecutionStartToCloseTimeout(ctx, 30 * time.Minute)
//        ctx := WithWorkflowTaskStartToCloseTimeout(ctx, time.Minute)
//	  ctx := WithWorkflowTaskList(ctx, "example-group")
//  wfn - workflow function. for new execution it can be different from the currently running.
//  args - arguments for the new workflow.
//
func NewContinueAsNewError(ctx Context, wfn interface{}, args ...interface{}) *ContinueAsNewError {
	return internal.NewContinueAsNewError(ctx, wfn, args)
}
