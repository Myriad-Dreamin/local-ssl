package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
)

func getBatchConfig(reader io.Reader) *BatchCertsConfig {
	var decoder = yaml.NewDecoder(reader)
	decoder.SetStrict(true)
	var conf BatchCertsConfig
	if err := decoder.Decode(&conf); err != nil {
		panicHelper(err)
	}
	return &conf
}

func colorWith(s, style string, ignored bool) string {
	if ignored {
		return s
	}
	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", style, s)
}

func printEvaluatedConf(conf *BatchCertsConfig, collected []MergedCertConf, ignored bool) {
	fmt.Printf("Certs Evaluated:\n")
	for _, merged := range collected {
		fmt.Printf("[")
		for i, ks := range merged.KeyPath {
			if i == 0 {
			} else {
				fmt.Printf(", ")
			}
			if len(ks) > 1 {
				fmt.Printf("%s:[%d elems]", colorWith(ks[0], "01;34", ignored), len(ks)-1)
			} else {
				fmt.Printf("%s", ks[0])
			}
		}
		if merged.Role == "$inline" {
			fmt.Printf("] => [...&%s:%p]\n", merged.NameInConf, merged.InlineConf)
		} else {
			fmt.Printf("] => [...%s:%p, ...&%s:%p]\n", merged.Role, merged.RoleConf, merged.NameInConf, merged.InlineConf)
		}
	}
	fmt.Printf("Asset Variables:\n")
	for k, a := range conf.Assets {
		fmt.Printf("where %s => %v\n", colorWith(k, "01;34", ignored), a)
	}
	fmt.Printf("Expanded Certs:\n")
	for _, merged := range collected {
		for _, e := range merged.ExpandPath {
			fmt.Printf("- %s:\n", e.Path)
			if len(e.Mapping) != 0 {
				fmt.Printf("where: { ")
				var beg = true
				for k, v := range e.Mapping {
					if beg {
						beg = false
					} else {
						fmt.Printf(", ")
					}
					fmt.Printf("%s: %s", colorWith(k, "01;34", ignored), colorWith(v, "01;32", ignored))
				}
				fmt.Printf(" }\n")
			}
		}
	}
}
