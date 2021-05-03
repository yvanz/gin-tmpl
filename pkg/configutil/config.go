package configutil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"

	"github.com/uber/kraken/utils/stringset"
	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v3"
)

var ErrCycleRef = fmt.Errorf("cyclic reference in configuration extends detected")

type Extends struct {
	Extends string `yaml:"extends"`
}

type ValidationError struct {
	errorMap validator.ErrorMap
}

func (e ValidationError) ErrForField(name string) error {
	return e.errorMap[name]
}

func (e ValidationError) Error() string {
	var w bytes.Buffer

	fmt.Fprintf(&w, "validation failed")
	for f, err := range e.errorMap {
		fmt.Fprintf(&w, "   %s: %v\n", f, err)
	}

	return w.String()
}

func Load(filename string, config interface{}) error {
	if filename == "" {
		return fmt.Errorf("no configuration file is specified")
	}
	filenames, err := resolveExtends(filename, readExtend)
	if err != nil {
		return err
	}
	return loadFiles(config, filenames)
}

type getExtend func(filename string) (extends string, err error)

func resolveExtends(filename string, extendReader getExtend) ([]string, error) {
	filenames := []string{filename}
	seen := make(stringset.Set)
	for {
		extends, err := extendReader(filename)
		if err != nil {
			return nil, err
		} else if extends == "" {
			break
		}

		if !filepath.IsAbs(extends) {
			extends = path.Join(filepath.Dir(filename), extends)
		}

		if seen.Has(extends) {
			return nil, ErrCycleRef
		}

		filenames = append([]string{extends}, filenames...)
		seen.Add(extends)
		filename = extends
	}
	return filenames, nil
}

func readExtend(configFile string) (string, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return "", err
	}

	var cfg Extends
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return "", fmt.Errorf("unmarshal %s: %s", configFile, err)
	}
	return cfg.Extends, nil
}

func loadFiles(config interface{}, fnames []string) error {
	for _, fname := range fnames {
		data, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}

		if err := yaml.Unmarshal(data, config); err != nil {
			return fmt.Errorf("unmarshal %s: %s", fname, err)
		}
	}

	if err := validator.Validate(config); err != nil {
		return ValidationError{
			errorMap: err.(validator.ErrorMap),
		}
	}
	return nil
}
