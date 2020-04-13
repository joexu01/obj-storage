package objectstream

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type TempPutStream struct {
	Server string
	Uuid string
}

func NewTempPutStream(server, object string, size int64) (*TempPutStream, error) {
	request, err := http.NewRequest(
		http.MethodPost, `http://`+server+`/temp/`+object, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("size", fmt.Sprintf("%d", size))
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	uuid, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &TempPutStream{Server: server, Uuid: string(uuid)}, nil
}

func (w *TempPutStream) Write(p []byte) (n int, err error) {
	// 对资源做局部更新
	request, err := http.NewRequest(
		http.MethodPatch, `http://`+w.Server+"/temp/"+w.Uuid, strings.NewReader(string(p)))
	if err != nil {
		return 0, err
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf(
			"data server: return http status code %d", resp.StatusCode)
	}
	return len(p), nil
}

func (w *TempPutStream) Commit(good bool) {
	method := http.MethodDelete
	if good {
		method = http.MethodPut
	}
	request, _ := http.NewRequest(method,
		`http://`+w.Server+`/temp/`+w.Uuid, nil)
	client := http.Client{}
	_, err := client.Do(request)
	if err != nil {
		log.Println("temp put stream commit error:", err)
	}
}

func NewTempGetStream(server, uuid string) (*GetStream, error) {
	return newGetStream("http://" + server + "/temp/" + uuid)
}
