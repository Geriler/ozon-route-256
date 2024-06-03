package model

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gopkg.in/validator.v2"
)

var (
	FieldUserID = "user_id"
	FieldSkuID  = "sku_id"
)

type UserRequest struct {
	UserID UserID `validate:"nonzero"`
}

type UserSKURequest struct {
	UserRequest
	SKU SkuID `validate:"nonzero"`
}

type UserSKUCountRequest struct {
	UserSKURequest
	Count uint16 `validate:"nonzero"`
}

func GetValidateUserRequest(r *http.Request) (*UserRequest, error) {
	req := &UserRequest{}

	userId, err := strconv.Atoi(r.PathValue(FieldUserID))
	if err != nil {
		return req, err
	}

	req.UserID = UserID(userId)

	if err = validator.Validate(req); err != nil {
		return req, err
	}

	return req, nil
}

func GetValidateUserSKURequest(r *http.Request) (*UserSKURequest, error) {
	req := &UserSKURequest{}

	userId, err := strconv.Atoi(r.PathValue(FieldUserID))
	if err != nil {
		return req, err
	}

	skuId, err := strconv.Atoi(r.PathValue(FieldSkuID))
	if err != nil {
		return req, err
	}

	req.UserID = UserID(userId)
	req.SKU = SkuID(skuId)

	if err = validator.Validate(req); err != nil {
		return req, err
	}

	return req, nil
}

func GetValidateUserSKUCountRequest(r *http.Request) (*UserSKUCountRequest, error) {
	req := &UserSKUCountRequest{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return req, err
	}

	userId, err := strconv.Atoi(r.PathValue(FieldUserID))
	if err != nil {
		return req, err
	}

	skuId, err := strconv.Atoi(r.PathValue(FieldSkuID))
	if err != nil {
		return req, err
	}

	req.UserID = UserID(userId)
	req.SKU = SkuID(skuId)

	if err = validator.Validate(req); err != nil {
		return req, err
	}

	return req, nil
}
