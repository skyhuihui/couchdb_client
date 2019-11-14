package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	And       = "$and"       //Array	Matches if all the selectors in the array match.
	Or        = "$or"        //Array	Matches if any of the selectors in the array match. All selectors must use the same index.
	Not       = "$not"       //Selector	Matches if the given selector does not match.
	Nor       = "$nor"       //Array	Matches if none of the selectors in the array match.
	All       = "$all"       //Array	Matches an array value if it contains all the elements of the argument array.
	ElemMatch = "$elemMatch" //Selector	Matches and returns all documents that contain an array field with at least one element that matches all the specified query criteria.
	AllMatch  = "$allMatch"  //Selector	Matches and returns all documents that contain an array field with all its elements matching all the specified query criteria.

	Lt     = "$lt"     //Any JSON	The field is less than the argument
	Lte    = "$lte"    //Any JSON	The field is less than or equal to the argument.
	Eq     = "$eq"     //Any JSON	The field is equal to the argument
	Ne     = "$ne"     //Any JSON	The field is not equal to the argument.
	Gte    = "$gte"    //Any JSON	The field is greater than or equal to the argument.
	Gt     = "$gt"     //Any JSON	The field is greater than the to the argument.
	Exists = "$exists" //Boolean	Check whether the field exists or not, regardless of its value.
	Type   = "$type"   //String	Check the document field’s type. Valid values are "null", "boolean", "number", "string", "array", and "object".
	In     = "$in"     //Array of JSON values	The document field must exist in the list provided.
	Nin    = "$nin"    //Array of JSON values	The document field not must exist in the list provided.
	Size   = "$size"   //Integer	Special condition to match the length of an array field in a document. Non-array fields cannot match this condition.
	Mod    = "$mod"    //[Divisor, Remainder]	Divisor and Remainder are both positive or negative integers. Non-integer values result in a 404. Matches documents where field % Divisor == Remainder is true, and only when the document field is an integer.
	Regex  = "$regex"  //String	A regular expression pattern to match against the document field. Only matches when the field is a string value and matches the supplied regular expression. The matching algorithms are based on the Perl Compatible Regular Expression (PCRE) library. For more information about what is implemented, see the see the Erlang Regular Expression )
	Asc    = "asc"
	Desc   = "desc"
)

type CouchDBClient interface {
	Find() (*CouchSelectorBody, error)
	CreateIndex() (*CouchIndexBody, error)
}

type CouchDB struct {
	Host     string `json:"Host,omitempty"`     // ip
	Port     string `json:"Port,omitempty"`     // port
	DataBase string `json:"DataBase,omitempty"` // DataBase name
	Args     Args   `json:"Args,omitempty"`     // args
	Index    Index  `json:"Index,omitempty"`    // args     Args   `json:"Args,omitempty"`     // args
}

type Args struct {
	// http://docs.couchdb.org/en/stable/api/database/find.html
	Selector       map[string]interface{} `json:"selector,omitempty"`        // – 选择器， 查询条件参数 JSON object describing criteria used to select documents. More information provided in the section on selector syntax. Required
	Limit          int64                  `json:"limit,omitempty"`           //  – 查询条数，配合bookmark 使用可以达到分页效果  Maximum number of results returned. Default is 25. Optional
	Skip           int64                  `json:"skip,omitempty"`            // – Skip the first ‘n’ results, where ‘n’ is the value specified. Optional
	Sort           []interface{}          `json:"sort,omitempty"`            // – 排序 JSON array following sort syntax. Optional
	Fields         []interface{}          `json:"fields,omitempty"`          // – 过滤 JSON array specifying which fields of each object should be returned. If it is omitted, the entire object is returned. More information provided in the section on filtering fields. Optional
	UseIndex       string                 `json:"use_index,omitempty"`       //  索引( string|array)  – Instruct a query to use a specific index. Specified either as "<design_document>" or ["<design_document>", "<index_name>"]. Optional
	R              int64                  `json:"r,omitempty"`               // (number)   – Read quorum needed for the result. This defaults to 1, in which case the document found in the index is returned. If set to a higher value, each document is read from at least that many replicas before it is returned in the results. This is likely to take more time than using only the document stored locally with the index. Optional, default: 1
	Bookmark       string                 `json:"bookmark,omitempty"`        //  – A string that enables you to specify which page of results you require. Used for paging through result sets. Every query returns an opaque string under the bookmark key that can then be passed back in a query to get the next page of results. If any part of the selector query changes between requests, the results are undefined. Optional, default: null
	Update         bool                   `json:"update,omitempty"`          //(boolean) – Whether to update the index prior to returning the result. Default is true. Optional
	Stable         bool                   `json:"stable,omitempty"`          //(boolean)  – Whether or not the view results should be returned from a “stable” set of shards. Optional
	Stale          string                 `json:"stale,omitempty"`           // – Combination of update=false and stable=true options. Possible options: "ok", false (default). Optional
	ExecutionStats bool                   `json:"execution_stats,omitempty"` //  在查询响应中包含执行统计信息 (boolean)  – Include execution statistics in the query response. Optional, default: ``false``\
}

type CouchSelectorBody struct {
	Docs           []interface{}  `json:"docs,omitempty"`            // – JSON object describing criteria used to select documents. More information provided in the section on selector syntax. Required
	BookMark       string         `json:"bookmark,omitempty"`        //  – 配合limit使用， 查询一次返回的， 下次再查询传入这个查询的就是下一页， 分页效果
	ExecutionStats ExecutionStats `json:"execution_stats,omitempty"` // —在查询响应中包含执行统计信息
	Warning        string         `json:"warning,omitempty"`         // – 异常信息
}

type ExecutionStats struct {
	ExecutionTimeMs         float64 `json:"execution_time_ms,omitempty"`          // 执行时间
	ResultsReturned         int     `json:"results_returned,omitempty"`           // 结果返回
	TotalDocsExamined       int     `json:"total_docs_examined,omitempty"`        // 已检查的文档总数
	TotalKeysExamined       int     `json:"total_keys_examined,omitempty"`        // 已检查的总键数
	TotalQuorumDocsExamined int     `json:"total_quorum_docs_examined,omitempty"` // 已审核的法定文档总数
}

/*
{
	"index": {
		"partial_filter_selector": {
			"Attribute.UserHospitalizedCases": {
				"$in": [{
					"Name": "刘"
				}]
			}
		},
		 "fields": ["UserKey"]
	},
	"name": "0x00ifhd72832jfsuajzi",
	"type": "json"
}
*/
type Index struct {
	Index interface{} `json:"index,omitempty"` //fields []string,  partial_filter_selector  interface{}               // 描述要创建的索引的json对象。(json) – JSON object describing the index to create.
	DDoc  string      `json:"ddoc,omitempty"`  // (string) – Name of the design document in which the index will be created. By default, each index will be created in its own design document. Indexes can be grouped into design documents for efficiency. However, a change to one index in a design document will invalidate all other indexes in the same document (similar to views). Optional
	Name  string      `json:"name,omitempty"`  // (string) – Name of the index. If no name is provided, a name will be generated automatically. Optional
	Type  string      `json:"type,omitempty"`  // (string) – Can be "json" or "text". Defaults to json. Geospatial indexes will be supported in the future. Optional Text indexes are supported via a third party library Optional
}

type CouchIndexBody struct {
	Result string `json:"result,omitempty"` //(string) – Flag to show whether the index was created or one already exists. Can be “created” or “exists”.
	Id     string `json:"id,omitempty"`     //(string) – Id of the design document the index was created in.
	Name   string `json:"name,omitempty"`   //(string) – Name of the index created.
}

func NewCouchDBClientService(c CouchDB) CouchDBClient {
	return &c
}

func (c CouchDB) Find() (*CouchSelectorBody, error) {
	cBytes, _ := json.Marshal(c.Args)
	url := c.Host + ":" + c.Port + "/" + c.DataBase + "/" + "_find"
	payload := strings.NewReader(string(cBytes))
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	couchSelectorBody := &CouchSelectorBody{}
	err = json.Unmarshal(body, couchSelectorBody)
	if err != nil {
		return nil, err
	}
	return couchSelectorBody, nil
}

func (c CouchDB) CreateIndex() (*CouchIndexBody, error) {
	cBytes, _ := json.Marshal(c.Index)
	url := c.Host + ":" + c.Port + "/" + c.DataBase + "/" + "_index"
	fmt.Println(string(cBytes))
	payload := strings.NewReader(string(cBytes))
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	couchIndexBody := &CouchIndexBody{}
	err = json.Unmarshal(body, couchIndexBody)
	if err != nil {
		return nil, err
	}
	return couchIndexBody, nil
}
