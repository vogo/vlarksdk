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
	"time"

	"github.com/vogo/vlarksdk/maparser"
)

func ParseMapFieldTextLink(val any) (string, error) {
	if val == nil {
		return "", nil
	}

	if s, ok := val.(string); ok {
		return s, nil
	}

	arr, arrok := val.([]interface{})
	if !arrok {
		return "", fmt.Errorf("MapFieldTextFieldParser: invalid type %T", val)
	}

	if len(arr) == 0 {
		return "", nil
	}

	var result string
	for _, v := range arr {
		m, ok := v.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("MapFieldTextFieldParser: invalid type %T", val)
		}

		text, ok := m["text"]
		if !ok {
			continue
		}
		s, ok := text.(string)

		if !ok {
			return "", fmt.Errorf("MapFieldTextFieldParser: invalid type %T", text)
		}
		link, ok := m["link"]
		if !ok {
			continue
		}
		l, ok := link.(string)

		if !ok {
			return "", fmt.Errorf("MapFieldTextFieldParser: invalid type %T", text)
		}
		result += fmt.Sprintf("%s(%s) ", s, l)

	}

	return result, nil
}

func MapFieldTextValueParser(val any) (any, error) {
	return ParseMapFieldText(val)
}

func ParseMapFieldText(val any) (string, error) {
	if val == nil {
		return "", nil
	}

	str, ok := val.(string)
	if ok {
		return str, nil
	}

	arr, ok := val.([]interface{})
	if !ok {
		return "", fmt.Errorf("ParseMapFieldText: invalid type %T", val)
	}

	if len(arr) == 0 {
		return "", nil
	}

	var textArr []string
	for _, v := range arr {
		s, err := parseMapTextField(v)
		if err != nil {
			return "", err
		}
		textArr = append(textArr, s)
	}

	return strings.Join(textArr, ","), nil
}

func parseMapTextField(v interface{}) (string, error) {
	str, ok := v.(string)
	if ok {
		return str, nil
	}

	m, ok := v.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("parseMapTextField: invalid type %T", v)
	}

	text, ok := m["text"]
	if !ok {
		return "", fmt.Errorf("parseMapTextField: %s", "text field not found")
	}
	s, ok := text.(string)

	if !ok {
		return "", fmt.Errorf("parseMapTextField: invalid type %T", text)
	}

	return s, nil
}

func ParseMapFieldAttachUrls(val any) (string, error) {
	if val == nil {
		return "", nil
	}

	arr, ok := val.([]interface{})
	if !ok {
		return "", fmt.Errorf("ParseMapFieldAttachUrls: invalid type %T", val)
	}

	if len(arr) == 0 {
		return "", nil
	}

	var urls []string
	for _, v := range arr {
		m, ok := v.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("ParseMapFieldAttachUrls: invalid type %T", val)
		}

		token, ok := m["file_token"]
		if !ok {
			return "", fmt.Errorf("ParseMapFieldAttachUrls: %s", "file_token field not found")
		}
		t, ok := token.(string)
		if !ok {
			return "", fmt.Errorf("ParseMapFieldAttachUrls: invalid type %T", token)
		}
		name, ok := m["name"]
		if !ok {
			return "", fmt.Errorf("ParseMapFieldAttachUrls: %s", "name field not found")
		}
		n, ok := name.(string)
		if !ok {
			return "", fmt.Errorf("ParseMapFieldAttachUrls: invalid type %T", name)
		}
		urls = append(urls, fmt.Sprintf("%s(%s)", n, t))
	}

	return strings.Join(urls, ","), nil
}

func MapFieldAttachParser(dest reflect.Value, val any) error {
	if val == nil {
		return nil
	}

	s, err := ParseMapFieldAttachUrls(val)
	if err != nil {
		return err
	}
	maparser.SetValue(dest, s)
	return nil
}

func MapFieldTextLinkParser(dest reflect.Value, val any) error {
	if val == nil {
		return nil
	}

	s, err := ParseMapFieldTextLink(val)
	if err != nil {
		return err
	}

	dest.SetString(s)
	return nil
}

func MapFieldTextFieldParser(dest reflect.Value, val any) error {
	if val == nil {
		return nil
	}

	s, err := ParseMapFieldText(val)
	if err != nil {
		return err
	}

	dest.SetString(s)
	return nil
}

func MapFieldTextDateParser(dest reflect.Value, val any) error {
	if val == nil {
		return nil
	}

	var timestamp int64
	switch v := val.(type) {
	case int:
		timestamp = int64(v)
	case int32:
		timestamp = int64(v)
	case int64:
		timestamp = v
	case float32:
		timestamp = int64(v)
	case float64:
		timestamp = int64(v)
	}
	if timestamp > 0 {
		date := time.UnixMilli(timestamp)
		maparser.SetValue(dest, date)
		return nil
	}

	s, err := ParseMapFieldText(val)
	if err != nil {
		return err
	}

	date, err := time.Parse("2006/01/02", s)
	if err != nil {
		date, err = time.Parse("2006-01-02", s)
		if err != nil {
			return fmt.Errorf("MapFieldTextDateParser: invalid time format %s, error:%s", s, err)
		}
	}

	maparser.SetValue(dest, date)
	return nil
}

func TimestampFieldParser(dest reflect.Value, val any) error {
	if val == nil {
		return nil
	}

	t, err := ParseTimestampValue(val)
	if err != nil {
		return err
	}
	maparser.SetValue(dest, t)
	return nil
}

func TimestampValueParser(val any) (any, error) {
	return ParseTimestampValue(val)
}

func ParseTimestampValue(val any) (time.Time, error) {
	if val == nil {
		return time.Time{}, nil
	}

	var timestamp int64

	switch v := val.(type) {
	case int:
		timestamp = int64(v)
	case int32:
		timestamp = int64(v)
	case int64:
		timestamp = v
	case float32:
		timestamp = int64(v)
	case float64:
		timestamp = int64(v)
	case []interface{}:
		if len(v) > 0 {
			return ParseTimestampValue(v[0])
		}
		return time.Time{}, nil
	default:
		return time.Time{}, fmt.Errorf("TimestampFieldParser: invalid type %T", val)
	}

	t := time.UnixMilli(timestamp)
	return t, nil
}

func LarkDaysParser(dest reflect.Value, val any) error {
	if val == nil {
		return nil
	}

	var days float64

	switch v := val.(type) {
	case int:
		days = float64(v)
	case int32:
		days = float64(v)
	case int64:
		days = float64(v)
	case float32:
		days = float64(v)
	case float64:
		days = v
	default:
		return fmt.Errorf("TimestampFieldParser: invalid type %T", val)
	}

	t, _ := time.Parse("2006/01/02", "1899/12/30")
	t = t.Add(time.Duration(days*24*3600*1000) * time.Millisecond)
	maparser.SetValue(dest, t)
	return nil
}

func FuncIntParser(dest reflect.Value, val any) error {
	if val == nil {
		return nil
	}

	_ = maparser.IntFieldParser(dest, val)

	return nil
}
