package main

//go:generate msgp

type Leaderboard struct {
	Name string `msg:"name"`
	Members map[string]*Member `msg:"members"`

	MembersIndex []*MemberIndex  `msg:"members_index"`
}

func (l Leaderboard) Len() int           { return len(l.MembersIndex) }
func (l Leaderboard) Swap(i, j int)      { l.MembersIndex[i], l.MembersIndex[j] = l.MembersIndex[j], l.MembersIndex[i] }
func (l Leaderboard) Less(i, j int) bool { return l.MembersIndex[i].Score > l.MembersIndex[j].Score }


func (l *Leaderboard) CheckOrAdd(member *Member) []*MemberIndex {
	for _, val := range l.MembersIndex {
		if val.Key == member.Key {
			return l.MembersIndex
		}
	}
	return append(l.MembersIndex, member.toIndex())
}