/*
 * @Autor: qing fu liu
 * @Email: 1805003304@qq.com
 * @Github: https://github.com/qingfuliu
 * @Date: 2022-03-21 09:29:07
 * @LastEditors: qingfu liu
 * @LastEditTime: 2022-03-21 20:07:27
 * @FilePath: \golang\Gee\main.go
 * @Description:
 */
package main

import (
	"./gee"
)

func handle1(c *gee.Context) {
	c.String(200, "%s", "hhhhhh")
}
func handle2(c *gee.Context) {
	c.String(200, "%s", "mohupipei")
}
func main() {
	e := gee.New()
	group := e.Group("/v1")
	group.Get("/*", handle2)
	group.Get("/*/hhhh", handle1)
	e.Post("/p/hhhh", handle1)
	e.Run("127.0.0.1:8080")
}
