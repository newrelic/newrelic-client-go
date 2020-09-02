// Code generated by tutone: DO NOT EDIT
package edge

// EdgeCreateTraceObserverResponseErrorType - Known error codes and messages for `CreateTraceObserverResponseError`.
type EdgeCreateTraceObserverResponseErrorType string

var EdgeCreateTraceObserverResponseErrorTypeTypes = struct {
	// A trace observer already exists for this account family and provider region.
	ALREADY_EXISTS EdgeCreateTraceObserverResponseErrorType
	// Trace observers aren’t available in provider region.
	NO_AVAILABILITY_IN_REGION EdgeCreateTraceObserverResponseErrorType
	// You don’t have permission to make this trace observer call.
	UNAUTHORIZED_USER EdgeCreateTraceObserverResponseErrorType
	// We couldn’t process this request.
	UNEXPECTED_ERROR EdgeCreateTraceObserverResponseErrorType
}{
	// A trace observer already exists for this account family and provider region.
	ALREADY_EXISTS: "ALREADY_EXISTS",
	// Trace observers aren’t available in provider region.
	NO_AVAILABILITY_IN_REGION: "NO_AVAILABILITY_IN_REGION",
	// You don’t have permission to make this trace observer call.
	UNAUTHORIZED_USER: "UNAUTHORIZED_USER",
	// We couldn’t process this request.
	UNEXPECTED_ERROR: "UNEXPECTED_ERROR",
}

// EdgeDeleteTraceObserverResponseErrorType - Known error codes and messages for `DeleteTraceObserverResponseError`.
type EdgeDeleteTraceObserverResponseErrorType string

var EdgeDeleteTraceObserverResponseErrorTypeTypes = struct {
	// The trace observer has already been deleted.
	ALREADY_DELETED EdgeDeleteTraceObserverResponseErrorType
	// No trace observer was found with the id given.
	NOT_FOUND EdgeDeleteTraceObserverResponseErrorType
	// You don’t have permission to make this trace observer call.
	UNAUTHORIZED_USER EdgeDeleteTraceObserverResponseErrorType
	// We couldn’t process this request.
	UNEXPECTED_ERROR EdgeDeleteTraceObserverResponseErrorType
}{
	// The trace observer has already been deleted.
	ALREADY_DELETED: "ALREADY_DELETED",
	// No trace observer was found with the id given.
	NOT_FOUND: "NOT_FOUND",
	// You don’t have permission to make this trace observer call.
	UNAUTHORIZED_USER: "UNAUTHORIZED_USER",
	// We couldn’t process this request.
	UNEXPECTED_ERROR: "UNEXPECTED_ERROR",
}

// EdgeEndpointStatus - Status of the endpoint.
type EdgeEndpointStatus string

var EdgeEndpointStatusTypes = struct {
	// The endpoint has been created and is available for use.
	CREATED EdgeEndpointStatus
	// The endpoint has been deleted and is no longer available for use.
	DELETED EdgeEndpointStatus
}{
	// The endpoint has been created and is available for use.
	CREATED: "CREATED",
	// The endpoint has been deleted and is no longer available for use.
	DELETED: "DELETED",
}

// EdgeEndpointType - Type of connection established with the trace observer. Currently, only `PUBLIC` is supported.
type EdgeEndpointType string

var EdgeEndpointTypeTypes = struct {
	// PUBLIC: the endpoint is reachable on the internet.
	PUBLIC EdgeEndpointType
}{
	// PUBLIC: the endpoint is reachable on the internet.
	PUBLIC: "PUBLIC",
}

// EdgeProviderRegion - Provider and region where the trace observer is located. Currently, only AWS regions are supported.
type EdgeProviderRegion string

var EdgeProviderRegionTypes = struct {
	// Provider: `AWS`, Region: `eu-west-1`
	AWS_EU_WEST_1 EdgeProviderRegion
	// Provider: `AWS`, Region: `us-east-1`
	AWS_US_EAST_1 EdgeProviderRegion
	// Provider: `AWS`, Region: `us-east-2`
	AWS_US_EAST_2 EdgeProviderRegion
	// Provider: `AWS`, Region: `us-west-2`
	AWS_US_WEST_2 EdgeProviderRegion
}{
	// Provider: `AWS`, Region: `eu-west-1`
	AWS_EU_WEST_1: "AWS_EU_WEST_1",
	// Provider: `AWS`, Region: `us-east-1`
	AWS_US_EAST_1: "AWS_US_EAST_1",
	// Provider: `AWS`, Region: `us-east-2`
	AWS_US_EAST_2: "AWS_US_EAST_2",
	// Provider: `AWS`, Region: `us-west-2`
	AWS_US_WEST_2: "AWS_US_WEST_2",
}

// EdgeTraceObserverResponseErrorType - Known error codes and messages for `TraceObserverResponseError`.
type EdgeTraceObserverResponseErrorType string

var EdgeTraceObserverResponseErrorTypeTypes = struct {
	// You don’t have permission to make this trace observer call.
	UNAUTHORIZED_USER EdgeTraceObserverResponseErrorType
	// We couldn’t process this request.
	UNEXPECTED_ERROR EdgeTraceObserverResponseErrorType
}{
	// You don’t have permission to make this trace observer call.
	UNAUTHORIZED_USER: "UNAUTHORIZED_USER",
	// We couldn’t process this request.
	UNEXPECTED_ERROR: "UNEXPECTED_ERROR",
}

// EdgeTraceObserverStatus - Status of the trace observer.
type EdgeTraceObserverStatus string

var EdgeTraceObserverStatusTypes = struct {
	// The trace observer has been created and is available for use.
	CREATED EdgeTraceObserverStatus
	// The trace observer has been deleted and is no longer available for use.
	DELETED EdgeTraceObserverStatus
}{
	// The trace observer has been created and is available for use.
	CREATED: "CREATED",
	// The trace observer has been deleted and is no longer available for use.
	DELETED: "DELETED",
}

// EdgeAccountStitchedFields -
type EdgeAccountStitchedFields struct {
	// Provides access to Tracing data.
	Tracing EdgeTracing `json:"tracing"`
}

// EdgeAgentEndpointDetail - All the details necessary to configure an agent to connect to an endoint.
type EdgeAgentEndpointDetail struct {
	// Full host name that is used to connect to the endpoint. This is the part that will be placed into an agent config named `infinite_tracing.trace_observer.host`.
	Host string `json:"host"`
	// Port that is used to connect to the endpoint. This is the part that will be placed into an agent config named `infinite_tracing.trace_observer.port`.
	Port int `json:"port"`
}

func (x *EdgeAgentEndpointDetail) ImplementsEdgeEndpointDetail() {}

// EdgeCreateTraceObserverInput - Data required to create a trace observer.
type EdgeCreateTraceObserverInput struct {
	// Name of the trace observer.
	Name string `json:"name"`
	// Provider and region where the trace observer must run. Currently, only AWS regions are supported.
	ProviderRegion EdgeProviderRegion `json:"providerRegion"`
}

// EdgeCreateTraceObserverResponse - Successfully created trace observers, or one or more error responses if there were issues.
type EdgeCreateTraceObserverResponse struct {
	// Errors that may occur when creating a `TraceObserver`. Defaults to `null` in case of success.
	Errors []EdgeCreateTraceObserverResponseError `json:"errors"`
	// The trace observer defined in `CreateTraceObserverInput`. Defaults to `null` in case of failure.
	TraceObserver EdgeTraceObserver `json:"traceObserver"`
}

// EdgeCreateTraceObserverResponseError - Description of errors that may occur while attempting to create a trace observer.
type EdgeCreateTraceObserverResponseError struct {
	// Error message, with further detail to help resolve the issue.
	Message string `json:"message"`
	// Error that may occur while attempting to create a trace observer.
	Type EdgeCreateTraceObserverResponseErrorType `json:"type"`
}

// EdgeCreateTraceObserverResponses - Array of responses, one for each trace observer creation request.
type EdgeCreateTraceObserverResponses struct {
	// Array of trace observer creation responses, one for each `CreateTraceObserverInput`.
	Responses []EdgeCreateTraceObserverResponse `json:"responses"`
}

// EdgeDeleteTraceObserverInput - Data required to delete a trace observer.
type EdgeDeleteTraceObserverInput struct {
	// Globally unique identifier of the trace observer being deleted.
	ID int `json:"id"`
}

// EdgeDeleteTraceObserverResponse - Successfully deleted trace observers, or one or more error responses if there were issues.
type EdgeDeleteTraceObserverResponse struct {
	// Errors that may occur when deleting a `TraceObserver`. Defaults to `null` in case of success.
	Errors []EdgeDeleteTraceObserverResponseError `json:"errors"`
	// The trace observer that was deleted. Defaults to `null` in case of failure.
	TraceObserver EdgeTraceObserver `json:"traceObserver"`
}

// EdgeDeleteTraceObserverResponseError - Description of errors that may occur while attempting to delete a trace observer.
type EdgeDeleteTraceObserverResponseError struct {
	// Error message, with further detail to help resolve the issue.
	Message string `json:"message"`
	// Error that may occur while attempting to delete a trace observer.
	Type EdgeDeleteTraceObserverResponseErrorType `json:"type"`
}

// EdgeDeleteTraceObserverResponses - Array of responses, one for each trace observer deletion request.
type EdgeDeleteTraceObserverResponses struct {
	// Array of trace observer deletion responses, one for each `DeleteTraceObserverInput`.
	Responses []EdgeDeleteTraceObserverResponse `json:"responses"`
}

// EdgeEndpoint - An `Endpoint` describes access to an endpoint pointing to a trace observer. Currently, only one endpoint per trace observer is supported.
type EdgeEndpoint struct {
	// Connection information related to the agent configuration.
	Agent EdgeAgentEndpointDetail `json:"agent"`
	// Type of the endpoint.
	EndpointType EdgeEndpointType `json:"endpointType"`
	// Connection information related to the Infinite Tracing Trace API (HTTP 1.1) configuration.
	Https EdgeHttpsEndpointDetail `json:"https"`
	// Status of the endpoint.
	Status EdgeEndpointStatus `json:"status"`
}

// EdgeHttpsEndpointDetail - All the details necessary to configure an integration to connect to the Infinite Tracing Trace API (HTTP 1.1) endpoint.
type EdgeHttpsEndpointDetail struct {
	// Full host name that is used to connect to the endpoint.
	Host string `json:"host"`
	// Port that is used to connect to the endpoint.
	Port int `json:"port"`
	// Full URL used to send data to the endpoint. For instance, if you were using the
	//  [Java Telemetry SDK](https://docs.newrelic.com/docs/data-ingest-apis/get-data-new-relic/new-relic-sdks/telemetry-sdks-send-custom-telemetry-data-new-relic)
	//  this is the data you would use to create a `URI` to pass to the [`uriOverride`](https://github.com/newrelic/newrelic-telemetry-sdk-java/blob/85e526cf6fbba0640f20d2d7a3ab0dab89f958b3/telemetry_core/src/main/java/com/newrelic/telemetry/AbstractSenderBuilder.java#L37-L48)
	//  method.
	Url string `json:"url"`
}

func (x *EdgeHttpsEndpointDetail) ImplementsEdgeEndpointDetail() {}

// EdgeTraceObserver - `TraceObserver` handles a group of tracing services for an account family.
type EdgeTraceObserver struct {
	// List of endpoints associated with this trace observer. Currently, only one endpoint per trace observer is supported.
	Endpoints []EdgeEndpoint `json:"endpoints"`
	// Globally unique identifier of this trace observer.
	ID int `json:"id"`
	// Human-readable name of this trace observer.
	Name string `json:"name"`
	// Provider-specific region of this endpoint (for example, `AWS_US_EAST_1`). Currently, only AWS regions are supported.
	ProviderRegion EdgeProviderRegion `json:"providerRegion"`
	// Status of the trace observer.
	Status EdgeTraceObserverStatus `json:"status"`
}

// EdgeTraceObserverResponse - Array of trace observers, or a list of errors for why they couldn't be retrieved.
type EdgeTraceObserverResponse struct {
	// All trace observer's response errors, if any.
	Errors []EdgeTraceObserverResponseError `json:"errors"`
	// All trace observers found, if any.
	TraceObservers []EdgeTraceObserver `json:"traceObservers"`
}

// EdgeTraceObserverResponseError - Description of errors that may occur while attempting to retrieve a trace observer.
type EdgeTraceObserverResponseError struct {
	// Error message, with further detail to help resolve the issue.
	Message string `json:"message"`
	// Error that may occur while attempting to retrieve a trace observer.
	Type EdgeTraceObserverResponseErrorType `json:"type"`
}

// EdgeTracing - This field provides access to Tracing data.
type EdgeTracing struct {
	// Lists the existing trace observers for this account family.
	TraceObservers EdgeTraceObserverResponse `json:"traceObservers"`
}
