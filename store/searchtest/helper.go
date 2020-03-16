// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package searchtest

import (
	"fmt"
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/store"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
)

type SearchTestHelper struct {
	Store                       store.Store
	Team                        *model.Team
	AnotherTeam                 *model.Team
	User                        *model.User
	User2                       *model.User
	ChannelA                    *model.Channel
	ChannelAlternate            *model.Channel
	ChannelPrivate              *model.Channel
	ChannelHyphen               *model.Channel
	ChannelWhiteSpace           *model.Channel
	ChannelWhiteSpaceAndHyphens *model.Channel
	ChannelPurpose              *model.Channel
	ChannelOtherTeam            *model.Channel
	ChannelDeleted              *model.Channel
}

func (th *SearchTestHelper) InitFixtures() error {
	// Create teams
	team, err := th.createTeam("test-team", "test tea", model.TEAM_OPEN)
	if err != nil {
		return err
	}
	anotherTeam, err := th.createTeam("another-test-team", "Another test team", model.TEAM_OPEN)
	if err != nil {
		return err
	}

	// Create users
	user, err := th.createUser("user1", "user1", "test", "user1")
	if err != nil {
		return err
	}
	user2, err := th.createUser("user2", "user2", "test", "user2")
	if err != nil {
		return err
	}

	// Create channels
	channelA, err := th.createChannel(team.Id, th.createChannelName(), "ChannelA", "", model.CHANNEL_OPEN, false)
	if err != nil {
		return err
	}
	channelAlternate, err := th.createChannel(team.Id, th.createChannelName(), "ChannelA (alternate)", "", model.CHANNEL_OPEN, false)
	if err != nil {
		return err
	}
	channelPrivate, err := th.createChannel(team.Id, th.createChannelName(), "ChannelB", "", model.CHANNEL_PRIVATE, false)
	if err != nil {
		return err
	}
	channelHyphen, err := th.createChannel(team.Id, "off-topic", "Off-Topic", "", model.CHANNEL_OPEN, false)
	if err != nil {
		return err
	}
	channelWhiteSpace, err := th.createChannel(team.Id, "town-square", "Town Square", "", model.CHANNEL_OPEN, false)
	if err != nil {
		return err
	}
	channelWhiteSpaceAndHyphens, err := th.createChannel(team.Id, "native-mobile-apps", "Native Mobile Apps", "", model.CHANNEL_OPEN, false)
	if err != nil {
		return err
	}
	channelPurpose, err := th.createChannel(team.Id, "another-town-square", "Town Square", "This can now be searchable!", model.CHANNEL_OPEN, false)
	if err != nil {
		return err
	}
	channelDeleted, err := th.createChannel(team.Id, th.createChannelName(), "ChannelA (deleted)", "", model.CHANNEL_OPEN, true)
	if err != nil {
		return err
	}

	err = th.addUserToTeams(user, []string{team.Id})
	if err != nil {
		return err
	}

	err = th.addUserToChannels(user, []string{channelA.Id, channelAlternate.Id, channelPrivate.Id,
		channelHyphen.Id, channelWhiteSpace.Id, channelWhiteSpaceAndHyphens.Id,
		channelPurpose.Id, channelDeleted.Id})
	if err != nil {
		return err
	}

	th.Team = team
	th.AnotherTeam = anotherTeam
	th.User = user
	th.User2 = user2
	th.ChannelA = channelA
	th.ChannelAlternate = channelAlternate
	th.ChannelPrivate = channelPrivate
	th.ChannelHyphen = channelHyphen
	th.ChannelWhiteSpace = channelWhiteSpace
	th.ChannelWhiteSpaceAndHyphens = channelWhiteSpaceAndHyphens
	th.ChannelPurpose = channelPurpose
	th.ChannelDeleted = channelDeleted

	return nil
}

func (th *SearchTestHelper) CleanFixtures() error {
	err := th.deleteChannels([]*model.Channel{
		th.ChannelA, th.ChannelAlternate, th.ChannelPrivate,
		th.ChannelHyphen, th.ChannelWhiteSpace, th.ChannelWhiteSpaceAndHyphens,
		th.ChannelPurpose, th.ChannelDeleted,
	})
	if err != nil {
		return err
	}

	err = th.deleteTeam(th.Team)
	if err != nil {
		return err
	}

	err = th.deleteTeam(th.AnotherTeam)
	if err != nil {
		return err
	}

	err = th.deleteUser(th.User)
	if err != nil {
		return err
	}

	err = th.deleteUser(th.User2)
	if err != nil {
		return err
	}

	return nil
}

func (th *SearchTestHelper) createTeam(name, displayName, teamType string) (*model.Team, error) {
	team, appError := th.Store.Team().Save(&model.Team{
		Name:        name,
		DisplayName: displayName,
		Type:        teamType,
	})
	if appError != nil {
		return nil, errors.New(appError.Error())
	}

	return team, nil
}

func (th *SearchTestHelper) deleteTeam(team *model.Team) error {
	appError := th.Store.Team().RemoveAllMembersByTeam(team.Id)
	if appError != nil {
		return errors.New(appError.Error())
	}
	appError = th.Store.Team().PermanentDelete(team.Id)
	if appError != nil {
		return errors.New(appError.Error())
	}

	return nil
}

func (th *SearchTestHelper) createTeamMember(teamID, userID string) *model.TeamMember {
	return &model.TeamMember{
		TeamId: teamID,
		UserId: userID,
	}
}

func (th *SearchTestHelper) makeEmail() string {
	return "success_" + model.NewId() + "@simulator.amazonseth.Store.com"
}

func (th *SearchTestHelper) createUser(username, nickname, firstName, lastName string) (*model.User, error) {
	user, appError := th.Store.User().Save(&model.User{
		Username:  username,
		Password:  username,
		Nickname:  nickname,
		FirstName: firstName,
		LastName:  lastName,
		Email:     th.makeEmail(),
	})
	if appError != nil {
		return nil, errors.New(appError.Error())
	}

	return user, nil
}

func (th *SearchTestHelper) deleteUser(user *model.User) error {
	appError := th.Store.User().PermanentDelete(user.Id)
	if appError != nil {
		return errors.New(appError.Error())
	}

	return nil
}

func (th *SearchTestHelper) createChannelName() string {
	return "zz" + model.NewId() + "b"
}

func (th *SearchTestHelper) createChannel(teamID, name, displayName, purpose, channelType string, deleted bool) (*model.Channel, error) {
	channel, appError := th.Store.Channel().Save(&model.Channel{
		TeamId:      teamID,
		DisplayName: displayName,
		Name:        name,
		Type:        channelType,
		Purpose:     purpose,
	}, 999)
	if appError != nil {
		return nil, errors.New(appError.Error())
	}

	if deleted {
		appError := th.Store.Channel().Delete(channel.Id, model.GetMillis())
		if appError != nil {
			return nil, errors.New(appError.Error())
		}
	}

	return channel, nil
}

func (th *SearchTestHelper) deleteChannel(channel *model.Channel) error {
	appError := th.Store.Channel().PermanentDeleteMembersByChannel(channel.Id)
	if appError != nil {
		return errors.New(appError.Error())
	}

	appError = th.Store.Channel().PermanentDelete(channel.Id)
	if appError != nil {
		return errors.New(appError.Error())
	}

	return nil
}

func (th *SearchTestHelper) saveChannels(channels []*model.Channel) error {
	for _, channel := range channels {
		_, err := th.Store.Channel().Save(channel, 100)
		if err != nil {
			return errors.New(err.Error())
		}
	}

	return nil
}

func (th *SearchTestHelper) deleteChannels(channels []*model.Channel) error {
	for _, channel := range channels {
		err := th.deleteChannel(channel)
		if err != nil {
			return err
		}
	}

	return nil
}

func (th *SearchTestHelper) createPost(userID, channelID, message, hashtags string, createAt int64) (*model.Post, error) {
	var creationTime int64 = 1000000
	if createAt > 0 {
		creationTime = createAt
	}
	post, appError := th.Store.Post().Save(&model.Post{
		Message:       message,
		ChannelId:     channelID,
		PendingPostId: model.NewId() + ":" + fmt.Sprint(model.GetMillis()),
		UserId:        userID,
		Hashtags:      hashtags,
		CreateAt:      creationTime,
	})
	if appError != nil {
		return nil, errors.New(appError.Error())
	}

	return post, nil
}

func (th *SearchTestHelper) addUserToTeams(user *model.User, teamIds []string) error {
	for _, teamId := range teamIds {
		_, err := th.Store.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: user.Id}, -1)
		if err != nil {
			return errors.New(err.Error())
		}
	}

	return nil
}

func (th *SearchTestHelper) addUserToChannels(user *model.User, channelIds []string) error {

	for _, channelId := range channelIds {
		_, err := th.Store.Channel().SaveMember(&model.ChannelMember{
			ChannelId:   channelId,
			UserId:      user.Id,
			NotifyProps: model.GetDefaultChannelNotifyProps(),
		})
		if err != nil {
			return errors.New(err.Error())
		}
	}

	return nil
}

func (th *SearchTestHelper) assertUsers(t *testing.T, expected, actual []*model.User) {
	expectedUsernames := make([]string, 0, len(expected))
	for _, user := range expected {
		expectedUsernames = append(expectedUsernames, user.Username)
	}

	actualUsernames := make([]string, 0, len(actual))
	for _, user := range actual {
		actualUsernames = append(actualUsernames, user.Username)
	}

	if assert.Equal(t, expectedUsernames, actualUsernames) {
		assert.Equal(t, expected, actual)
	}
}

func (th *SearchTestHelper) assertUsersMatchInAnyOrder(t *testing.T, expected, actual []*model.User) {
	expectedUsernames := make([]string, 0, len(expected))
	for _, user := range expected {
		expectedUsernames = append(expectedUsernames, user.Username)
	}

	actualUsernames := make([]string, 0, len(actual))
	for _, user := range actual {
		actualUsernames = append(actualUsernames, user.Username)
	}

	if assert.ElementsMatch(t, expectedUsernames, actualUsernames) {
		assert.ElementsMatch(t, expected, actual)
	}
}

func (th *SearchTestHelper) checkPostInSearchResults(t *testing.T, postId string, searchResults []string) {
	t.Helper()
	assert.Contains(t, searchResults, postId, "Did not find expected post in search results.")
}

func (th *SearchTestHelper) checkPostNotInSearchResults(t *testing.T, postId string, searchResults []string) {
	t.Helper()
	assert.NotContains(t, searchResults, postId, "Found post in search results that should not be there.")
}

func (th *SearchTestHelper) checkMatchesEqual(t *testing.T, expected model.PostSearchMatches, actual map[string][]string) {
	a := assert.New(t)

	a.Len(actual, len(expected), "Received matches for a different number of posts")

	for postId, expectedMatches := range expected {
		a.ElementsMatch(expectedMatches, actual[postId], fmt.Sprintf("%v: expected %v, got %v", postId, expectedMatches, actual[postId]))
	}
}

type ByChannelDisplayName model.ChannelList

func (s ByChannelDisplayName) Len() int { return len(s) }
func (s ByChannelDisplayName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByChannelDisplayName) Less(i, j int) bool {
	if s[i].DisplayName != s[j].DisplayName {
		return s[i].DisplayName < s[j].DisplayName
	}

	return s[i].Id < s[j].Id
}
