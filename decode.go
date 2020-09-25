package json

type DecoderFunc func(v Value, data []byte) error
type DecoderWithSchemaFunc func(v Value, data []byte, schema *DecodeSchema) error

var Decoder DecoderFunc
var DecoderWithSchema DecoderWithSchemaFunc

func (v *valueImpl)Decode(str []byte) error {
    return Decoder(v, str)
}

func (v *valueImpl)DecodeWithSchema(str []byte, schema *DecodeSchema) error {
    return DecoderWithSchema(v, str, schema)
}

func init() {
    //Decoder = JsonParserDecoder
    Decoder = SimdJsonDecoder
    DecoderWithSchema = SchemaJsonDecoder
}
