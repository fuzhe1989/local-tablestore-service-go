package main

// OTSError wraps HTTPStatus and Error in pb.
type OTSError struct {
	HTTPStatus int
	Code       string
	Message    string
}

func newErrorAccessKeyIDNotExist() OTSError {
	return OTSError{
		HTTPStatus: 403,
		Code:       "OTSAuthFailed",
		Message:    "The AccessKeyID does not exist.",
	}
}

func newErrorAccessKeyIDDisabled() OTSError {
	return OTSError{
		HTTPStatus: 403,
		Code:       "OTSAuthFailed",
		Message:    "The AccessKeyID is disabled.",
	}
}

func newErrorUserNotExist() OTSError {
	return OTSError{
		HTTPStatus: 403,
		Code:       "OTSAuthFailed",
		Message:    "The user does not exist.",
	}
}
