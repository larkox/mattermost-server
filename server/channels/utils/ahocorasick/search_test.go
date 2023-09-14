package ahocosarick

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	tree := NewNode()
	tree.AddString(KeywordString{KeywordTerm{Type: "string", Term: "hello"}}, LeafObject{UserID: "user1", MentionType: 1})
	tree.AddString(KeywordString{KeywordTerm{Type: "string", Term: "hello my dear friend"}}, LeafObject{UserID: "user2", MentionType: 2})
	tree.AddString(KeywordString{KeywordTerm{Type: "string", Term: "bye"}}, LeafObject{UserID: "user3", MentionType: 3})
	tree.AddString(KeywordString{
		KeywordTerm{Type: "string", Term: "some"},
		KeywordTerm{Type: "alphanumeric"},
		KeywordTerm{Type: "string", Term: "thing"},
	}, LeafObject{UserID: "user4", MentionType: 4})
	tree.AddString(KeywordString{
		KeywordTerm{Type: "string", Term: "zet"},
		KeywordTerm{Type: "alphanumeric*"},
		KeywordTerm{Type: "string", Term: "zet"},
	}, LeafObject{UserID: "user5", MentionType: 5})
	tree.AddString(KeywordString{
		KeywordTerm{Type: "wordseparator"},
		KeywordTerm{Type: "string", Term: "foo"},
		KeywordTerm{Type: "wordseparator"},
	}, LeafObject{UserID: "user6", MentionType: 6})
	tree.AddString(KeywordString{
		KeywordTerm{Type: "string", Term: "foo"},
		KeywordTerm{Type: "wordseparator*"},
		KeywordTerm{Type: "string", Term: "foo"},
	}, LeafObject{UserID: "user7", MentionType: 7})
	tree.AddString(KeywordString{
		KeywordTerm{Type: "string", Term: "hi-world", CaseInsensitive: true},
	}, LeafObject{UserID: "user8", MentionType: 8})

	leaves := tree.Search("hello my dear friend and good bye")
	require.Len(t, leaves, 3, "all three strings are found once")

	leaves = tree.Search("goodbye my darling")
	require.Len(t, leaves, 1)
	require.Equal(t, "bye", leaves[0].Term)
	require.Equal(t, "user3", leaves[0].Values[0].UserID)

	leaves = tree.Search("byebye my darling")
	require.Len(t, leaves, 1)
	require.Equal(t, "bye", leaves[0].Term)
	require.Equal(t, "user3", leaves[0].Values[0].UserID)

	leaves = tree.Search("Bye my darling")
	require.Len(t, leaves, 0)

	leaves = tree.Search("some2thing my darling")
	require.Len(t, leaves, 1)
	require.Equal(t, "some{{alphanumeric}}thing", leaves[0].Term)
	require.Equal(t, "user4", leaves[0].Values[0].UserID)

	leaves = tree.Search("some thing my darling")
	require.Len(t, leaves, 0)

	leaves = tree.Search("asdfasdfzetasdfasdfasdfasfzetdasdfasdf and something else")
	require.Len(t, leaves, 1)
	require.Equal(t, "zet{{alphanumeric*}}zet", leaves[0].Term)
	require.Equal(t, "user5", leaves[0].Values[0].UserID)

	leaves = tree.Search("something foo my friend")
	require.Len(t, leaves, 1)
	require.Equal(t, "{{wordseparator}}foo{{wordseparator}}", leaves[0].Term)
	require.Equal(t, "user6", leaves[0].Values[0].UserID)

	leaves = tree.Search("foo my friend")
	require.Len(t, leaves, 1)
	require.Equal(t, "{{wordseparator}}foo{{wordseparator}}", leaves[0].Term)
	require.Equal(t, "user6", leaves[0].Values[0].UserID)

	leaves = tree.Search("my friend foo")
	require.Len(t, leaves, 1)
	require.Equal(t, "{{wordseparator}}foo{{wordseparator}}", leaves[0].Term)
	require.Equal(t, "user6", leaves[0].Values[0].UserID)

	leaves = tree.Search("myfoo !.-! foofoo")
	require.Len(t, leaves, 1)
	require.Equal(t, "foo{{wordseparator*}}foo", leaves[0].Term)
	require.Equal(t, "user7", leaves[0].Values[0].UserID)

	leaves = tree.Search("hi-world")
	require.Len(t, leaves, 1)
	require.Equal(t, "hi-world", leaves[0].Term)
	require.Equal(t, "user8", leaves[0].Values[0].UserID)

	leaves = tree.Search("Hi-WoRlD")
	require.Len(t, leaves, 1)
	require.Equal(t, "Hi-WoRlD", leaves[0].Term)
	require.Equal(t, "user8", leaves[0].Values[0].UserID)
}
