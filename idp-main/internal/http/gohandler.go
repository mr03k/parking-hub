package http

import (
	"net/http"
	"reflect"

	"application/internal/http/handler"

	_ "application/docs/swagger"

	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/swgui/v5emb"
)

func NewHTTPHandler(r *openapi3.Reflector, o handler.OAPI, svcs ...handler.Handler) (http.Handler, error) {
	mux := http.NewServeMux()

	// for _, svc := range svcs {
	// 	svc.RegisterMuxRouter(mux)
	// }

	for _, svc := range svcs {

		if _, err := r.Reflect(svc); err == nil {
			svc.RegisterMuxRouter(mux)
		} else {
			return nil, err
		}

		if reflect.TypeOf(svc).Implements(reflect.TypeFor[handler.OpenApiHandler]()) {
			svc.(handler.OpenApiHandler).OpenApiSpec(o)
		}

	}

	// wd, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// swagDir := path.Join(wd, "docs", "swagger")
	// _, err = os.ReadDir(swagDir)
	// if err == nil {
	// 	fs := http.FileServer(http.Dir(swagDir))

	// 	mux.Handle("/docs/swagger/", http.StripPrefix("/docs/swagger/", fs))

	// 	mux.Handle("/swagger/", v5emb.New(
	// 		"swagger",
	// 		"/docs/swagger/swagger.json",
	// 		"/swagger/",
	// 	))

	// } else {
	// }

	sw := NewSwagger(o)
	mux.HandleFunc("/docs/swagger/swagger.json", sw.swagerjson)

	mux.Handle("/swagger/", v5emb.New(
		"swagger",
		"/docs/swagger/swagger.json",
		"/swagger/",
	))

	return mux, nil
}

type Swagger struct {
	OpenAPI handler.OAPI
}

func NewSwagger(o handler.OAPI) *Swagger {
	return &Swagger{
		OpenAPI: o,
	}
}

func (s *Swagger) swagerjson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json, err := s.OpenAPI.GetJsonData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json) //nolint // ignore error
}
