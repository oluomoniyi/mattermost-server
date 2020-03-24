package searchtest

import (
	"testing"

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
	return
}
