package main

import "net/http"

// Status wraps HTTPStatus and Error in pb.
type Status struct {
	HTTPStatus int
	Code       string
	Message    string
}

// IsOK ...
func (err *Status) IsOK() bool {
	return err.HTTPStatus == http.StatusOK
}

func newStatusOK() Status {
	return Status{
		HTTPStatus: http.StatusOK,
	}
}

func newAccessKeyIDNotExistError() Status {
	return Status{
		HTTPStatus: 403,
		Code:       "OTSAuthFailed",
		Message:    "The AccessKeyID does not exist.",
	}
}

func newAccessKeyIDDisabledError() Status {
	return Status{
		HTTPStatus: 403,
		Code:       "OTSAuthFailed",
		Message:    "The AccessKeyID is disabled.",
	}
}

func newUserNotExistError() Status {
	return Status{
		HTTPStatus: 403,
		Code:       "OTSAuthFailed",
		Message:    "The user does not exist.",
	}
}

func newInstanceNotFoundError() Status {
	return Status{
		HTTPStatus: 403,
		Code:       "OTSAuthFailed",
		Message:    "The instance is not found.",
	}
}

func newInternalServerError() Status {
	return Status{
		HTTPStatus: 500,
		Code:       "OTSInternalServerError",
		Message:    "Internal server error.",
	}
}
