package searchtest

import (
	"testing"

	"github.com/mattermost/mattermost-server/v5/store"
	"github.com/stretchr/testify/require"
)

var searchChannelStoreTests = []searchTest{
	{
		"Should be able to autocomplete a channel by name",
		testAutocompleteChannelByName,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to autocomplete a channel by display name",
		testAutocompleteChannelByDisplayName,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to autocomplete a channel by a part of its name when has parts splitted by - character",
		testAutocompleteChannelByNameSplittedWithDashChar,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to autocomplete a channel by a part of its name when has parts splitted by , character",
		testAutocompleteChannelByNameSplittedWithCommaChar,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to autocomplete a channel by a part of its name when has parts splitted by _ character",
		testAutocompleteChannelByNameSplittedWithUnderscoreChar,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to autocomplete a channel by a part of its display name when has parts splitted by whitespace character",
		testAutocompleteChannelByDisplayNameSplittedByWhitespaces,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to autocomplete retrieving all channels if the term is empty",
		testAutocompleteChannelByDisplayNameSplittedByWhitespaces,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to autocomplete channels in a case insensitive manner",
		testSearchChannelsInCaseInsensitiveManner,
		[]string{ENGINE_ALL},
	},
	{
		"Should autocomplete only returning public channels",
		testSearchOnlyPublicChannels,
		[]string{ENGINE_ALL},
	},
	{
		"Should support to autocomplete having a hyphen as the last character",
		testSearchShouldSupportHavingHyphenAsLastCharacter,
		[]string{ENGINE_ALL},
	},
	{
		"Should support to autocomplete with archived channels",
		testSearchShouldSupportAutocompleteWithArchivedChannels,
		[]string{ENGINE_ALL},
	},
}

func TestSearchChannelStore(t *testing.T, s store.Store, testEngine *SearchTestEngine) {
	th := &SearchTestHelper{
		Store: s,
	}
	err := th.InitFixtures()
	require.Nil(t, err)
	defer th.CleanFixtures()
	runTestSearch(t, testEngine, searchChannelStoreTests, th)
}

func testAutocompleteChannelByName(t *testing.T, th *SearchTestHelper) {
	return
}

func testAutocompleteChannelByDisplayName(t *testing.T, th *SearchTestHelper) {
	return
}

func testAutocompleteChannelByNameSplittedWithDashChar(t *testing.T, th *SearchTestHelper) {
	return
}

func testAutocompleteChannelByNameSplittedWithCommaChar(t *testing.T, th *SearchTestHelper) {
	return
}

func testAutocompleteChannelByNameSplittedWithUnderscoreChar(t *testing.T, th *SearchTestHelper) {
	return
}

func testAutocompleteChannelByDisplayNameSplittedByWhitespaces(t *testing.T, th *SearchTestHelper) {
	return
}

func testSearchChannelsInCaseInsensitiveManner(t *testing.T, th *SearchTestHelper) {
	return
}

func testSearchOnlyPublicChannels(t *testing.T, th *SearchTestHelper) {
	return
}

func testSearchShouldSupportWildcardAfterHyphen(t *testing.T, th *SearchTestHelper) {
	return
}

func testSearchShouldSupportHavingHyphenAsLastCharacter(t *testing.T, th *SearchTestHelper) {
	return
}

func testSearchShouldSupportAutocompleteWithArchivedChannels(t *testing.T, th *SearchTestHelper) {
	return
}
