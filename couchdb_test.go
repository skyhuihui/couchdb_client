package db

import (
	"fmt"
	"testing"
)

func TestSelectorFind(t *testing.T) {
	/*
		{
		"selector": {
			"$and": [{
				"Attribute": {
					"Name": {
						"$eq": "董华凌（受体）"
					}
				}
			}],
			"Attribute": {
				"UserKey": "003a09f30f151d56d550d60ab0fad3ef6d399677b8beaab0d498456d2f30780f",
				"Age": "51",
				"InvitationCode": {
					"$eq": "8M2F84K"
				}
			},
			"_id": "003a09f30f151d56d550d60ab0fad3ef6d399677b8beaab0d498456d2f30780f"
		}
		}
	*/
	//var a []map[string]interface{}
	args := Args{
		Selector: map[string]interface{}{
			And: []interface{}{
				map[string]interface{}{
					"Attribute": map[string]interface{}{
						"Name": map[string]interface{}{
							Eq: "董华凌（受体）",
						},
					},
				},
			},
			"Attribute": map[string]interface{}{
				"UserKey": "003a09f30f151d56d550d60ab0fad3ef6d399677b8beaab0d498456d2f30780f",
				"Age":     "51",
				"InvitationCode": map[string]interface{}{
					Eq: "8M2F84K",
				},
			},
			"_id": "003a09f30f151d56d550d60ab0fad3ef6d399677b8beaab0d498456d2f30780f",
		},
		Limit: 10,
	}
	couchCli := CouchDB{
		Host:     "http://112.126.97.169",
		Port:     "5984",
		DataBase: "contract_$user$c$c",
		Args:     args,
	}
	couchDBClient := NewCouchDBClientService(couchCli)
	body, _ := couchDBClient.Find()
	fmt.Println(body)
}

func TestCreateIndex(t *testing.T) {
	/*
		{
			"index": {
				"partial_filter_selector": {
					"Attribute.UserBasics.IDType":{
		          "$in": [
		          	"0"
		          ]
		        }
				},
				 "fields": ["UserKey",
				 "Attribute",
				 "Password",
				 "_id",
				 "_rev",
				 "~version"]
			},
			"name": "0x00ifhd728320xjfsuajzi",
			"type": "json"
		}
	*/
	//var a []map[string]interface{}
	index := Index{
		Index: map[string]interface{}{
			"fields": []string{"UserKey",
				"Attribute",
				"Password",
				"_id",
				"_rev",
				"~version"},
			"partial_filter_selector": map[string]interface{}{
				"Attribute.UserBasics.IDType": map[string]interface{}{
					In: []string{"0"},
				},
			},
		},
		Name: "0x00ifhd728320xjfsuajzi",
		Type: "json",
	}
	couchCli := CouchDB{
		Host:     "http://112.126.97.169",
		Port:     "5984",
		DataBase: "contract_$user$c$c",
		Index:    index,
	}
	couchDBClient := NewCouchDBClientService(couchCli)
	body, _ := couchDBClient.CreateIndex()
	fmt.Println(body)
}
