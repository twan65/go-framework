package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"path/filepath"
	"text/template"
)

type Context struct {
	Params map[string]interface{}

	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

type HandlerFunc func(*Context)

// コンテキストにJSONフォーマットにデータをレンダリング
func (c *Context) RenderJson(v interface{}) {
	// HTTP StatusをStatusOKで指定
	c.ResponseWriter.WriteHeader(http.StatusOK)

	// ContentーTypeをapplication/jsonにする。
	c.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")

	// vをjsonで出力
	if err := json.NewEncoder(c.ResponseWriter).Encode(v); err != nil {
		c.RenderErr(http.StatusInternalServerError, err)
	}
}

func (c *Context) RenderXml(v interface{}) {
	c.ResponseWriter.WriteHeader(http.StatusOK)
	c.ResponseWriter.Header().Set("Content-Type", "application/xml; charset=utf-8")

	if err := xml.NewEncoder(c.ResponseWriter).Encode(v); err != nil {
		c.RenderErr(http.StatusInternalServerError, err)
	}
}

// エラー状態を適切なHTTP STATUSにレンダリング
func (c *Context) RenderErr(code int, err error) {
	if err != nil {
		if code > 0 {
			// 正常なコードの場合
			http.Error(c.ResponseWriter, http.StatusText(code), code)
		} else {
			// 正常のコードではない場合
			defaultErr := http.StatusInternalServerError
			http.Error(c.ResponseWriter, http.StatusText(defaultErr), defaultErr)
		}
	}
}

// テンプレートオブジェクトを保存するためのmap
var templates = map[string]*template.Template{}

func (c *Context) RenderTemplate(path string, v interface{}) {
	// pathに該当テンプレートがあるかをチェック
	t, ok := templates[path]
	if !ok {
		// テンプレートオブジェクトを生成
		t = template.Must(template.ParseFiles(filepath.Join(".", path)))
		templates[path] = t
	}

	// 最終結果をResponseWriterに出力
	t.Execute(c.ResponseWriter, v)
}

func (c *Context) Redirect(url string) {
	http.Redirect(c.ResponseWriter, c.Request, url, http.StatusMovedPermanently)
}
