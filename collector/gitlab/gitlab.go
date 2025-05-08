package gitlabcollector

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"vislab/libs/check"
	"vislab/sources/gitlab"
	gitlabTypes "vislab/sources/gitlab/types"
	"vislab/storage"
	"vislab/types"
)

func getNeededGroups(chosenGroups []string, git *gitlab.Client) ([]*gitlabTypes.Group, error) {
	groups, _, err := git.Groups.ListAll(context.Background(), &gitlabTypes.ListGroupsOptions{})
	if err != nil {
		return nil, err
	}
	slog.Debug("received all gitlab groups", "groups", fmt.Sprintf("%s", groups))

	neededGroups := groups

	if len(chosenGroups) > 0 {
		neededGroups = []*gitlabTypes.Group{}

		for _, flagGroup := range chosenGroups {
			if !slices.ContainsFunc(groups, func(group *gitlabTypes.Group) bool {
				if group.Path == flagGroup {
					neededGroups = append(neededGroups, group)
					return true
				}
				return false
			}) {
				slog.Error("group not found in gitlab", "group", flagGroup)
			}
		}
	}

	slog.Debug("got needed groups", "groups", fmt.Sprintf("%s", neededGroups))
	return neededGroups, nil
}

func getNeededProjects(neededGroups []*gitlabTypes.Group, git *gitlab.Client) ([]*gitlabTypes.Project, error) {
	neededProjects := []*gitlabTypes.Project{}

	for _, group := range neededGroups {
		projects, _, err := git.Groups.ListAllProjects(context.Background(), group.ID, &gitlabTypes.ListProjectsOptions{})
		if err != nil {
			slog.Error("failed to get projects for group", "err", err, "group", group.Path)
		}

		neededProjects = append(neededProjects, projects...)
	}

	slog.Debug("got needed projects", "projects", fmt.Sprintf("%s", neededProjects))
	return neededProjects, nil
}

func isAlreadyExist(ctx context.Context, service *types.Service, storage storage.Storage) (bool, error) {
	dbService, err := storage.Service().Get(ctx, *service.Name)
	if err != nil {
		if strings.Contains(err.Error(), "not found") { // TODO: add cool err handle
			return false, nil

		}
		return false, err
	}

	if check.ComparePointers(service.LatestTag, dbService.LatestTag) {
		return true, nil
	}

	return false, nil
}
