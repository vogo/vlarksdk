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

package vbitable

import "github.com/vogo/vlarksdk/maparser"

func init() {
	maparser.SetFieldParser("single_user_name_email", nil, SingleUserNameEmail)
	maparser.SetFieldParser("single_user_email", nil, SingleUserEmail)
	maparser.SetFieldParser("single_user_name", nil, SingleUserName)
	maparser.SetFieldParser("single_user_id", nil, SingleUserId)
	maparser.SetFieldParser("multiple_user_name_email", nil, MultipleUserNameEmail)
	maparser.SetFieldParser("single_user", nil, SingleUser)
	maparser.SetFieldParser("multiple_users", nil, MultipleUsers)
	maparser.SetFieldParser("map_field_text", MapFieldTextValueParser, MapFieldTextFieldParser)
	maparser.SetFieldParser("map_field_text_link", nil, MapFieldTextLinkParser)
	maparser.SetFieldParser("map_field_text_date", nil, MapFieldTextDateParser)
	maparser.SetFieldParser("timestamp", TimestampValueParser, TimestampFieldParser)
	maparser.SetFieldParser("lark_days", nil, LarkDaysParser)
	maparser.SetFieldParser("func_int", nil, FuncIntParser)
	maparser.SetFieldParser("map_field_attach", nil, MapFieldAttachParser)
	maparser.SetFieldParser("file_array", FileArrayValueParser, FileArrayFieldParser)
}
