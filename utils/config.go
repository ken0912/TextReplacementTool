package util

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Config struct {
	Mymap       map[string]string
	SectionName string
}

func (c *Config) InitConfig(path string) {
	c.Mymap = make(map[string]string)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		linestr, err := r.ReadString('\n')

		linestr = strings.TrimSpace(linestr)
		if len(linestr) == 0 {
			continue
		}
		//ignore the ones that start with '#' or ';'
		if linestr[0] == '#' || linestr[0] == ';' {
			continue
		}
		//find the SectionName
		if strings.HasPrefix(linestr, "[") && strings.HasSuffix(linestr, "]") {
			c.SectionName = linestr[1 : len(linestr)-1]
		}

		n := strings.Index(linestr, "=")
		if n < 0 {
			continue
		}
		nb := strings.TrimSpace(linestr[0:n])
		ne := strings.TrimSpace(linestr[n+1 : len(linestr)])

		if np := strings.Index(linestr, "#"); np > -1 {
			ne = strings.TrimSpace(linestr[n+1 : np])
		}
		key := c.SectionName + "=" + nb
		c.Mymap[key] = ne

		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}
	}
}

func (c *Config) Read(node, key string) string {
	key = node + "=" + key
	v, found := c.Mymap[key]
	if !found {
		return ""
	}
	return v
}
