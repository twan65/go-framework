package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type User struct {
	Id        string
	AddressId string
}

func main() {
	s := NewServer()

	s.HandleFunc("GET", "/", func(c *Context) {
		c.RenderTemplate("/public/index.html",
			map[string]interface{}{"time": time.Now()})
	})

	s.HandleFunc("GET", "/about", func(c *Context) {
		fmt.Fprintln(c.ResponseWriter, "about")
	})

	s.HandleFunc("GET", "/users/:id", func(c *Context) {
		u := User{Id: c.Params["id"].(string)}
		c.RenderXml(u)
	})

	s.HandleFunc("GET", "/users/:user_id/addresses/:address_id", func(c *Context) {
		u := User{c.Params["user_id"].(string), c.Params["address_id"].(string)}
		c.RenderJson(u)
	})

	s.HandleFunc("POST", "/users", func(c *Context) {
		c.RenderJson(c.Params)
	})

	s.HandleFunc("GET", "/login", func(c *Context) {
		// "login.html"レンダリング
		c.RenderTemplate("/public/login.html",
			map[string]interface{}{"message": "ログインが必要です"})
	})

	s.HandleFunc("POST", "/login", func(c *Context) {
		// ログイン情報を確認してクッキーに認証トークンを保存
		if CheckLogin(c.Params["username"].(string), c.Params["password"].(string)) {
			http.SetCookie(c.ResponseWriter, &http.Cookie{
				Name:  "X_AUTH",
				Value: Sign(VerifyMessage),
				Path:  "/",
			})
			c.Redirect("/")
		}
		// IDとPWが合致しないとid"/login"をレンダリング
		c.RenderTemplate("/public/login.html",
			map[string]interface{}{"message": "idまたはpasswordを確認してください。"})

	})

	s.Use(AuthHandler)

	s.Run(":8080")
}

const VerifyMessage = "verified"

func AuthHandler(next HandlerFunc) HandlerFunc {
	ignore := []string{"/login", "public/index.html"}
	return func(c *Context) {
		// URL prefixが"/login", "public/index.html"だったらautをチェックしない
		for _, s := range ignore {
			if strings.HasPrefix(c.Request.URL.Path, s) {
				next(c)
				return
			}
		}

		if v, err := c.Request.Cookie("X_AUTH"); err == http.ErrNoCookie {
			// "X_AUTH"クッキーの値がない場合"/login"
			c.Redirect("/login")
			return
		} else if err != nil {
			// エラー処理
			c.RenderErr(http.StatusInternalServerError, err)
			return
		} else if Verify(VerifyMessage, v.Value) {

			next(c)
			return
		}

		c.Redirect("/login")
	}
}

func CheckLogin(username, password string) bool {
	// ログイン処理
	const (
		USERNAME = "twan"
		PASSWORD = "12345"
	)

	return username == USERNAME && password == PASSWORD
}

// 認証トークン確認
func Verify(message, sig string) bool {
	return hmac.Equal([]byte(sig), []byte(Sign(message)))
}

// 認証トークン生成
func Sign(message string) string {
	secretKey := []byte("golang-book-secret-key2")
	if len(secretKey) == 0 {
		return ""
	}
	mac := hmac.New(sha1.New, secretKey)
	io.WriteString(mac, message)
	return hex.EncodeToString(mac.Sum(nil))
}
