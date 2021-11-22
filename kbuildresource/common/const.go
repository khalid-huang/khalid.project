package common

const (
	BuildJobCreateRequestType string = "buildjob_create"
	BuildJobDeleteRequestType string = "buildjob_delete"
	BuildJobPrefix string = "buildjob"

	// request status
	RequestStatusPending string = "pending"
	RequestStatusExecuting string = "executing"
	RequestStatusFailed string = "failed"
	RequestStatusSuccess string = "success"
)
