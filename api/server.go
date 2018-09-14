package api

import (
	"bytes"
	_ "encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"httpschedule2/engine"
	"io/ioutil"
	_ "io/ioutil"
	"net"
	"net/http"
	_ "net/http/pprof"
	_ "strconv"
	"strings"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

type HttpApiFunc func(eng *engine.Engine, w http.ResponseWriter, r *http.Request, vars map[string]string) error

func createRouter(eng *engine.Engine, enableCors bool) (*mux.Router, error) {
	r := mux.NewRouter()

	m := map[string]map[string]HttpApiFunc{
		"GET": {
			"/services":               listServices,
			"/services/{servicename}": getServices,
		},
		"POST": {
			"/selectnode": selectnode,
		},
		"OPTIONS": {
			"": optionsHandler,
		},
	}

	for method, routes := range m {
		for route, fct := range routes {
			localRoute := route
			localFct := fct
			localMethod := method

			f := makeHttpHandler(eng, localMethod, localRoute, localFct, enableCors)

			if localRoute == "" {
				r.Methods(localMethod).HandlerFunc(f)
			} else {
				r.Path(localRoute).Methods(localMethod).HandlerFunc(f)
			}
		}
	}

	addViewHandler(r)
	addPprofHandler(r)

	return r, nil

}

func ServerApi(job *engine.Job) engine.Status {

	if len(job.Args) == 0 {

		return job.Errorf("%s args cann't space", job.Name)
	}

	protoAddr := job.Args[0]

	parts := strings.SplitN(protoAddr, "://", 2)
	if len(parts) != 2 {
		return job.Errorf("%s args is wrong", job.Name)
	}

	err := ServerAndListen(parts[0], parts[1], job)

	if err != nil {
		return job.Errorf("%v", err)
	}

	return 0

}

func ServerAndListen(proto, addr string, job *engine.Job) error {
	var l net.Listener

	enableCors := job.GetEnvBool("EnableCors")

	r, err := createRouter(job.Eng, enableCors)

	if err != nil {
		return err
	}
	l, err = net.Listen(proto, addr)

	if err != nil {
		return err
	}

	//return http.ListenAndServe(addr, r)

	httpSrv := http.Server{Addr: addr, Handler: r}
	return httpSrv.Serve(l)

}

//tackle error

func makeHttpHandler(eng *engine.Engine, localMethod string, localRoute string, handlerFunc HttpApiFunc, enableCors bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if enableCors {
			writeCorsHeaders(w, r)
		}

		if err := handlerFunc(eng, w, r, mux.Vars(r)); err != nil {
			fmt.Printf("%v\n", err)
		}
	}
}

func writeCorsHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
}

func writeJson(w http.ResponseWriter, buffer *bytes.Buffer) {
	text := buffer.Bytes()
	w.Header().Set("Content-Type", "application/json")
	w.Write(text)
}

func optionsHandler(eng *engine.Engine, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	w.WriteHeader(http.StatusOK)
	return nil
}

func listServices(eng *engine.Engine, w http.ResponseWriter, r *http.Request, vars map[string]string) error {

	/*	session, err := store.Get(r, "user")

		if err != nil {
			return err
		}

		val, ok := session.Values["name"]

		if !ok {
			return fmt.Errorf("you are not login, cann't views nodes informations")
		}

		username, _ := val.(string)

		if username != "admin" {
			return fmt.Errorf("you are not admin,  cann't views nodes informations")
		}*/

	job := eng.Job("listServices")

	w.Header().Set("Content-type", "application/json")
	job.Stdout.Add(w)

	if err := job.Run(); err != nil {
		return err
	}

	return nil
}

func selectnode(eng *engine.Engine, w http.ResponseWriter, r *http.Request, vars map[string]string) error {

	/*	session, err := store.Get(r, "user")

		if err != nil {
			return err
		}

		val, ok := session.Values["name"]

		if !ok {
			return fmt.Errorf("you are not login, cann't views nodes informations")
		}

		username, _ := val.(string)

		if username != "admin" {
			return fmt.Errorf("you are not admin,  cann't views nodes informations")
		}*/

	job := eng.Job("selectnode")
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	job.SetEnv("data", string(data))
	w.Header().Set("Content-type", "application/json")
	job.Stdout.Add(w)

	if err := job.Run(); err != nil {
		return err
	}

	return nil
}
func getServices(eng *engine.Engine, w http.ResponseWriter, r *http.Request, vars map[string]string) error {

	/*	session, err := store.Get(r, "user")

		if err != nil {
			return err
		}

		val, ok := session.Values["name"]

		if !ok {
			return fmt.Errorf("you are not login, cann't views nodes informations")
		}

		username, _ := val.(string)

		if username != "admin" {
			return fmt.Errorf("you are not admin,  cann't views nodes informations")
		}*/

	serviceName := vars["servicename"]

	/*	q := r.URL.Query()

		signType := q.Get("signType")
		sign := q.Get("sign")*/

	if serviceName == "" {
		return fmt.Errorf("servicename cann't null")
	}
	/*	if signType == "" {
			return fmt.Errorf("signType cann't null")
		}
		if sign == "" {
			return fmt.Errorf("sign cann't null")
		}

		validateMap := make(map[string]string)
		validateMap["servicename"] = serviceName

		if !validateByMap(validateMap, signType, []byte(sign)) {
			return fmt.Errorf("validate fail")
		}*/

	job := eng.Job("getService")
	job.SetEnv("name", serviceName)
	job.Stdout.Add(w)

	if err := job.Run(); err != nil {
		return err
	}

	return nil
}
