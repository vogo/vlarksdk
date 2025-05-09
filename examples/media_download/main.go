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
	"io"
	"log"
	"os"

	drivev1 "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
	"github.com/vogo/vlarksdk"
	"github.com/vogo/vogo/vos"
)

func main() {
	appId := vos.EnvString("LARK_APP_ID")
	appSecret := vos.EnvString("LARK_APP_SECRET")

	if appId == "" || appSecret == "" {
		panic("appId or appSecret is empty")
	}

	testClientDownload(appId, appSecret)
}

func testClientDownload(appId string, appSecret string) {
	vlarksdk.InitLarkService(appId, appSecret)

	fileToken := "MS6tbq0ZPomK8CxsOhtcgiO9n6e"
	downloadReq := drivev1.NewDownloadMediaReqBuilder().FileToken(fileToken).Build()
	downloadResp, err := vlarksdk.LarkCli.Drive.V1.Media.Download(context.Background(), downloadReq)
	if err != nil {
		log.Printf("err: %v", err)
		return
	}

	outputFile, err := os.Create("/Users/hk/temp/test.png")
	if err != nil {
		log.Printf("create file error: %s", err)
		return
	}
	defer outputFile.Close()
	fileSize, fileErr := io.Copy(outputFile, downloadResp.File)
	if fileErr != nil {
		log.Printf("copy file error: %s", fileErr)
		return
	}
	log.Printf("file download success, file size: %d", fileSize)
}
