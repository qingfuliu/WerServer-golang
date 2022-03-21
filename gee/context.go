/*
 * @Autor: qing fu liu
 * @Email: 1805003304@qq.com
 * @Github: https://github.com/qingfuliu
 * @Date: 2022-03-21 15:09:36
 * @LastEditors: qingfu liu
 * @LastEditTime: 2022-03-21 21:13:31
 * @FilePath: \golang\Gee\gee\context.go
 * @Description:
 */
package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer       http.ResponseWriter
	Req          *http.Request
	Path, Method string
	StatusCode   int
	index        int
	handles      []HandleFunc
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{Writer: w, Req: r, Path: r.URL.Path, Method: r.Method, index: -1, handles: make([]HandleFunc, 0)}
}
func (c *Context) next() {
	s := len(c.handles)
	for ; c.index < s; c.index++ {
		c.handles[c.index](c)
	}
}
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) Json(code int, jsonData interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "application/json")
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(jsonData); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
