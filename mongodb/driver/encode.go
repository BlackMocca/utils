package driver

import (
	"reflect"

	helperModel "github.com/Blackmocca/utils"
	frsUUID "github.com/gofrs/uuid"
	googleUUID "github.com/google/uuid"
	"github.com/guregu/null/zero"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

func DateEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != reflect.TypeOf(helperModel.Date{}) {
		return bsoncodec.ValueEncoderError{Name: "DateTimeEncodeValue", Types: []reflect.Type{reflect.TypeOf(helperModel.Date{})}, Received: val}
	}
	dataTi := val.Interface().(helperModel.Date)

	return vw.WriteString(dataTi.String())
}

func TimestampEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != reflect.TypeOf(helperModel.Timestamp{}) {
		return bsoncodec.ValueEncoderError{Name: "DateTimeEncodeValue", Types: []reflect.Type{reflect.TypeOf(helperModel.Timestamp{})}, Received: val}
	}
	dataTi := val.Interface().(helperModel.Timestamp)

	return vw.WriteString(dataTi.String())
}

func GoogleUUIDEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != reflect.TypeOf(googleUUID.UUID{}) {
		return bsoncodec.ValueEncoderError{Name: "GoogleUUIDEncodeValue", Types: []reflect.Type{reflect.TypeOf(googleUUID.UUID{})}, Received: val}
	}
	b := val.Interface().(googleUUID.UUID)
	return vw.WriteBinaryWithSubtype(b[:], uuidSubtype)
}

func FrsUUIDEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != reflect.TypeOf(frsUUID.UUID{}) {
		return bsoncodec.ValueEncoderError{Name: "uuidEncodeValue", Types: []reflect.Type{reflect.TypeOf(frsUUID.UUID{})}, Received: val}
	}

	b := val.Interface().(frsUUID.UUID)
	return vw.WriteBinaryWithSubtype(b.Bytes(), uuidSubtype)
}

func ZeroStringEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != reflect.TypeOf(zero.String{}) {
		return bsoncodec.ValueEncoderError{Name: "zeroStringEncodeValue", Types: []reflect.Type{reflect.TypeOf(zero.String{})}, Received: val}
	}

	b := val.Interface()
	return vw.WriteString(b.(zero.String).ValueOrZero())
}

func ZeroIntEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	ptr := zero.Int{}
	if !val.IsValid() || val.Type() != reflect.TypeOf(ptr) {
		return bsoncodec.ValueEncoderError{Name: "zeroIntEncodeValue", Types: []reflect.Type{reflect.TypeOf(ptr)}, Received: val}
	}
	b := val.Interface()

	_, err := cast.ToInt32E(b.(zero.Int).ValueOrZero())
	if err == nil {
		return vw.WriteInt32(cast.ToInt32(b.(zero.Int).ValueOrZero()))
	}

	return vw.WriteInt64(b.(zero.Int).ValueOrZero())
}

func ZeroFloatEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	ptr := zero.Float{}
	if !val.IsValid() || val.Type() != reflect.TypeOf(ptr) {
		return bsoncodec.ValueEncoderError{Name: "zeroFloatEncodeValue", Types: []reflect.Type{reflect.TypeOf(ptr)}, Received: val}
	}

	b := val.Interface()
	return vw.WriteDouble(b.(zero.Float).ValueOrZero())
}

func ZeroBoolEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	ptr := zero.Bool{}
	if !val.IsValid() || val.Type() != reflect.TypeOf(ptr) {
		return bsoncodec.ValueEncoderError{Name: "zeroBoolEncodeValue", Types: []reflect.Type{reflect.TypeOf(ptr)}, Received: val}
	}

	b := val.Interface()
	return vw.WriteBoolean(b.(zero.Bool).ValueOrZero())
}
