package featurex

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
)

type (
	SubConfig    map[string]string
	ConfigLoader struct {
		yamlConfig map[string]SubConfig
	}
)

var configFile = flag.String("fx", "./config/app.local.yaml", "configuration file")

func NewConfigLoader() *ConfigLoader {
	flag.Parse()
	cl := &ConfigLoader{}
	cl.yamlConfig = make(map[string]SubConfig)
	maybeReadConfig(*configFile, cl.yamlConfig)
	return cl
}

// Load config from yaml and env
func (cl *ConfigLoader) Load(cfg interface{}) {
	cl.loadConfig(cfg)
}

func maybeReadConfig(filepath string, config map[string]SubConfig) {
	if _, err := os.Stat(filepath); err != nil { // abort if file does not exist
		return
	}
	fp, err := os.Open(filepath)
	if err != nil {
		logx.Info("Failed to open app config file")
		return
	} else {
		defer fp.Close()
	}

	bytes, err := ioutil.ReadAll(fp)
	if err != nil {
		logx.Error("Failed to read config file", "error", err.Error())
	}
	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		logx.Error("Failed to unmarshal config file", "error", err.Error())
	}
}

func (cl *ConfigLoader) loadConfig(featureConfig interface{}) {
	val := reflect.ValueOf(featureConfig)
	elm := reflect.Indirect(val)
	loadConfigTree(featureConfig, cl.yamlConfig[strings.ToLower(elm.Type().Name())])
}

func loadConfigTree(configTree interface{}, yamlSubConfig SubConfig) {
	val := reflect.ValueOf(configTree)
	elm := reflect.Indirect(val)

	if elm.Kind() != reflect.Struct {
		panic(fmt.Errorf("expected a <*struct> of %v", val))
	}

	typ := elm.Type()

	// loop fields of each struct field
	for i := 0; i < typ.NumField(); i++ {
		field := elm.Field(i)
		structField := typ.Field(i)

		if !field.CanSet() {
			continue
		}

		envTag := structField.Tag.Get("env")
		if value := os.Getenv(envTag); value != "" {
			setField(field, value)
			continue
		}
		yamlTag := structField.Tag.Get("yaml")
		if yamlTag != "" {
			if yamlSubConfig == nil {
				continue
			}
			if value, ok := yamlSubConfig[yamlTag]; ok {
				setField(field, value)
			}
			continue
		}

		// go deeper
		if field.CanInterface() {
			if field.Kind() != reflect.Struct {
				continue
			}
			// restore to pointer struct
			field = field.Addr()
			loadConfigTree(field.Interface(), yamlSubConfig)
		}
	}
}

func setField(field reflect.Value, value string) {
	typ := field.Type()

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		if field.IsNil() {
			field.Set(reflect.New(typ))
		}
		field = field.Elem()
	}

	switch typ.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var (
			val int64
			err error
		)
		if field.Kind() == reflect.Int64 && typ.PkgPath() == "time" && typ.Name() == "Duration" {
			var d time.Duration
			d, err = time.ParseDuration(value)
			val = int64(d)
		} else {
			val, err = strconv.ParseInt(value, 0, typ.Bits())
		}
		if err != nil {
			panic(err)
		}

		field.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(value, 0, typ.Bits())
		if err != nil {
			panic(err)
		}
		field.SetUint(val)
	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			panic(err)
		}
		field.SetBool(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(value, typ.Bits())
		if err != nil {
			panic(err)
		}
		field.SetFloat(val)
	case reflect.Slice:
		vals := strings.Split(value, ";")
		sl := reflect.MakeSlice(typ, len(vals), len(vals))
		for i, val := range vals {
			setField(sl.Index(i), val)
		}
		field.Set(sl)
	case reflect.Map:
		mp := reflect.MakeMap(typ)
		if len(strings.TrimSpace(value)) != 0 {
			pairs := strings.Split(value, ";")
			for _, pair := range pairs {
				kvpair := strings.Split(pair, ":")
				if len(kvpair) != 2 {
					panic(fmt.Errorf("invalid map item: %q", pair))
				}
				k := reflect.New(typ.Key()).Elem()
				setField(k, kvpair[0])
				v := reflect.New(typ.Elem()).Elem()
				setField(v, kvpair[1])
				mp.SetMapIndex(k, v)
			}
		}
		field.Set(mp)
	default:
		panic(fmt.Errorf("invalid type: %v", typ))
	}
}
