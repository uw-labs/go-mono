// Copyright 2019 Google LLC.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.22.0
// 	protoc        v3.11.2
// source: google/api/client.proto

package annotations

import (
	reflect "reflect"

	proto "github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

var file_google_api_client_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptor.MethodOptions)(nil),
		ExtensionType: ([]string)(nil),
		Field:         1051,
		Name:          "google.api.method_signature",
		Tag:           "bytes,1051,rep,name=method_signature",
		Filename:      "google/api/client.proto",
	},
	{
		ExtendedType:  (*descriptor.ServiceOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         1049,
		Name:          "google.api.default_host",
		Tag:           "bytes,1049,opt,name=default_host",
		Filename:      "google/api/client.proto",
	},
	{
		ExtendedType:  (*descriptor.ServiceOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         1050,
		Name:          "google.api.oauth_scopes",
		Tag:           "bytes,1050,opt,name=oauth_scopes",
		Filename:      "google/api/client.proto",
	},
}

// Extension fields to descriptor.MethodOptions.
var (
	// A definition of a client library method signature.
	//
	// In client libraries, each proto RPC corresponds to one or more methods
	// which the end user is able to call, and calls the underlying RPC.
	// Normally, this method receives a single argument (a struct or instance
	// corresponding to the RPC request object). Defining this field will
	// add one or more overloads providing flattened or simpler method signatures
	// in some languages.
	//
	// The fields on the method signature are provided as a comma-separated
	// string.
	//
	// For example, the proto RPC and annotation:
	//
	//   rpc CreateSubscription(CreateSubscriptionRequest)
	//       returns (Subscription) {
	//     option (google.api.method_signature) = "name,topic";
	//   }
	//
	// Would add the following Java overload (in addition to the method accepting
	// the request object):
	//
	//   public final Subscription createSubscription(String name, String topic)
	//
	// The following backwards-compatibility guidelines apply:
	//
	//   * Adding this annotation to an unannotated method is backwards
	//     compatible.
	//   * Adding this annotation to a method which already has existing
	//     method signature annotations is backwards compatible if and only if
	//     the new method signature annotation is last in the sequence.
	//   * Modifying or removing an existing method signature annotation is
	//     a breaking change.
	//   * Re-ordering existing method signature annotations is a breaking
	//     change.
	//
	// repeated string method_signature = 1051;
	E_MethodSignature = &file_google_api_client_proto_extTypes[0]
)

// Extension fields to descriptor.ServiceOptions.
var (
	// The hostname for this service.
	// This should be specified with no prefix or protocol.
	//
	// Example:
	//
	//   service Foo {
	//     option (google.api.default_host) = "foo.googleapi.com";
	//     ...
	//   }
	//
	// optional string default_host = 1049;
	E_DefaultHost = &file_google_api_client_proto_extTypes[1]
	// OAuth scopes needed for the client.
	//
	// Example:
	//
	//   service Foo {
	//     option (google.api.oauth_scopes) = \
	//       "https://www.googleapis.com/auth/cloud-platform";
	//     ...
	//   }
	//
	// If there is more than one scope, use a comma-separated string:
	//
	// Example:
	//
	//   service Foo {
	//     option (google.api.oauth_scopes) = \
	//       "https://www.googleapis.com/auth/cloud-platform,"
	//       "https://www.googleapis.com/auth/monitoring";
	//     ...
	//   }
	//
	// optional string oauth_scopes = 1050;
	E_OauthScopes = &file_google_api_client_proto_extTypes[2]
)

var File_google_api_client_proto protoreflect.FileDescriptor

var file_google_api_client_proto_rawDesc = []byte{
	0x0a, 0x17, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x61, 0x70, 0x69, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3a, 0x4a, 0x0a, 0x10, 0x6d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x1e, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x9b, 0x08, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x0f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x3a, 0x43, 0x0a, 0x0c, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x68,
	0x6f, 0x73, 0x74, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0x99, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x66,
	0x61, 0x75, 0x6c, 0x74, 0x48, 0x6f, 0x73, 0x74, 0x3a, 0x43, 0x0a, 0x0c, 0x6f, 0x61, 0x75, 0x74,
	0x68, 0x5f, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x73, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x9a, 0x08, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x53, 0x63, 0x6f, 0x70, 0x65, 0x73, 0x42, 0x69, 0x0a,
	0x0e, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x42,
	0x0b, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x41,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2e, 0x6f, 0x72,
	0x67, 0x2f, 0x67, 0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x3b, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0xa2, 0x02, 0x04, 0x47, 0x41, 0x50, 0x49, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_google_api_client_proto_goTypes = []interface{}{
	(*descriptor.MethodOptions)(nil),  // 0: google.protobuf.MethodOptions
	(*descriptor.ServiceOptions)(nil), // 1: google.protobuf.ServiceOptions
}
var file_google_api_client_proto_depIdxs = []int32{
	0, // 0: google.api.method_signature:extendee -> google.protobuf.MethodOptions
	1, // 1: google.api.default_host:extendee -> google.protobuf.ServiceOptions
	1, // 2: google.api.oauth_scopes:extendee -> google.protobuf.ServiceOptions
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	0, // [0:3] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_google_api_client_proto_init() }
func file_google_api_client_proto_init() {
	if File_google_api_client_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_google_api_client_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 3,
			NumServices:   0,
		},
		GoTypes:           file_google_api_client_proto_goTypes,
		DependencyIndexes: file_google_api_client_proto_depIdxs,
		ExtensionInfos:    file_google_api_client_proto_extTypes,
	}.Build()
	File_google_api_client_proto = out.File
	file_google_api_client_proto_rawDesc = nil
	file_google_api_client_proto_goTypes = nil
	file_google_api_client_proto_depIdxs = nil
}
