package json

import (
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"testing"

	//"github.com/stretchr/testify/assert"
)


func BenchmarkDecodeWithSimdJson(b *testing.B) {
	Decoder = SimdJsonDecoder
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := NewNull()
		if err := val.Decode(testData); err != nil {
			b.Fatal(err)
		}
		_ = val.Encode(nil)
		val.DecRef()
	}
}

func BenchmarkDecodeWithSchemaJson(b *testing.B) {
	Decoder = SimdJsonDecoder

	schema := &DecodeSchema{
		Type: Object,
		Property: map[string]*DecodeSchema{
			"data": {
				Type: Array,
				ElementSchema: &DecodeSchema{
					Type: Object,
					Property: map[string]*DecodeSchema{
						"base": {
							Type: Object,
							Property: map[string]*DecodeSchema{
								"album_audio_id": {
									Type: Int,
								},
							},
						},
					},
					DefaultAction: SaveToSource,
				},
			},
		},
		DefaultAction: SaveToSource,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := NewNull()
		if err := val.DecodeWithSchema(testData, schema); err != nil {
			b.Fatal(err)
		}
		_ = val.Encode(nil)
		val.DecRef()
	}
}

func BenchmarkStdUnmarshal(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var res map[string]interface{}
		if err := json.Unmarshal(testData, &res); err != nil {
			b.Fatal(err)
		}
		//if _, err := json.Marshal(res); err != nil {
		//	b.Fatal(err)
		//}
	}
}


func BenchmarkIteratorUnmarshal(b *testing.B) {
	var json1 = jsoniter.ConfigCompatibleWithStandardLibrary
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var res map[string]interface{}
		if err := json1.Unmarshal(testData, &res); err != nil {
			b.Fatal(err)
		}
		if _, err := json.Marshal(res); err != nil {
			b.Fatal(err)
		}
	}
}
