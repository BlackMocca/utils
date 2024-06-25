package mongo

import (
	"reflect"

	helperModel "github.com/Blackmocca/utils"
	"github.com/Blackmocca/utils/mongodb/driver"
	frsUUID "github.com/gofrs/uuid"
	googleUUID "github.com/google/uuid"
	"github.com/guregu/null/zero"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

/*
สำหรับ package mongodb-driver
วิธีใช้ คือ
mongo.Connect(ctx, options.Client().ApplyURI(uri).SetRetryReads(true).SetRetryWrites(true).SetRegistry(builder.Build()))
*/
func Registry() *bsoncodec.Registry {
	builder := bson.NewRegistry()

	/* models */
	date := helperModel.Date{}
	timestamp := helperModel.Timestamp{}
	googleUUID := googleUUID.UUID{}
	frsUUID := frsUUID.UUID{}
	zs := zero.String{}
	zi := zero.Int{}
	zb := zero.Bool{}
	zf := zero.Float{}

	/* encoding */
	builder.RegisterTypeEncoder(reflect.TypeOf(date), bsoncodec.ValueEncoderFunc(driver.DateEncodeValue))
	builder.RegisterTypeEncoder(reflect.TypeOf(timestamp), bsoncodec.ValueEncoderFunc(driver.TimestampEncodeValue))
	builder.RegisterTypeEncoder(reflect.TypeOf(googleUUID), bsoncodec.ValueEncoderFunc(driver.GoogleUUIDEncodeValue))
	builder.RegisterTypeEncoder(reflect.TypeOf(frsUUID), bsoncodec.ValueEncoderFunc(driver.FrsUUIDEncodeValue))
	builder.RegisterTypeEncoder(reflect.TypeOf(zs), bsoncodec.ValueEncoderFunc(driver.ZeroStringEncodeValue))
	builder.RegisterTypeEncoder(reflect.TypeOf(zi), bsoncodec.ValueEncoderFunc(driver.ZeroIntEncodeValue))
	builder.RegisterTypeEncoder(reflect.TypeOf(zb), bsoncodec.ValueEncoderFunc(driver.ZeroBoolEncodeValue))
	builder.RegisterTypeEncoder(reflect.TypeOf(zf), bsoncodec.ValueEncoderFunc(driver.ZeroFloatEncodeValue))

	/* decoding */
	builder.RegisterTypeDecoder(reflect.TypeOf(date), bsoncodec.ValueDecoderFunc(driver.DateDecodeValue))
	builder.RegisterTypeDecoder(reflect.TypeOf(timestamp), bsoncodec.ValueDecoderFunc(driver.TimestampDecodeValue))
	builder.RegisterTypeDecoder(reflect.TypeOf(googleUUID), bsoncodec.ValueDecoderFunc(driver.GoogleUUIDDecodeValue))
	builder.RegisterTypeDecoder(reflect.TypeOf(frsUUID), bsoncodec.ValueDecoderFunc(driver.FrsUUIDDecodeValue))
	builder.RegisterTypeDecoder(reflect.TypeOf(zs), bsoncodec.ValueDecoderFunc(driver.ZeroStringDecodeValue))
	builder.RegisterTypeDecoder(reflect.TypeOf(zi), bsoncodec.ValueDecoderFunc(driver.ZeroIntDecodeValue))
	builder.RegisterTypeDecoder(reflect.TypeOf(zb), bsoncodec.ValueDecoderFunc(driver.ZeroBoolDecodeValue))
	builder.RegisterTypeDecoder(reflect.TypeOf(zf), bsoncodec.ValueDecoderFunc(driver.ZeroFloatDecodeValue))

	return builder
}
