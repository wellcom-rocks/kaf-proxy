package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	coreruleset "github.com/corazawaf/coraza-coreruleset/v4"
	"github.com/corazawaf/coraza/v3"
	txhttp "github.com/corazawaf/coraza/v3/http"
	"github.com/corazawaf/coraza/v3/types"
	"github.com/jcchavezs/mergefs"
	"github.com/jcchavezs/mergefs/io"
)

func main() {
	waf := createWAF()

	http.Handle("/", txhttp.WrapHandler(waf, http.HandlerFunc(exampleHandler)))

	fmt.Println("Server is running. Listening port: 8090")

	log.Fatal(http.ListenAndServe(":8090", nil))
}

func createWAF() coraza.WAF {
	// directivesFile := "./default.conf"
	// if s := os.Getenv("DIRECTIVES_FILE"); s != "" {
	// 	directivesFile = s
	// }

	waf, err := coraza.NewWAF(
		coraza.NewWAFConfig().
			WithErrorCallback(logError).
			WithDirectives(`
				Include ./default.conf
                Include @owasp_crs/REQUEST-911-METHOD-ENFORCEMENT.conf
				Include @owasp_crs/REQUEST-942-APPLICATION-ATTACK-SQLI.conf
				Include @owasp_crs/REQUEST-913-SCANNER-DETECTION.conf
				Include @owasp_crs/RESPONSE-951-DATA-LEAKAGES-SQL.conf
            `).
			WithRootFS(mergefs.Merge(coreruleset.FS, io.OSFS)),
	)
	if err != nil {
		log.Fatal(err)
	}
	return waf
}

func exampleHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	resBody := "Hello world, transaction not disrupted."

	if body := os.Getenv("RESPONSE_BODY"); body != "" {
		resBody = body
	}

	if h := os.Getenv("RESPONSE_HEADERS"); h != "" {
		key, val, _ := strings.Cut(h, ":")
		w.Header().Set(key, val)
	}

	// The server generates the response
	w.Write([]byte(resBody))
}

func logError(error types.MatchedRule) {
	msg := error.ErrorLog()
	fmt.Printf("[logError][%s] %s\n", error.Rule().Severity(), msg)
}
