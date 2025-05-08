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

package maparser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type (
	FieldParser func(dest reflect.Value, val any) error
	ValueParser func(val any) (any, error)
)

var fieldParserMap = map[string]FieldParser{}

var valueParserMap = map[string]ValueParser{}

func init() {
	SetFieldParser("string", StringValueParser, StringFieldParser)
	SetFieldParser("int", IntValueParser, IntFieldParser)
	SetFieldParser("float", FloatValueParser, FloatFieldParser)
	SetFieldParser("array_to_string", ArrayToStringValueParser, ArrayToStringFieldParser)
}

func SetFieldParser(name string, valueParser ValueParser, fieldParser FieldParser) {
	valueParserMap[name] = valueParser
	fieldParserMap[name] = fieldParser
}

type fieldConfig struct {
	Key    string
	Parser FieldParser
}

// PackPtr pack a Ptr value
func PackPtr(v reflect.Value) reflect.Value {
	vv := reflect.New(v.Type())
	vv.Elem().Set(v)
	return vv
}

func SetValue(dest reflect.Value, val any) {
	v := reflect.ValueOf(val)
	if dest.Kind() == reflect.Ptr && v.Kind() != reflect.Ptr {
		v = PackPtr(v)
	}
	dest.Set(v)
}

var typeMapFieldConfigMap = map[reflect.Type]map[string]fieldConfig{}

func getTypeMapFieldConfig(t reflect.Type) (map[string]fieldConfig, error) {
	c, ok := typeMapFieldConfigMap[t]
	if !ok {
		c = make(map[string]fieldConfig)
		typeMapFieldConfigMap[t] = c

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			key, tagOk := field.Tag.Lookup("key")
			if !tagOk {
				continue
			}
			parserName, parserOk := field.Tag.Lookup("parser")
			if !parserOk {
				parserName = "string"
			}

			parser, parserOk := fieldParserMap[parserName]
			if !parserOk {
				return nil, fmt.Errorf("invalid parser %s", parserName)
			}

			c[field.Name] = fieldConfig{
				Key:    key,
				Parser: parser,
			}
		}
	}

	return c, nil
}

func Parse(dest any, m map[string]any) (perr error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			perr = fmt.Errorf("parse error: %v", panicErr)
		}
	}()

	destValue := reflect.ValueOf(dest)
	destType := reflect.TypeOf(dest)
	for destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
		destValue = destValue.Elem()
	}

	fieldConfigMap, err := getTypeMapFieldConfig(destType)
	if err != nil {
		return err
	}

	var parseErr error
	for fieldName, config := range fieldConfigMap {
		val := m[config.Key]
		field := destValue.FieldByName(fieldName)
		if !field.IsValid() {
			continue
		}

		if parseErr = config.Parser(field, val); parseErr != nil {
			return fmt.Errorf("parse error: %v, field: %s, value: %v ", parseErr, fieldName, val)
		}
	}

	return nil
}

func StringValueParser(val any) (any, error) {
	if val == nil {
		return nil, nil
	}

	switch v := val.(type) {
	case []interface{}:
		if len(v) > 0 {
			val = v[0]
		}
	}

	s, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("invalid type %T", val)
	}

	return s, nil
}

func StringFieldParser(dest reflect.Value, val any) error {
	s, err := StringValueParser(val)
	if err != nil {
		return err
	}
	if s == nil {
		return nil
	}

	dest.SetString(s.(string))

	return nil
}

func ParseStringField(val any) string {
	if val == nil {
		return ""
	}
	s, ok := val.(string)
	if ok {
		return s
	}

	return fmt.Sprintf("%v", val)
}

func ArrayToStringValueParser(val any) (any, error) {
	if val == nil {
		return nil, nil
	}

	anyArr, ok := val.([]interface{})
	if !ok || !isArrAllString(anyArr) {
		switch v := val.(type) {
		case []interface{}:
			if len(v) > 0 {
				val = v[0]
			}
		}
	}

	arr, ok := val.([]string)
	if ok {
		return strings.Join(arr, ","), nil
	}

	faceArr, ok := val.([]interface{})
	if ok {
		var s string
		for _, v := range faceArr {
			s, ok = v.(string)
			if !ok {
				return nil, fmt.Errorf("ArrayToStringValueParser: invalid array item type %T", v)
			}

			arr = append(arr, s)
		}

		return strings.Join(arr, ","), nil
	}

	return nil, fmt.Errorf("ArrayToStringValueParser: invalid type %T", val)
}

func isArrAllString(arr []interface{}) bool {
	for _, v := range arr {
		if _, ok := v.(string); !ok {
			return false
		}
	}

	return true
}

func ArrayToStringFieldParser(dest reflect.Value, val any) error {
	s, err := ArrayToStringValueParser(val)
	if err != nil {
		return err
	}
	if s == nil {
		return nil
	}
	dest.SetString(s.(string))
	return nil
}

func IntValueParser(val any) (any, error) {
	if val == nil {
		return nil, nil
	}

	switch v := val.(type) {
	case []interface{}:
		if len(v) > 0 {
			val = v[0]
		}
	}

	switch v := val.(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return nil, fmt.Errorf("invalid type %T", val)
	}
}

func IntFieldParser(dest reflect.Value, val any) error {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case int:
		dest.SetInt(int64(v))
	case int32:
		dest.SetInt(int64(v))
	case int64:
		dest.SetInt(v)
	case float32:
		dest.SetInt(int64(v))
	case float64:
		dest.SetInt(int64(v))
	case string:
		if v == "" {
			return nil
		}

		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		dest.SetInt(i)
	default:
		return fmt.Errorf("invalid type %T", val)
	}

	return nil
}

func ParseIntField(val any) (int64, error) {
	if val == nil {
		return 0, nil
	}

	switch v := val.(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("invalid type: %T", val)
	}
}

func FloatValueParser(val any) (any, error) {
	if val == nil {
		return nil, nil
	}

	switch v := val.(type) {
	case []interface{}:
		if len(v) > 0 {
			val = v[0]
		}
	}

	switch v := val.(type) {
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		return strToFloat(v)
	default:
		return nil, fmt.Errorf("invalid type %T", val)
	}
}

func FloatFieldParser(dest reflect.Value, val any) error {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case int:
		dest.SetFloat(float64(v))
	case int32:
		dest.SetFloat(float64(v))
	case int64:
		dest.SetFloat(float64(v))
	case float32:
		dest.SetFloat(float64(v))
	case float64:
		dest.SetFloat(v)
	case string:
		f, err := strToFloat(v)
		if err != nil {
			return fmt.Errorf("FloatFieldParser: invalid flaot %s, error: %v", val, err)
		}
		dest.SetFloat(f)
	default:
		str := fmt.Sprintf("%v", val)

		f, err := strToFloat(str)
		if err != nil {
			return fmt.Errorf("FloatFieldParser: invalid flaot %s, error: %v", str, err)
		}

		dest.SetFloat(f)
	}

	return nil
}

func ParseFloatField(val any) (float64, error) {
	switch v := val.(type) {
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		f, err := strToFloat(v)
		if err != nil {
			return 0, fmt.Errorf("FloatFieldParser: invalid flaot %s, error: %v", val, err)
		}
		return f, nil
	default:
		str := fmt.Sprintf("%v", val)

		f, err := strToFloat(str)
		if err != nil {
			return 0, fmt.Errorf("FloatFieldParser: invalid flaot %s, error: %v", str, err)
		}

		return f, nil
	}
}

func strToFloat(v string) (float64, error) {
	if strings.HasSuffix(v, "%") {
		v = strings.TrimSuffix(v, "%")
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid float %s", v)
		}
		return f / 100, nil
	}

	return strconv.ParseFloat(v, 64)
}
