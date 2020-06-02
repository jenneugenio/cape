// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.11.4
// source: data_connector.proto

package proto

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type RecordType int32

const (
	RecordType_DOCUMENT RecordType = 0
)

// Enum value maps for RecordType.
var (
	RecordType_name = map[int32]string{
		0: "DOCUMENT",
	}
	RecordType_value = map[string]int32{
		"DOCUMENT": 0,
	}
)

func (x RecordType) Enum() *RecordType {
	p := new(RecordType)
	*p = x
	return p
}

func (x RecordType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RecordType) Descriptor() protoreflect.EnumDescriptor {
	return file_data_connector_proto_enumTypes[0].Descriptor()
}

func (RecordType) Type() protoreflect.EnumType {
	return &file_data_connector_proto_enumTypes[0]
}

func (x RecordType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RecordType.Descriptor instead.
func (RecordType) EnumDescriptor() ([]byte, []int) {
	return file_data_connector_proto_rawDescGZIP(), []int{0}
}

// We need to know the exact field type so the client can convert the data into the appropriate field type.
type FieldType int32

const (
	FieldType_BIGINT    FieldType = 0
	FieldType_INT       FieldType = 1
	FieldType_SMALLINT  FieldType = 2
	FieldType_BOOL      FieldType = 3
	FieldType_BYTEA     FieldType = 4
	FieldType_CHAR      FieldType = 5
	FieldType_VARCHAR   FieldType = 6
	FieldType_TEXT      FieldType = 7
	FieldType_DOUBLE    FieldType = 8
	FieldType_REAL      FieldType = 9
	FieldType_TIMESTAMP FieldType = 10
)

// Enum value maps for FieldType.
var (
	FieldType_name = map[int32]string{
		0:  "BIGINT",
		1:  "INT",
		2:  "SMALLINT",
		3:  "BOOL",
		4:  "BYTEA",
		5:  "CHAR",
		6:  "VARCHAR",
		7:  "TEXT",
		8:  "DOUBLE",
		9:  "REAL",
		10: "TIMESTAMP",
	}
	FieldType_value = map[string]int32{
		"BIGINT":    0,
		"INT":       1,
		"SMALLINT":  2,
		"BOOL":      3,
		"BYTEA":     4,
		"CHAR":      5,
		"VARCHAR":   6,
		"TEXT":      7,
		"DOUBLE":    8,
		"REAL":      9,
		"TIMESTAMP": 10,
	}
)

func (x FieldType) Enum() *FieldType {
	p := new(FieldType)
	*p = x
	return p
}

func (x FieldType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FieldType) Descriptor() protoreflect.EnumDescriptor {
	return file_data_connector_proto_enumTypes[1].Descriptor()
}

func (FieldType) Type() protoreflect.EnumType {
	return &file_data_connector_proto_enumTypes[1]
}

func (x FieldType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FieldType.Descriptor instead.
func (FieldType) EnumDescriptor() ([]byte, []int) {
	return file_data_connector_proto_rawDescGZIP(), []int{1}
}

type QueryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DataSource string `protobuf:"bytes,1,opt,name=data_source,json=dataSource,proto3" json:"data_source,omitempty"` // Label of the data source we’re targeting
	Query      string `protobuf:"bytes,2,opt,name=query,proto3" json:"query,omitempty"`
	Limit      int64  `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset     int64  `protobuf:"varint,4,opt,name=offset,proto3" json:"offset,omitempty"`
}

func (x *QueryRequest) Reset() {
	*x = QueryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_connector_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryRequest) ProtoMessage() {}

func (x *QueryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_data_connector_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryRequest.ProtoReflect.Descriptor instead.
func (*QueryRequest) Descriptor() ([]byte, []int) {
	return file_data_connector_proto_rawDescGZIP(), []int{0}
}

func (x *QueryRequest) GetDataSource() string {
	if x != nil {
		return x.DataSource
	}
	return ""
}

func (x *QueryRequest) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

func (x *QueryRequest) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *QueryRequest) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

type Record struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Schema is only available on the first message
	Schema *Schema `protobuf:"bytes,1,opt,name=schema,proto3" json:"schema,omitempty"`
	// Each field value goes over the wire and the number of values maps to the number of fields in the Definition.
	Fields []*Field `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
}

func (x *Record) Reset() {
	*x = Record{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_connector_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Record) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Record) ProtoMessage() {}

func (x *Record) ProtoReflect() protoreflect.Message {
	mi := &file_data_connector_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Record.ProtoReflect.Descriptor instead.
func (*Record) Descriptor() ([]byte, []int) {
	return file_data_connector_proto_rawDescGZIP(), []int{1}
}

func (x *Record) GetSchema() *Schema {
	if x != nil {
		return x.Schema
	}
	return nil
}

func (x *Record) GetFields() []*Field {
	if x != nil {
		return x.Fields
	}
	return nil
}

type Field struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Value:
	//	*Field_Int32
	//	*Field_Int64
	//	*Field_Timestamp
	//	*Field_Float
	//	*Field_Double
	//	*Field_String_
	//	*Field_Bytes
	//	*Field_Bool
	Value isField_Value `protobuf_oneof:"value"`
}

func (x *Field) Reset() {
	*x = Field{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_connector_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Field) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Field) ProtoMessage() {}

func (x *Field) ProtoReflect() protoreflect.Message {
	mi := &file_data_connector_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Field.ProtoReflect.Descriptor instead.
func (*Field) Descriptor() ([]byte, []int) {
	return file_data_connector_proto_rawDescGZIP(), []int{2}
}

func (m *Field) GetValue() isField_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *Field) GetInt32() int32 {
	if x, ok := x.GetValue().(*Field_Int32); ok {
		return x.Int32
	}
	return 0
}

func (x *Field) GetInt64() int64 {
	if x, ok := x.GetValue().(*Field_Int64); ok {
		return x.Int64
	}
	return 0
}

func (x *Field) GetTimestamp() *timestamp.Timestamp {
	if x, ok := x.GetValue().(*Field_Timestamp); ok {
		return x.Timestamp
	}
	return nil
}

func (x *Field) GetFloat() float32 {
	if x, ok := x.GetValue().(*Field_Float); ok {
		return x.Float
	}
	return 0
}

func (x *Field) GetDouble() float64 {
	if x, ok := x.GetValue().(*Field_Double); ok {
		return x.Double
	}
	return 0
}

func (x *Field) GetString_() string {
	if x, ok := x.GetValue().(*Field_String_); ok {
		return x.String_
	}
	return ""
}

func (x *Field) GetBytes() []byte {
	if x, ok := x.GetValue().(*Field_Bytes); ok {
		return x.Bytes
	}
	return nil
}

func (x *Field) GetBool() bool {
	if x, ok := x.GetValue().(*Field_Bool); ok {
		return x.Bool
	}
	return false
}

type isField_Value interface {
	isField_Value()
}

type Field_Int32 struct {
	Int32 int32 `protobuf:"varint,1,opt,name=int32,proto3,oneof"`
}

type Field_Int64 struct {
	Int64 int64 `protobuf:"varint,2,opt,name=int64,proto3,oneof"`
}

type Field_Timestamp struct {
	Timestamp *timestamp.Timestamp `protobuf:"bytes,3,opt,name=timestamp,proto3,oneof"`
}

type Field_Float struct {
	Float float32 `protobuf:"fixed32,4,opt,name=float,proto3,oneof"`
}

type Field_Double struct {
	Double float64 `protobuf:"fixed64,5,opt,name=double,proto3,oneof"`
}

type Field_String_ struct {
	String_ string `protobuf:"bytes,6,opt,name=string,proto3,oneof"`
}

type Field_Bytes struct {
	Bytes []byte `protobuf:"bytes,7,opt,name=bytes,proto3,oneof"`
}

type Field_Bool struct {
	Bool bool `protobuf:"varint,8,opt,name=bool,proto3,oneof"`
}

func (*Field_Int32) isField_Value() {}

func (*Field_Int64) isField_Value() {}

func (*Field_Timestamp) isField_Value() {}

func (*Field_Float) isField_Value() {}

func (*Field_Double) isField_Value() {}

func (*Field_String_) isField_Value() {}

func (*Field_Bytes) isField_Value() {}

func (*Field_Bool) isField_Value() {}

type Schema struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DataSource string       `protobuf:"bytes,1,opt,name=data_source,json=dataSource,proto3" json:"data_source,omitempty"` // label of data source data is being returned from
	Target     string       `protobuf:"bytes,2,opt,name=target,proto3" json:"target,omitempty"`                           // the target of the schema (e.g. a postgres table)
	Type       RecordType   `protobuf:"varint,3,opt,name=type,proto3,enum=cape.RecordType" json:"type,omitempty"`
	Fields     []*FieldInfo `protobuf:"bytes,4,rep,name=fields,proto3" json:"fields,omitempty"`
}

func (x *Schema) Reset() {
	*x = Schema{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_connector_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Schema) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Schema) ProtoMessage() {}

func (x *Schema) ProtoReflect() protoreflect.Message {
	mi := &file_data_connector_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Definition.ProtoReflect.Descriptor instead.
func (*Schema) Descriptor() ([]byte, []int) {
	return file_data_connector_proto_rawDescGZIP(), []int{3}
}

func (x *Schema) GetDataSource() string {
	if x != nil {
		return x.DataSource
	}
	return ""
}

func (x *Schema) GetTarget() string {
	if x != nil {
		return x.Target
	}
	return ""
}

func (x *Schema) GetType() RecordType {
	if x != nil {
		return x.Type
	}
	return RecordType_DOCUMENT
}

func (x *Schema) GetFields() []*FieldInfo {
	if x != nil {
		return x.Fields
	}
	return nil
}

type SchemaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DataSource string `protobuf:"bytes,1,opt,name=data_source,json=dataSource,proto3" json:"data_source,omitempty"`
}

func (x *SchemaRequest) Reset() {
	*x = SchemaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_connector_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SchemaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SchemaRequest) ProtoMessage() {}

func (x *SchemaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_data_connector_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SchemaRequest.ProtoReflect.Descriptor instead.
func (*SchemaRequest) Descriptor() ([]byte, []int) {
	return file_data_connector_proto_rawDescGZIP(), []int{4}
}

func (x *SchemaRequest) GetDataSource() string {
	if x != nil {
		return x.DataSource
	}
	return ""
}

// Returned when you request database schema
type SchemaResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Schemas []*Schema `protobuf:"bytes,1,rep,name=schemas,proto3" json:"schemas,omitempty"`
}

func (x *SchemaResponse) Reset() {
	*x = SchemaResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_connector_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SchemaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SchemaResponse) ProtoMessage() {}

func (x *SchemaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_data_connector_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SchemaResponse.ProtoReflect.Descriptor instead.
func (*SchemaResponse) Descriptor() ([]byte, []int) {
	return file_data_connector_proto_rawDescGZIP(), []int{5}
}

func (x *SchemaResponse) GetSchemas() []*Schema {
	if x != nil {
		return x.Schemas
	}
	return nil
}

// FieldInfo represents all information about a field including its type, the number of bits or bytes, and the fields name.
type FieldInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Field FieldType `protobuf:"varint,1,opt,name=field,proto3,enum=cape.FieldType" json:"field,omitempty"`
	Size  int64     `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	Name  string    `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *FieldInfo) Reset() {
	*x = FieldInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_connector_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FieldInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldInfo) ProtoMessage() {}

func (x *FieldInfo) ProtoReflect() protoreflect.Message {
	mi := &file_data_connector_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FieldInfo.ProtoReflect.Descriptor instead.
func (*FieldInfo) Descriptor() ([]byte, []int) {
	return file_data_connector_proto_rawDescGZIP(), []int{6}
}

func (x *FieldInfo) GetField() FieldType {
	if x != nil {
		return x.Field
	}
	return FieldType_BIGINT
}

func (x *FieldInfo) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *FieldInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type VersionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *VersionRequest) Reset() {
	*x = VersionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_connector_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VersionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersionRequest) ProtoMessage() {}

func (x *VersionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_data_connector_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersionRequest.ProtoReflect.Descriptor instead.
func (*VersionRequest) Descriptor() ([]byte, []int) {
	return file_data_connector_proto_rawDescGZIP(), []int{7}
}

type VersionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version   string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	BuildDate string `protobuf:"bytes,2,opt,name=build_date,json=buildDate,proto3" json:"build_date,omitempty"`
}

func (x *VersionResponse) Reset() {
	*x = VersionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_connector_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VersionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersionResponse) ProtoMessage() {}

func (x *VersionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_data_connector_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersionResponse.ProtoReflect.Descriptor instead.
func (*VersionResponse) Descriptor() ([]byte, []int) {
	return file_data_connector_proto_rawDescGZIP(), []int{8}
}

func (x *VersionResponse) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *VersionResponse) GetBuildDate() string {
	if x != nil {
		return x.BuildDate
	}
	return ""
}

var File_data_connector_proto protoreflect.FileDescriptor

var file_data_connector_proto_rawDesc = []byte{
	0x0a, 0x14, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63, 0x61, 0x70, 0x65, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x73, 0x0a,
	0x0c, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a,
	0x0b, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71,
	0x75, 0x65, 0x72, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66,
	0x66, 0x73, 0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73,
	0x65, 0x74, 0x22, 0x53, 0x0a, 0x06, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x24, 0x0a, 0x06,
	0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x63,
	0x61, 0x70, 0x65, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x06, 0x73, 0x63, 0x68, 0x65,
	0x6d, 0x61, 0x12, 0x23, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x63, 0x61, 0x70, 0x65, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x52,
	0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x22, 0xf6, 0x01, 0x0a, 0x05, 0x46, 0x69, 0x65, 0x6c,
	0x64, 0x12, 0x16, 0x0a, 0x05, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x48, 0x00, 0x52, 0x05, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x12, 0x16, 0x0a, 0x05, 0x69, 0x6e, 0x74,
	0x36, 0x34, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x05, 0x69, 0x6e, 0x74, 0x36,
	0x34, 0x12, 0x3a, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x48, 0x00, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x16, 0x0a,
	0x05, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x48, 0x00, 0x52, 0x05,
	0x66, 0x6c, 0x6f, 0x61, 0x74, 0x12, 0x18, 0x0a, 0x06, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x01, 0x48, 0x00, 0x52, 0x06, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x12,
	0x18, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x00, 0x52, 0x06, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12, 0x16, 0x0a, 0x05, 0x62, 0x79, 0x74,
	0x65, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x05, 0x62, 0x79, 0x74, 0x65,
	0x73, 0x12, 0x14, 0x0a, 0x04, 0x62, 0x6f, 0x6f, 0x6c, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x48,
	0x00, 0x52, 0x04, 0x62, 0x6f, 0x6f, 0x6c, 0x42, 0x07, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x90, 0x01, 0x0a, 0x06, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x1f, 0x0a, 0x0b, 0x64,
	0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x12, 0x24, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x10, 0x2e, 0x63, 0x61, 0x70, 0x65, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x27, 0x0a, 0x06, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x63, 0x61, 0x70,
	0x65, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x06, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x73, 0x22, 0x30, 0x0a, 0x0d, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x53,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x38, 0x0a, 0x0e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x07, 0x73, 0x63, 0x68, 0x65, 0x6d,
	0x61, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x63, 0x61, 0x70, 0x65, 0x2e,
	0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x07, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x73, 0x22,
	0x5a, 0x0a, 0x09, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x25, 0x0a, 0x05,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x63, 0x61,
	0x70, 0x65, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x54, 0x79, 0x70, 0x65, 0x52, 0x05, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x10, 0x0a, 0x0e, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x4a, 0x0a,
	0x0f, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x75,
	0x69, 0x6c, 0x64, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x62, 0x75, 0x69, 0x6c, 0x64, 0x44, 0x61, 0x74, 0x65, 0x2a, 0x1a, 0x0a, 0x0a, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x4f, 0x43, 0x55, 0x4d,
	0x45, 0x4e, 0x54, 0x10, 0x00, 0x2a, 0x89, 0x01, 0x0a, 0x09, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x42, 0x49, 0x47, 0x49, 0x4e, 0x54, 0x10, 0x00, 0x12,
	0x07, 0x0a, 0x03, 0x49, 0x4e, 0x54, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x53, 0x4d, 0x41, 0x4c,
	0x4c, 0x49, 0x4e, 0x54, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x42, 0x4f, 0x4f, 0x4c, 0x10, 0x03,
	0x12, 0x09, 0x0a, 0x05, 0x42, 0x59, 0x54, 0x45, 0x41, 0x10, 0x04, 0x12, 0x08, 0x0a, 0x04, 0x43,
	0x48, 0x41, 0x52, 0x10, 0x05, 0x12, 0x0b, 0x0a, 0x07, 0x56, 0x41, 0x52, 0x43, 0x48, 0x41, 0x52,
	0x10, 0x06, 0x12, 0x08, 0x0a, 0x04, 0x54, 0x45, 0x58, 0x54, 0x10, 0x07, 0x12, 0x0a, 0x0a, 0x06,
	0x44, 0x4f, 0x55, 0x42, 0x4c, 0x45, 0x10, 0x08, 0x12, 0x08, 0x0a, 0x04, 0x52, 0x45, 0x41, 0x4c,
	0x10, 0x09, 0x12, 0x0d, 0x0a, 0x09, 0x54, 0x49, 0x4d, 0x45, 0x53, 0x54, 0x41, 0x4d, 0x50, 0x10,
	0x0a, 0x32, 0xa9, 0x01, 0x0a, 0x0d, 0x44, 0x61, 0x74, 0x61, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x12, 0x2b, 0x0a, 0x05, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x12, 0x2e, 0x63,
	0x61, 0x70, 0x65, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x0c, 0x2e, 0x63, 0x61, 0x70, 0x65, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x30, 0x01,
	0x12, 0x33, 0x0a, 0x06, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x13, 0x2e, 0x63, 0x61, 0x70,
	0x65, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x14, 0x2e, 0x63, 0x61, 0x70, 0x65, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x14, 0x2e, 0x63, 0x61, 0x70, 0x65, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x63, 0x61, 0x70, 0x65, 0x2e, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x09, 0x5a,
	0x07, 0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_data_connector_proto_rawDescOnce sync.Once
	file_data_connector_proto_rawDescData = file_data_connector_proto_rawDesc
)

func file_data_connector_proto_rawDescGZIP() []byte {
	file_data_connector_proto_rawDescOnce.Do(func() {
		file_data_connector_proto_rawDescData = protoimpl.X.CompressGZIP(file_data_connector_proto_rawDescData)
	})
	return file_data_connector_proto_rawDescData
}

var file_data_connector_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_data_connector_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_data_connector_proto_goTypes = []interface{}{
	(RecordType)(0),             // 0: cape.RecordType
	(FieldType)(0),              // 1: cape.FieldType
	(*QueryRequest)(nil),        // 2: cape.QueryRequest
	(*Record)(nil),              // 3: cape.Record
	(*Field)(nil),               // 4: cape.Field
	(*Schema)(nil),              // 5: cape.Definition
	(*SchemaRequest)(nil),       // 6: cape.SchemaRequest
	(*SchemaResponse)(nil),      // 7: cape.SchemaResponse
	(*FieldInfo)(nil),           // 8: cape.FieldInfo
	(*VersionRequest)(nil),      // 9: cape.VersionRequest
	(*VersionResponse)(nil),     // 10: cape.VersionResponse
	(*timestamp.Timestamp)(nil), // 11: google.protobuf.Timestamp
}
var file_data_connector_proto_depIdxs = []int32{
	5,  // 0: cape.Record.schema:type_name -> cape.Definition
	4,  // 1: cape.Record.fields:type_name -> cape.Field
	11, // 2: cape.Field.timestamp:type_name -> google.protobuf.Timestamp
	0,  // 3: cape.Definition.type:type_name -> cape.RecordType
	8,  // 4: cape.Definition.fields:type_name -> cape.FieldInfo
	5,  // 5: cape.SchemaResponse.schemas:type_name -> cape.Definition
	1,  // 6: cape.FieldInfo.field:type_name -> cape.FieldType
	2,  // 7: cape.DataConnector.Query:input_type -> cape.QueryRequest
	6,  // 8: cape.DataConnector.Definition:input_type -> cape.SchemaRequest
	9,  // 9: cape.DataConnector.Version:input_type -> cape.VersionRequest
	3,  // 10: cape.DataConnector.Query:output_type -> cape.Record
	7,  // 11: cape.DataConnector.Definition:output_type -> cape.SchemaResponse
	10, // 12: cape.DataConnector.Version:output_type -> cape.VersionResponse
	10, // [10:13] is the sub-list for method output_type
	7,  // [7:10] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_data_connector_proto_init() }
func file_data_connector_proto_init() {
	if File_data_connector_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_data_connector_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_data_connector_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Record); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_data_connector_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Field); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_data_connector_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Schema); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_data_connector_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SchemaRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_data_connector_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SchemaResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_data_connector_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FieldInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_data_connector_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VersionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_data_connector_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VersionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_data_connector_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*Field_Int32)(nil),
		(*Field_Int64)(nil),
		(*Field_Timestamp)(nil),
		(*Field_Float)(nil),
		(*Field_Double)(nil),
		(*Field_String_)(nil),
		(*Field_Bytes)(nil),
		(*Field_Bool)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_data_connector_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_data_connector_proto_goTypes,
		DependencyIndexes: file_data_connector_proto_depIdxs,
		EnumInfos:         file_data_connector_proto_enumTypes,
		MessageInfos:      file_data_connector_proto_msgTypes,
	}.Build()
	File_data_connector_proto = out.File
	file_data_connector_proto_rawDesc = nil
	file_data_connector_proto_goTypes = nil
	file_data_connector_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// DataConnectorClient is the client API for DataConnector service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DataConnectorClient interface {
	Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (DataConnector_QueryClient, error)
	Schema(ctx context.Context, in *SchemaRequest, opts ...grpc.CallOption) (*SchemaResponse, error)
	Version(ctx context.Context, in *VersionRequest, opts ...grpc.CallOption) (*VersionResponse, error)
}

type dataConnectorClient struct {
	cc grpc.ClientConnInterface
}

func NewDataConnectorClient(cc grpc.ClientConnInterface) DataConnectorClient {
	return &dataConnectorClient{cc}
}

func (c *dataConnectorClient) Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (DataConnector_QueryClient, error) {
	stream, err := c.cc.NewStream(ctx, &_DataConnector_serviceDesc.Streams[0], "/cape.DataConnector/Query", opts...)
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

func (c *dataConnectorClient) Schema(ctx context.Context, in *SchemaRequest, opts ...grpc.CallOption) (*SchemaResponse, error) {
	out := new(SchemaResponse)
	err := c.cc.Invoke(ctx, "/cape.DataConnector/Definition", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataConnectorClient) Version(ctx context.Context, in *VersionRequest, opts ...grpc.CallOption) (*VersionResponse, error) {
	out := new(VersionResponse)
	err := c.cc.Invoke(ctx, "/cape.DataConnector/Version", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataConnectorServer is the server API for DataConnector service.
type DataConnectorServer interface {
	Query(*QueryRequest, DataConnector_QueryServer) error
	Schema(context.Context, *SchemaRequest) (*SchemaResponse, error)
	Version(context.Context, *VersionRequest) (*VersionResponse, error)
}

// UnimplementedDataConnectorServer can be embedded to have forward compatible implementations.
type UnimplementedDataConnectorServer struct {
}

func (*UnimplementedDataConnectorServer) Query(*QueryRequest, DataConnector_QueryServer) error {
	return status.Errorf(codes.Unimplemented, "method Query not implemented")
}
func (*UnimplementedDataConnectorServer) Schema(context.Context, *SchemaRequest) (*SchemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Definition not implemented")
}
func (*UnimplementedDataConnectorServer) Version(context.Context, *VersionRequest) (*VersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Version not implemented")
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

func _DataConnector_Schema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataConnectorServer).Schema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cape.DataConnector/Definition",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataConnectorServer).Schema(ctx, req.(*SchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataConnector_Version_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataConnectorServer).Version(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cape.DataConnector/Version",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataConnectorServer).Version(ctx, req.(*VersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DataConnector_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cape.DataConnector",
	HandlerType: (*DataConnectorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Definition",
			Handler:    _DataConnector_Schema_Handler,
		},
		{
			MethodName: "Version",
			Handler:    _DataConnector_Version_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Query",
			Handler:       _DataConnector_Query_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "data_connector.proto",
}
