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

	"github.com/vogo/vlarksdk/vutil"
)

type FileInfo struct {
	FileToken string  `json:"file_token"`
	Name      string  `json:"name"`
	Size      float64 `json:"size"`
	TempUrl   string  `json:"temp_url"`
	Type      string  `json:"type"`
	Url       string  `json:"url"`
}

func ParseFileInfoMap(m map[string]any) (info *FileInfo, err error) {
	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			err = fmt.Errorf("parse file info map failed: %v", recoverErr)
		}
	}()

	if m == nil {
		return nil, nil
	}

	fileInfo := &FileInfo{}
	fileInfo.FileToken = vutil.EnsureMapString(m, "file_token")
	fileInfo.Name = vutil.EnsureMapString(m, "name")
	fileInfo.Size = vutil.EnsureMapFloat64(m, "size")
	fileInfo.TempUrl = vutil.EnsureMapString(m, "temp_url")
	fileInfo.Type = vutil.EnsureMapString(m, "type")
	fileInfo.Url = vutil.EnsureMapString(m, "url")

	return fileInfo, nil
}

func FileArrayValueParser(val any) (any, error) {
	if val == nil {
		return nil, nil
	}

	switch v := val.(type) {
	case []interface{}:
		infoArr := make([]*FileInfo, 0, len(v))
		for i := 0; i < len(v); i++ {
			m, ok := v[i].(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid type %T", v[i])
			}

			info, err := ParseFileInfoMap(m)
			if err != nil {
				return nil, err
			}
			infoArr = append(infoArr, info)
		}
		return infoArr, nil
	case map[string]interface{}:
		info, err := ParseFileInfoMap(v)
		if err != nil {
			return nil, err
		}
		return []*FileInfo{info}, nil
	default:
		return nil, fmt.Errorf("invalid type %T", val)
	}
}

func FileArrayFieldParser(dest reflect.Value, val any) error {
	if val == nil {
		return nil
	}

	fileInfo, err := FileArrayValueParser(val)
	if err != nil {
		return err
	}
	if fileInfo == nil {
		return nil
	}

	dest.Set(reflect.ValueOf(fileInfo))

	return nil
}
