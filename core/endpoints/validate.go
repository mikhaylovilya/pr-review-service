package endpoints

import "errors"

func (u *UserDto) ValidateUserDto() error {
	return nil
}
func (m *TeamMemberDto) ValidateTeamMemberDto() error {
	if m.Id == "" {
		return errors.New("Id is required in TeamMember object")
	}
	if m.Name == "" {
		return errors.New("Name is required in TeamMember object")
	}

	return nil
}
func (t *TeamDto) ValidateTeamDto() error {
	if t.TeamName == "" {
		return errors.New("TeamName is required in Team object")
	}

	if len(t.Members) == 0 {
		return errors.New("[]Members object is nil or it's len is 0")
	}
	for _, m := range t.Members {
		if err := m.ValidateTeamMemberDto(); err != nil {
			return err
		}
	}
	return nil
}
func (u *PullRequestShortDto) ValidatePullRequestDto() error {
	return nil
}
