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
		"Should be able to search inactive users",
		testShouldBeAbleToSearchInactiveUsers,
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
	t.Run("Autocomplete users with channel restrictions", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "", options)
		require.Nil(t, apperr)
		th.User.Sanitize(map[string]bool{})
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{th.User, userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
	t.Run("Autocomplete users with term and channel restrictions", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "alt", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
	t.Run("Autocomplete users with all channels restricted", func(t *testing.T) {
		options.ViewRestrictions = &model.ViewUsersRestrictions{Channels: []string{}}
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "", options)
		require.Nil(t, apperr)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
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
	t.Run("Should return results for users in the team", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "", options)
		require.Nil(t, apperr)
		th.User.Sanitize(map[string]bool{})
		th.User2.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{th.User}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{th.User2}, users.OutOfChannel)
	})
	t.Run("Should return empty because we're filtering all the teams", func(t *testing.T) {
		options.ViewRestrictions = &model.ViewUsersRestrictions{Teams: []string{}}
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "", options)
		require.Nil(t, apperr)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
}
func testShouldReturnNothingWithoutProperAccess(t *testing.T, th *SearchTestHelper) {
	options := &model.UserSearchOptions{
		AllowFullNames:        true,
		Limit:                 model.USER_SEARCH_DEFAULT_LIMIT,
		ListOfAllowedChannels: []string{th.ChannelBasic.Id},
	}
	t.Run("Should return results users for the defined channel in the list", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "", options)
		require.Nil(t, apperr)
		th.User.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{th.User}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
	t.Run("Should return empty because we're filtering all the channels", func(t *testing.T) {
		options.ListOfAllowedChannels = []string{}
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "", options)
		require.Nil(t, apperr)
		th.User.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
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
	t.Run("Should return results when the first name is unique", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "altfirstname", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
	t.Run("Should return results for in the channel and out of the channel with the same first name", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "basicfirstname", options)
		require.Nil(t, apperr)
		th.User.Sanitize(map[string]bool{})
		th.User2.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{th.User}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{th.User2}, users.OutOfChannel)
	})
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
	t.Run("Should return results when the last name is unique", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "altlastname", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
	t.Run("Should return results for in the channel and out of the channel with the same last name", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "basiclastname", options)
		require.Nil(t, apperr)
		th.User.Sanitize(map[string]bool{})
		th.User2.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{th.User}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{th.User2}, users.OutOfChannel)
	})
}
func testAutocompleteUserByNickName(t *testing.T, th *SearchTestHelper) {
	userAlternate, err := th.createUser("alternateusername", "alternatenickname", "firstname", "altlastname")
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
	t.Run("Should return results when the nickname is unique", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "alternatenickname", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
	t.Run("Should return users that share the same part of the nickname", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "basicnickname", options)
		require.Nil(t, apperr)
		th.User.Sanitize(map[string]bool{})
		th.User2.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{th.User}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{th.User2}, users.OutOfChannel)
	})
}
func testAutocompleteUserByEmail(t *testing.T, th *SearchTestHelper) {
	userAlternate, err := th.createUser("alternateusername", "alternatenickname", "firstname", "altlastname")
	require.Nil(t, err)
	userAlternate.Email = "useralt@test.email.com"
	_, apperr := th.Store.User().Update(userAlternate, false)
	require.Nil(t, apperr)
	defer th.deleteUser(userAlternate)
	err = th.addUserToTeams(userAlternate, []string{th.Team.Id})
	require.Nil(t, err)
	err = th.addUserToChannels(userAlternate, []string{th.ChannelBasic.Id})
	require.Nil(t, err)
	options := &model.UserSearchOptions{
		AllowEmails: true,
		Limit:       model.USER_SEARCH_DEFAULT_LIMIT,
	}
	// Should return results when the email is unique
	t.Run("Should return results when the email is unique", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "useralt@test.email.com", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
	// Should return users that share the same email user prefix
	t.Run("Should return users that share the same email user prefix", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "success_", options)
		require.Nil(t, apperr)
		th.User.Sanitize(map[string]bool{})
		th.User2.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{th.User}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{th.User2}, users.OutOfChannel)
	})
	// Should return users that share the same domain
	t.Run("Should return users that share the same email domain", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "simulator.amazon.com", options)
		require.Nil(t, apperr)
		th.User.Sanitize(map[string]bool{})
		th.User2.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{th.User}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{th.User2}, users.OutOfChannel)
	})
}
func testShouldNotMatchSpecificQueriesEmail(t *testing.T, th *SearchTestHelper) {
	options := &model.UserSearchOptions{
		AllowEmails: false,
		Limit:       model.USER_SEARCH_DEFAULT_LIMIT,
	}
	users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "success_", options)
	require.Nil(t, apperr)
	th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
}
func testAutocompleteUserByUsernameWithPoint(t *testing.T, th *SearchTestHelper) {
	userAlternate, err := th.createUser("alternate.username", "alternatenickname", "firstname", "altlastname")
	require.Nil(t, err)
	defer th.deleteUser(userAlternate)
	err = th.addUserToTeams(userAlternate, []string{th.Team.Id})
	require.Nil(t, err)
	err = th.addUserToChannels(userAlternate, []string{th.ChannelBasic.Id})
	require.Nil(t, err)
	options := &model.UserSearchOptions{
		Limit: model.USER_SEARCH_DEFAULT_LIMIT,
	}
	// Should return results when searching for the whole username with point
	t.Run("Should return results when searching for the whole username with point", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "alternate.username", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
	t.Run("Should return results when searching for part of the username including the point", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, ".username", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
}
func testAutocompleteUserByUsernameWithWhiteSpace(t *testing.T, th *SearchTestHelper) {
	userAlternate, err := th.createUser("alternate username", "alternatenickname", "firstname", "altlastname")
	require.Nil(t, err)
	defer th.deleteUser(userAlternate)
	err = th.addUserToTeams(userAlternate, []string{th.Team.Id})
	require.Nil(t, err)
	err = th.addUserToChannels(userAlternate, []string{th.ChannelBasic.Id})
	require.Nil(t, err)
	options := &model.UserSearchOptions{
		Limit: model.USER_SEARCH_DEFAULT_LIMIT,
	}
	t.Run("Should return results when searching for the whole username with point", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "alternate username", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
	t.Run("Should return results when searching for part of the username including the whitespace", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, " username", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
}
func testAutocompleteUserByUsernameWithUnderscore(t *testing.T, th *SearchTestHelper) {
	userAlternate, err := th.createUser("alternate_username", "alternatenickname", "firstname", "altlastname")
	require.Nil(t, err)
	defer th.deleteUser(userAlternate)
	err = th.addUserToTeams(userAlternate, []string{th.Team.Id})
	require.Nil(t, err)
	err = th.addUserToChannels(userAlternate, []string{th.ChannelBasic.Id})
	require.Nil(t, err)
	options := &model.UserSearchOptions{
		Limit: model.USER_SEARCH_DEFAULT_LIMIT,
	}
	t.Run("Should return results when searching for the whole username with underscore", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "alternate_username", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
	t.Run("Should return results when searching for part of the username including the underscore", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "_username", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
}
func testAutocompleteUserByUsernameWithHyphen(t *testing.T, th *SearchTestHelper) {
	userAlternate, err := th.createUser("alternate-username", "alternatenickname", "firstname", "altlastname")
	require.Nil(t, err)
	defer th.deleteUser(userAlternate)
	err = th.addUserToTeams(userAlternate, []string{th.Team.Id})
	require.Nil(t, err)
	err = th.addUserToChannels(userAlternate, []string{th.ChannelBasic.Id})
	require.Nil(t, err)
	options := &model.UserSearchOptions{
		Limit: model.USER_SEARCH_DEFAULT_LIMIT,
	}
	t.Run("Should return results when searching for the whole username with hyphen", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "alternate-username", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
	t.Run("Should return results when searching for part of the username including the hyphen", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "-username", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
}
func testShouldEscapePercentageCharacter(t *testing.T, th *SearchTestHelper) {
	userAlternate, err := th.createUser("alternateusername", "alternate%nickname", "firstname", "altlastname")
	require.Nil(t, err)
	defer th.deleteUser(userAlternate)
	err = th.addUserToTeams(userAlternate, []string{th.Team.Id})
	require.Nil(t, err)
	err = th.addUserToChannels(userAlternate, []string{th.ChannelBasic.Id})
	require.Nil(t, err)
	options := &model.UserSearchOptions{
		Limit: model.USER_SEARCH_DEFAULT_LIMIT,
	}
	users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "alternate%", options)
	require.Nil(t, apperr)
	userAlternate.Sanitize(map[string]bool{})
	th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
}
func testShouldEscapeDashCharacter(t *testing.T, th *SearchTestHelper) {
	userAlternate, err := th.createUser("alternate_username", "alternatenickname", "firstname", "altlastname")
	require.Nil(t, err)
	defer th.deleteUser(userAlternate)
	err = th.addUserToTeams(userAlternate, []string{th.Team.Id})
	require.Nil(t, err)
	err = th.addUserToChannels(userAlternate, []string{th.ChannelBasic.Id})
	require.Nil(t, err)
	options := &model.UserSearchOptions{
		Limit: model.USER_SEARCH_DEFAULT_LIMIT,
	}
	users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "alternate_", options)
	require.Nil(t, apperr)
	userAlternate.Sanitize(map[string]bool{})
	th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
}

func testShouldBeAbleToSearchInactiveUsers(t *testing.T, th *SearchTestHelper) {
	userAlternate, err := th.createUser("alternate-username", "alternatenickname", "firstname", "altlastname")
	require.Nil(t, err)
	userAlternate.DeleteAt = model.GetMillis()
	_, apperr := th.Store.User().Update(userAlternate, false)
	require.Nil(t, apperr)
	defer th.deleteUser(userAlternate)
	err = th.addUserToTeams(userAlternate, []string{th.Team.Id})
	require.Nil(t, err)
	err = th.addUserToChannels(userAlternate, []string{th.ChannelBasic.Id})
	require.Nil(t, err)
	options := &model.UserSearchOptions{
		AllowInactive: true,
		Limit:         model.USER_SEARCH_DEFAULT_LIMIT,
	}
	t.Run("Should return inactive users if we allow it", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "alternate-username", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
	t.Run("Shouldn't return return inactive users if we don't allow it", func(t *testing.T) {
		options.AllowInactive = false
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "alternate-username", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
}

func testShouldIgnoreLeadingAtSymbols(t *testing.T, th *SearchTestHelper) {
	options := &model.UserSearchOptions{
		Limit: model.USER_SEARCH_DEFAULT_LIMIT,
	}
	users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "@basicusername", options)
	require.Nil(t, apperr)
	th.User.Sanitize(map[string]bool{})
	th.User2.Sanitize(map[string]bool{})
	th.assertUsersMatchInAnyOrder(t, []*model.User{th.User}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{th.User2}, users.OutOfChannel)
}

func testSearchUsersShouldBeCaseInsensitive(t *testing.T, th *SearchTestHelper) {
	options := &model.UserSearchOptions{
		Limit: model.USER_SEARCH_DEFAULT_LIMIT,
	}
	users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "BaSiCUsErNaMe", options)
	require.Nil(t, apperr)
	th.User.Sanitize(map[string]bool{})
	th.User2.Sanitize(map[string]bool{})
	th.assertUsersMatchInAnyOrder(t, []*model.User{th.User}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{th.User2}, users.OutOfChannel)
}

func testSearchOneTwoCharUsersnameAndFirstLastNames(t *testing.T, th *SearchTestHelper) {
	userAlternate, err := th.createUser("ho", "alternatenickname", "zi", "k")
	require.Nil(t, err)
	userAlternate.Email = "x@example.com"
	_, apperr := th.Store.User().Update(userAlternate, false)
	require.Nil(t, apperr)
	defer th.deleteUser(userAlternate)
	err = th.addUserToTeams(userAlternate, []string{th.Team.Id})
	require.Nil(t, err)
	err = th.addUserToChannels(userAlternate, []string{th.ChannelBasic.Id})
	require.Nil(t, err)
	options := &model.UserSearchOptions{
		AllowFullNames: true,
		Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
	}
	t.Run("Should support two characters in the full name", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "zi", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
	t.Run("Should support two characters in the username", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "ho", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
	t.Run("Should support one character in the email", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "x", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
}

func testShouldSupportKoreanCharacters(t *testing.T, th *SearchTestHelper) {
	userAlternate, err := th.createUser("alternate-username", "alternate-nickname", "서강준", "안신원")
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
	t.Run("Should support hanja korean characters", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "서강준", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
	t.Run("Should support hangul korean characters", func(t *testing.T) {
		users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "안신원", options)
		require.Nil(t, apperr)
		userAlternate.Sanitize(map[string]bool{})
		th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
		th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
	})
}

func testSupportSearchUsersByContainingTerms(t *testing.T, th *SearchTestHelper) {
	options := &model.UserSearchOptions{
		Limit: model.USER_SEARCH_DEFAULT_LIMIT,
	}
	users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "nick", options)
	require.Nil(t, apperr)
	th.User.Sanitize(map[string]bool{})
	th.User2.Sanitize(map[string]bool{})
	th.assertUsersMatchInAnyOrder(t, []*model.User{th.User}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{th.User2}, users.OutOfChannel)
}
func testSearchWithHyphenAtTheEndOfTheTerm(t *testing.T, th *SearchTestHelper) {
	userAlternate, err := th.createUser("alternate-username", "alternate-nickname", "altfirst", "altlast")
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
	users, apperr := th.Store.User().AutocompleteUsersInChannel(th.Team.Id, th.ChannelBasic.Id, "alternate-", options)
	require.Nil(t, apperr)
	userAlternate.Sanitize(map[string]bool{})
	th.assertUsersMatchInAnyOrder(t, []*model.User{userAlternate}, users.InChannel)
	th.assertUsersMatchInAnyOrder(t, []*model.User{}, users.OutOfChannel)
}
