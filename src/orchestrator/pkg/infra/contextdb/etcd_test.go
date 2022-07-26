// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package contextdb

import (
	"context"

	pkgerrors "github.com/pkg/errors"

	"strings"
	"testing"

	mvccpb "go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type kv struct {
	Key   []byte
	Value []byte
}

// MockEtcdClient for mocking etcd
type MockEtcdClient struct {
	Kvs   []*mvccpb.KeyValue
	Count int64
	Err   error
}

// Mocking only Single Value
// Put function
func (e *MockEtcdClient) Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	var m mvccpb.KeyValue
	m.Key = []byte(key)
	m.Value = []byte(val)
	e.Count = e.Count + 1
	e.Kvs = append(e.Kvs, &m)
	return &clientv3.PutResponse{}, e.Err
}

// Get function
func (e *MockEtcdClient) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	var g clientv3.GetResponse
	g.Kvs = e.Kvs
	g.Count = e.Count
	return &g, e.Err
}

// Delete function
func (e *MockEtcdClient) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	return &clientv3.DeleteResponse{}, e.Err
}

type testStruct struct {
	Name string `json:"name"`
	Num  int    `json:"num"`
}

// TestPut test Put
func TestPut(t *testing.T) {
	testCases := []struct {
		label         string
		mockEtcd      *MockEtcdClient
		expectedError string
		key           string
		value         *testStruct
	}{
		{
			label:    "Success Case",
			mockEtcd: &MockEtcdClient{},
			key:      "test1",
			value:    &testStruct{Name: "test", Num: 5},
		},
		{
			label:         "Key is null",
			mockEtcd:      &MockEtcdClient{},
			key:           "",
			expectedError: "Key is null",
		},
		{
			label:         "Value is nil",
			mockEtcd:      &MockEtcdClient{},
			key:           "test1",
			value:         nil,
			expectedError: "Value is nil",
		},
		{
			label:         "Error creating etcd entry",
			mockEtcd:      &MockEtcdClient{Err: pkgerrors.New("DB Error")},
			key:           "test1",
			value:         &testStruct{Name: "test", Num: 5},
			expectedError: "Error creating etcd entry: DB Error",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			ctx := context.Background()
			cli, _ := NewEtcdClient(&clientv3.Client{}, EtcdConfig{})
			getEtcd = func(e *EtcdClient) Etcd {
				return testCase.mockEtcd
			}
			err := cli.Put(ctx, testCase.key, testCase.value)
			if err != nil {
				if testCase.expectedError == "" {
					t.Fatalf("Method returned an un-expected (%s)", err)
				}
				if !strings.Contains(string(err.Error()), testCase.expectedError) {
					t.Fatalf("Method returned an error (%s)", err)
				}
			}

		})
	}
}

func TestGet(t *testing.T) {
	testCases := []struct {
		label         string
		mockEtcd      *MockEtcdClient
		expectedError string
		key           string
		value         *testStruct
	}{
		{
			label:         "Key is null",
			mockEtcd:      &MockEtcdClient{},
			key:           "",
			value:         nil,
			expectedError: "Key is null",
		},
		{
			label:         "Key doesn't exist",
			mockEtcd:      &MockEtcdClient{},
			key:           "test1",
			value:         &testStruct{},
			expectedError: "Key doesn't exist",
		},
		{
			label:         "Error getting etcd entry",
			mockEtcd:      &MockEtcdClient{Err: pkgerrors.New("DB Error")},
			key:           "test1",
			value:         &testStruct{},
			expectedError: "Error getting etcd entry: DB Error",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			ctx := context.Background()
			cli, _ := NewEtcdClient(&clientv3.Client{}, EtcdConfig{})
			getEtcd = func(e *EtcdClient) Etcd {
				return testCase.mockEtcd
			}
			err := cli.Get(ctx, testCase.key, testCase.value)
			if err != nil {
				if testCase.expectedError == "" {
					t.Fatalf("Method returned an un-expected (%s)", err)
				}
				if !strings.Contains(string(err.Error()), testCase.expectedError) {
					t.Fatalf("Method returned an error (%s)", err)
				}
			}

		})
	}
}

func TestGetString(t *testing.T) {
	testCases := []struct {
		label         string
		mockEtcd      *MockEtcdClient
		expectedError string
		value         string
	}{
		{
			label:    "Success Case",
			mockEtcd: &MockEtcdClient{},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			ctx := context.Background()
			cli, _ := NewEtcdClient(&clientv3.Client{}, EtcdConfig{})
			getEtcd = func(e *EtcdClient) Etcd {
				return testCase.mockEtcd
			}
			err := cli.Put(ctx, "test", "test1")
			if err != nil {
				t.Error("Test failed", err)
			}
			var s string
			err = cli.Get(ctx, "test", &s)
			if err != nil {
				t.Error("Test failed", err)
			}
			if "test1" != s {
				t.Error("Get Failed")
			}
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		label         string
		mockEtcd      *MockEtcdClient
		expectedError string
	}{
		{
			label:    "Success Case",
			mockEtcd: &MockEtcdClient{},
		},
		{
			label:         "Delete failed etcd entry",
			mockEtcd:      &MockEtcdClient{Err: pkgerrors.New("DB Error")},
			expectedError: "Delete failed etcd entry: DB Error",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			ctx := context.Background()
			cli, _ := NewEtcdClient(&clientv3.Client{}, EtcdConfig{})
			getEtcd = func(e *EtcdClient) Etcd {
				return testCase.mockEtcd
			}
			err := cli.Delete(ctx, "test")
			if err != nil {
				if testCase.expectedError == "" {
					t.Fatalf("Method returned an un-expected (%s)", err)
				}
				if !strings.Contains(string(err.Error()), testCase.expectedError) {
					t.Fatalf("Method returned an error (%s)", err)
				}
			}

		})
	}
}

func TestGetAll(t *testing.T) {
	testCases := []struct {
		label         string
		mockEtcd      *MockEtcdClient
		expectedError string
	}{
		{
			label:         "Key doesn't exist",
			mockEtcd:      &MockEtcdClient{},
			expectedError: "Key doesn't exist",
		},
		{
			label:         "Error getting etcd entry",
			mockEtcd:      &MockEtcdClient{Err: pkgerrors.New("DB Error")},
			expectedError: "Error getting etcd entry: DB Error",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			ctx := context.Background()
			cli, _ := NewEtcdClient(&clientv3.Client{}, EtcdConfig{})
			getEtcd = func(e *EtcdClient) Etcd {
				return testCase.mockEtcd
			}
			_, err := cli.GetAllKeys(ctx, "test")
			if err != nil {
				if testCase.expectedError == "" {
					t.Fatalf("Method returned an un-expected (%s)", err)
				}
				if !strings.Contains(string(err.Error()), testCase.expectedError) {
					t.Fatalf("Method returned an error (%s)", err)
				}
			}
		})
	}
}
