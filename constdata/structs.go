package constdata


type Sender struct {
	Id         uint64 `json:"id"`
	MemberName string `json:"memberName"`


	Permission PermissionLevel `json:"permission"`

	GroupIn Group `json:"group"`
}

type Group struct {
	Id   uint64 `json:"id"`
	Name string`json:"name"`

	Permission PermissionLevel `json:"permission"`
}
