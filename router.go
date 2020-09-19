package main

import (
	"net/http"
	"strings"
)

type router struct {
	// Key: httpメソッド
	// Value: URLパターン毎に実行するHandlerFunc
	// ２次元map
	handlers map[string]map[string]HandlerFunc
}

func (r *router) HandleFunc(method, pattern string, h HandlerFunc) {

	// httpメソッドとして登録されているmapがあるか確認
	m, ok := r.handlers[method]
	if !ok {
		// 新しいmapを作成
		m = make(map[string]HandlerFunc)
		r.handlers[method] = m
	}

	// httpメソッドで登録されているmapにURLパターンとハンドラ関数を登録
	m[pattern] = h
}

// WebからのリクエストのhttpメソッドとURL経路を分析して該当ハンドラを動作させる。
// func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
// 	// httpメソッドに合うhandlersを繰り返ししてリクエストURLに該当ハンドラを見つける
// 	for pattern, handler := range r.handlers[req.Method] {
// 		if ok, params := match(pattern, req.URL.Path); ok {
// 			// Context生成
// 			c := Context{
// 				Params:         make(map[string]interface{}),
// 				ResponseWriter: w,
// 				Request:        req,
// 			}

// 			for k, v := range params {
// 				c.Params[k] = v
// 			}

// 			// リクエストURLに該当ハンドラを実行
// 			handler(&c)
// 			return
// 		}
// 	}

// 	// ハンドラが見つからなかった場合NotFoundエラー処理
// 	http.NotFound(w, req)
// 	return
// }

func (r *router) handler() HandlerFunc {
	return func(c *Context) {
		// httpメソッドに合うhandlersを繰り返ししてリクエストURLに該当ハンドラを見つける
		for pattern, handler := range r.handlers[c.Request.Method] {
			if ok, params := match(pattern, c.Request.URL.Path); ok {
				for k, v := range params {
					c.Params[k] = v
				}

				handler(c)
				return
			}
		}

		http.NotFound(c.ResponseWriter, c.Request)
		return
	}
}

func match(pattern, path string) (bool, map[string]string) {

	if pattern == path {
		return true, nil
	}

	// パターンとパスを”/”単位で区分
	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")

	if len(patterns) != len(paths) {
		return false, nil
	}

	// パターンに一致するURLのパラメータを保存するためのparams mapを作成
	params := make(map[string]string)

	// ”/”で区分されているパターン/パスの各文字列を繰り返ししながら比較
	for i := 0; i < len(patterns); i++ {
		switch {
		case patterns[i] == paths[i]:
		case len(patterns[i]) > 0 && patterns[i][0] == ':':
			params[patterns[i][1:]] = paths[i]
		default:
			return false, nil
		}
	}

	return true, params
}
