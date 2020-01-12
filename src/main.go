package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/buaazp/fasthttprouter"
	"github.com/fuzhe1989/local-tablestore-service-go/src/protocol"
	"github.com/golang/protobuf/proto"
	"github.com/valyala/fasthttp"
)

var (
	addr                = flag.String("addr", ":8080", "TCP address to listen to")
	predefinedInstances = flag.String("predefined-instances", "test", "Split by comma(',') if predefine multiple instances")
)

const (
	xOtsDate                = "x-ots-date"
	xOtsApiversion          = "x-ots-apiversion"
	xOtsAccesskeyid         = "x-ots-accesskeyid"
	xOtsContentmd5          = "x-ots-contentmd5"
	xOtsHeaderStsToken      = "x-ots-ststoken"
	xOtsHeaderChargeAdmin   = "x-ots-charge-for-admin"
	xOtsSignature           = "x-ots-signature"
	xOtsRequestCompressType = "x-ots-request-compress-type"
	xOtsRequestCompressSize = "x-ots-request-compress-size"
	xOtsResponseCompressTye = "x-ots-response-compress-type"
	xOtsPrefix              = "x-ots"
	xOtsDateFormat          = "2006-01-02T15:04:05.123Z"
	xOtsInstanceName        = "x-ots-instancename"
	xOtsRequestID           = "x-ots-requestid"
)

func main() {
	flag.Parse()

	server := NewServer(strings.Split(*predefinedInstances, ","))

	router := fasthttprouter.New()
	router.GET("/", requestHandler)
	router.GET("/index", index)
	router.GET("/hello/:name", hello)

	router.POST("/", requestHandler)
	router.POST("/CreateTable", wrap("CreateTable", server.GetHandler("CreateTable")))

	log.Fatal(fasthttp.ListenAndServe(*addr, router.Handler))
}

// TODO: demo code, remove later
func requestHandler(ctx *fasthttp.RequestCtx) {
	s := fmt.Sprintf("Hello, world!\n\n")

	s += fmt.Sprintf("Request method is %q\n", ctx.Method())
	s += fmt.Sprintf("RequestURI is %q\n", ctx.RequestURI())
	s += fmt.Sprintf("Requested path is %q\n", ctx.Path())
	s += fmt.Sprintf("Host is %q\n", ctx.Host())
	s += fmt.Sprintf("Query string is %q\n", ctx.QueryArgs())
	s += fmt.Sprintf("User-Agent is %q\n", ctx.UserAgent())
	s += fmt.Sprintf("Connection has been established at %s\n", ctx.ConnTime())
	s += fmt.Sprintf("Request has been started at %s\n", ctx.Time())
	s += fmt.Sprintf("Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
	s += fmt.Sprintf("Your ip is %q\n\n", ctx.RemoteIP())
	s += fmt.Sprintf("Raw request is:\n---CUT---\n%s\n---CUT---", &ctx.Request)

	fmt.Fprint(ctx, s)
	fmt.Println(s)

	ctx.SetContentType("text/plain; charset=utf8")

	// Set arbitrary headers
	ctx.Response.Header.Set("X-My-Header", "my-header-value")

	// Set cookies
	var c fasthttp.Cookie
	c.SetKey("cookie-name")
	c.SetValue("cookie-value")
	ctx.Response.Header.SetCookie(&c)
}

// RPCMethodHandler ...
type RPCMethodHandler = func(*fasthttp.RequestCtx)

func wrap(methodName string, f ServerMethodHandler) RPCMethodHandler {
	return func(ctx *fasthttp.RequestCtx) {
		reqID := NewRequestID()
		ctxHeader := &ctx.Request.Header
		req := Request{
			Method:         methodName,
			InstanceName:   ctxHeader.Peek(xOtsInstanceName),
			APIVersion:     ctxHeader.Peek(xOtsApiversion),
			Contentmd5:     ctxHeader.Peek(xOtsContentmd5),
			ClientDateTime: ctxHeader.Peek(xOtsDate),
			Signature:      ctxHeader.Peek(xOtsSignature),
			Data:           ctx.Request.Body(),
		}
		resp := f(req)
		ctx.SetBody(resp.Data)
		ctx.SetStatusCode(resp.HTTPStatus)
		ctx.Response.Header.Set(xOtsRequestID, reqID)
	}
}

// TODO: demo code, remove later
func handleCreateTable(ctx *fasthttp.RequestCtx) {
	s := fmt.Sprintf("Receive /CreateTable request\n\n")
	s += fmt.Sprintf("Host is %q\n", ctx.Host())
	s += fmt.Sprintf("User-Agent is %q\n", ctx.UserAgent())
	s += fmt.Sprintf("Connection has been established at %s\n", ctx.ConnTime())
	s += fmt.Sprintf("Request has been started at %s\n", ctx.Time())
	s += fmt.Sprintf("Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
	s += fmt.Sprintf("Client ip is %q\n\n", ctx.RemoteIP())

	otsHeader := parseHeader(&ctx.Request.Header)
	s += fmt.Sprintf("Instance is %s\n", otsHeader.InstanceName)
	s += fmt.Sprintf("API version is %s\n", otsHeader.APIVersion)
	s += fmt.Sprintf("Content MD5 is %s\n", otsHeader.Contentmd5)
	s += fmt.Sprintf("Date is %s\n", otsHeader.Date)
	s += fmt.Sprintf("Signature is %s\n", otsHeader.Signature)

	body := ctx.PostBody()
	request := protocol.CreateTableRequest{}
	err := proto.Unmarshal(body, &request)
	if err != nil {
		s += fmt.Sprintf("Parse Body Error: %s", err)
	} else {
		tableName := request.TableMeta.TableName
		s += fmt.Sprintf("TableName is %s", *tableName)
	}

	fmt.Fprint(ctx, s)
	fmt.Println(s)

	ctx.SetContentType("text/plain; charset=utf8")

	// Set arbitrary headers
	ctx.Response.Header.Set("X-My-Header", "my-header-value")

	// Set cookies
	var c fasthttp.Cookie
	c.SetKey("cookie-name")
	c.SetValue("cookie-value")
	ctx.Response.Header.SetCookie(&c)
}

// TODO: demo code, remove later
func index(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Welcome!\n")
}

// TODO: demo code, remove later
func hello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("name"))
}

// TablestoreRequestHeader common fields in header of each Tablestore request
// TODO: demo code, remove later
type TablestoreRequestHeader struct {
	InstanceName []byte
	APIVersion   []byte
	Contentmd5   []byte
	Date         []byte
	Signature    []byte
}

// TODO: demo code, remove later
func parseHeader(ctxHeader *fasthttp.RequestHeader) TablestoreRequestHeader {
	return TablestoreRequestHeader{
		InstanceName: ctxHeader.Peek(xOtsInstanceName),
		APIVersion:   ctxHeader.Peek(xOtsApiversion),
		Contentmd5:   ctxHeader.Peek(xOtsContentmd5),
		Date:         ctxHeader.Peek(xOtsDate),
		Signature:    ctxHeader.Peek(xOtsSignature),
	}
}
