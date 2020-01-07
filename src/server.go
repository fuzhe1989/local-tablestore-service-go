package main

import (
	"net/http"

	"github.com/golang/protobuf/proto"
)

// Server stores data and provides all of handlers
type Server struct {
}

// ServerMethodHandler ...
type ServerMethodHandler = func(Request) Response

// GetHandler returns method handlers
func (server *Server) GetHandler(method string) ServerMethodHandler {
	switch method {
	case "CreateTable":
		return server.handleCreateTable
	}
	panic("Unknown method:" + method)
}

func (server *Server) handleCreateTable(req Request) Response {
	reqPB := CreateTableRequest{}
	err := proto.Unmarshal(req.Data, &reqPB)
	if err != nil {
		// TODO: use correct error
		return respondError(newErrorUserNotExist())
	}
	respPB := CreateTableResponse{}
	return respondProto(&respPB)
}

func respondError(err OTSError) Response {
	errPB := Error{
		Code:    &err.Code,
		Message: &err.Message,
	}
	data, marshalErr := proto.Marshal(&errPB)
	if marshalErr != nil {
		panic("Error on marshal PB:" + marshalErr.Error())
	}
	return Response{
		Data:       data,
		HTTPStatus: err.HTTPStatus,
	}
}

func respondProto(message proto.Message) Response {
	data, marshalErr := proto.Marshal(message)
	if marshalErr != nil {
		panic("Error on marshal PB:" + marshalErr.Error())
	}
	return Response{
		Data:       data,
		HTTPStatus: http.StatusOK,
	}
}
