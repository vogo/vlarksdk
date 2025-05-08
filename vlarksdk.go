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

package vlarksdk

import (
	"log"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"

	_ "github.com/vogo/vlarksdk/maparser"
	_ "github.com/vogo/vlarksdk/vbitable/vbitableparser"
)

var LarkCli *lark.Client

func InitLarkService(appId, appSecret string) {
	log.Printf("lark app id: %s", appId)

	LarkCli = lark.NewClient(appId, appSecret, lark.WithLogLevel(larkcore.LogLevelInfo))
}
