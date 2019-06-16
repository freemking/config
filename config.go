package config

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Config struct {
	ConfigMap map[string]map[string]string
	node      string
}

func Load (path string) (config *Config, err error){
	config = new(Config)
	err = config.InitConfig(path)
	return
}

func (c *Config) InitConfig(path string) error {
	c.ConfigMap = make(map[string]map[string]string)
	temp_map := make(map[string]string)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		s := strings.TrimSpace(string(b))
		if strings.Index(s, "#") == 0 {
			continue
		}

		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			c.node = strings.TrimSpace(s[n1+1 : n2])
			temp_map = make(map[string]string)
			continue
		}

		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		first := strings.TrimSpace(s[:index])
		if len(first) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])

		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}
		temp_map[first] = strings.TrimSpace(second)
		c.ConfigMap[c.node] = temp_map
	}
	return nil
}

func (c *Config) Read(node, key string) string {
	v, v_found := c.ConfigMap[node][key]
	if !v_found {
		return ""
	}
	return v
}