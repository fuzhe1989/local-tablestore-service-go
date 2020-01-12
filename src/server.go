package main

import (
	"net/http"
	"sync"

	"github.com/fuzhe1989/local-tablestore-service-go/src/protocol"
	"github.com/golang/protobuf/proto"
)

// Server stores data and provides all of handlers
type Server struct {
	MetaManager MetaManager
}

// NewServer constructs a new Server
func NewServer(instances []string) *Server {
	server := new(Server)
	server.MetaManager.InstanceMap = make(map[string]TableMap)
	for _, instance := range instances {
		server.MetaManager.InstanceMap[instance] = make(TableMap)
	}
	return server
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
	reqPB := protocol.CreateTableRequest{}
	err := proto.Unmarshal(req.Data, &reqPB)
	if err != nil {
		// TODO: use correct error
		return respondError(newUserNotExistError())
	}
	tableMeta, err := toServerTableMeta()
	if err != nil {
		return respondError(newInternalServerError())
	}
	status := server.MetaManager.TryCreateTable(string(req.InstanceName), tableMeta)
	if !status.IsOK() {
		return respondError(status)
	}
	respPB := protocol.CreateTableResponse{}
	return respondProto(&respPB)
}

func respondError(err Status) Response {
	errPB := protocol.Error{
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

// TableMeta differs from the pb version.
type TableMeta struct {
}

func toServerTableMeta() (*TableMeta, error) {
	return &TableMeta{}, nil
}

// TableMap ...
type TableMap = map[string]*TableMeta

// MetaManager manages meta of all tables
type MetaManager struct {
	Lock sync.Mutex
	// instanceName -> tableMap
	InstanceMap map[string]TableMap
}

// TryCreateTable ...
func (manager *MetaManager) TryCreateTable(instanceName string, tableMeta *TableMeta) Status {
	manager.Lock.Lock()
	defer manager.Lock.Unlock()

	tables := manager.InstanceMap[instanceName]
	if tables == nil {
		return newInstanceNotFoundError()
	}

	// TODO
	return newStatusOK()
}
