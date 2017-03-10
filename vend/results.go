package vend

// Result represents the result of an attempted vend operation
type Result uint

const (
	// NoResult is if the operation did not result in a vend
	NoResult Result = iota
	// ResultSuccess is when everything is fine
	ResultSuccess
	// ResultEmpty is when you vend something that's not there
	ResultEmpty
	// ResultBusy is when the vending machine is already vending
	ResultBusy
	// ResultOffline is when the vending machine isn't talking to us right now
	ResultOffline
	// ResultJammed is when the vending machine physically can't vend
	ResultJammed
	// ResultHardwareFailure is when the vending machine is on fire
	ResultHardwareFailure
	// ResultUnknownFailure is when something else went wrong
	ResultUnknownFailure
	// ResultInvalidRequest is when you do something stupid
	ResultInvalidRequest
	// ResultAborted is when the application shuts down mid-vend
	ResultAborted
)

// AllResults is a handy array of all results such as you might want in a template
var AllResults = map[string]Result{
	"ResultSuccess":         ResultSuccess,
	"ResultEmpty":           ResultEmpty,
	"ResultBusy":            ResultBusy,
	"ResultOffline":         ResultOffline,
	"ResultJammed":          ResultJammed,
	"ResultHardwareFailure": ResultHardwareFailure,
	"ResultUnknownFailure":  ResultUnknownFailure,
	"ResultInvalidRequest":  ResultInvalidRequest,
}
