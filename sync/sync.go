package sync

import (
	"errors"
	"os"
	"sort"
	"strings"

	"github.com/y1rn/nsync/neko"
	"gopkg.in/yaml.v3"
)

func Sync(url, token, configPath, nodePrefix string) error {

	configs, err := loadConfig(configPath)
	if err != nil {
		return err
	}

	ruleResp, err := neko.LoadRule(url, token)
	if err != nil {
		return err
	}
	if len(ruleResp.Data) == 0 {
		return errors.New("no neko rule data")
	}
	for _, config := range configs {
		err = sync(config, ruleResp.Data, nodePrefix)
		if err != nil {
			return err
		}
	}

	return nil
}

func sync(config Config, data []neko.Data, nodePrefix string) error {
	proxies := []map[string]any{}
	pg := ProxyGroup{
		Name:    config.Group.Name,
		Type:    config.Group.Type,
		Proxies: config.Group.Proxies,
		Url:     config.Group.Url,
	}
	if pg.Proxies == nil {
		pg.Proxies = []string{}
	}

	for _, proxy := range config.Proxies {
		port := proxy[KEY_PORT]
		// t := proxy[KEY_TYPE]

		for _, d := range data {
			if d.Rport == port {
				proxyName := nodePrefix + "-" + d.Name
				proxies = append(proxies, newProxy(proxy, d, proxyName))

				//append proxy to group
				exist := false
				for _, pName := range pg.Proxies {
					if pName == proxyName {
						exist = true
					}
				}
				if !exist {
					pg.Proxies = append(pg.Proxies, proxyName)
				}
			}
		}

	}

	//sort group proxies by name
	sort.Strings(pg.Proxies)

	temp, err := config.LoadTemplate()
	if err != nil {
		return err
	}

	temp[KEY_PROXIES] = proxyAppend(temp[KEY_PROXIES].([]any), proxies, KEY_NAME, nodePrefix)

	//replace template group
	pgs := temp[KEY_GROUP].([]any)
	exist := false

	for i, g := range pgs {
		tg := g.(map[string]any)
		if tg[KEY_NAME] == pg.Name {
			pgs[i] = pg
			exist = true
			break
		}
	}
	if !exist {
		pgs = append(pgs, pg)
	}

	temp[KEY_GROUP] = pgs

	return writeYaml(config.OutPut, temp)
}

func newProxy(p map[string]any, d neko.Data, proxyName string) map[string]any {
	proxy := make(map[string]any)
	for k, v := range p {
		proxy[k] = v
	}
	proxy[KEY_NAME] = proxyName
	proxy[KEY_PORT] = d.Port
	proxy[KEY_SERVER] = d.Host
	return proxy
}

func proxyAppend(dest []any, src []map[string]any, key, nodePrefix string) []any {
	d := []any{}
	for _, p := range dest {
		proxy := p.(map[string]any)
		name := proxy[key].(string)
		if !strings.HasPrefix(name, nodePrefix) {
			d = append(d, p)
		}
	}

	for _, p := range src {
		d = append(d, p)
	}
	return d
}

func writeYaml(path string, v any) error {
	file, err := os.Create(path)
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		return err
	}

	en := yaml.NewEncoder(file)
	en.SetIndent(1)

	return en.Encode(v)
}
