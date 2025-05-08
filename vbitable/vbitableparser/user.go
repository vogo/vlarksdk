/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package vbitableparser

import (
	"fmt"
	"reflect"
	"strings"
)

func ParseUserNameEmail(val any) (string, error) {
	if val == nil {
		return "", nil
	}

	m, ok := val.(map[string]interface{})
	if !ok {
		var arr []interface{}
		var mm map[string]interface{}

		if arr, ok = val.([]interface{}); ok {
			if mm, ok = arr[0].(map[string]interface{}); ok {
				m = mm
			}
		}

		if !ok {
			return "", fmt.Errorf("ParseUserNameEmail: invalid type %T", val)
		}
	}

	email, ok := m["email"]
	if !ok {
		return "", fmt.Errorf("ParseUserNameEmail: %s", "email field not found")
	}

	name, ok := m["name"]
	if !ok {
		return "", fmt.Errorf("ParseUserNameEmail: %s", "name field not found")
	}

	emailPrefix, _, _ := strings.Cut(email.(string), "@")

	return fmt.Sprintf("%s(%s)", name, emailPrefix), nil
}

func ParseUserName(val any) (string, error) {
	if val == nil {
		return "", nil
	}

	m, ok := val.(map[string]interface{})
	if !ok {
		var arr []interface{}
		var mm map[string]interface{}

		if arr, ok = val.([]interface{}); ok {
			if mm, ok = arr[0].(map[string]interface{}); ok {
				m = mm
			}
		}

		if !ok {
			return "", fmt.Errorf("ParseUserNameEmail: invalid type %T", val)
		}
	}

	name, ok := m["name"]
	if !ok {
		return "", fmt.Errorf("ParseUserNameEmail: %s", "name field not found")
	}

	return name.(string), nil
}

func ParseUserEmail(val any) (string, error) {
	if val == nil {
		return "", nil
	}

	m, ok := val.(map[string]interface{})
	if !ok {
		var arr []interface{}
		var mm map[string]interface{}

		if arr, ok = val.([]interface{}); ok {
			if mm, ok = arr[0].(map[string]interface{}); ok {
				m = mm
			}
		}

		if !ok {
			return "", fmt.Errorf("ParseUserNameEmail: invalid type %T", val)
		}
	}

	email, ok := m["email"]
	if !ok {
		return "", fmt.Errorf("ParseUserEmail: %s", "email field not found")
	}

	return email.(string), nil
}

func ParseMultipleUserUnionIds(val any) ([]string, error) {
	if val == nil {
		return nil, nil
	}

	arr, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("ParseMultipleUserUnionIds: invalid type %T", val)
	}
	var users []string
	for _, v := range arr {
		unionId, err := ParseUserId(v)
		if err != nil {
			return nil, err
		}
		users = append(users, unionId)
	}

	return users, nil
}

func ParseUserId(val any) (string, error) {
	if val == nil {
		return "", nil
	}

	m, ok := val.(map[string]interface{})
	if !ok {
		var arr []interface{}
		var mm map[string]interface{}

		if arr, ok = val.([]interface{}); ok {
			if mm, ok = arr[0].(map[string]interface{}); ok {
				m = mm
			}
		}

		if !ok {
			return "", fmt.Errorf("ParseUserId: invalid type %T", val)
		}
	}

	unionId, ok := m["id"]
	if !ok {
		return "", fmt.Errorf("ParseUserId: %s", "id field not found")
	}

	return unionId.(string), nil
}

func SingleUserNameEmail(dest reflect.Value, val any) error {
	if val == nil {
		return nil
	}

	nameEmail, err := ParseUserNameEmail(val)
	if err != nil {
		return err
	}

	dest.SetString(nameEmail)
	return nil
}

func SingleUserEmail(dest reflect.Value, val any) error {
	if val == nil {
		return nil
	}

	email, err := ParseUserEmail(val)
	if err != nil {
		return err
	}

	dest.SetString(email)
	return nil
}

func SingleUserName(dest reflect.Value, val any) error {
	if val == nil {
		return nil
	}

	id, err := ParseUserName(val)
	if err != nil {
		return err
	}

	dest.SetString(id)
	return nil
}

func SingleUserId(dest reflect.Value, val any) error {
	if val == nil {
		return nil
	}

	id, err := ParseUserId(val)
	if err != nil {
		return err
	}

	dest.SetString(id)
	return nil
}

func MultipleUserNameEmail(dest reflect.Value, val any) error {
	if val == nil {
		return nil
	}

	arr, ok := val.([]interface{})
	if !ok {
		return fmt.Errorf("MultipleUserNameEmail: invalid type %T", val)
	}
	var users []string
	for _, v := range arr {
		nameEmail, err := ParseUserNameEmail(v)
		if err != nil {
			return err
		}
		users = append(users, nameEmail)
	}

	dest.SetString(strings.Join(users, ","))
	return nil
}

type LarkUser struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	OpenId string `json:"open_id"`
}

func SingleUser(dest reflect.Value, val any) error {
	if val == nil {
		return nil
	}

	u, err := ParseUser(val)
	if err != nil {
		return err
	}

	dest.Set(reflect.ValueOf(u))
	return nil
}

func MultipleUsers(dest reflect.Value, val any) error {
	if val == nil {
		return nil
	}

	u, err := ParseMultipleUsers(val)
	if err != nil {
		return err
	}

	dest.Set(reflect.ValueOf(u))
	return nil
}

func ParseMultipleUsers(val any) ([]*LarkUser, error) {
	if val == nil {
		return nil, nil
	}

	arr, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("ParseMultipleUsers: invalid type %T", val)
	}
	var users []*LarkUser
	for _, v := range arr {
		u, err := ParseUser(v)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func ParseUser(val any) (*LarkUser, error) {
	if val == nil {
		return nil, nil
	}

	m, ok := val.(map[string]interface{})
	if !ok {
		var arr []interface{}
		var mm map[string]interface{}

		if arr, ok = val.([]interface{}); ok {
			if mm, ok = arr[0].(map[string]interface{}); ok {
				m = mm
			}
		}

		if !ok {
			return nil, fmt.Errorf("ParseUser: invalid type %T", val)
		}
	}

	u := &LarkUser{}
	openId, ok := m["id"]
	if !ok {
		return nil, fmt.Errorf("ParseUser: %s", "id field not found")
	}

	u.OpenId = openId.(string)

	name, ok := m["name"]
	if !ok {
		return nil, fmt.Errorf("ParseUser: %s", "name field not found")
	}
	u.Name = name.(string)

	email, ok := m["email"]
	if !ok {
		return nil, fmt.Errorf("ParseUser: %s", "email field not found")
	}
	u.Email = email.(string)

	return u, nil
}
