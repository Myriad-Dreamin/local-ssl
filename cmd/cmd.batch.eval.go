package main

import (
	"fmt"
	"strings"
)

type MergedCertConf struct {
	RoleConf   *CertConfBase `yaml:"roleConf"`
	InlineConf *CertConfBase `yaml:"inlineConf"`
	CertConf   CertConfBase  `yaml:"conf"`
	RawKey     string        `yaml:"rawKey"`
	Empty      bool          `yaml:"emptyAssets"`
	KeyPath    [][]string    `yaml:"keyPath"`
	ExpandPath []*ExpandItem `yaml:"expandPath"`
	Role       string        `yaml:"role"`
	NameInConf string        `yaml:"rawName"`
	Name       string        `yaml:"name"`
}

type MergedCertConfCmpImpl []MergedCertConf

func (cmp MergedCertConfCmpImpl) Less(i, j int) bool {
	if cmp[i].RawKey != cmp[j].RawKey {
		return cmp[i].RawKey < cmp[j].RawKey
	}
	if cmp[i].Name != cmp[j].Name {
		return cmp[i].Name < cmp[j].Name
	}
	if cmp[i].Role != cmp[j].Role {
		return cmp[i].Role < cmp[j].Role
	}
	return false
}

func (cmp MergedCertConfCmpImpl) Len() int {
	return len(cmp)
}

func (cmp MergedCertConfCmpImpl) Swap(i, j int) {
	cmp[i], cmp[j] = cmp[j], cmp[i]
}

func evaluateBatchConfig(conf *BatchCertsConfig) (collected []MergedCertConf, errors []error) {
	for k, vs := range conf.Certs {
		for _, v := range vs {
			var merged = MergedCertConf{
				RawKey:     k,
				NameInConf: v.Name,
				Name:       v.Name,
				Role:       v.Role,
				InlineConf: &v.CertConf,
				RoleConf:   nil,
			}
			if len(merged.NameInConf) == 0 {
				merged.Name = merged.Role
			}
			for _, p := range strings.Split(merged.RawKey, "/") {
				if strings.HasPrefix(p, "$") {
					p = p[1:]
					if _, ok := conf.Assets[p]; !ok {
						errors = append(errors, fmt.Errorf("asset name not found %s", p))
						merged.KeyPath = append(merged.KeyPath, []string{p})
					} else {
						var a = conf.Assets[p]
						if len(a) == 0 {
							merged.Empty = true
						}
						merged.KeyPath = append(merged.KeyPath, append([]string{p}, a...))
					}
				} else {
					merged.KeyPath = append(merged.KeyPath, []string{p})
				}
			}
			if !merged.Empty {
				merged.ExpandPath = expandPath(merged.KeyPath)
			}
			switch merged.Role {
			case "$inline":
				break
			case "":
				merged.Role = "$default"
				fallthrough
			default:
				if len(conf.Roles) != 0 {
					if role, ok := conf.Roles[merged.Role]; ok {
						merged.RoleConf = &role.CertConf
					} else {
						errors = append(errors, fmt.Errorf("role not found %s", merged.Role))
					}
				}
			}
			collected = append(collected, merged)
		}
	}

	return
}
