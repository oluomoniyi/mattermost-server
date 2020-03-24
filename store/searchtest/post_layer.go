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
	post, err := th.createPost(th.User.Id, direct.Id, "dm test", "", 0)
	require.Nil(t, err)
	_, err = th.createPost(th.User.Id, direct.Id, "dm other", "", 0)
	require.Nil(t, err)
	defer th.deleteUserPosts(th.User.Id)
	params := &model.SearchParams{Terms: "test"}
	results, err := th.Store.Post().Search(th.Team.Id, th.User.Id, params)
	require.Nil(t, err)
	require.Len(t, results.Posts, 1)
	th.checkPostInSearchResults(t, post.Id, results.Posts)
}
