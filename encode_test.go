package json

import (
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"testing"

	//"github.com/stretchr/testify/assert"
)

var testData []byte

func TestMain(m *testing.M) {
	//data, err := ioutil.ReadFile("testdata/data.json")
	data, err := ioutil.ReadFile("testdata/response.json")
	if err != nil {
		log.Fatal(err)
	}
	testData = data

	m.Run()
}


func BenchmarkEncodeWithSimdJson(b *testing.B) {
	Decoder = SimdJsonDecoder
	val := NewNull()
	if err := val.Decode(testData); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = val.Encode(nil)
	}
	val.DecRef()
}

func BenchmarkEncodeWithSchemaJson(b *testing.B) {
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
	val := NewNull()
	if err := val.DecodeWithSchema(testData, schema); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = val.Encode(nil)
	}
	val.DecRef()
}


func BenchmarkStdMarshal(b *testing.B) {
	var res map[string]interface{}
	if err := json.Unmarshal(testData, &res); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(res); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkIteratorMarshal(b *testing.B) {
	var json1 = jsoniter.ConfigCompatibleWithStandardLibrary
	var res map[string]interface{}
	if err := json1.Unmarshal(testData, &res); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(res); err != nil {
			b.Fatal(err)
		}
	}
}