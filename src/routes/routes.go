package routes

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HelloWorldHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hello World!")
}
