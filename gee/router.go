/*
 * @Autor: qing fu liu
 * @Email: 1805003304@qq.com
 * @Github: https://github.com/qingfuliu
 * @Date: 2022-03-21 15:24:33
 * @LastEditors: qingfu liu
 * @LastEditTime: 2022-03-21 20:20:57
 * @FilePath: \golang\Gee\gee\router.go
 * @Description:
 */
package gee

import (
	"log"
	"strings"
)

type trieNode struct {
	isEnd, isWild bool
	value         string
	children      map[string]*trieNode
	handle        HandleFunc
	numsOfWild    int
}

func parserPath(path string) []string {
	temp := strings.Split(path, "/")
	for index := 0; index < len(temp); {
		if temp[index] == "" {
			temp = append(temp[:index], temp[index+1:]...)
		} else {
			index++
		}
	}
	return temp
}

func (n *trieNode) insert(parts []string, handle HandleFunc) {
	temp := n
	for len(parts) > 0 {
		value, ok := temp.children["*"]
		if !ok {
			for key, val := range temp.children {
				if key[0] == ':' && strings.Contains(parts[0], key) {
					ok = true
					value = val
					break
				}
			}
		}
		if !ok {
			value, ok = temp.children[parts[0]]
		}
		if !ok {
			value = &trieNode{isEnd: len(parts) == 1, isWild: parts[0] == "*" || parts[0][0] == ':', children: make(map[string]*trieNode)}
			temp.children[parts[0]] = value
		}
		parts = parts[1:]
		if value.isWild {
			value.numsOfWild++
		}
		temp = value
	}
	temp.isEnd = true
	temp.handle = handle
}

func (n *trieNode) search(path string) (HandleFunc, bool) {
	parts := parserPath(path)
	temp := n
	for len(parts) > 0 {
		value, ok := temp.children[parts[0]]
		if !ok {
			value, ok = temp.children["*"]
			if !ok {
				for key, val := range temp.children {
					if key[0] == ':' && strings.Contains(parts[0], key) {
						ok = true
						value = val
						break
					}
				}
			}
			if !ok {
				return nil, false
			}
		}
		temp = value
		parts = parts[1:]
	}
	return temp.handle, true
}

func (n *trieNode) delete(path []string) bool {
	if len(path) == 0 {
		return true
	}

	value, ok := n.children[path[0]]

	if !ok {
		return false
	}

	ok = value.delete(path[1:])

	if ok {
		if len(n.children) <= 1 && (!value.isWild || value.isWild) && value.numsOfWild == 1 {
			delete(n.children, path[1])
			return true
		} else if value.isWild {
			value.numsOfWild--
		}
	}
	return false
}

type router struct {
	trieNode *trieNode
}

func newRouter() *router {
	return &router{trieNode: &trieNode{isEnd: false, isWild: false, children: make(map[string]*trieNode), handle: nil}}
}

func (r *router) addRouter(method, pattern string, handle HandleFunc) {
	method = strings.ToUpper(method)
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "/" + pattern
	parts := parserPath(key)
	r.trieNode.insert(parts, handle)
}

func (r *router) deleteRoute(method, pattern string) {
	method = strings.ToUpper(method)
	log.Printf("Route delete %4s - %s", method, pattern)
	key := method + "/" + pattern
	parts := parserPath(key)
	r.trieNode.delete(parts)
}

func (r *router) handle(c *Context) {
	key := c.Method + "/" + c.Path
	if handler, ok := r.trieNode.search(key); ok {
		handler(c)
	} else {
		c.String(404, "404 NOT FOUND :%s", key)
	}
}
