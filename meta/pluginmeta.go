package meta

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type PluginMeta struct {
	Name        string
	Version     string
	URI         string
	Description string
	Author      string
	AuthorURI   string
}

var (
	comment_open   *regexp.Regexp
	comment_close  *regexp.Regexp
	plugin_name_re *regexp.Regexp
	plugin_uri_re  *regexp.Regexp
	description_re *regexp.Regexp
	author_re      *regexp.Regexp
	version_re     *regexp.Regexp
	author_uri_re  *regexp.Regexp
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
	plugin_name_re, err = regexp.Compile("\\* +Plugin Name: *(.*)")
	if err != nil {
		panic(err)
	}
	version_re, err = regexp.Compile("\\* +Version: *(.*)")
	if err != nil {
		panic(err)
	}
	plugin_uri_re, err = regexp.Compile("\\* +Plugin URI: *(.*)")
	if err != nil {
		panic(err)
	}
	description_re, err = regexp.Compile("\\* +Description: *(.*)")
	if err != nil {
		panic(err)
	}
	author_re, err = regexp.Compile("\\* +Author: *(.*)")
	if err != nil {
		panic(err)
	}
	author_uri_re, err = regexp.Compile("\\* +Author URI: *(.*)")
	if err != nil {
		panic(err)
	}
}

func New() *PluginMeta {
	return &PluginMeta{}
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

func (pm *PluginMeta) ParseMetaLine(line string) error {
	match := plugin_name_re.FindStringSubmatch(line)
	if match != nil {
		pm.Name = strings.TrimRight(match[1], " \t")
		return nil
	}

	match = version_re.FindStringSubmatch(line)
	if match != nil {
		pm.Version = strings.TrimRight(match[1], " \t")
		return nil
	}

	match = plugin_uri_re.FindStringSubmatch(line)
	if match != nil {
		pm.URI = strings.TrimRight(match[1], " \t")
		return nil
	}

	match = description_re.FindStringSubmatch(line)
	if match != nil {
		pm.Description = strings.TrimRight(match[1], " \t")
		return nil
	}

	match = author_re.FindStringSubmatch(line)
	if match != nil {
		pm.Author = strings.TrimRight(match[1], " \t")
		return nil
	}

	match = author_uri_re.FindStringSubmatch(line)
	if match != nil {
		pm.AuthorURI = strings.TrimRight(match[1], " \t")
		return nil
	}

	return fmt.Errorf("unknown meta field: %s", line)
}
