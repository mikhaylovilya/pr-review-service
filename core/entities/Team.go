package entities

import (
	"errors"
)

type Team struct {
	TeamName string
	Members  []User
}

func NewTeam(teamName string, members []User) (*Team, error) {
	if len(members) == 0 {
		return &Team{}, errors.New("Members slice is nil or it's len is 0")
	}

	return &Team{
		TeamName: teamName,
		Members:  members,
	}, nil
}
