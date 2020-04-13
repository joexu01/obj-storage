package objectstream

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
)

type GetStream struct {
	reader io.Reader
}

type PutStream struct {
	writer *io.PipeWriter
	c      chan error
}

// Get Stream

func newGetStream(url string) (*GetStream, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"data server: return http status code %d", resp.StatusCode)
	}

	return &GetStream{resp.Body}, nil
}

func NewGetStream(server, object string) (*GetStream, error) {
	if server == "" || object == "" {
		return nil, fmt.Errorf("invalid server %s object %s", server, object)
	}
	return newGetStream(`http://` + server + `/objects/` + object)
}

func (r *GetStream) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

// Put Stream

func NewPutStream(server, object string) *PutStream {
	reader, writer := io.Pipe()
	c := make(chan error)
	go func() {
		request, err := http.NewRequest(http.MethodPut, `http://`+server+`/objects/`+object, reader)
		if err != nil {
			c <- err
			runtime.Gosched()  // Debug时调用了一下，一般new request不会出错
		}
		client := http.Client{}
		response, err := client.Do(request)
		if err == nil && response.StatusCode != http.StatusOK {
			// 状态码非200也是error
			err = fmt.Errorf("data server: return http code %d", response.StatusCode)
		}
		c <- err
	}()
	return &PutStream{writer: writer, c: c}
}

func (w *PutStream) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

// Close用于关闭writer Close使得管道另一边的
// reader读到文件尾，否则另一goroutine中的
// client.Do() 始终阻塞 最后从c中读取错误
func (w *PutStream) Close() error {
	_ = w.writer.Close()
	return <- w.c
}
