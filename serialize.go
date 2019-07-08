package hessian

type Serializer interface {
	serializeObject(*Encoder, POJO) error
}
type DeSerializer interface {
	deserializeObject(*Decoder) (interface{}, error)
}

var SerializerMap = make(map[string]Serializer, 16)
var DeSerializerMap = make(map[string]DeSerializer, 16)

func SetSerializer(key string, codec Serializer) {
	SerializerMap[key] = codec
}

func GetSerializer(key string) (Serializer, bool) {
	codec, ok := SerializerMap[key]
	return codec, ok
}

func SetDeSerializer(key string, codec DeSerializer) {
	DeSerializerMap[key] = codec
}

func GetDeSerializer(key string) (DeSerializer, bool) {
	codec, ok := DeSerializerMap[key]
	return codec, ok
}
