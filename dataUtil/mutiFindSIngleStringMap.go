package datautil

type MutliToOneMap struct {
	basicData  map[string]string
	noNameList []string
}

func NewMutliToOneMap() MutliToOneMap {
	return MutliToOneMap{
		basicData:  make(map[string]string),
		noNameList: make([]string, 0),
	}
}
func (c *MutliToOneMap) IsEmpty() bool {
	return len(c.basicData) == 0 && len(c.noNameList) == 0
}
func (c *MutliToOneMap) Put(key, value string) {
	c.basicData[key] = value
}
func (c *MutliToOneMap) PutNoName(value string) {
	c.noNameList = append(c.noNameList, value)
}
func (c *MutliToOneMap) SetNoNameCmdOrder(cmdNames ...TargetValues) {
	var index = 0
	for _, v := range cmdNames { 
		_, ok := c.Get(v...)
		//没找到指令
		if !ok && index < len(c.noNameList) {
			//按照顺序识别指令
			c.Put(v.GetSign(), c.noNameList[index])
			index++
		}
	}
}

func (c *MutliToOneMap) Get(keys ...string) (string, bool) {
	//找到一个就是成功
	for _, key := range keys {
		data, ok := (c.basicData)[key]
		if ok {
			return data, true
		}
	}
	return "", false
}

func (c *MutliToOneMap) GetWithDefault(defaultValue string, keys ...string) (string, bool) {
	//找到一个就是成功
	v, ok := c.Get(keys...)
	if ok {
		return v, ok
	}
	return defaultValue, false
}