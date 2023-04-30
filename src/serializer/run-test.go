package main

import (
	"errors"
	"fmt"
	"os"
	"time"
)

func runSingleSerializerTest(s Serializer) (time.Duration, time.Duration, int, error) {
	dummy := makeDummy()

	// Serialize
	start := time.Now()
	serialized, err := s.serialize(dummy)
	if err != nil {
		return 0, 0, 0, err
	}
	end := time.Now()
	serializeTime := end.Sub(start)

	// Deserialize
	start = time.Now()
	_, err = s.deserialize(serialized)
	if err != nil {
		return 0, 0, 0, err
	}
	end = time.Now()
	deserializeTime := end.Sub(start)

	return serializeTime, deserializeTime, len(serialized), nil
}

func runSerializerTests(s Serializer, nTests uint64) (string, error) {
	if nTests == 0 {
		return "", errors.New("nTests must be positive")
	}

	var serializeTimeTotal time.Duration = 0
	var deserializeTimeTotal time.Duration = 0
	var serializedBytes int = -1

	for i := 0; i < int(nTests); i++ {
		serializeTime, deserializeTime, bytes, err := runSingleSerializerTest(s)
		if err != nil {
			return "", err
		}

		serializeTimeTotal += serializeTime
		deserializeTimeTotal += deserializeTime
		if serializedBytes == -1 {
			serializedBytes = bytes
		}

		if serializedBytes != bytes {
			return "", errors.New("serialized object size is inconsistent between test runs")
		}
	}

	return fmt.Sprintf("%v – %v – %v – %v\n", s.getFormat(), serializedBytes, serializeTimeTotal/time.Duration(nTests), deserializeTimeTotal/time.Duration(nTests)), nil
}

func calculateResult() string {
	var s Serializer
	switch os.Getenv("S_FORMAT") {
	case "gob":
		s = GobSerializer{}
	case "xml":
		s = XmlSerializer{}
	case "json":
		s = JsonSerializer{}
	case "protobuf":
		s = ProtobufSerializer{}
	case "avro":
		s = AvroSerializer{}
	case "yaml":
		s = YamlSerializer{}
	case "message-pack":
		s = MessagePackSerializer{}
	default:
		panic(fmt.Sprintf("invalid S_FORMAT: %q", os.Getenv("S_FORMAT")))
	}

	result, err := runSerializerTests(s, 1000)
	if err != nil {
		return fmt.Sprintf("unexpected error: %q\n", err.Error())
	}
	return result
}
