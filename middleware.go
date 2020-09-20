package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

type Middleware func(next HandlerFunc) HandlerFunc

func logHandler(next HandlerFunc) HandlerFunc {
	return func(c *Context) {

		t := time.Now()

		// 次のハンドラを実行
		next(c)

		// Webリクエスト情報と全体所要時間を残す
		log.Printf("[%s] %q %v¥n",
			c.Request.Method,
			c.Request.URL.String(),
			time.Now().Sub(t))
	}
}

func recoverHandler(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(c.ResponseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next(c)
	}
}

// PostリクエストのFormデータをContextのParamsに保存
func parseFormHandler(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		c.Request.ParseForm()

		fmt.Println(c.Request.PostForm)

		for k, v := range c.Request.PostForm {
			if len(v) > 0 {
				c.Params[k] = v[0]
			}
		}

		next(c)
	}
}

// JSONデータ解析してContextのParamsに保存
func parseJsonBodyHandler(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		var m map[string]interface{}

		if json.NewDecoder(c.Request.Body).Decode(&m); len(m) > 0 {
			for k, v := range m {
				c.Params[k] = v
			}
		}

		next(c)
	}
}

func staticHandler(next HandlerFunc) HandlerFunc {
	var (
		dir       = http.Dir(".")
		indexFile = "index.html"
	)
	return func(c *Context) {
		// httpメソッドがGETかHEADではない場合は次のハンドラを実行する。
		if c.Request.Method != "GET" && c.Request.Method != "HEAD" {
			next(c)
			return
		}

		file := c.Request.URL.Path
		// URL経路に該当のファイルを開く
		f, err := dir.Open(file)
		if err != nil {
			// 失敗時、次のハンドラを実行する。
			next(c)
			return
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			// ファイルが正常ではない場合、次のハンドラを実行
			next(c)
			return
		}

		// URL経路がディレクトリの場合はindexFileを利用
		if fi.IsDir() {
			// ディレクトリ経路をURLで使う時は最後に”/”を付けなければならない
			if !strings.HasSuffix(c.Request.URL.Path, "/") {
				http.Redirect(c.ResponseWriter, c.Request, c.Request.URL.Path+"/", http.StatusFound)
				return
			}

			file = path.Join(file, indexFile)

			// indexFileを開く
			f, err = dir.Open(file)
			if err != nil {
				next(c)
				return
			}
			defer f.Close()

			fi, err = f.Stat()
			if err != nil || fi.IsDir() {
				next(c)
				return
			}
		}

		// fileの内容を渡す。(nextハンドラの制御権は渡さずリクエスト処理を終了)
		http.ServeContent(c.ResponseWriter, c.Request, file, fi.ModTime(), f)
	}
}
