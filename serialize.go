package hessian

type Serializer interface {
	serializeObject(*Encoder, POJO) error
	deserializeObject(*Decoder) (interface{}, error)
}

var CodecMap = make(map[string]Serializer, 16)

func SetSerializer(key string, codec Serializer) {
	CodecMap[key] = codec
}

func GetSerializer(key string) (Serializer, bool) {
	codec, ok := CodecMap[key]
	return codec, ok
}
