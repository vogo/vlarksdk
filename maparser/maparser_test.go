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
	"testing"

	"github.com/stretchr/testify/assert"
)

type Obj struct {
	Str string `json:"str" key:"str" parser:"string"`
	Int int    `json:"int" key:"int" parser:"int"`
}

func TestParse(t *testing.T) {
	obj := &Obj{}
	m := map[string]any{
		"str": "test",
		"int": 123,
	}

	err := Parse(obj, m)
	assert.Nil(t, err)

	assert.Equal(t, "test", obj.Str)
	assert.Equal(t, 123, obj.Int)
}
