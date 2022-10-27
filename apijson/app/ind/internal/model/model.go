package model

type ReqList struct {
	PageNum   int      `p:"_pn" d:"1" v:"integer#_pn必须为int" json:"_pn,omitempty"`
	PageSize  int      `p:"_ps" d:"10" v:"integer#_pn必须为int" json:"_ps,omitempty"`
	CreatedAt []string `json:"createdAt" q:"TimeBetween"`
}
