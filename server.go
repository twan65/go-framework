package main

import "net/http"

type Server struct {
	*router
	middlewares  []Middleware
	startHandler HandlerFunc
}

func NewServer() *Server {
	r := &router{make(map[string]map[string]HandlerFunc)}
	s := &Server{router: r}

	// 基本middlewareチェイン登録
	s.middlewares = []Middleware{
		logHandler,
		recoverHandler,
		staticHandler,
		parseFormHandler,
		parseJsonBodyHandler}

	return s
}

func (s *Server) Run(addr string) {

	// startHandlerをラウターハンドラに指定
	s.startHandler = s.router.handler()

	for i := len(s.middlewares) - 1; i >= 0; i-- {
		s.startHandler = s.middlewares[i](s.startHandler)
	}

	// Web Server起動
	if err := http.ListenAndServe(addr, s); err != nil {
		panic(err)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Context生成
	c := &Context{
		Params:         make(map[string]interface{}),
		ResponseWriter: w,
		Request:        r,
	}

	for k, v := range r.URL.Query() {
		c.Params[k] = v[0]
	}

	s.startHandler(c)
}

func (s *Server) Use(middlewares ...Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}
