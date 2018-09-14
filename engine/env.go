package engine

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Env []string

func (env *Env) Get(key string) (value string) {
	for _, kv := range *env {
		if strings.Index(kv, "=") == -1 {
			continue
		}

		parts := strings.SplitN(kv, "=", 2)
		if parts[0] != key {
			continue
		} else {
			value = parts[1]
		}
	}
	return
}

func (env *Env) Set(key, value string) {
	*env = append(*env, key+"="+value)
}

func (env *Env) GetBool(key string) bool {

	s := strings.ToLower(strings.Trim(env.Get(key), " \t"))
	if s == "" || s == "0" || s == "no" || s == "false" || s == "none" {
		return false
	}
	return true
}

func (env *Env) SetBool(key string, value bool) {
	if value {
		env.Set(key, "1")
	} else {
		env.Set(key, "0")
	}
}

func (env *Env) SetInt(key string, value int) {
	env.Set(key, fmt.Sprintf("%d", value))
}

func (env *Env) SetInt64(key string, value int64) {
	env.Set(key, fmt.Sprintf("%d", value))
}

func (env *Env) GetInt64(key string) int64 {
	s := strings.Trim(env.Get(key), " \t")
	val, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		return 0
	}
	return val
}

func (env *Env) GetInt(key string) int {
	return int(env.GetInt64(key))
}

func (env *Env) GetJson(name string, iface interface{}) error {
	sval := env.Get(name)

	if sval == "" {
		return nil
	}
	return json.Unmarshal([]byte(sval), iface)
}

func (env *Env) SetJson(name string, value interface{}) error {
	sval, err := json.Marshal(value)
	if err != nil {
		return err
	}

	env.Set(name, string(sval))
	return nil
}

func (env *Env) Decode(src io.Reader) error {
	m := make(map[string]interface{})

	if err := json.NewDecoder(src).Decode(&m); err != nil {
		return err
	}

	for k, v := range m {
		env.SetAuto(k, v)
	}

	return nil
}

func (env *Env) SetAuto(k string, v interface{}) {
	if val, ok := v.(string); ok {
		env.Set(k, val)
	} else if val, err := json.Marshal(v); err != nil {
		env.Set(k, string(val))
	} else {
		env.Set(k, fmt.Sprintf("%v", v))
	}
}
