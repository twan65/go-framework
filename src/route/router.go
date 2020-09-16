package route

import (
	"net/http"
	"strings"
)

type Router struct {
	// Key: httpメソッド
	// Value: URLパターン毎に実行するHandlerFunc
	// ２次元map
	Handlers map[string]map[string]http.HandlerFunc
}

func (r *Router) HandleFunc(method, pattern string, h http.HandlerFunc) {

	// httpメソッドとして登録されているmapがあるか確認
	m, ok := r.Handlers[method]
	if !ok {
		// 新しいmapを作成
		m = make(map[string]http.HandlerFunc)
		r.Handlers[method] = m
	}

	// httpメソッドで登録されているmapにURLパターンとハンドラ関数を登録
	m[pattern] = h
}

// WebからのリクエストのhttpメソッドとURL経路を分析して該当ハンドラを動作させる。
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// httpメソッドに合うhandlersを繰り返ししてリクエストURLに該当ハンドラを見つける
	for pattern, handler := range r.Handlers[req.Method] {
		if ok, _ := match(pattern, req.URL.Path); ok {
			// リクエストURLに該当ハンドラを実行
			handler(w, req)
			return
		}
	}

	// ハンドラが見つからなかった場合NotFoundエラー処理
	http.NotFound(w, req)
	return
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
