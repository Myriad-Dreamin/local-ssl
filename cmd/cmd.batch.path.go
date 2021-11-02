package main

import (
	"fmt"
	"strings"
)

type ExpandItem struct {
	Mapping map[string]string `yaml:"mapping"`
	Path    string            `yaml:"path"`
}

func (e *ExpandItem) Clone() *ExpandItem {
	if e == nil {
		return &ExpandItem{
			Mapping: make(map[string]string),
			Path:    "",
		}
	}

	var c = &ExpandItem{
		Mapping: make(map[string]string),
		Path:    e.Path,
	}
	for k, v := range e.Mapping {
		c.Mapping[k] = v
	}
	return c
}

func (e *ExpandItem) Extend(path string) *ExpandItem {
	var c = e.Clone()
	c.Path = strings.Join([]string{c.Path, path}, "/")
	return c
}

func (e *ExpandItem) ExtendKV(key, path string) *ExpandItem {
	var c = e.Clone()
	c.Mapping[key] = path
	c.Path = strings.Join([]string{c.Path, path}, "/")
	return c
}

func expandPath_(e [][]string, paths []*ExpandItem) []*ExpandItem {
	if len(e) == 0 {
		return paths
	}
	x := e[0]
	switch len(x) {
	case 0:
		panic(fmt.Errorf("invalid path %v", e))
	case 1:
		for i := range paths {
			paths[i] = paths[i].Extend(x[0])
		}
	default:
		var l, pl = len(x) - 2, len(paths)
		for j, y := range x[1:] {
			if j == l {
				for i := 0; i < pl; i++ {
					paths[i] = paths[i].ExtendKV(x[0], y)
				}
			} else {
				for i := 0; i < pl; i++ {
					paths = append(paths, paths[i].ExtendKV(x[0], y))
				}
			}
		}
	}
	return expandPath_(e[1:], paths)
}

func expandPath(p [][]string) []*ExpandItem {
	var es = expandPath_(p, []*ExpandItem{nil})
	if len(es) == 1 && es[0] == nil {
		panic(fmt.Errorf("invalid expansion: %v", p))
	}
	for _, e := range es {
		e.Path = strings.TrimPrefix(e.Path, "/")
	}
	return es
}
