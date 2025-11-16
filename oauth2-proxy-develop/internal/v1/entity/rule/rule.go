package rule

type Action string

const (
	ActionAllow Action = "allow"
	ActionDeny  Action = "deny"
	ActionAuth  Action = "auth"
)

//------------routes------------
// name: str
// path: str path
// type: exact,prefix
// action: allow,deny,auth
//------------------------------

type Rule struct {
	Name   string
	Path   string
	Action Action
}

func NewRule(name string, path string, action Action) *Rule {
	return &Rule{Name: name, Path: path, Action: action}
}
