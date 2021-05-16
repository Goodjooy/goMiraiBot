package collection

type FilterHandle func(interface{})bool

func Filter(handle FilterHandle,collect []interface{})[]interface{}