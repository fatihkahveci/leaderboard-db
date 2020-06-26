package main

import (
	"errors"
	"sort"
	"sync"
)

//go:generate msgp

type Store struct {
	mu sync.RWMutex
	Leaderboards map[string]*Leaderboard `msg:"leaderboards"`
}


func NewStore() *Store {
	store := &Store{}
	store.Leaderboards = make(map[string]*Leaderboard)

	return store
}

func (s *Store) AddMember(leaderboardName string,key string, score float64, fields map[string]string)  {
	s.mu.Lock()
	defer s.mu.Unlock()

	member := &Member{
		Key: key,
		Score: score,
		Fields: fields,
	}
	leaderboard := s.Leaderboards[leaderboardName]
	if leaderboard == nil {
		//Append now leaderboard
		leaderboard = &Leaderboard{}
		leaderboard.Members = make(map[string]*Member)
		s.Leaderboards[leaderboardName] = leaderboard
		s.Leaderboards[leaderboardName].MembersIndex = append(s.Leaderboards[leaderboardName].MembersIndex, member.toIndex())
	}


	check := leaderboard.Members[key]
	leaderboard.Members[key] = member

	if check != nil {
		for i, v := range leaderboard.MembersIndex {
			if key == v.Key  {
				leaderboard.MembersIndex[i] = member.toIndex()
				//sort.Sort(leaderboard)
				return
			}
		}
	}


	leaderboard.MembersIndex = append(leaderboard.MembersIndex, member.toIndex())
	//sort.Sort(leaderboard)
}

func (s *Store) DeleteMember(leaderboardName, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	leaderboard, ok := s.Leaderboards[leaderboardName]
	if !ok {
		return errors.New("Leaderboard not found")
	}
	member, ok := leaderboard.Members[key]
	if !ok {
		return errors.New("Member not found")
	}

	delete(leaderboard.Members, member.Key)

	var newMemberIndex []*MemberIndex
	for _,val := range leaderboard.MembersIndex {
		if val.Key == member.Key {
			continue
		}
		newMemberIndex = append(newMemberIndex, val)
	}

	leaderboard.MembersIndex = newMemberIndex

	if len(leaderboard.Members) == 0 {
		//Delete leaderboard if empty
		delete(s.Leaderboards, leaderboardName)
		return nil
	}

	sort.Sort(leaderboard)

	return nil
}

func (s *Store) DeleteLeaderboard(leaderboardName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.Leaderboards[leaderboardName]
	if !ok {
		return errors.New("Leaderboard not found")
	}

	delete(s.Leaderboards, leaderboardName)

	return nil
}


func (s *Store) MemberScore(leaderboardName, key string) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	leaderboard, ok := s.Leaderboards[leaderboardName]
	if !ok {
		return 0,errors.New("Leaderboard not found")
	}
	member, ok := leaderboard.Members[key]
	if !ok {
		return 0,errors.New("Member not found")
	}

	return member.Score, nil

}

func (s *Store) MemberRank(leaderboardName, key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	leaderboard, ok := s.Leaderboards[leaderboardName]
	if !ok {
		return 0,errors.New("Leaderboard not found")
	}
	member, ok := leaderboard.Members[key]
	if !ok {
		return 0,errors.New("Member not found")
	}

	rank := 0
	for index,val := range leaderboard.MembersIndex {
		if val.Key == member.Key {
			rank = index+1
		}
	}

	return rank, nil

}
func (s *Store) UpdateMemberScore(leaderboardName, key string, score float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	leaderboard, ok := s.Leaderboards[leaderboardName]
	if !ok {
		return errors.New("Leaderboard not found")
	}
	member, ok := leaderboard.Members[key]
	if !ok {
		return errors.New("Member not found")
	}

	member.Score = member.Score + score

	for index,val := range leaderboard.MembersIndex {
		if val.Key == member.Key {
			leaderboard.MembersIndex[index].Score = leaderboard.MembersIndex[index].Score + score
		}
	}

	sort.Sort(leaderboard)

	return nil
}

func (s *Store) MemberAddField(leaderboardName, key string, fields map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	leaderboard, ok := s.Leaderboards[leaderboardName]
	if !ok {
		return errors.New("Leaderboard not found")
	}
	member, ok := leaderboard.Members[key]
	if !ok {
		return errors.New("Member not found")
	}

	oldFields := member.Fields
	newFields := fields

	for key,val := range oldFields {
		newFields[key] = val
	}

	member.Fields = newFields
	return nil
}

func (s *Store) GetSort(leaderboardName string, fields map[string]string) ([]string,error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var results []string

	leaderboard, ok := s.Leaderboards[leaderboardName]
	if !ok {
		return results, errors.New("Leaderboard not found")
	}

	if len (fields) > 0 {
		tmpLeaderboard := &Leaderboard{
			Name: leaderboardName,
		}
		tmpLeaderboard.Members = make(map[string]*Member)
		for memberKey, member  := range leaderboard.Members {
			found := true
			for key,val := range fields{
				data,ok := member.Fields[key]
				if !ok || data != val {
					found = false
					continue
				}
			}

			if found {
				tmpLeaderboard.Members[memberKey] = member
				tmpLeaderboard.MembersIndex = append(tmpLeaderboard.MembersIndex, member.toIndex())
			}
		}

		sort.Sort(tmpLeaderboard)

		for _, val := range tmpLeaderboard.MembersIndex{
			results = append(results, val.Key)
		}

		return results,nil
	}


	sort.Sort(leaderboard)

	for _, val := range leaderboard.MembersIndex{
		results = append(results, val.Key)
	}

	return results,nil
}