package ahocorasick

import (
	"fmt"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/require"
)

func addUserListToTree(tree *Node, userList []*model.User) {
	for _, user := range userList {
		tree.AddString(KeywordString{
			{Type: "wordseparator"},
			{Type: "string", Term: "@" + user.Username, CaseInsensitive: true},
			{Type: "wordseparator"},
		}, LeafObject{UserID: user.Id, MentionType: 1})
		if len(user.FirstName) > 0 {
			tree.AddString(KeywordString{{Type: "wordseparator"}, {Type: "string", Term: user.FirstName}, {Type: "wordseparator"}}, LeafObject{UserID: user.Id, MentionType: 2})
		}
		keywords := user.GetMentionKeys()
		for _, k := range keywords {
			tree.AddString(KeywordString{{Type: "wordseparator"}, {Type: "string", Term: k, CaseInsensitive: true}, {Type: "wordseparator"}}, LeafObject{UserID: user.Id, MentionType: 3})
		}
	}
}

func BenchmarkCorasick(b *testing.B) {
	const USER_COUNT = 50000
	users := make([]*model.User, USER_COUNT)
	for i := 0; i < USER_COUNT; i++ {
		newUser := model.User{
			Id:        model.NewId(),
			Username:  model.NewId(),
			FirstName: model.NewId(),
		}
		newUser.AddNotifyProp(model.MentionKeysNotifyProp, fmt.Sprintf("%s,%s,%s", model.NewId(), model.NewId(), model.NewId()))
		users[i] = &newUser
	}

	b.Run("Build", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			tree := NewNode()
			addUserListToTree(tree, users)
		}
	})

	tree := NewNode()
	addUserListToTree(tree, users)

	b.Run("Search", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			tree.Search(`
		Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam lobortis laoreet arcu ac maximus. Fusce vitae velit massa. Aenean ac tortor eros. Duis pellentesque orci non nisl venenatis, eu maximus ligula semper. Sed ligula enim, suscipit in tempor ut, pulvinar eget mauris. Mauris porttitor nibh quis pharetra blandit. Fusce finibus mollis ante vel dictum. Vestibulum consequat pretium molestie. Integer blandit mi et tellus placerat, sed iaculis mauris dapibus. Ut eu purus accumsan metus ultricies hendrerit. Mauris at eleifend nisi, vitae malesuada purus. In porta quam et pharetra accumsan.
		
		In hac habitasse platea dictumst. Nullam ultrices porta consequat. Nunc vestibulum dolor eget augue commodo efficitur. Etiam vel leo consequat, pellentesque est a, rutrum mi. Vivamus bibendum rhoncus metus, eget sollicitudin dui eleifend posuere. Curabitur pretium, diam eu ornare pretium, dui magna vulputate leo, id faucibus ex ligula ut augue. Aliquam sit amet lobortis tortor. Nam quis malesuada massa.
		
		Suspendisse sed felis iaculis est tempor convallis in eget augue. Praesent vehicula nisl id laoreet bibendum. In venenatis congue ultrices. Fusce accumsan ut massa et laoreet. Morbi vitae interdum libero. In malesuada enim vitae pretium aliquet. Donec volutpat sollicitudin lacinia. Integer vitae porta enim, id auctor nibh. Donec fringilla, arcu id aliquam consectetur, leo eros vulputate enim, eget faucibus urna purus at tortor. Donec ac condimentum mi, vel porta nisl. Maecenas imperdiet placerat luctus. Quisque finibus, ex et elementum dapibus, sem justo mollis nisi, nec pharetra purus sem id odio. Curabitur rhoncus nunc eget est interdum, eu tempus ligula aliquet.
		
		Phasellus elit mauris, tempus ut neque a, dictum posuere metus. Fusce in urna quis tellus finibus lobortis id sed lorem. Donec congue egestas metus, eu fringilla ipsum tincidunt nec. Duis et egestas dolor. Cras tristique eu magna quis sodales. Nulla non leo eget sem imperdiet vulputate id eu purus. Nunc venenatis quam eu gravida varius.
		
		Duis elementum arcu non quam tempus, id imperdiet quam dapibus. In hac habitasse platea dictumst. Nam congue pulvinar odio, eu elementum tortor tempus vel. Duis placerat dictum felis, quis imperdiet risus consectetur non. Cras at leo a magna volutpat sollicitudin vitae quis erat. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Nullam facilisis rutrum lacus, vel aliquet ex. Nulla cursus semper ex. Aliquam eu tincidunt orci. Nunc at consectetur quam. Etiam sodales quis nisi sed luctus. Donec et lacus purus. Pellentesque quis malesuada enim. Vestibulum lobortis auctor consectetur.
		
		Integer in est et dolor suscipit ullamcorper. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Nulla facilisi. Suspendisse ultrices, mi at rhoncus posuere, ex odio laoreet ante, in dictum nunc justo ac nibh. Sed malesuada, massa id sollicitudin viverra, diam lacus ultrices risus, quis imperdiet ante velit ullamcorper libero. Aliquam molestie nulla posuere ante porta aliquet. Aenean finibus vitae ligula non pretium. Nam sem enim, cursus ultricies egestas sit amet, lobortis eu odio. Nunc sed sagittis mauris.
		
		Etiam urna ante, volutpat ac nisi molestie, ornare condimentum augue. Praesent tortor urna, laoreet nec augue quis, suscipit pulvinar ligula. Cras hendrerit nisl non urna fermentum finibus sit amet eget orci. Fusce ligula lacus, dapibus at nulla placerat, euismod vestibulum felis. Ut ac magna lobortis tortor volutpat aliquam. Etiam dignissim orci nec erat feugiat gravida. Praesent imperdiet venenatis sem non tempor. Vivamus sit amet bibendum ante. Nunc quis laoreet erat, a fermentum lectus. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Pellentesque lacinia nisi eget ex rhoncus interdum. Nullam euismod lectus volutpat turpis pulvinar posuere. Donec lacus purus, hendrerit sed porta fermentum, ornare et nulla.
		
		Fusce vitae scelerisque nisi. Integer sit amet tortor purus. Nulla commodo nisl non purus fermentum, sed pharetra ligula egestas. Aliquam tempus libero vel metus consequat, eu aliquet erat placerat. Sed mattis dapibus maximus. Etiam posuere, lorem non fermentum venenatis, magna nisl dapibus urna, sit amet feugiat urna elit at velit. Curabitur non pulvinar diam, id iaculis tellus. Praesent pretium, tellus eget viverra pellentesque, diam eros porttitor nisi, ut placerat elit lorem sed turpis. Aliquam non arcu cursus, elementum purus at, consectetur ante. Integer facilisis rhoncus suscipit. Maecenas et elit tellus. Nam vestibulum dictum lacus a laoreet. Morbi varius sollicitudin auctor. Aliquam quis lacus lorem. Nam lobortis posuere urna, in efficitur risus posuere at. Ut porta in arcu id venenatis.
		
		Nam tristique, nisi nec maximus vehicula, libero ante lobortis quam, sit amet feugiat metus mauris at quam. Sed ultricies ipsum ligula. In porttitor pharetra magna, id tincidunt risus euismod vel. Maecenas a justo eget orci hendrerit semper. Etiam porta dui in eros vehicula, ut feugiat quam mattis. Sed at felis eu tellus faucibus dapibus nec eu risus. Phasellus bibendum lacus quis tempor pretium. Phasellus pharetra pellentesque augue, quis maximus arcu tempor sed. Ut eu gravida dolor, in iaculis lorem. Nulla vulputate porttitor ligula eu iaculis. Vestibulum accumsan a erat in dapibus. Fusce eget nunc commodo, consectetur arcu eget, pretium purus.
		
		Phasellus pellentesque cursus arcu non varius. Ut lorem massa, eleifend vehicula pellentesque eu, faucibus tristique leo. Quisque dui elit, laoreet sed malesuada a, iaculis at ante. Vivamus lectus justo, sodales sit amet ipsum ut, pellentesque suscipit nisl. Praesent quis leo nec lorem tempus elementum vel id lorem. In pulvinar, sem quis scelerisque feugiat, ante nisl auctor lorem, quis convallis nulla massa convallis dolor. Duis sed ligula placerat libero tempor luctus. Donec commodo ipsum sed arcu eleifend, nec cursus nisi venenatis. Curabitur pulvinar sapien libero, eu gravida sapien ornare non. Curabitur non venenatis lorem. Suspendisse nec tincidunt purus. Suspendisse potenti. Vestibulum congue lorem in tellus rutrum faucibus. Vivamus ac lectus faucibus, vestibulum neque in, posuere sapien. Nam fermentum sem id nunc facilisis, quis molestie erat euismod.
		
		Pellentesque orci risus, aliquam id hendrerit et, tincidunt ac magna. Quisque elit turpis, mattis in cursus et, vulputate nec erat. Morbi aliquet maximus sodales. Aliquam ac arcu libero. Integer mollis accumsan nisl faucibus mattis. Curabitur non commodo sapien. Aliquam hendrerit gravida libero, quis faucibus erat pellentesque eget. Cras vel ultrices mauris, convallis lobortis ligula. Suspendisse rutrum nec elit non molestie. Nullam vel odio massa. Etiam ultrices lacus ac luctus faucibus.
		` + users[0].FirstName)
		}
	})
}

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
	require.Equal(t, "xhxi-xwxoxrxlxd", leaves[0].Term)
	require.Equal(t, "user8", leaves[0].Values[0].UserID)

	leaves = tree.Search("Hi-WoRlD")
	require.Len(t, leaves, 1)
	require.Equal(t, "xhxi-xwxoxrxlxd", leaves[0].Term)
	require.Equal(t, "user8", leaves[0].Values[0].UserID)
}
