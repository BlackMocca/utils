package driver

import (
	"fmt"
	"reflect"
	"time"

	helperModel "github.com/Blackmocca/utils"
	frsUUID "github.com/gofrs/uuid"
	googleUUID "github.com/google/uuid"
	"github.com/guregu/null/zero"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

func DateDecodeValue(ec bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != reflect.TypeOf(helperModel.Date{}) {
		return bsoncodec.ValueDecoderError{Name: "DateTimeEncodeValue", Types: []reflect.Type{reflect.TypeOf(helperModel.Date{})}, Received: val}
	}

	switch vrType := vr.Type(); vrType {
	case bsontype.Null:
		break
	case bsontype.Undefined:
		break
	case bsontype.String:
		str, err := vr.ReadString()
		if err != nil {
			return err
		}
		val.Set(reflect.ValueOf(helperModel.NewDateFromString(str)))
	case bsontype.DateTime:
		dt, err := vr.ReadDateTime()
		if err != nil {
			return err
		}
		ti := time.Unix(dt/1000, 0)
		val.Set(reflect.ValueOf(helperModel.NewDateFromTime(ti)))
	default:
		return fmt.Errorf("cannot decode %v into a Timestamp", vrType)
	}

	return nil
}

func TimestampDecodeValue(ec bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != reflect.TypeOf(helperModel.Timestamp{}) {
		return bsoncodec.ValueDecoderError{Name: "DateTimeEncodeValue", Types: []reflect.Type{reflect.TypeOf(helperModel.Timestamp{})}, Received: val}
	}

	switch vrType := vr.Type(); vrType {
	case bsontype.Null:
		break
	case bsontype.Undefined:
		break
	case bsontype.String:
		str, err := vr.ReadString()
		if err != nil {
			return err
		}
		val.Set(reflect.ValueOf(helperModel.NewTimestampFromString(str)))
	case bsontype.DateTime:
		dt, err := vr.ReadDateTime()
		if err != nil {
			return err
		}
		ti := time.Unix(dt/1000, 0)
		val.Set(reflect.ValueOf(helperModel.NewTimestampFromTime(ti)))
	default:
		return fmt.Errorf("cannot decode %v into a Timestamp", vrType)
	}

	return nil
}

func GoogleUUIDDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != reflect.TypeOf(googleUUID.UUID{}) {
		return bsoncodec.ValueDecoderError{Name: "GoogleUUIDDecodeValue", Types: []reflect.Type{reflect.TypeOf(googleUUID.UUID{})}, Received: val}
	}

	var data []byte
	var subtype byte
	var err error
	switch vrType := vr.Type(); vrType {
	case bsontype.Binary:
		data, subtype, err = vr.ReadBinary()
		if subtype != uuidSubtype {
			return fmt.Errorf("unsupported binary subtype %v for UUID", subtype)
		}
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return fmt.Errorf("cannot decode %v into a UUID", vrType)
	}

	if err != nil {
		return err
	}
	uuid2, err := googleUUID.FromBytes(data)
	if err != nil {
		return err
	}
	val.Set(reflect.ValueOf(uuid2))
	return nil
}

func FrsUUIDDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != reflect.TypeOf(frsUUID.UUID{}) {
		return bsoncodec.ValueDecoderError{Name: "uuidDecodeValue", Types: []reflect.Type{reflect.TypeOf(frsUUID.UUID{})}, Received: val}
	}

	switch vrType := vr.Type(); vrType {
	case bsontype.Binary:
		data, subtype, err := vr.ReadBinary()
		if err != nil {
			return err
		}
		if subtype != uuidSubtype {
			return fmt.Errorf("unsupported binary subtype %v for UUID", subtype)
		}

		uuid2, err := frsUUID.FromBytes(data)
		if err != nil {
			return err
		}
		val.Set(reflect.ValueOf(uuid2))

	case bsontype.String:
		str, err := vr.ReadString()
		if err != nil {
			return err
		}
		uuid2, err := frsUUID.FromString(str)
		if err != nil {
			return err
		}

		val.Set(reflect.ValueOf(uuid2))
	case bsontype.Null:
		return nil
	case bsontype.Undefined:
		return nil
	default:
		return fmt.Errorf("cannot decode %v into a UUID", vrType)
	}

	return nil
}

func ZeroStringDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != reflect.TypeOf(zero.String{}) {
		return bsoncodec.ValueDecoderError{Name: "zeroStringDecodeValue", Types: []reflect.Type{reflect.TypeOf(zero.String{})}, Received: val}
	}

	switch vrType := vr.Type(); vrType {
	case bsontype.String:
		str, err := vr.ReadString()
		if err != nil {
			return err
		}

		val.Set(reflect.ValueOf(zero.StringFrom(str)))
	case bsontype.Null:
		return nil
	case bsontype.Undefined:
		return nil
	default:
		return fmt.Errorf("cannot decode %v into a zeroString", vrType)
	}

	return nil
}

func ZeroIntDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	ptr := zero.Int{}
	if !val.CanSet() || val.Type() != reflect.TypeOf(ptr) {
		return bsoncodec.ValueDecoderError{Name: "zeroIntDecodeValue", Types: []reflect.Type{reflect.TypeOf(ptr)}, Received: val}
	}

	switch vrType := vr.Type(); vrType {
	case bsontype.Int32:
		data, err := vr.ReadInt32()
		if err != nil {
			return err
		}
		val.Set(reflect.ValueOf(zero.IntFrom(cast.ToInt64(data))))
	case bsontype.Int64:
		data, err := vr.ReadInt64()
		if err != nil {
			return err
		}
		val.Set(reflect.ValueOf(zero.IntFrom(data)))
	case bsontype.Null:
		return nil
	case bsontype.Undefined:
		return nil
	default:
		return fmt.Errorf("cannot decode %v into a zeroInt", vrType)
	}

	return nil
}

func ZeroFloatDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	ptr := zero.Float{}
	if !val.CanSet() || val.Type() != reflect.TypeOf(ptr) {
		return bsoncodec.ValueDecoderError{Name: "zeroFloatDecodeValue", Types: []reflect.Type{reflect.TypeOf(ptr)}, Received: val}
	}

	switch vrType := vr.Type(); vrType {
	case bsontype.Double:
		data, err := vr.ReadDouble()
		if err != nil {
			return err
		}
		val.Set(reflect.ValueOf(zero.FloatFrom(data)))
	case bsontype.Null:
		return nil
	case bsontype.Undefined:
		return nil
	default:
		return fmt.Errorf("cannot decode %v into a zeroFloat", vrType)
	}

	return nil
}

func ZeroBoolDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	ptr := zero.Bool{}
	if !val.CanSet() || val.Type() != reflect.TypeOf(ptr) {
		return bsoncodec.ValueDecoderError{Name: "zeroBoolDecodeValue", Types: []reflect.Type{reflect.TypeOf(ptr)}, Received: val}
	}

	switch vrType := vr.Type(); vrType {
	case bsontype.Boolean:
		data, err := vr.ReadBoolean()
		if err != nil {
			return err
		}
		val.Set(reflect.ValueOf(zero.BoolFrom(data)))
	case bsontype.Null:
		return nil
	case bsontype.Undefined:
		return nil
	default:
		return fmt.Errorf("cannot decode %v into a zeroBool", vrType)
	}

	return nil
}
