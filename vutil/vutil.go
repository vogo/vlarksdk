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

package vutil

import "fmt"

func EnsureMapString(m map[string]any, key string) string {
	val, ok := m[key]
	if !ok {
		panic("map key not found: " + key)
	}

	s, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("map key not string, key: %s, val: %v, type: %T", key, val, val))
	}

	return s
}

func EnsureMapInt64(m map[string]any, key string) int64 {
	val, ok := m[key]
	if !ok {
		panic("map key not found: " + key)
	}
	i, ok := val.(int64)
	if !ok {
		panic(fmt.Sprintf("map key not int64, key: %s, val: %v, type: %T", key, val, val))
	}
	return i
}

func EnsureMapFloat64(m map[string]any, key string) float64 {
	val, ok := m[key]
	if !ok {
		panic("map key not found: " + key)
	}
	i, ok := val.(float64)
	if !ok {
		panic(fmt.Sprintf("map key not float64, key: %s, val: %v, type: %T", key, val, val))
	}
	return i
}
