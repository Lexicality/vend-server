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
	"ResultJammed":          ResultJammed,
	"ResultHardwareFailure": ResultHardwareFailure,
	"ResultUnknownFailure":  ResultUnknownFailure,
	"ResultInvalidRequest":  ResultInvalidRequest,
}

func (res Result) String() string {
	switch res {
	case ResultSuccess:
		return "Success"
	case ResultEmpty:
		return "Empty"
	case ResultJammed:
		return "Jammed"
	case ResultHardwareFailure:
		return "Hardware Failure"
	case ResultUnknownFailure:
		return "Unknown Failure"
	case ResultAborted:
		return "Aborted before vend start"
	case ResultInvalidRequest:
		return "Invalid Location"
	default:
		return "Invalid Result"
	}
}
