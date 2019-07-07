package hessian

type Serializer interface {
	serializeObject(*Encoder, POJO) error
	deserializeObject(*Decoder) (interface{}, error)
}

var CodecMap = make(map[string]Serializer, 16)

func SetCodec(key string, codec Serializer) {
	CodecMap[key] = codec
}

func GetCodec(key string) (Serializer, bool) {
	codec, ok := CodecMap[key]
	return codec, ok
}
