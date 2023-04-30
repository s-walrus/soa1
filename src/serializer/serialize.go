package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"

	"github.com/s-walrus/soa1/gen"

	"github.com/linkedin/goavro"
	"github.com/vmihailenco/msgpack"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

type Serializer interface {
	serialize(d Dummy) ([]byte, error)
	deserialize(data []byte) (Dummy, error)
	getFormat() string
}

// GOB

type GobSerializer struct{}

func (GobSerializer) serialize(d Dummy) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(d)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (GobSerializer) deserialize(data []byte) (Dummy, error) {
	var d Dummy
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&d)
	if err != nil {
		return Dummy{}, err
	}
	return d, nil
}

func (GobSerializer) getFormat() string {
	return "Gob"
}

// XML

type XmlSerializer struct{}

func (XmlSerializer) serialize(d Dummy) ([]byte, error) {
	return xml.Marshal(d)
}

func (XmlSerializer) deserialize(data []byte) (Dummy, error) {
	var d Dummy
	err := xml.Unmarshal(data, &d)
	if err != nil {
		return Dummy{}, err
	}
	return d, nil
}

func (XmlSerializer) getFormat() string {
	return "XML"
}

// JSON

type JsonSerializer struct{}

func (JsonSerializer) serialize(d Dummy) ([]byte, error) {
	return json.Marshal(d)
}

func (JsonSerializer) deserialize(data []byte) (Dummy, error) {
	var d Dummy
	err := json.Unmarshal(data, &d)
	if err != nil {
		return Dummy{}, err
	}
	return d, nil
}

func (JsonSerializer) getFormat() string {
	return "JSON"
}

// Protobuf

type ProtobufSerializer struct{}

func (ProtobufSerializer) serialize(d Dummy) ([]byte, error) {
	pb := &gen.Dummy{
		Name: d.Name,
		Age:  int32(d.Age),
	}
	return proto.Marshal(pb)
}

func (ProtobufSerializer) deserialize(data []byte) (Dummy, error) {
	pb := &gen.Dummy{}
	err := proto.Unmarshal(data, pb)
	if err != nil {
		return Dummy{}, err
	}
	return Dummy{
		Name: pb.GetName(),
		Age:  int(pb.GetAge()),
	}, nil
}

func (ProtobufSerializer) getFormat() string {
	return "Protobuf"
}

// Avro

type AvroSerializer struct{}

func (AvroSerializer) getSchema() string {
	return `
        {
            "type": "record",
            "name": "Dummy",
            "fields": [
                { "name": "name", "type": "string" },
                { "name": "limbs", "type": { "type": "array", "items": "string" } },
                { "name": "stats", "type": { "type": "map", "values": "int" } },
                { "name": "age", "type": "int" },
                { "name": "mood", "type": "float" }
            ]
        }
    `
}

func (s AvroSerializer) serialize(d Dummy) ([]byte, error) {
	native := map[string]interface{}{
		"name":  d.Name,
		"limbs": d.Limbs,
		"stats": d.Stats,
		"age":   d.Age,
		"mood":  d.Mood,
	}
	codec, err := goavro.NewCodec(s.getSchema())
	if err != nil {
		return nil, err
	}
	textual, err := codec.TextualFromNative(nil, native)
	if err != nil {
		return nil, err
	}
	return textual, nil
}

func (s AvroSerializer) deserialize(data []byte) (Dummy, error) {
	codec, err := goavro.NewCodec(s.getSchema())
	if err != nil {
		return Dummy{}, err
	}
	avroData, _, err := codec.NativeFromTextual(data)
	if err != nil {
		return Dummy{}, err
	}
	native := avroData.(map[string]interface{})
	d := makeDummy()
	d.Name = native["name"].(string)
	d.Age = int(native["age"].(int32))
	d.Mood = native["mood"].(float64)
	for _, limb := range native["limbs"].([]interface{}) {
		d.Limbs = append(d.Limbs, limb.(string))
	}
	for stat, value := range native["stats"].(map[string]interface{}) {
		d.Stats[stat] = int(value.(int32))
	}
	return d, nil
}

func (AvroSerializer) getFormat() string {
	return "Avro"
}

// YAML

type YamlSerializer struct{}

func (YamlSerializer) serialize(d Dummy) ([]byte, error) {
	return yaml.Marshal(d)
}

func (YamlSerializer) deserialize(data []byte) (Dummy, error) {
	var d Dummy
	err := yaml.Unmarshal(data, &d)
	if err != nil {
		return Dummy{}, err
	}
	return d, nil
}

func (YamlSerializer) getFormat() string {
	return "YAML"
}

// Message pack

type MessagePackSerializer struct{}

func (MessagePackSerializer) serialize(d Dummy) ([]byte, error) {
	return msgpack.Marshal(d)
}

func (MessagePackSerializer) deserialize(data []byte) (Dummy, error) {
	var d Dummy
	err := msgpack.Unmarshal(data, &d)
	if err != nil {
		return Dummy{}, err
	}
	return d, nil
}

func (MessagePackSerializer) getFormat() string {
	return "MessagePack"
}
