package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/go-ini/ini"
)

var (
	root    = flag.String("root", "/etc/netdata", "the root folder to search for file")
	prefix  = flag.String("prefix", "ND_", "the prefix for config loading")
	allKeys = make(map[string][]*key)
)

type key struct {
	file    string
	section string
	key     string
	value   string
}

func newKey(env string) *key {
	parts := strings.SplitN(env, "|", 2)
	if len(parts) != 2 {
		return nil
	}
	// first split for section
	res := key{file: strings.Trim(parts[0], "\t ")}
	parts = strings.SplitN(parts[1], "/", 2)
	if len(parts) != 2 {
		return nil
	}
	res.section = strings.Trim(parts[0], " \t")
	kv := strings.SplitN(parts[1], "=", 2)
	if len(kv) != 2 {
		return nil
	}
	res.key = strings.Trim(kv[0], " \t")
	res.value = strings.Trim(kv[1], " \t")

	return &res
}

func main() {
	flag.Parse()

	all := os.Environ()
	for i := range all {
		tmp := strings.SplitN(all[i], "=", 2)
		if len(tmp) != 2 {
			continue
		}
		key := strings.Trim(tmp[0], " \t")
		if strings.HasPrefix(key, *prefix) {
			k := newKey(strings.Trim(tmp[1], " \n"))
			if k == nil {
				log.Printf("The key is invalid : %s", all[i])
				continue
			}
			allKeys[k.file] = append(allKeys[k.file], k)
		}
	}

	for i, v := range allKeys {
		target := filepath.Join(*root, i)
		f, err := ini.Load(target)
		if err != nil {
			log.Printf("file not found, create an empty one: %s", target)
			f = ini.Empty()
		}

		for k := range v {
			sec, err := f.NewSection(v[k].section)
			if err != nil {
				log.Fatal(err)
			}
			_, err = sec.NewKey(v[k].key, v[k].value)
			if err != nil {
				log.Fatal(err)
			}
		}
		if err := f.SaveTo(target); err != nil {
			log.Fatal(err)
		}
	}
}
