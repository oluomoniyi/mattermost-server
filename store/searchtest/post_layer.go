package searchtest

import (
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/store"
	"github.com/stretchr/testify/require"
)

var searchPostStoreTests = []searchTest{
	{
		"Should be able to search posts including results from DMs",
		testSearchPostsIncludingDMs,
		[]string{ENGINE_ALL},
	},
	{
		"Should return pinned and unpinned posts",
		testSearchReturnPinnedAndUnpinned,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to search for exact phrases in quotes",
		testSearchExactPhraseInQuotes,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to search for email addresses with or without quotes",
		testSearchEmailAddresses,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to search when markdown underscores are applied",
		testSearchMarkdownUnderscores,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to search for non-latin words",
		testSearchNonLatinWords,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to search for alternative spellings of words",
		testSearchAlternativeSpellings,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to search for alternative spellings of words with and without accents",
		testSearchAlternativeSpellingsAccents,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to search or exclude messages written by a specific user",
		testSearchOrExcludePostsBySpecificUser,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to search or exclude messages written in a specific channel",
		testSearchOrExcludePostsInChannel,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to search or exclude messages written in a DM or GM",
		testSearchOrExcludePostsInDMGM,
		[]string{ENGINE_ALL},
	},
}

func TestSearchPostStore(t *testing.T, s store.Store, testEngine *SearchTestEngine) {
	th := &SearchTestHelper{
		Store: s,
	}
	err := th.SetupBasicFixtures()
	require.Nil(t, err)
	defer th.CleanFixtures()

	runTestSearch(t, testEngine, searchPostStoreTests, th)
}

func testSearchPostsIncludingDMs(t *testing.T, th *SearchTestHelper) {
	direct, err := th.createDirectChannel(th.Team.Id, "direct", "direct", []*model.User{th.User, th.User2})
	require.Nil(t, err)
	defer th.deleteChannel(direct)

	p1, err := th.createPost(th.User.Id, direct.Id, "dm test", "", 0, false)
	require.Nil(t, err)
	_, err = th.createPost(th.User.Id, direct.Id, "dm other", "", 0, false)
	require.Nil(t, err)
	p2, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "channel test", "", 0, false)
	require.Nil(t, err)
	defer th.deleteUserPosts(th.User.Id)

	params := &model.SearchParams{Terms: "test"}
	results, err := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
	require.Nil(t, err)

	require.Len(t, results.Posts, 2)
	th.checkPostInSearchResults(t, p1.Id, results.Posts)
	th.checkPostInSearchResults(t, p2.Id, results.Posts)
}

func testSearchReturnPinnedAndUnpinned(t *testing.T, th *SearchTestHelper) {
	p1, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "channel test unpinned", "", 0, false)
	require.Nil(t, err)
	p2, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "channel test pinned", "", 0, true)
	require.Nil(t, err)
	defer th.deleteUserPosts(th.User.Id)

	params := &model.SearchParams{Terms: "test"}
	results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
	require.Nil(t, apperr)

	require.Len(t, results.Posts, 2)
	th.checkPostInSearchResults(t, p1.Id, results.Posts)
	th.checkPostInSearchResults(t, p2.Id, results.Posts)
}

func testSearchExactPhraseInQuotes(t *testing.T, th *SearchTestHelper) {
	p1, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "channel test 1 2 3", "", 0, false)
	require.Nil(t, err)
	_, err = th.createPost(th.User.Id, th.ChannelBasic.Id, "channel test 123", "", 0, false)
	require.Nil(t, err)
	defer th.deleteUserPosts(th.User.Id)

	params := &model.SearchParams{Terms: "\"channel test 1 2 3\""}
	results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
	require.Nil(t, apperr)

	require.Len(t, results.Posts, 1)
	th.checkPostInSearchResults(t, p1.Id, results.Posts)
}

func testSearchEmailAddresses(t *testing.T, th *SearchTestHelper) {
	p1, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "test email test@test.com", "", 0, false)
	require.Nil(t, err)
	_, err = th.createPost(th.User.Id, th.ChannelBasic.Id, "test email test2@test.com", "", 0, false)
	require.Nil(t, err)
	defer th.deleteUserPosts(th.User.Id)

	t.Run("Should search email addresses enclosed by quotes", func(t *testing.T) {
		params := &model.SearchParams{Terms: "\"test@test.com\""}
		results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
		require.Nil(t, apperr)

		require.Len(t, results.Posts, 1)
		th.checkPostInSearchResults(t, p1.Id, results.Posts)
	})

	t.Run("Should search email addresses without quotes", func(t *testing.T) {
		params := &model.SearchParams{Terms: "test@test.com"}
		results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
		require.Nil(t, apperr)

		require.Len(t, results.Posts, 1)
		th.checkPostInSearchResults(t, p1.Id, results.Posts)
	})
}

func testSearchMarkdownUnderscores(t *testing.T, th *SearchTestHelper) {
	p1, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "_start middle end_ _both_", "", 0, false)
	require.Nil(t, err)
	defer th.deleteUserPosts(th.User.Id)

	t.Run("Should search the start inside the markdown underscore", func(t *testing.T) {
		params := &model.SearchParams{Terms: "start"}
		results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
		require.Nil(t, apperr)

		require.Len(t, results.Posts, 1)
		th.checkPostInSearchResults(t, p1.Id, results.Posts)
	})

	t.Run("Should search a word in the middle of the markdown underscore", func(t *testing.T) {
		params := &model.SearchParams{Terms: "middle"}
		results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
		require.Nil(t, apperr)

		require.Len(t, results.Posts, 1)
		th.checkPostInSearchResults(t, p1.Id, results.Posts)
	})

	t.Run("Should search in the end of the markdown underscore", func(t *testing.T) {
		params := &model.SearchParams{Terms: "end"}
		results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
		require.Nil(t, apperr)

		require.Len(t, results.Posts, 1)
		th.checkPostInSearchResults(t, p1.Id, results.Posts)
	})

	t.Run("Should search inside markdown underscore", func(t *testing.T) {
		params := &model.SearchParams{Terms: "both"}
		results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
		require.Nil(t, apperr)

		require.Len(t, results.Posts, 1)
		th.checkPostInSearchResults(t, p1.Id, results.Posts)
	})
}

func testSearchNonLatinWords(t *testing.T, th *SearchTestHelper) {
	t.Run("Should be able to search chinese words", func(t *testing.T) {
		p1, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "你好", "", 0, false)
		require.Nil(t, err)
		p2, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "你", "", 0, false)
		require.Nil(t, err)
		defer th.deleteUserPosts(th.User.Id)

		t.Run("Should search one word", func(t *testing.T) {
			params := &model.SearchParams{Terms: "你"}
			results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
			require.Nil(t, apperr)

			require.Len(t, results.Posts, 1)
			th.checkPostInSearchResults(t, p2.Id, results.Posts)
		})
		t.Run("Should search two words", func(t *testing.T) {
			params := &model.SearchParams{Terms: "你好"}
			results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
			require.Nil(t, apperr)

			require.Len(t, results.Posts, 1)
			th.checkPostInSearchResults(t, p1.Id, results.Posts)
		})
		t.Run("Should search with wildcard", func(t *testing.T) {
			params := &model.SearchParams{Terms: "你*"}
			results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
			require.Nil(t, apperr)

			require.Len(t, results.Posts, 2)
			th.checkPostInSearchResults(t, p1.Id, results.Posts)
			th.checkPostInSearchResults(t, p2.Id, results.Posts)
		})
	})
	t.Run("Should be able to search cyrillic words", func(t *testing.T) {
		p1, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "слово test", "", 0, false)
		require.Nil(t, err)
		defer th.deleteUserPosts(th.User.Id)

		t.Run("Should search one word", func(t *testing.T) {
			params := &model.SearchParams{Terms: "слово"}
			results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
			require.Nil(t, apperr)

			require.Len(t, results.Posts, 1)
			th.checkPostInSearchResults(t, p1.Id, results.Posts)
		})
		t.Run("Should search using wildcard", func(t *testing.T) {
			params := &model.SearchParams{Terms: "слов*"}
			results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
			require.Nil(t, apperr)

			require.Len(t, results.Posts, 1)
			th.checkPostInSearchResults(t, p1.Id, results.Posts)
		})
	})

	t.Run("Should be able to search japanese words", func(t *testing.T) {
		p1, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "本", "", 0, false)
		require.Nil(t, err)
		p2, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "本木", "", 0, false)
		require.Nil(t, err)
		defer th.deleteUserPosts(th.User.Id)

		t.Run("Should search one word", func(t *testing.T) {
			params := &model.SearchParams{Terms: "本"}
			results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
			require.Nil(t, apperr)

			require.Len(t, results.Posts, 2)
			th.checkPostInSearchResults(t, p1.Id, results.Posts)
			th.checkPostInSearchResults(t, p2.Id, results.Posts)
		})
		t.Run("Should search two words", func(t *testing.T) {
			params := &model.SearchParams{Terms: "本木"}
			results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
			require.Nil(t, apperr)

			require.Len(t, results.Posts, 1)
			th.checkPostInSearchResults(t, p2.Id, results.Posts)
		})
		t.Run("Should search with wildcard", func(t *testing.T) {
			params := &model.SearchParams{Terms: "本*"}
			results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
			require.Nil(t, apperr)

			require.Len(t, results.Posts, 2)
			th.checkPostInSearchResults(t, p1.Id, results.Posts)
			th.checkPostInSearchResults(t, p2.Id, results.Posts)
		})
	})

	t.Run("Should be able to search korean words", func(t *testing.T) {
		p1, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "불", "", 0, false)
		require.Nil(t, err)
		p2, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "불다", "", 0, false)
		require.Nil(t, err)
		defer th.deleteUserPosts(th.User.Id)

		t.Run("Should search one word", func(t *testing.T) {
			params := &model.SearchParams{Terms: "불"}
			results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
			require.Nil(t, apperr)

			require.Len(t, results.Posts, 1)
			th.checkPostInSearchResults(t, p1.Id, results.Posts)
		})
		t.Run("Should search two words", func(t *testing.T) {
			params := &model.SearchParams{Terms: "불다"}
			results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
			require.Nil(t, apperr)

			require.Len(t, results.Posts, 1)
			th.checkPostInSearchResults(t, p2.Id, results.Posts)
		})
		t.Run("Should search with wildcard", func(t *testing.T) {
			params := &model.SearchParams{Terms: "불*"}
			results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
			require.Nil(t, apperr)

			require.Len(t, results.Posts, 2)
			th.checkPostInSearchResults(t, p1.Id, results.Posts)
			th.checkPostInSearchResults(t, p2.Id, results.Posts)
		})
	})
}

func testSearchAlternativeSpellings(t *testing.T, th *SearchTestHelper) {
	p1, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "Straße test", "", 0, false)
	require.Nil(t, err)
	p2, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "Strasse test", "", 0, false)
	require.Nil(t, err)
	defer th.deleteUserPosts(th.User.Id)

	params := &model.SearchParams{Terms: "Straße"}
	results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
	require.Nil(t, apperr)

	require.Len(t, results.Posts, 2)
	th.checkPostInSearchResults(t, p1.Id, results.Posts)
	th.checkPostInSearchResults(t, p2.Id, results.Posts)

	params = &model.SearchParams{Terms: "Strasse"}
	results, apperr = th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
	require.Nil(t, apperr)

	require.Len(t, results.Posts, 2)
	th.checkPostInSearchResults(t, p1.Id, results.Posts)
	th.checkPostInSearchResults(t, p2.Id, results.Posts)
}

func testSearchAlternativeSpellingsAccents(t *testing.T, th *SearchTestHelper) {
	p1, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "café", "", 0, false)
	require.Nil(t, err)
	p2, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "café", "", 0, false)
	require.Nil(t, err)
	defer th.deleteUserPosts(th.User.Id)

	params := &model.SearchParams{Terms: "café"}
	results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
	require.Nil(t, apperr)

	require.Len(t, results.Posts, 2)
	th.checkPostInSearchResults(t, p1.Id, results.Posts)
	th.checkPostInSearchResults(t, p2.Id, results.Posts)

	params = &model.SearchParams{Terms: "café"}
	results, apperr = th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
	require.Nil(t, apperr)

	require.Len(t, results.Posts, 2)
	th.checkPostInSearchResults(t, p1.Id, results.Posts)
	th.checkPostInSearchResults(t, p2.Id, results.Posts)

	params = &model.SearchParams{Terms: "cafe"}
	results, apperr = th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
	require.Nil(t, apperr)

	require.Len(t, results.Posts, 0)
}

func testSearchOrExcludePostsBySpecificUser(t *testing.T, th *SearchTestHelper) {
	p1, err := th.createPost(th.User.Id, th.ChannelPrivate.Id, "test fromuser", "", 0, false)
	require.Nil(t, err)
	_, err = th.createPost(th.User2.Id, th.ChannelPrivate.Id, "test fromuser 2", "", 0, false)
	require.Nil(t, err)
	defer th.deleteUserPosts(th.User.Id)

	params := &model.SearchParams{
		Terms:     "fromuser",
		FromUsers: []string{th.User.Id},
	}
	results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
	require.Nil(t, apperr)

	require.Len(t, results.Posts, 1)
	th.checkPostInSearchResults(t, p1.Id, results.Posts)
}

func testSearchOrExcludePostsInChannel(t *testing.T, th *SearchTestHelper) {
	p1, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "test fromuser", "", 0, false)
	require.Nil(t, err)
	_, err = th.createPost(th.User2.Id, th.ChannelPrivate.Id, "test fromuser 2", "", 0, false)
	require.Nil(t, err)
	defer th.deleteUserPosts(th.User.Id)

	params := &model.SearchParams{
		Terms:      "fromuser",
		InChannels: []string{th.ChannelBasic.Id},
	}
	results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
	require.Nil(t, apperr)

	require.Len(t, results.Posts, 1)
	th.checkPostInSearchResults(t, p1.Id, results.Posts)
}

func testSearchOrExcludePostsInDMGM(t *testing.T, th *SearchTestHelper) {
	direct, err := th.createDirectChannel(th.Team.Id, "direct", "direct", []*model.User{th.User, th.User2})
	require.Nil(t, err)
	defer th.deleteChannel(direct)

	group, err := th.createGroupChannel(th.Team.Id, "test group", []*model.User{th.User, th.User2})
	require.Nil(t, err)
	defer th.deleteChannel(group)

	p1, err := th.createPost(th.User.Id, direct.Id, "test fromuser", "", 0, false)
	require.Nil(t, err)
	p2, err := th.createPost(th.User2.Id, group.Id, "test fromuser 2", "", 0, false)
	require.Nil(t, err)
	defer th.deleteUserPosts(th.User.Id)

	t.Run("Should be able to search in both DM and GM channels", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:      "fromuser",
			InChannels: []string{direct.Id, group.Id},
		}
		results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
		require.Nil(t, apperr)

		require.Len(t, results.Posts, 2)
		th.checkPostInSearchResults(t, p1.Id, results.Posts)
		th.checkPostInSearchResults(t, p2.Id, results.Posts)
	})

	t.Run("Should be able to search only in DM channel", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:      "fromuser",
			InChannels: []string{direct.Id},
		}
		results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
		require.Nil(t, apperr)

		require.Len(t, results.Posts, 1)
		th.checkPostInSearchResults(t, p1.Id, results.Posts)
	})

	t.Run("Should be able to search only in GM channel", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:      "fromuser",
			InChannels: []string{group.Id},
		}
		results, apperr := th.Store.Post().SearchPostsInTeamForUser([]*model.SearchParams{params}, th.User.Id, th.Team.Id, false, false, 0, 20)
		require.Nil(t, apperr)

		require.Len(t, results.Posts, 1)
		th.checkPostInSearchResults(t, p2.Id, results.Posts)
	})
}
