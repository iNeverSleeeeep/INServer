// Code generated by protoc-gen-go-json. DO NOT EDIT.
// source: math.proto

package engine

import (
	"bytes"

	"github.com/golang/protobuf/jsonpb"
)

// MarshalJSON implements json.Marshaler
func (msg *Vector2) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	err := (&jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: false,
		OrigName:     false,
	}).Marshal(&buf, msg)
	return buf.Bytes(), err
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *Vector2) UnmarshalJSON(b []byte) error {
	return jsonpb.Unmarshal(bytes.NewReader(b), msg)
}

// MarshalJSON implements json.Marshaler
func (msg *Vector3) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	err := (&jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: false,
		OrigName:     false,
	}).Marshal(&buf, msg)
	return buf.Bytes(), err
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *Vector3) UnmarshalJSON(b []byte) error {
	return jsonpb.Unmarshal(bytes.NewReader(b), msg)
}

// MarshalJSON implements json.Marshaler
func (msg *Rect) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	err := (&jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: false,
		OrigName:     false,
	}).Marshal(&buf, msg)
	return buf.Bytes(), err
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *Rect) UnmarshalJSON(b []byte) error {
	return jsonpb.Unmarshal(bytes.NewReader(b), msg)
}

// MarshalJSON implements json.Marshaler
func (msg *Quaternion) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	err := (&jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: false,
		OrigName:     false,
	}).Marshal(&buf, msg)
	return buf.Bytes(), err
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *Quaternion) UnmarshalJSON(b []byte) error {
	return jsonpb.Unmarshal(bytes.NewReader(b), msg)
}
