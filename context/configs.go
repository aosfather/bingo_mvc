package context

import (
	"fmt"
	log "github.com/aosfather/bingo_utils"
	"github.com/aosfather/bingo_utils/files"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type YamlConfig struct {
	config map[interface{}]interface{}
}

func (this *YamlConfig) LoadFromFile(file string) {
	if file != "" && files.IsFileExist(file) {
		f, err := os.Open(file)
		if err == nil {
			txt, _ := ioutil.ReadAll(f)
			err := yaml.Unmarshal(txt, &this.config)
			if err != nil {
				log.Err(err.Error())
				panic("load config file error!")
			}
		}

	}
	if this.config == nil {
		this.config = make(map[interface{}]interface{})
	}
}

//不能获取bingo自身的属性，只能获取应用自身的扩展属性
func (this *YamlConfig) GetPropertyForCustom(key string) string {
	if strings.HasPrefix(key, "bingo.") {
		return ""
	}
	return this.GetProperty(key)
}

func (this *YamlConfig) GetProperty(key string) string {
	if this.config == nil {
		return ""
	}
	v, ok := this.config[key]
	if ok {
		return v.(string)
	} else {
		if strings.Index(key, ".") > 0 {
			keys := strings.Split(key, ".")
			return this.getvalue(this.config, keys, 0)
		}
	}
	return ""
}

func (this *YamlConfig) getvalue(m map[interface{}]interface{}, keys []string, index int) string {
	v, ok := m[keys[index]]
	if ok {
		if value, ok := v.(string); ok {
			return value
		}

		if value, ok := v.(map[interface{}]interface{}); ok {
			return this.getvalue(value, keys, index+1)
		}
		if v != nil {
			return fmt.Sprintf("%v", v)
		}

	}
	return ""
}
