package es

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
)

type Metadata struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
	Size    int64  `json:"size"`
	Hash    string `json:"hash"`
}

type hit struct {
	Source Metadata `json:"_source"`
}

type searchResult struct {
	Hits struct {
		Total int
		Hits  []hit
	}
}

func getMetadata(name string, versionId int) (meta Metadata, err error) {
	// "http://ES_SERVER/metadata/objects/%s_%d/_source"
	client, err := elastic.NewClient(
		elastic.SetSniff(false), elastic.SetURL(os.Getenv("ES_SERVER")))
	if err != nil {
		return
	}

	resp, err := client.Get().Index("metadata").
		Id(fmt.Sprintf("%s_%d", name, versionId)).
		Do(context.Background())
	if err != nil {
		return
	}

	err = json.Unmarshal(resp.Source, &meta)
	if err != nil {
		return
	}
	return
}

//SearchLatestVersion用来获取文件最新版本的
//元数据
func SearchLatestVersion(name string) (meta Metadata, err error) {
	client, err := elastic.NewClient(
		elastic.SetSniff(false), elastic.SetURL(os.Getenv("ES_SERVER")))
	if err != nil {
		return
	}
	query := elastic.NewTermQuery(
		"name", fmt.Sprintf("%s", url.PathEscape(name)))
	resp, err := client.Search().Index("metadata").
		Query(query).Sort("version", false).
		Size(1).Do(context.Background())
	if err != nil {
		return
	}

	var result Metadata
	if resp.Hits.TotalHits.Value == 0 {
		return
	}
	err = json.Unmarshal(resp.Hits.Hits[0].Source, &result)
	if err != nil {
		return
	}
	return
}

//GetMetadata根据文件名和版本号
//获取文件元数据，当版本号为0时
//获取最新版本的元数据
func GetMetadata(name string, version int) (Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetadata(name, version)
}

func PutMetadata(name string, version int, size int64, hash string) error {
	// request URL:
	// "http://ES_SERVER/metadata/objects/%name_%version?op_type=create"
	client, err := elastic.NewClient(
		elastic.SetSniff(false), elastic.SetURL(os.Getenv("ES_SERVER")))
	if err != nil {
		return err
	}

	meta := Metadata{
		Name:    name,
		Version: version,
		Size:    size,
		Hash:    hash,
	}

	resp, err := client.Index().Index("metadata").
		Id(fmt.Sprintf("%s_%d", name, version)).
		OpType("create").
		BodyJson(meta).
		Do(context.Background())
	if err != nil {
		if elastic.IsStatusCode(err, http.StatusConflict) {
			return PutMetadata(name, version+1, size, hash)
		}
		return fmt.Errorf("fail to put metadata: %v %v", err, resp)
	}
	return nil
}

func AddVersion(name, hash string, size int64) error {
	version, err := SearchLatestVersion(name)
	if err != nil {
		return err
	}
	return PutMetadata(name, version.Version+1, size, hash)
}

func SearchAllVersions(name string, from, size int) ([]Metadata, error) {
	client, err := elastic.NewClient(
		elastic.SetSniff(false), elastic.SetURL(os.Getenv("ES_SERVER")))
	if err != nil {
		return nil, err
	}

	request := client.Search().Index("metadata").
		Sort("name", true).
		Sort("version", false).
		From(from).Size(size)
	if name != "" {
		query := elastic.NewTermQuery("name", name)
		request = request.Query(query)
	}
	resp, err := request.Do(context.Background())
	if err != nil {
		return nil, err
	}

	var meta Metadata
	metas := make([]Metadata, 0)
	for _, item := range resp.Each(reflect.TypeOf(meta)) {
		if m, ok := item.(Metadata); ok {
			metas = append(metas, m)
		}
	}

	return metas, nil
}

func DelMetadata(name string, version int) {
	client, err := elastic.NewClient(
		elastic.SetSniff(false), elastic.SetURL(os.Getenv("ES_SERVER")))
	if err != nil {
		log.Println("error deleting metadata:", err)
		return
	}
	_, err = client.Delete().
		Index("metadata").
		Id(fmt.Sprintf("%s_%d", name, version)).
		Do(context.Background())
	if err != nil {
		log.Println("error deleting metadata:", err)
		return
	}
}

type Bucket struct {
	Key         string
	Doc_count   int
	Min_version struct {
		Value float32
	}
}

type aggregateResult struct {
	Aggregations struct {
		Group_by_name struct {
			Buckets []Bucket
		}
	}
}

func SearchVersionStatus(min_doc_count int) ([]Bucket, error) {
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/_search", os.Getenv("ES_SERVER"))
	body := fmt.Sprintf(`
        {
          "size": 0,
          "aggs": {
            "group_by_name": {
              "terms": {
                "field": "name",
                "min_doc_count": %d
              },
              "aggs": {
                "min_version": {
                  "min": {
                    "field": "version"
                  }
                }
              }
            }
          }
        }`, min_doc_count)
	request, _ := http.NewRequest("GET", url, strings.NewReader(body))
	r, e := client.Do(request)
	if e != nil {
		return nil, e
	}
	b, _ := ioutil.ReadAll(r.Body)
	var ar aggregateResult
	json.Unmarshal(b, &ar)
	return ar.Aggregations.Group_by_name.Buckets, nil
}

func HasHash(hash string) (bool, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=hash:%s&size=0", os.Getenv("ES_SERVER"), hash)
	r, e := http.Get(url)
	if e != nil {
		return false, e
	}
	b, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(b, &sr)
	return sr.Hits.Total != 0, nil
}

func SearchHashSize(hash string) (size int64, e error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=hash:%s&size=1",
		os.Getenv("ES_SERVER"), hash)
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to search hash size: %d", r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		size = sr.Hits.Hits[0].Source.Size
	}
	return
}
