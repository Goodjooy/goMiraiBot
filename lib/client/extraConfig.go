package client

import "reflect"

func GetExtraConfig(key string) (interface{}, bool) {
	i, ok := cfg.ExtraConfig[key]

	return i, ok
}

func ParseExtraConfigToTarget(key string, target interface{}) bool {
	config, ok := GetExtraConfig(key)
	if !ok {
		return ok
	}
	cr := config.(map[interface{}]interface{})

	targetType := reflect.TypeOf(target).Elem()
	targetValure := reflect.ValueOf(target).Elem()

	count := targetType.NumField()
	for i := 0; i < count; i++ {
		feild := targetType.Field(i)
		feildType := feild.Type

		tags := feild.Tag.Get("config")

		data, ok := cr[tags]

		if ok {
			dataValue := reflect.ValueOf(data)
			dataType := dataValue.Type()

			//如果为数字
			if dataType != feildType {
				if dataType.Kind() == reflect.Float64 {
					if feildType.Kind() <= reflect.Uint64 &&
						feildType.Kind() >= reflect.Uint {
						dataValue = reflect.ValueOf(uint64(data.(float64)))
					} else if feildType.Kind() >= reflect.Int &&
						feildType.Kind() <= reflect.Int64 {
						dataValue = reflect.ValueOf(int64(data.(float64)))
					}
				}
			} else if dataType.Kind() == reflect.Struct {
				continue

			} else {
				dataValue = reflect.ValueOf(data.(string))
			}

			vFeild := targetValure.Field(i)
			if vFeild.CanSet() {
				vFeild.Set(dataValue)
			}
		}
	}

	return true
}
