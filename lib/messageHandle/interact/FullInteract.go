package interact


type FullSingleInteract interface {
	InteractSideInformation
	Interact
}

/*
FullContextInteract 有上下文关系的信息交互部分,
能够提供连续的信息交互.优先级高于普通信息交互
*/
type FullContextInteract interface {
	InteractSideInformation
	ContextInteract
}