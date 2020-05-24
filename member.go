package main

//go:generate msgp

type Member struct {
	Key string `msg:"key"`
	Score float64 `msg:"score"`
	Fields map[string]string `msg:"fields"`
}

type MemberIndex struct {
	Key string `msg:"key"`
	Score float64 `msg:"score"`
}


func (m *Member) toIndex() *MemberIndex {
	return &MemberIndex{
		Key: m.Key,
		Score: m.Score,
	}
}