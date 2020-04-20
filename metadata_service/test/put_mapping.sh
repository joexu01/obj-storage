#!/bin/bash

curl 10.29.102.173:9200/metadata -XDELETE

# elasticsearch v7
curl 10.29.102.173:9200/metadata -XPUT -d'{
  "mappings":{
    "properties":{
      "name":{
        "type":"keyword"
      },
      "version":{
        "type":"integer"
      },
      "size":{
        "type":"integer"
      },
      "hash":{
        "type":"text"
      }
    }
  }
}' -H'Accept: */*' -H'Content-Type: application/json'
