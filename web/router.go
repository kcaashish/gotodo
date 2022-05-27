package web

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var routes = []route{
		newRoute("GET", "/users/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})", s.getUser()),
		newRoute("GET", "/users", s.getUsers()),
		newRoute("POST", "/users", s.createUser()),
		newRoute("PATCH", "/users/update", s.updateUser()),
		newRoute("DELETE", "/users/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})", s.deleteUser()),

		newRoute("GET", "/todo/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})", s.getTodoList()),
		newRoute("GET", "/todo", s.getTodoLists()),
		newRoute("POST", "/todo", s.createTodoList()),
		newRoute("PATCH", "/todo/update", s.updateTodoList()),
		newRoute("DELETE", "/todo/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})", s.deleteTodoList()),
	}

	var allow []string

	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}

	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

type ctxKey struct{}

func getField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}
