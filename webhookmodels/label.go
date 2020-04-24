package webhookmodels

// Label represents a github label
type Label struct {
	ID     int64  `json:"id"`
	NodeID string `json:"node_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

type Labels []Label

func (l *Labels) Names() []string {
	ret := []string{}
	for _, i := range *l {
		ret = append(ret, i.Name)
	}
	return ret
}
