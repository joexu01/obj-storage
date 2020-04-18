package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

type object struct {
	Name    string `json:"name"`
	Teacher string `json:"teacher"`
}

func main() {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	result, err := client.Search().Index("object").
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", result.Hits.Hits[0].Source)
	fmt.Printf("%s", result.Hits.Hits[1].Source)
	//result, err := client.Get().Index("object").Id("1").Do(context.Background())
	//if err != nil {
	//	panic(err)
	//}
	//
	//var store object
	//
	//err = json.Unmarshal(result.Source, &store)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(store)
	//fmt.Printf("%v\n", result)
	//
	//obj := object{
	//	Name:    "Python",
	//	Teacher: "abc",
	//}
	//
	//resp, err := client.Index().Index("object").BodyJson(obj).Do(context.Background())
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%+v\n", resp)
}
