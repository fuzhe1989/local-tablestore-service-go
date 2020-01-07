package main

// Request contains all information in a request to a Tablestore server.
type Request struct {
	Method         string
	InstanceName   []byte
	APIVersion     []byte
	Contentmd5     []byte
	ClientDateTime []byte
	Signature      []byte
	Data           []byte
}

// Response contains all information need send back to client
type Response struct {
	Data       []byte
	HTTPStatus int
}
