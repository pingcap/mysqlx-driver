// Code generated by protoc-gen-go.
// source: mysqlx_sql.proto
// DO NOT EDIT!

/*
Package Mysqlx_Sql is a generated protocol buffer package.

Messages of the MySQL Package

It is generated from these files:
	mysqlx_sql.proto

It has these top-level messages:
	StmtExecute
	StmtExecuteOk
*/
package Mysqlx_Sql

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import "github.com/sjmudd/go-mysqlx-driver/Mysqlx_Datatypes"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

// execute a statement in the given namespace
//
// .. uml::
//
//   client -> server: StmtExecute
//   ... zero or more Resultsets ...
//   server --> client: StmtExecuteOk
//
// Notices:
//   This message may generate a notice containing WARNINGs generated by its execution.
//   This message may generate a notice containing INFO messages generated by its execution.
//
// :param namespace: namespace of the statement to be executed
// :param stmt: statement that shall be executed.
// :param args: values for wildcard replacements
// :param compact_metadata: send only type information for :protobuf:msg:`Mysqlx.Resultset::ColumnMetadata`, skipping names and others
// :returns:
//    * zero or one :protobuf:msg:`Mysqlx.Resultset::` followed by :protobuf:msg:`Mysqlx.Sql::StmtExecuteOk`
type StmtExecute struct {
	Namespace        *string                 `protobuf:"bytes,3,opt,name=namespace,def=sql" json:"namespace,omitempty"`
	Stmt             []byte                  `protobuf:"bytes,1,req,name=stmt" json:"stmt,omitempty"`
	Args             []*Mysqlx_Datatypes.Any `protobuf:"bytes,2,rep,name=args" json:"args,omitempty"`
	CompactMetadata  *bool                   `protobuf:"varint,4,opt,name=compact_metadata,def=0" json:"compact_metadata,omitempty"`
	XXX_unrecognized []byte                  `json:"-"`
}

func (m *StmtExecute) Reset()                    { *m = StmtExecute{} }
func (m *StmtExecute) String() string            { return proto.CompactTextString(m) }
func (*StmtExecute) ProtoMessage()               {}
func (*StmtExecute) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

const Default_StmtExecute_Namespace string = "sql"
const Default_StmtExecute_CompactMetadata bool = false

func (m *StmtExecute) GetNamespace() string {
	if m != nil && m.Namespace != nil {
		return *m.Namespace
	}
	return Default_StmtExecute_Namespace
}

func (m *StmtExecute) GetStmt() []byte {
	if m != nil {
		return m.Stmt
	}
	return nil
}

func (m *StmtExecute) GetArgs() []*Mysqlx_Datatypes.Any {
	if m != nil {
		return m.Args
	}
	return nil
}

func (m *StmtExecute) GetCompactMetadata() bool {
	if m != nil && m.CompactMetadata != nil {
		return *m.CompactMetadata
	}
	return Default_StmtExecute_CompactMetadata
}

// statement executed successful
type StmtExecuteOk struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *StmtExecuteOk) Reset()                    { *m = StmtExecuteOk{} }
func (m *StmtExecuteOk) String() string            { return proto.CompactTextString(m) }
func (*StmtExecuteOk) ProtoMessage()               {}
func (*StmtExecuteOk) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*StmtExecute)(nil), "Mysqlx.Sql.StmtExecute")
	proto.RegisterType((*StmtExecuteOk)(nil), "Mysqlx.Sql.StmtExecuteOk")
}

var fileDescriptor0 = []byte{
	// 199 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0xc8, 0xad, 0x2c, 0x2e,
	0xcc, 0xa9, 0x88, 0x07, 0x12, 0x7a, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0x5c, 0xbe, 0x60, 0x11,
	0xbd, 0xe0, 0xc2, 0x1c, 0x29, 0x31, 0xa8, 0x6c, 0x4a, 0x62, 0x49, 0x62, 0x49, 0x65, 0x41, 0x6a,
	0x31, 0x44, 0x8d, 0x52, 0x25, 0x17, 0x77, 0x70, 0x49, 0x6e, 0x89, 0x6b, 0x45, 0x6a, 0x72, 0x69,
	0x49, 0xaa, 0x90, 0x18, 0x17, 0x67, 0x5e, 0x62, 0x6e, 0x6a, 0x71, 0x41, 0x62, 0x72, 0xaa, 0x04,
	0xb3, 0x02, 0xa3, 0x06, 0xa7, 0x15, 0x33, 0x50, 0x9f, 0x10, 0x0f, 0x17, 0x4b, 0x31, 0x50, 0x99,
	0x04, 0xa3, 0x02, 0x93, 0x06, 0x8f, 0x90, 0x32, 0x17, 0x4b, 0x62, 0x51, 0x7a, 0xb1, 0x04, 0x93,
	0x02, 0xb3, 0x06, 0xb7, 0x91, 0xa8, 0x1e, 0xd4, 0x1e, 0x17, 0xb8, 0xd9, 0x8e, 0x79, 0x95, 0x42,
	0xf2, 0x5c, 0x02, 0xc9, 0xf9, 0xb9, 0x40, 0x83, 0x4a, 0xe2, 0x73, 0x53, 0x4b, 0x12, 0x41, 0x16,
	0x4b, 0xb0, 0x00, 0x4d, 0xe4, 0xb0, 0x62, 0x4d, 0x4b, 0xcc, 0x29, 0x4e, 0x55, 0xe2, 0xe7, 0xe2,
	0x45, 0xb2, 0xda, 0x3f, 0xdb, 0x49, 0x8e, 0x4b, 0x06, 0xa8, 0x43, 0x0f, 0xec, 0x52, 0xbd, 0xe4,
	0x2c, 0x08, 0xa3, 0x02, 0xe2, 0xd0, 0xa4, 0xd2, 0x34, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc5,
	0x35, 0x9a, 0x03, 0xe2, 0x00, 0x00, 0x00,
}
