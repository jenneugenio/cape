// Code generated by protoc-gen-go. DO NOT EDIT.
// source: data_connector.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	data_connector.proto

It has these top-level messages:
	QueryRequest
	Record
	Schema
	Field
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type RecordType int32

const (
	RecordType_DOCUMENT RecordType = 0
)

var RecordType_name = map[int32]string{
	0: "DOCUMENT",
}
var RecordType_value = map[string]int32{
	"DOCUMENT": 0,
}

func (x RecordType) String() string {
	return proto1.EnumName(RecordType_name, int32(x))
}
func (RecordType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// We need to know the exact field type so the client can convert the data into the appropriate field type.
type FieldType int32

const (
	FieldType_BIGINT       FieldType = 0
	FieldType_BIGSERIAL    FieldType = 1
	FieldType_BIT          FieldType = 2
	FieldType_VARBIT       FieldType = 3
	FieldType_BOOL         FieldType = 4
	FieldType_BOX          FieldType = 5
	FieldType_BYTEA        FieldType = 6
	FieldType_CHAR         FieldType = 7
	FieldType_VARCHAR      FieldType = 8
	FieldType_CIDR         FieldType = 9
	FieldType_CIRCLE       FieldType = 10
	FieldType_DATE         FieldType = 11
	FieldType_DOUBLE       FieldType = 12
	FieldType_INET         FieldType = 13
	FieldType_INT          FieldType = 14
	FieldType_INTERVAL     FieldType = 15
	FieldType_JSON         FieldType = 16
	FieldType_JSONB        FieldType = 17
	FieldType_LINE         FieldType = 18
	FieldType_LSEG         FieldType = 19
	FieldType_MACADDR      FieldType = 20
	FieldType_MACADDR8     FieldType = 21
	FieldType_MONEY        FieldType = 22
	FieldType_NUMERIC      FieldType = 23
	FieldType_PATH         FieldType = 24
	FieldType_PGLSN        FieldType = 25
	FieldType_POINT        FieldType = 26
	FieldType_POLYGON      FieldType = 27
	FieldType_REAL         FieldType = 28
	FieldType_SMALLINT     FieldType = 29
	FieldType_SMALLSERIAL  FieldType = 30
	FieldType_SERIAL       FieldType = 31
	FieldType_TEXT         FieldType = 32
	FieldType_TIME         FieldType = 33
	FieldType_TIMETZ       FieldType = 34
	FieldType_TIMESTAMP    FieldType = 35
	FieldType_TIMESTAMPTZ  FieldType = 36
	FieldType_TSQUERY      FieldType = 37
	FieldType_TSVECTOR     FieldType = 38
	FieldType_TXIDSNAPSHOT FieldType = 39
	FieldType_UUID         FieldType = 40
	FieldType_XML          FieldType = 41
)

var FieldType_name = map[int32]string{
	0:  "BIGINT",
	1:  "BIGSERIAL",
	2:  "BIT",
	3:  "VARBIT",
	4:  "BOOL",
	5:  "BOX",
	6:  "BYTEA",
	7:  "CHAR",
	8:  "VARCHAR",
	9:  "CIDR",
	10: "CIRCLE",
	11: "DATE",
	12: "DOUBLE",
	13: "INET",
	14: "INT",
	15: "INTERVAL",
	16: "JSON",
	17: "JSONB",
	18: "LINE",
	19: "LSEG",
	20: "MACADDR",
	21: "MACADDR8",
	22: "MONEY",
	23: "NUMERIC",
	24: "PATH",
	25: "PGLSN",
	26: "POINT",
	27: "POLYGON",
	28: "REAL",
	29: "SMALLINT",
	30: "SMALLSERIAL",
	31: "SERIAL",
	32: "TEXT",
	33: "TIME",
	34: "TIMETZ",
	35: "TIMESTAMP",
	36: "TIMESTAMPTZ",
	37: "TSQUERY",
	38: "TSVECTOR",
	39: "TXIDSNAPSHOT",
	40: "UUID",
	41: "XML",
}
var FieldType_value = map[string]int32{
	"BIGINT":       0,
	"BIGSERIAL":    1,
	"BIT":          2,
	"VARBIT":       3,
	"BOOL":         4,
	"BOX":          5,
	"BYTEA":        6,
	"CHAR":         7,
	"VARCHAR":      8,
	"CIDR":         9,
	"CIRCLE":       10,
	"DATE":         11,
	"DOUBLE":       12,
	"INET":         13,
	"INT":          14,
	"INTERVAL":     15,
	"JSON":         16,
	"JSONB":        17,
	"LINE":         18,
	"LSEG":         19,
	"MACADDR":      20,
	"MACADDR8":     21,
	"MONEY":        22,
	"NUMERIC":      23,
	"PATH":         24,
	"PGLSN":        25,
	"POINT":        26,
	"POLYGON":      27,
	"REAL":         28,
	"SMALLINT":     29,
	"SMALLSERIAL":  30,
	"SERIAL":       31,
	"TEXT":         32,
	"TIME":         33,
	"TIMETZ":       34,
	"TIMESTAMP":    35,
	"TIMESTAMPTZ":  36,
	"TSQUERY":      37,
	"TSVECTOR":     38,
	"TXIDSNAPSHOT": 39,
	"UUID":         40,
	"XML":          41,
}

func (x FieldType) String() string {
	return proto1.EnumName(FieldType_name, int32(x))
}
func (FieldType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type QueryRequest struct {
	DataSource string `protobuf:"bytes,1,opt,name=data_source,json=dataSource" json:"data_source,omitempty"`
	Query      string `protobuf:"bytes,2,opt,name=query" json:"query,omitempty"`
}

func (m *QueryRequest) Reset()                    { *m = QueryRequest{} }
func (m *QueryRequest) String() string            { return proto1.CompactTextString(m) }
func (*QueryRequest) ProtoMessage()               {}
func (*QueryRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *QueryRequest) GetDataSource() string {
	if m != nil {
		return m.DataSource
	}
	return ""
}

func (m *QueryRequest) GetQuery() string {
	if m != nil {
		return m.Query
	}
	return ""
}

type Record struct {
	// Schema is only available on the first message
	Schema *Schema `protobuf:"bytes,1,opt,name=schema" json:"schema,omitempty"`
	// Each field value goes over the wire and the number of values maps to the number of fields in the Schema.
	Value [][]byte `protobuf:"bytes,2,rep,name=value,proto3" json:"value,omitempty"`
}

func (m *Record) Reset()                    { *m = Record{} }
func (m *Record) String() string            { return proto1.CompactTextString(m) }
func (*Record) ProtoMessage()               {}
func (*Record) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Record) GetSchema() *Schema {
	if m != nil {
		return m.Schema
	}
	return nil
}

func (m *Record) GetValue() [][]byte {
	if m != nil {
		return m.Value
	}
	return nil
}

// Schema contains information
type Schema struct {
	DataSource string     `protobuf:"bytes,1,opt,name=data_source,json=dataSource" json:"data_source,omitempty"`
	Target     string     `protobuf:"bytes,2,opt,name=target" json:"target,omitempty"`
	Type       RecordType `protobuf:"varint,3,opt,name=type,enum=cape.proto.RecordType" json:"type,omitempty"`
	Fields     []*Field   `protobuf:"bytes,4,rep,name=fields" json:"fields,omitempty"`
}

func (m *Schema) Reset()                    { *m = Schema{} }
func (m *Schema) String() string            { return proto1.CompactTextString(m) }
func (*Schema) ProtoMessage()               {}
func (*Schema) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Schema) GetDataSource() string {
	if m != nil {
		return m.DataSource
	}
	return ""
}

func (m *Schema) GetTarget() string {
	if m != nil {
		return m.Target
	}
	return ""
}

func (m *Schema) GetType() RecordType {
	if m != nil {
		return m.Type
	}
	return RecordType_DOCUMENT
}

func (m *Schema) GetFields() []*Field {
	if m != nil {
		return m.Fields
	}
	return nil
}

// Field represents all information about a field including its type, the number of bits or bytes, and the fields name.
type Field struct {
	Field FieldType `protobuf:"varint,1,opt,name=field,enum=cape.proto.FieldType" json:"field,omitempty"`
	Size  int64     `protobuf:"varint,2,opt,name=size" json:"size,omitempty"`
	Name  string    `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
}

func (m *Field) Reset()                    { *m = Field{} }
func (m *Field) String() string            { return proto1.CompactTextString(m) }
func (*Field) ProtoMessage()               {}
func (*Field) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Field) GetField() FieldType {
	if m != nil {
		return m.Field
	}
	return FieldType_BIGINT
}

func (m *Field) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *Field) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto1.RegisterType((*QueryRequest)(nil), "cape.proto.QueryRequest")
	proto1.RegisterType((*Record)(nil), "cape.proto.Record")
	proto1.RegisterType((*Schema)(nil), "cape.proto.Schema")
	proto1.RegisterType((*Field)(nil), "cape.proto.Field")
	proto1.RegisterEnum("cape.proto.RecordType", RecordType_name, RecordType_value)
	proto1.RegisterEnum("cape.proto.FieldType", FieldType_name, FieldType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for DataConnector service

type DataConnectorClient interface {
	Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (DataConnector_QueryClient, error)
}

type dataConnectorClient struct {
	cc *grpc.ClientConn
}

func NewDataConnectorClient(cc *grpc.ClientConn) DataConnectorClient {
	return &dataConnectorClient{cc}
}

func (c *dataConnectorClient) Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (DataConnector_QueryClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_DataConnector_serviceDesc.Streams[0], c.cc, "/cape.proto.DataConnector/Query", opts...)
	if err != nil {
		return nil, err
	}
	x := &dataConnectorQueryClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DataConnector_QueryClient interface {
	Recv() (*Record, error)
	grpc.ClientStream
}

type dataConnectorQueryClient struct {
	grpc.ClientStream
}

func (x *dataConnectorQueryClient) Recv() (*Record, error) {
	m := new(Record)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for DataConnector service

type DataConnectorServer interface {
	Query(*QueryRequest, DataConnector_QueryServer) error
}

func RegisterDataConnectorServer(s *grpc.Server, srv DataConnectorServer) {
	s.RegisterService(&_DataConnector_serviceDesc, srv)
}

func _DataConnector_Query_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(QueryRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DataConnectorServer).Query(m, &dataConnectorQueryServer{stream})
}

type DataConnector_QueryServer interface {
	Send(*Record) error
	grpc.ServerStream
}

type dataConnectorQueryServer struct {
	grpc.ServerStream
}

func (x *dataConnectorQueryServer) Send(m *Record) error {
	return x.ServerStream.SendMsg(m)
}

var _DataConnector_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cape.proto.DataConnector",
	HandlerType: (*DataConnectorServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Query",
			Handler:       _DataConnector_Query_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "data_connector.proto",
}

func init() { proto1.RegisterFile("data_connector.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 643 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x53, 0xdd, 0x4f, 0x1a, 0x4f,
	0x14, 0xfd, 0x21, 0xb0, 0xc8, 0x05, 0xf5, 0x3a, 0x3f, 0xb5, 0xd4, 0x7e, 0x48, 0xe9, 0x17, 0xd2,
	0xc4, 0x34, 0xf4, 0xa1, 0x7d, 0xdd, 0x8f, 0x29, 0x8c, 0xd9, 0x0f, 0x9c, 0x1d, 0x08, 0x98, 0x26,
	0xcd, 0x16, 0xa7, 0xad, 0x89, 0x0a, 0xc2, 0xd2, 0xc4, 0xfe, 0x25, 0xfd, 0x5f, 0xfb, 0xd2, 0xdc,
	0x61, 0xab, 0x24, 0x3e, 0xf4, 0x89, 0x73, 0xcf, 0x3d, 0xf7, 0xdc, 0xc3, 0xec, 0x0c, 0xec, 0x9c,
	0x25, 0x69, 0xf2, 0x79, 0x3c, 0xb9, 0xba, 0xd2, 0xe3, 0x74, 0x32, 0x3b, 0x9a, 0xce, 0x26, 0xe9,
	0x84, 0xc1, 0x38, 0x99, 0xea, 0x25, 0x6e, 0x70, 0xa8, 0x9e, 0x2c, 0xf4, 0xec, 0x46, 0xea, 0xeb,
	0x85, 0x9e, 0xa7, 0xec, 0x00, 0x2a, 0x66, 0x66, 0x3e, 0x59, 0xcc, 0xc6, 0xba, 0x96, 0xab, 0xe7,
	0x9a, 0x65, 0x09, 0x44, 0xc5, 0x86, 0x61, 0x3b, 0x50, 0xbc, 0xa6, 0x81, 0xda, 0x9a, 0x69, 0x2d,
	0x8b, 0xc6, 0x31, 0x58, 0x52, 0x8f, 0x27, 0xb3, 0x33, 0xd6, 0x02, 0x6b, 0x3e, 0xfe, 0xae, 0x2f,
	0x13, 0x33, 0x5b, 0x69, 0xb3, 0xa3, 0xbb, 0x6d, 0x47, 0xb1, 0xe9, 0xc8, 0x4c, 0x41, 0x5e, 0x3f,
	0x92, 0x8b, 0x85, 0xae, 0xad, 0xd5, 0xf3, 0xcd, 0xaa, 0x5c, 0x16, 0x8d, 0x5f, 0x39, 0xb0, 0x96,
	0xc2, 0x7f, 0xa7, 0xd9, 0x03, 0x2b, 0x4d, 0x66, 0xdf, 0x74, 0x9a, 0xc5, 0xc9, 0x2a, 0xd6, 0x82,
	0x42, 0x7a, 0x33, 0xd5, 0xb5, 0x7c, 0x3d, 0xd7, 0xdc, 0x6c, 0xef, 0xad, 0x66, 0x58, 0xe6, 0x54,
	0x37, 0x53, 0x2d, 0x8d, 0x86, 0x1d, 0x82, 0xf5, 0xf5, 0x5c, 0x5f, 0x9c, 0xcd, 0x6b, 0x85, 0x7a,
	0xbe, 0x59, 0x69, 0x6f, 0xaf, 0xaa, 0x3f, 0x52, 0x47, 0x66, 0x82, 0xc6, 0x27, 0x28, 0x1a, 0x82,
	0xbd, 0x81, 0xa2, 0xa1, 0x4c, 0xa4, 0xcd, 0xf6, 0xee, 0xbd, 0x11, 0xe3, 0xbf, 0xd4, 0x30, 0x06,
	0x85, 0xf9, 0xf9, 0x4f, 0x6d, 0x22, 0xe6, 0xa5, 0xc1, 0xc4, 0x5d, 0x25, 0x97, 0xcb, 0x80, 0x65,
	0x69, 0x70, 0x6b, 0x1f, 0xe0, 0x2e, 0x1c, 0xab, 0xc2, 0xba, 0x17, 0xb9, 0xfd, 0x80, 0x87, 0x0a,
	0xff, 0x6b, 0xfd, 0xce, 0x43, 0xf9, 0xd6, 0x98, 0x01, 0x58, 0x8e, 0xe8, 0x08, 0xea, 0xb0, 0x0d,
	0x28, 0x3b, 0xa2, 0x13, 0x73, 0x29, 0x6c, 0x1f, 0x73, 0xac, 0x04, 0x79, 0x47, 0x28, 0x5c, 0x23,
	0xcd, 0xc0, 0x96, 0x84, 0xf3, 0x6c, 0x1d, 0x0a, 0x4e, 0x14, 0xf9, 0x58, 0x30, 0xed, 0x68, 0x88,
	0x45, 0x56, 0x86, 0xa2, 0x33, 0x52, 0xdc, 0x46, 0x8b, 0xba, 0x6e, 0xd7, 0x96, 0x58, 0x62, 0x15,
	0x28, 0x0d, 0x6c, 0x69, 0x8a, 0x75, 0x43, 0x0b, 0x4f, 0x62, 0x99, 0xac, 0x5c, 0x21, 0x5d, 0x9f,
	0x23, 0x10, 0xeb, 0xd9, 0x8a, 0x63, 0x85, 0x58, 0x2f, 0xea, 0x3b, 0x3e, 0xc7, 0x2a, 0xb1, 0x22,
	0xe4, 0x0a, 0x37, 0x68, 0x01, 0xe5, 0xda, 0xa4, 0xfc, 0x22, 0x54, 0x5c, 0x0e, 0x6c, 0x1f, 0xb7,
	0x48, 0x70, 0x1c, 0x47, 0x21, 0x22, 0x2d, 0x26, 0xe4, 0xe0, 0x36, 0x91, 0xbe, 0x08, 0x39, 0x32,
	0x83, 0x62, 0xde, 0xc1, 0xff, 0x29, 0x42, 0x60, 0xbb, 0xb6, 0xe7, 0x49, 0xdc, 0x21, 0x8f, 0xac,
	0xf8, 0x80, 0xbb, 0x34, 0x19, 0x44, 0x21, 0x1f, 0xe1, 0x1e, 0xa9, 0xc2, 0x7e, 0xc0, 0xa5, 0x70,
	0xf1, 0x01, 0x0d, 0xf7, 0x6c, 0xd5, 0xc5, 0x1a, 0x29, 0x7a, 0x1d, 0x3f, 0x0e, 0xf1, 0xa1, 0x81,
	0x11, 0x25, 0xd9, 0x27, 0x71, 0x2f, 0xf2, 0x47, 0x9d, 0x28, 0xc4, 0x47, 0x24, 0x96, 0xdc, 0xf6,
	0xf1, 0x31, 0x99, 0xc7, 0x81, 0xed, 0xfb, 0x24, 0x7a, 0xc2, 0xb6, 0xa0, 0x62, 0xaa, 0xec, 0x20,
	0x9f, 0xd2, 0xdf, 0xcb, 0xf0, 0x01, 0x0d, 0x29, 0x3e, 0x54, 0x58, 0x37, 0x48, 0x04, 0x1c, 0x9f,
	0x51, 0x9f, 0x90, 0x3a, 0xc5, 0x06, 0x7d, 0x03, 0xc2, 0xb1, 0xb2, 0x83, 0x1e, 0x3e, 0x27, 0xaf,
	0xdb, 0x52, 0x9d, 0xe2, 0x0b, 0x4a, 0xa0, 0xe2, 0x93, 0x3e, 0x97, 0x23, 0x7c, 0x49, 0x7b, 0x55,
	0x3c, 0xe0, 0xae, 0x8a, 0x24, 0xbe, 0x62, 0x08, 0x55, 0x35, 0x14, 0x5e, 0x1c, 0xda, 0xbd, 0xb8,
	0x1b, 0x29, 0x7c, 0x4d, 0x2b, 0xfa, 0x7d, 0xe1, 0x61, 0x93, 0xce, 0x72, 0x18, 0xf8, 0x78, 0xd8,
	0xee, 0xc2, 0x86, 0x97, 0xa4, 0x89, 0xfb, 0xf7, 0x21, 0xb3, 0xf7, 0x50, 0x34, 0xcf, 0x96, 0xd5,
	0x56, 0x6f, 0xde, 0xea, 0x4b, 0xde, 0x67, 0xf7, 0x2f, 0xfd, 0xdb, 0x9c, 0x53, 0x3a, 0x2d, 0x1a,
	0xe6, 0x8b, 0x65, 0x7e, 0xde, 0xfd, 0x09, 0x00, 0x00, 0xff, 0xff, 0x25, 0x42, 0x88, 0xc2, 0x23,
	0x04, 0x00, 0x00,
}
