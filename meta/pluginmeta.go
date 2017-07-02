package meta

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

type PluginMeta struct {
	Fields map[string]string
}

var (
	comment_open  *regexp.Regexp
	comment_close *regexp.Regexp
	meta_re       *regexp.Regexp
)

func init() {
	var err error

	comment_open, err = regexp.Compile(" */\\*+")
	if err != nil {
		panic(err)
	}
	comment_close, err = regexp.Compile(".*\\*/")
	if err != nil {
		panic(err)
	}
	meta_re, err = regexp.Compile("\\*\\s*(.*?)\\s*:\\s*(.*?)\\s*$")
	if err != nil {
		panic(err)
	}
}

func New() *PluginMeta {
	return &PluginMeta{map[string]string{}}
}

func normalizeKey(key string) string {
	return strings.Title(strings.ToLower(key))
}

func (pm *PluginMeta) Get(key string) string {
	return pm.Fields[normalizeKey(key)]
}

func (pm *PluginMeta) Set(key, value string) {
	pm.Fields[normalizeKey(key)] = value
}

func (meta *PluginMeta) Scan(in io.Reader) error {
	open := false

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		if !open {
			open = comment_open.MatchString(scanner.Text())
		} else {
			line := scanner.Text()
			if comment_close.MatchString(line) {
				open = false
			} else {
				meta.ParseMetaLine(line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (pm *PluginMeta) ParseMetaLine(line string) bool {
	match := meta_re.FindStringSubmatch(line)
	if match != nil {
		pm.Set(match[1], match[2])
	}

	return match != nil
}
