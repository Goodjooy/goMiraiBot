package interactprefab

import "goMiraiQQBot/lib/constdata"

type InteractPerfabBulider struct {
	perfab *InteractPerfab
}

func NewInteractPerfab() *InteractPerfabBulider {
	return &InteractPerfabBulider{
		perfab: &InteractPerfab{},
	}
}

func (i *InteractPerfabBulider) AddInitFn(fn func()) *InteractPerfabBulider {
	i.perfab.initFn = fn
	return i
}
func (i *InteractPerfabBulider) SetUseage(useage string) *InteractPerfabBulider {
	i.perfab.useage = useage
	return i
}
func (i *InteractPerfabBulider) AddActivateSigns(signs ...string) *InteractPerfabBulider {
	i.perfab.activateSigns = append(i.perfab.activateSigns, signs...)
	return i
}

func (i *InteractPerfabBulider) AddActivateSource(sources ...constdata.MessageType) *InteractPerfabBulider {
	i.perfab.activateSource = append(i.perfab.activateSource, sources...)
	return i
}

func (i *InteractPerfabBulider) BuildPtr() *InteractPerfab {
	return i.perfab
}
func (i *InteractPerfabBulider) Build() InteractPerfab {
	return *i.perfab
}
