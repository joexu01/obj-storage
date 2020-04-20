package versions

import (
	"encoding/json"
	"github.com/joexu01/obj-storage/lib/es"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	from := 0
	size := 100
	name := strings.Split(r.URL.EscapedPath(), `/`)[2]
	for {
		metas, err := es.SearchAllVersions(name, from, size)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for i := range metas {
			bytes, _ := json.Marshal(metas[i])
			_, _ = w.Write(bytes)
			_, _ = w.Write([]byte("\n"))
		}
		if len(metas) != size {
			return
		}
		from += size
	}
}
