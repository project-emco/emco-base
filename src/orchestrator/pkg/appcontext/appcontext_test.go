// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package appcontext

import (
	"context"
	"fmt"
	"strings"
	"testing"

	pkgerrors "github.com/pkg/errors"
)

// Mock run time context
type MockRunTimeContext struct {
	Items map[string]interface{}
	Err   error
}

type MockCompositeAppMeta struct {
	Project      string
	CompositeApp string
	Version      string
	Release      string
}

func (c *MockRunTimeContext) RtcCreate(ctx context.Context) (interface{}, error) {
	var key string = "/context/9345674458787728/"

	if c.Items == nil {
		c.Items = make(map[string]interface{})
	}
	c.Items[key] = "9345674458787728"
	return interface{}(key), c.Err

}

func (c *MockRunTimeContext) RtcAddMeta(ctx context.Context, meta interface{}) error {
	var cid string = "/context/9345674458787728/"
	key := cid + "meta" + "/"
	if c.Items == nil {
		c.Items = make(map[string]interface{})
	}
	c.Items[key] = meta
	return nil
}

func (c *MockRunTimeContext) RtcInit() (interface{}, error) {
	var id string = "9345674458787728"
	return id, c.Err
}

func (c *MockRunTimeContext) RtcInitWithValue(i interface{}) (interface{}, error) {
	var id string = fmt.Sprintf("%v", i)
	return id, c.Err
}

func (c *MockRunTimeContext) RtcLoad(ctx context.Context, id interface{}) (interface{}, error) {
	str := "/context/" + fmt.Sprintf("%v", id) + "/"
	return interface{}(str), c.Err
}

func (c *MockRunTimeContext) RtcGet(ctx context.Context) (interface{}, error) {
	var key string = "/context/9345674458787728/"
	return key, c.Err
}

func (c *MockRunTimeContext) RtcGetMeta(ctx context.Context) (interface{}, error) {
	meta := CompositeAppMeta{Project: "pn", CompositeApp: "ca", Version: "v", Release: "rName"}
	return meta, nil
}

func (c *MockRunTimeContext) RtcAddLevel(ctx context.Context, handle interface{}, level string, value string) (interface{}, error) {
	str := fmt.Sprintf("%v", handle) + level + "/" + value + "/"
	c.Items[str] = value
	return nil, c.Err

}

func (c *MockRunTimeContext) RtcAddOneLevel(ctx context.Context, handle interface{}, level string, value interface{}) (interface{}, error) {
	str := fmt.Sprintf("%v", handle) + level + "/"
	c.Items[str] = value
	return nil, c.Err

}

func (c *MockRunTimeContext) RtcAddResource(ctx context.Context, handle interface{}, resname string, value interface{}) (interface{}, error) {
	str := fmt.Sprintf("%v", handle) + "resource" + "/" + resname + "/"
	c.Items[str] = value
	return nil, c.Err

}

func (c *MockRunTimeContext) RtcAddInstruction(ctx context.Context, handle interface{}, level string, insttype string, value interface{}) (interface{}, error) {
	str := fmt.Sprintf("%v", handle) + level + "/" + insttype + "/"
	c.Items[str] = value
	return nil, c.Err
}

func (c *MockRunTimeContext) RtcDeletePair(ctx context.Context, handle interface{}) error {
	str := fmt.Sprintf("%v", handle)
	delete(c.Items, str)
	return c.Err
}

func (c *MockRunTimeContext) RtcDeletePrefix(ctx context.Context, handle interface{}) error {
	for k := range c.Items {
		delete(c.Items, k)
	}
	return c.Err
}

func (c *MockRunTimeContext) RtcGetHandles(ctx context.Context, handle interface{}) ([]interface{}, error) {
	var keys []interface{}

	for k := range c.Items {
		keys = append(keys, string(k))
	}
	return keys, c.Err
}

func (c *MockRunTimeContext) RtcGetValue(ctx context.Context, handle interface{}, value interface{}) error {
	key := fmt.Sprintf("%v", handle)
	var s *string
	s = value.(*string)
	for kvKey, kvValue := range c.Items {
		if kvKey == key {
			*s = kvValue.(string)
			return c.Err
		}
	}
	return c.Err
}

func (c *MockRunTimeContext) RtcUpdateValue(ctx context.Context, handle interface{}, value interface{}) error {
	key := fmt.Sprintf("%v", handle)
	c.Items[key] = value
	return c.Err
}

func TestCreateCompositeApp(t *testing.T) {
	var ac = AppContext{}
	testCases := []struct {
		label         string
		mockRtcontext *MockRunTimeContext
		expectedError string
		meta          interface{}
	}{
		{
			label:         "Success case",
			mockRtcontext: &MockRunTimeContext{},
			meta:          interface{}(MockCompositeAppMeta{Project: "Testproject", CompositeApp: "TestCompApp", Version: "CompAppVersion", Release: "TestRelease"}),
		},
		{
			label:         "Create returns error case",
			mockRtcontext: &MockRunTimeContext{Err: pkgerrors.Errorf("Error creating run time context:")},
			expectedError: "Error creating run time context:",
			meta:          interface{}(MockCompositeAppMeta{Project: "Testproject", CompositeApp: "TestCompApp", Version: "CompAppVersion", Release: "TestRelease"}),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			ctx := context.Background()
			ac.rtc = testCase.mockRtcontext
			_, err := ac.CreateCompositeApp(ctx)
			if err != nil {
				if !strings.Contains(string(err.Error()), testCase.expectedError) {
					t.Fatalf("Method returned an error (%s)", err)
				}
			}

		})
	}
}

func TestGetCompositeApp(t *testing.T) {
	var ac = AppContext{}
	testCases := []struct {
		label         string
		mockRtcontext *MockRunTimeContext
		expectedError string
	}{
		{
			label:         "Success case",
			mockRtcontext: &MockRunTimeContext{},
		},
		{
			label:         "Get returns error case",
			mockRtcontext: &MockRunTimeContext{Err: pkgerrors.Errorf("Error getting run time context:")},
			expectedError: "Error getting run time context:",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			ctx := context.Background()
			ac.rtc = testCase.mockRtcontext
			_, err := ac.GetCompositeAppHandle(ctx)
			if err != nil {
				if !strings.Contains(string(err.Error()), testCase.expectedError) {
					t.Fatalf("Method returned an error (%s)", err)
				}
			}

		})
	}
}

func TestDeleteCompositeApp(t *testing.T) {
	var ac = AppContext{}
	testCases := []struct {
		label         string
		mockRtcontext *MockRunTimeContext
		expectedError string
	}{
		{
			label:         "Success case",
			mockRtcontext: &MockRunTimeContext{},
		},
		{
			label:         "Delete returns error case",
			mockRtcontext: &MockRunTimeContext{Err: pkgerrors.Errorf("Error deleting run time context:")},
			expectedError: "Error deleting run time context:",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			ctx := context.Background()
			ac.rtc = testCase.mockRtcontext
			err := ac.DeleteCompositeApp(ctx)
			if err != nil {
				if !strings.Contains(string(err.Error()), testCase.expectedError) {
					t.Fatalf("Method returned an error (%s)", err)
				}
			}

		})
	}
}

func TestAddApp(t *testing.T) {
	var ac = AppContext{}
	testCases := []struct {
		label         string
		mockRtcontext *MockRunTimeContext
		key           interface{}
		expectedError string
		meta          interface{}
	}{
		{
			label:         "Success case",
			mockRtcontext: &MockRunTimeContext{},
			key:           "/context/9345674458787728/",
			meta:          interface{}(MockCompositeAppMeta{Project: "Testproject", CompositeApp: "TestCompApp", Version: "CompAppVersion", Release: "TestRelease"}),
		},
		{
			label:         "Error case for adding app",
			mockRtcontext: &MockRunTimeContext{Err: pkgerrors.Errorf("Error adding app to run time context:")},
			key:           "/context/9345674458787728/",
			expectedError: "Error adding app to run time context:",
			meta:          interface{}(MockCompositeAppMeta{Project: "Testproject", CompositeApp: "TestCompApp", Version: "CompAppVersion", Release: "TestRelease"}),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			ctx := context.Background()
			ac.rtc = testCase.mockRtcontext
			_, err := ac.CreateCompositeApp(ctx)
			_, err = ac.AddApp(ctx, testCase.key, "testapp1")
			if err != nil {
				if !strings.Contains(string(err.Error()), testCase.expectedError) {
					t.Fatalf("Method returned an error (%s)", err)
				}
			}

		})
	}
}

func TestGetAppHandle(t *testing.T) {
	var ac = AppContext{}
	testCases := []struct {
		label         string
		mockRtcontext *MockRunTimeContext
		key           interface{}
		appname       string
		expectedError string
	}{
		{
			label:         "Success case",
			mockRtcontext: &MockRunTimeContext{},
			key:           "/context/9345674458787728/",
			appname:       "testapp1",
		},
		{
			label:         "Invalid app name case",
			mockRtcontext: &MockRunTimeContext{},
			key:           "/context/9345674458787728/",
			appname:       "",
		},
		{
			label:         "Delete returns error case",
			mockRtcontext: &MockRunTimeContext{Err: pkgerrors.Errorf("Error getting app handle from run time context:")},
			key:           "/context/9345674458787728/",
			appname:       "testapp1",
			expectedError: "Error getting app handle from run time context:",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			ctx := context.Background()
			ac.rtc = testCase.mockRtcontext
			_, err := ac.GetAppHandle(ctx, testCase.appname)
			if err != nil {
				if !strings.Contains(string(err.Error()), testCase.expectedError) {
					t.Fatalf("Method returned an error (%s)", err)
				}
			}

		})
	}
}
