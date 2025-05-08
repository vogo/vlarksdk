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

package main

import (
	"context"
	"log"
	"time"

	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"github.com/vogo/vlarksdk"
	"github.com/vogo/vlarksdk/maparser"
	"github.com/vogo/vogo/vos"
)

type Record struct {
	Id       int       `json:"id" key:"用户编号" parser:"int"`
	Name     string    `json:"name" key:"姓名" parser:"string"`
	Date     time.Time `json:"date" key:"记录最后更新时间" parser:"timestamp"`
	ModifyBy string    `json:"modify_by" key:"修改人" parser:"single_user_name"`
}

func main() {
	appId := vos.EnvString("LARK_APP_ID")
	appSecret := vos.EnvString("LARK_APP_SECRET")
	tableId := vos.EnvString("BITABLE_TABLE_ID")
	tableAppToken := vos.EnvString("BITABLE_APP_TOKEN")

	if appId == "" || appSecret == "" || tableId == "" || tableAppToken == "" {
		log.Fatalf("LARK_APP_ID, LARK_APP_SECRET, BITABLE_TABLE_ID, BITABLE_APP_TOKEN must be set")
	}

	vlarksdk.InitLarkService(appId, appSecret)

	queryTableRecord(tableId, tableAppToken)
}

func queryTableRecord(tableId, tableAppToken string) {
	queryReq := larkbitable.NewListAppTableRecordReqBuilder().TableId(tableId).
		AppToken(tableAppToken).
		Limit(1).
		Build()

	ctx := context.Background()
	queryResp, err := vlarksdk.LarkCli.Bitable.AppTableRecord.List(ctx, queryReq)
	if err != nil {
		log.Printf("query table data error: %s", err)
		return
	}
	if queryResp.Code > 0 {
		log.Printf("query table data error: %s", queryResp.CodeError)
		return
	}
	if len(queryResp.Data.Items) == 0 {
		log.Printf("query table data error: %s", "no data")
		return
	}
	item := queryResp.Data.Items[0]
	data := &Record{}
	pErr := maparser.Parse(data, item.Fields)
	if pErr != nil {
		log.Printf("parse table data error: %s", pErr)
		return
	}
	log.Printf("data: %+v", data)
}
