package googlephotos

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
)

func LoginServer(data chan *string, httpListenPort int, redirectPath string) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", httpListenPort))
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc(redirectPath, func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		writer.Write([]byte("    <script>\n      window.close();\n    </script>"))
		queryParts, _ := url.ParseQuery(request.URL.RawQuery)

		// Use the authorization code that is pushed to the redirect
		// URL.
		code := queryParts["code"][0]

		data <- &code

		go func() {
			l.Close()
			waitGroup.Wait()
		}()
	})

	go func() {
		defer waitGroup.Done()

		log.Printf("Start listening on http://0.0.0.0:%d%s for Google Login Response", httpListenPort, redirectPath)
		if err := http.Serve(l, nil); err != http.ErrServerClosed {
			//log.Panicf("Failed to start Server for local redirect with error %v", err)
		}
	}()
}
