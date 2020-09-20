package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
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
