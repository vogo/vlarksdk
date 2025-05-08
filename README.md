# vlarksdk

github.com/larksuite/oapi-sdk-go extension


## parse bitable record to struct
```go
type Record struct {
	Id       int       `json:"id" key:"用户编号" parser:"int"`
	Name     string    `json:"name" key:"姓名" parser:"string"`
	Date     time.Time `json:"date" key:"记录最后更新时间" parser:"timestamp"`
	ModifyBy string    `json:"modify_by" key:"修改人" parser:"single_user_name"`
}

func queryTableRecord(tableId, tableAppToken string) {
	queryReq := larkbitable.NewListAppTableRecordReqBuilder().TableId(tableId).
		AppToken(tableAppToken).
		Limit(1).
		Build()

	ctx := context.Background()
	queryResp, err := vlarksdk.LarkCli.Bitable.AppTableRecord.List(ctx, queryReq)
	item := queryResp.Data.Items[0]
	data := &Record{}
	pErr := maparser.Parse(data, item.Fields)
	log.Printf("data: %+v", data)
}
```
