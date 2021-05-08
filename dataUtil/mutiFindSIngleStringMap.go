package datautil
type MutliToOneMap map[string]string

func NewMutliToOneMap()MutliToOneMap{
	return make(MutliToOneMap)
}

func (c *MutliToOneMap) Get(keys ...string) (string, bool) {
	//找到一个就是成功
	for _, key := range keys {
		data, ok := (*c)[key]
		if ok {
			return data, true
		}
	}
	return "", false
}

func (c *MutliToOneMap) GetWithDefault(defaultValue string, keys ...string) (string, bool) {
	//找到一个就是成功
	for _, key := range keys {
		data, ok := (*c)[key]
		if ok {
			return data, true
		}
	}
	return defaultValue, false
}