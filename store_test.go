package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	store = NewStore()
	fields := map[string]string{
		"character":"monk",
	}

	store.AddMember("diablo_3","user_2", 200, fields)
	store.AddMember("diablo_3","user_3", 300, fields)

	fields2 := map[string]string{
		"character":"wizard",
	}

	store.AddMember("diablo_3","user_4", 400, fields2)
	store.AddMember("diablo_3","user_5", 500, fields2)
}

func TestStore_AddMember(t *testing.T) {
	fields := map[string]string{
		"character":"monk",
	}

	store.AddMember("diablo_3","user_1", 100, fields)
	check := store.Leaderboards["diablo_3"].Members["user_1"]

	assert.Equal(t, "user_1", check.Key)
}

func TestStore_MemberRank(t *testing.T) {
	rank, err := store.MemberRank("diablo_3","user_1")

	assert.Nil(t, err)
	assert.Equal(t, 5, rank)
}

func TestStore_MemberScore(t *testing.T) {
	rank, err := store.MemberScore("diablo_3","user_1")

	assert.Nil(t, err)
	assert.Equal(t, float64(100), rank)
}

func TestStore_GetSort(t *testing.T) {
	fields := make(map[string]string)
	resp,err := store.GetSort("diablo_3",fields)

	assert.Nil(t, err)
	assert.ElementsMatch(t,resp, []string{"user_5","user_4","user_3","user_2","user_1"})

}

func TestStore_GetSortWithFilter(t *testing.T) {
	fields := map[string]string{
		"character":"wizard",
	}
	resp,err := store.GetSort("diablo_3",fields)

	assert.Nil(t, err)
	assert.ElementsMatch(t,resp, []string{"user_5","user_4"})
}

func TestStore_UpdateMemberScore(t *testing.T) {
	err := store.UpdateMemberScore("diablo_3","user_1", -500)
	assert.Nil(t, err)
}

func BenchmarkStore_AddMember(b *testing.B) {
	fields := map[string]string{
		"character":"monk",
	}
	for n := 0; n < b.N; n++ {
		store.AddMember("diablo_3","user_1", 100, fields)
	}
}

func BenchmarkStore_UpdateMemberScore(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = store.UpdateMemberScore("diablo_3","user_1", 100)
	}
}

func BenchmarkStore_GetSort(b *testing.B) {
	fields := make(map[string]string)
	for n := 0; n < b.N; n++ {
		 store.GetSort("diablo_3",fields)
	}
}

func BenchmarkStore_GetSortWithFilter(b *testing.B) {
	fields := map[string]string{
		"character":"wizard",
	}
	for n := 0; n < b.N; n++ {
		store.GetSort("diablo_3",fields)
	}
}

func BenchmarkStore_MemberRank(b *testing.B) {
	for n := 0; n < b.N; n++ {
		store.MemberRank("diablo_3","user_1")
	}
}

func BenchmarkStore_MemberScore(b *testing.B) {
	for n := 0; n < b.N; n++ {
		store.MemberScore("diablo_3","user_1")
	}
}