// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package db

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

func TestCreateDBClient(t *testing.T) {
	t.Run("Successfully create DB client", func(t *testing.T) {
		expected := &MongoStore{}

		ctx := context.Background()
		err := createDBClient(ctx, "mongo", "testdb")
		if err != nil {
			t.Fatalf("CreateDBClient returned an error (%s)", err)
		}
		if reflect.TypeOf(DBconn) != reflect.TypeOf(expected) {
			t.Fatalf("CreateDBClient set DBconn as:\n result=%T\n expected=%T", DBconn, expected)
		}
	})
	t.Run("Fail to create client for unsupported DB", func(t *testing.T) {
		ctx := context.Background()
		err := createDBClient(ctx, "fakeDB", "testdb2")
		if err == nil {
			t.Fatal("CreateDBClient didn't return an error")
		}
		if !strings.Contains(string(err.Error()), "DB not supported") {
			t.Fatalf("CreateDBClient method returned an error (%s)", err)
		}
	})
}

func TestSerialize(t *testing.T) {

	inp := map[string]interface{}{
		"UUID":   "123e4567-e89b-12d3-a456-426655440000",
		"Data":   "sdaijsdiodalkfjsdlagf",
		"Number": 23,
		"Float":  34.4,
		"Map": map[string]interface{}{
			"m1": "m1",
			"m2": 2,
			"m3": 3.0,
		},
	}

	got, err := Serialize(inp)
	if err != nil {
		t.Fatal(err)
	}

	expected := "{\"Data\":\"sdaijsdiodalkfjsdlagf\"," +
		"\"Float\":34.4,\"Map\":{\"m1\":\"m1\",\"m2\":2,\"m3\":3}," +
		"\"Number\":23,\"UUID\":\"123e4567-e89b-12d3-a456-426655440000\"}"

	if expected != got {
		t.Errorf("Serialize returned unexpected string: %s;"+
			" expected %sv", got, expected)
	}
}

func TestDeSerialize(t *testing.T) {
	testCases := []struct {
		label    string
		input    string
		expected map[string]interface{}
		errMsg   string
	}{
		{
			label: "Sucessful deserialize entry",
			input: "{\"Data\":\"sdaijsdiodalkfjsdlagf\"," +
				"\"Float\":34.4,\"Map\":{\"m1\":\"m1\",\"m3\":3}," +
				"\"UUID\":\"123e4567-e89b-12d3-a456-426655440000\"}",
			expected: map[string]interface{}{
				"UUID":  "123e4567-e89b-12d3-a456-426655440000",
				"Data":  "sdaijsdiodalkfjsdlagf",
				"Float": 34.4,
				"Map": map[string]interface{}{
					"m1": "m1",
					"m3": 3.0,
				},
			},
		},
		{
			label:  "Fail to deserialize invalid entry",
			input:  "{invalid}",
			errMsg: "Error deSerializing {invalid}: invalid character 'i' looking for beginning of object key string",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			got := make(map[string]interface{})
			err := DeSerialize(testCase.input, &got)
			if err != nil {
				if testCase.errMsg == "" {
					t.Fatalf("DeSerialize method return an un-expected (%s)", err)
				}
				if !strings.Contains(string(err.Error()), testCase.errMsg) {
					t.Fatalf("DeSerialize method returned an error (%s)", err)
				}
			} else {
				if !reflect.DeepEqual(testCase.expected, got) {
					t.Errorf("Serialize returned unexpected : %v;"+
						" expected %v", got, testCase.expected)
				}
			}
		})
	}
}
