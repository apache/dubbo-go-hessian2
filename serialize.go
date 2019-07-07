package hessian

type serializer interface {
	serializeObject(*Encoder, POJO) error
	deserializeObject(*Decoder) (interface{}, error)
}

var CodecMap = make(map[string]serializer, 16)

func SetCodec(key string, codec serializer) {
	CodecMap[key] = codec
}

func GetCodec(key string) (serializer, bool) {
	codec, ok := CodecMap[key]
	return codec, ok
}
