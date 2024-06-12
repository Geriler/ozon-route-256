package model

import (
	"bytes"
	"net/http"
	"testing"
)

func TestGetValidateUserRequest(t *testing.T) {
	cases := []struct {
		name     string
		userId   string
		expected *UserRequest
		hasErr   bool
	}{
		{
			name:     "success",
			userId:   "1",
			expected: &UserRequest{UserID: 1},
			hasErr:   false,
		},
		{
			name:     "symbol error",
			userId:   "a",
			expected: nil,
			hasErr:   true,
		},
		{
			name:     "zero error",
			userId:   "0",
			expected: nil,
			hasErr:   true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := &http.Request{}
			r.SetPathValue(FieldUserID, tc.userId)
			req, err := GetValidateUserRequest(r)
			if tc.hasErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Error(err.Error())
			}
			if req.UserID != tc.expected.UserID {
				t.Errorf("expected %v, got %v", tc.expected, req)
			}
		})
	}
}

func TestGetValidateUserSKURequest(t *testing.T) {
	cases := []struct {
		name     string
		userId   string
		skuId    string
		expected *UserSKURequest
		hasErr   bool
	}{
		{
			name:     "success",
			userId:   "1",
			skuId:    "1",
			expected: &UserSKURequest{UserID: 1, SKU: 1},
			hasErr:   false,
		},
		{
			name:     "symbol error in user id",
			userId:   "a",
			skuId:    "1",
			expected: nil,
			hasErr:   true,
		},
		{
			name:     "zero error in user id",
			userId:   "0",
			skuId:    "1",
			expected: nil,
			hasErr:   true,
		},
		{
			name:     "symbol error in sku id",
			userId:   "1",
			skuId:    "a",
			expected: nil,
			hasErr:   true,
		},
		{
			name:     "zero error in sku id",
			userId:   "1",
			skuId:    "0",
			expected: nil,
			hasErr:   true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := &http.Request{}
			r.SetPathValue(FieldUserID, tc.userId)
			r.SetPathValue(FieldSkuID, tc.skuId)
			req, err := GetValidateUserSKURequest(r)
			if tc.hasErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Error(err.Error())
			}
			if req.UserID != tc.expected.UserID || req.SKU != tc.expected.SKU {
				t.Errorf("expected %v, got %v", tc.expected, req)
			}
		})
	}
}

func TestGetValidateUserSKUCountRequest(t *testing.T) {
	cases := []struct {
		name     string
		userId   string
		skuId    string
		count    string
		expected *UserSKUCountRequest
		hasErr   bool
	}{
		{
			name:     "success",
			userId:   "1",
			skuId:    "1",
			count:    "1",
			expected: &UserSKUCountRequest{UserID: 1, SKU: 1, Count: 1},
			hasErr:   false,
		},
		{
			name:     "symbol error in user id",
			userId:   "a",
			skuId:    "1",
			count:    "1",
			expected: nil,
			hasErr:   true,
		},
		{
			name:     "zero error in user id",
			userId:   "0",
			skuId:    "1",
			count:    "1",
			expected: nil,
			hasErr:   true,
		},
		{
			name:     "symbol error in sku id",
			userId:   "1",
			skuId:    "a",
			count:    "1",
			expected: nil,
			hasErr:   true,
		},
		{
			name:     "zero error in sku id",
			userId:   "1",
			skuId:    "0",
			count:    "1",
			expected: nil,
			hasErr:   true,
		},
		{
			name:     "zero error in count",
			userId:   "1",
			skuId:    "1",
			count:    "0",
			expected: nil,
			hasErr:   true,
		},
		{
			name:     "failed to decode request",
			userId:   "1",
			skuId:    "1",
			count:    "a",
			expected: nil,
			hasErr:   true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body := bytes.NewReader([]byte(`{"count":` + tc.count + `}`))
			r, _ := http.NewRequest(http.MethodPost, "https://example.com", body)
			r.SetPathValue(FieldUserID, tc.userId)
			r.SetPathValue(FieldSkuID, tc.skuId)
			req, err := GetValidateUserSKUCountRequest(r)
			if tc.hasErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Error(err.Error())
			}
			if req.UserID != tc.expected.UserID || req.SKU != tc.expected.SKU || req.Count != tc.expected.Count {
				t.Errorf("expected %v, got %v", tc.expected, req)
			}
		})
	}
}
