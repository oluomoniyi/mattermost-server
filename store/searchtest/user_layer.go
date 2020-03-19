package searchtest

import (
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/store"
	"github.com/stretchr/testify/require"
)

var searchUserStoreTests = []searchTest{
	{
		"Should retrieve all users in a channel if the search term is empty",
		testGetAllUsersInChannelWithEmptyTerm,
		[]string{ENGINE_ALL},
	},
	{
		"Should honor channel restrictions when autocompleting users",
		testHonorChannelRestrictionsAutocompletingUsers,
		[]string{ENGINE_ALL},
	},
	{
		"Should honor team restrictions when autocompleting users",
		testHonorTeamRestrictionsAutocompletingUsers,
		[]string{ENGINE_ALL},
	},
	{
		"Should return nothing if the user can't access the channels of a given search",
		testShouldReturnNothingWithoutProperAccess,
		[]string{ENGINE_ALL},
	},
	{
		"Should autocomplete for user using username",
		testAutocompleteUserByUsername,
		[]string{ENGINE_ALL},
	},
	{
		"Should autocomplete user searching by first name",
		testAutocompleteUserByFirstName,
		[]string{ENGINE_ALL},
	},
	{
		"Should autocomplete user searching by last name",
		testAutocompleteUserByLastName,
		[]string{ENGINE_ALL},
	},
	{
		"Should autocomplete for user using nickname",
		testAutocompleteUserByNickName,
		[]string{ENGINE_ALL},
	},
	{
		"Should autocomplete for user using email",
		testAutocompleteUserByEmail,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able not to match specific queries with mail",
		testShouldNotMatchSpecificQueriesEmail,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to autocomplete a user by part of its username splitted by point",
		testAutocompleteUserByUsernameWithPoint,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to autocomplete a user by part of its username splitted by comma",
		testAutocompleteUserByUsernameWithComma,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to autocomplete a user by part of its username splitted by whitespace",
		testAutocompleteUserByUsernameWithWhiteSpace,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to autocomplete a user by part of its username splitted by underscore",
		testAutocompleteUserByUsernameWithUnderscore,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to autocomplete a user by part of its username splitted by hyphen",
		testAutocompleteUserByUsernameWithHyphen,
		[]string{ENGINE_ALL},
	},
	{
		"Should escape the percentage character",
		testShouldEscapePercentageCharacter,
		[]string{ENGINE_ALL},
	},
	{
		"Should escape the dash character",
		testShouldEscapeDashCharacter,
		[]string{ENGINE_ALL},
	},
	{
		"Should be able to search deactivated users",
		testShouldBeAbleToSearchDeactivatedUsers,
		[]string{ENGINE_ALL},
	},
	{
		"Should ignore leading @ when searching users",
		testShouldIgnoreLeadingAtSymbols,
		[]string{ENGINE_ALL},
	},
	{
		"Should search users in a case insensitive manner",
		testSearchUsersShouldBeCaseInsensitive,
		[]string{ENGINE_ALL},
	},
	{
		"Should support one or two character usernames and first/last names in search",
		testSearchOneTwoCharUsersnameAndFirstLastNames,
		[]string{ENGINE_ALL},
	},
	{
		"Should support Korean characters",
		testShouldSupportKoreanCharacters,
		[]string{ENGINE_ALL},
	},
	{
		"Should support searching for users containing the term not only starting with it",
		testSupportSearchUsersByContainingTerms,
		[]string{ENGINE_ALL},
	},
	{
		"Should support search with a hyphen at the end of the term",
		testSearchWithHyphenAtTheEndOfTheTerm,
		[]string{ENGINE_ALL},
	},
}

func TestSearchUserStore(t *testing.T, s store.Store, testEngine *SearchTestEngine) {
	th := &SearchTestHelper{
		Store: s,
	}
	err := th.SetupBasicFixtures()
	require.Nil(t, err)
	defer th.CleanFixtures()
	runTestSearch(t, testEngine, searchUserStoreTests, th)
}

func testGetAllUsersInChannelWithEmptyTerm(t *testing.T, th *SearchTestHelper) {
	options := &model.UserSearchOptions{
		AllowFullNames: true,
		Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
	}
	users, err := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "", options)
	require.Nil(t, err)
	th.User.Sanitize(map[string]bool{})
	th.User2.Sanitize(map[string]bool{})
	th.assertUsersMatchInAnyOrder(t, []*model.User{th.User}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{th.User2}, users.OutOfChannel)
}
func testHonorChannelRestrictionsAutocompletingUsers(t *testing.T, th *SearchTestHelper) {
	userAlternate, err := th.createUser("user-alternate", "user-alternate", "user", "alternate")
	require.Nil(t, err)
	defer th.deleteUser(userAlternate)
	err = th.addUserToTeams(userAlternate, []string{th.Team.Id})
	require.Nil(t, err)
	err = th.addUserToChannels(userAlternate, []string{th.ChannelBasic.Id})
	require.Nil(t, err)
	options := &model.UserSearchOptions{
		AllowFullNames:   true,
		Limit:            model.USER_SEARCH_DEFAULT_LIMIT,
		ViewRestrictions: &model.ViewUsersRestrictions{Channels: []string{th.ChannelBasic.Id}},
	}
	// Autocomplete users with channel restrictions
	users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "", options)
	require.Nil(t, apperr)
	th.User.Sanitize(map[string]bool{})
	userAlternate.Sanitize(map[string]bool{})
	th.assertUsersMatchInAnyOrder(t, []*model.User{th.User, userAlternate}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	// Autocomplete users with term and channel restrictions
	users, apperr = th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "alt", options)
	require.Nil(t, apperr)
	userAlternate.Sanitize(map[string]bool{})
	th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	// Autocomplete users with all channels restricted
	options.ViewRestrictions = &model.ViewUsersRestrictions{Channels: []string{}}
	users, apperr = th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "", options)
	require.Nil(t, apperr)
	th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
}
func testHonorTeamRestrictionsAutocompletingUsers(t *testing.T, th *SearchTestHelper) {
	userAlternate, err := th.createUser("user-alternate", "user-alternate", "user", "alternate")
	defer th.deleteUser(userAlternate)
	require.Nil(t, err)
	err = th.addUserToTeams(userAlternate, []string{th.AnotherTeam.Id})
	require.Nil(t, err)
	err = th.addUserToChannels(userAlternate, []string{th.ChannelAnotherTeam.Id})
	require.Nil(t, err)
	options := &model.UserSearchOptions{
		AllowFullNames:   true,
		Limit:            model.USER_SEARCH_DEFAULT_LIMIT,
		ViewRestrictions: &model.ViewUsersRestrictions{Teams: []string{th.Team.Id}},
	}
	// Should return results for users in the team
	users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "", options)
	require.Nil(t, apperr)
	th.User.Sanitize(map[string]bool{})
	th.User2.Sanitize(map[string]bool{})
	th.assertUsersMatchInAnyOrder(t, []*model.User{th.User}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{th.User2}, users.OutOfChannel)
	// Should return empty because we're filtering all the teams
	options.ViewRestrictions = &model.ViewUsersRestrictions{Teams: []string{}}
	users, apperr = th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "", options)
	require.Nil(t, apperr)
	th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
}
func testShouldReturnNothingWithoutProperAccess(t *testing.T, th *SearchTestHelper) {
	options := &model.UserSearchOptions{
		AllowFullNames:        true,
		Limit:                 model.USER_SEARCH_DEFAULT_LIMIT,
		ListOfAllowedChannels: []string{th.ChannelBasic.Id},
	}
	// Should return results users for the defined channel in the list
	users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "", options)
	require.Nil(t, apperr)
	th.User.Sanitize(map[string]bool{})
	th.assertUsersMatchInAnyOrder(t, []*model.User{th.User}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	options.ListOfAllowedChannels = []string{}
	// Should return empty because we're filtering all the channels
	users, apperr = th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "", options)
	require.Nil(t, apperr)
	th.User.Sanitize(map[string]bool{})
	th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
}
func testAutocompleteUserByUsername(t *testing.T, th *SearchTestHelper) {
	userAlternate, err := th.createUser("alternateusername", "alternatenick", "user", "alternate")
	require.Nil(t, err)
	defer th.deleteUser(userAlternate)
	err = th.addUserToTeams(userAlternate, []string{th.Team.Id})
	require.Nil(t, err)
	err = th.addUserToChannels(userAlternate, []string{th.ChannelBasic.Id})
	require.Nil(t, err)
	options := &model.UserSearchOptions{
		AllowFullNames: false,
		Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
	}
	users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "basicusername", options)
	require.Nil(t, apperr)
	th.User.Sanitize(map[string]bool{})
	th.User2.Sanitize(map[string]bool{})
	th.assertUsersMatchInAnyOrder(t, []*model.User{th.User}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{th.User2}, users.OutOfChannel)
}
func testAutocompleteUserByFirstName(t *testing.T, th *SearchTestHelper) {
	userAlternate, err := th.createUser("user-alternate", "user-alternate", "altfirstname", "lastname")
	require.Nil(t, err)
	defer th.deleteUser(userAlternate)
	err = th.addUserToTeams(userAlternate, []string{th.Team.Id})
	require.Nil(t, err)
	err = th.addUserToChannels(userAlternate, []string{th.ChannelBasic.Id})
	require.Nil(t, err)
	options := &model.UserSearchOptions{
		AllowFullNames: true,
		Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
	}
	// Should return results when the first name is unique
	users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "altfirstname", options)
	require.Nil(t, apperr)
	userAlternate.Sanitize(map[string]bool{})
	th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	// Should return results for in the channel and out of the channel with the same first name
	users, apperr = th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "basicfirstname", options)
	require.Nil(t, apperr)
	th.User.Sanitize(map[string]bool{})
	th.User2.Sanitize(map[string]bool{})
	th.assertUsersMatchInAnyOrder(t, []*model.User{th.User}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{th.User2}, users.OutOfChannel)
}
func testAutocompleteUserByLastName(t *testing.T, th *SearchTestHelper) {
	userAlternate, err := th.createUser("user-alternate", "user-alternate", "firstname", "altlastname")
	require.Nil(t, err)
	defer th.deleteUser(userAlternate)
	err = th.addUserToTeams(userAlternate, []string{th.Team.Id})
	require.Nil(t, err)
	err = th.addUserToChannels(userAlternate, []string{th.ChannelBasic.Id})
	require.Nil(t, err)
	options := &model.UserSearchOptions{
		AllowFullNames: true,
		Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
	}
	// Should return results when the first name is unique
	users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "altlastname", options)
	require.Nil(t, apperr)
	userAlternate.Sanitize(map[string]bool{})
	th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	// Should return results for in the channel and out of the channel with the same first name
	users, apperr = th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "basiclastname", options)
	require.Nil(t, apperr)
	th.User.Sanitize(map[string]bool{})
	th.User2.Sanitize(map[string]bool{})
	th.assertUsersMatchInAnyOrder(t, []*model.User{th.User}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{th.User2}, users.OutOfChannel)
}
func testAutocompleteUserByNickName(t *testing.T, th *SearchTestHelper) {
	return
}
func testAutocompleteUserByEmail(t *testing.T, th *SearchTestHelper) {
	return
}
func testShouldNotMatchSpecificQueriesEmail(t *testing.T, th *SearchTestHelper) {
	return
}
func testAutocompleteUserByUsernameWithPoint(t *testing.T, th *SearchTestHelper) {
	return
}
func testAutocompleteUserByUsernameWithComma(t *testing.T, th *SearchTestHelper) {
	return
}
func testAutocompleteUserByUsernameWithWhiteSpace(t *testing.T, th *SearchTestHelper) {
	return
}
func testAutocompleteUserByUsernameWithUnderscore(t *testing.T, th *SearchTestHelper) {
	return
}
func testAutocompleteUserByUsernameWithHyphen(t *testing.T, th *SearchTestHelper) {
	return
}
func testShouldEscapePercentageCharacter(t *testing.T, th *SearchTestHelper) {
	return
}
func testShouldEscapeDashCharacter(t *testing.T, th *SearchTestHelper) {
	return
}
func testShouldBeAbleToSearchInactiveUsers(t *testing.T, th *SearchTestHelper) {
	return
}
func testShouldBeAbleToSearchDeactivatedUsers(t *testing.T, th *SearchTestHelper) {
	return
}
func testShouldIgnoreLeadingAtSymbols(t *testing.T, th *SearchTestHelper) {
	return
}
func testSearchUsersShouldBeCaseInsensitive(t *testing.T, th *SearchTestHelper) {
	return
}
func testSearchOneTwoCharUsersnameAndFirstLastNames(t *testing.T, th *SearchTestHelper) {
	return
}
func testShouldSupportKoreanCharacters(t *testing.T, th *SearchTestHelper) {
	return
}
func testSupportSearchUsersByContainingTerms(t *testing.T, th *SearchTestHelper) {
	return
}
func testSearchWithHyphenAtTheEndOfTheTerm(t *testing.T, th *SearchTestHelper) {
	return
}
