package endpoints

import "errors"

//	func (u *UserDto) Validate() error {
//		return nil
//	}
func (m *TeamMemberDto) Validate() error {
	if m.Id == "" {
		return errors.New("Id is required in TeamMember object")
	}
	if m.Name == "" {
		return errors.New("Name is required in TeamMember object")
	}

	return nil
}
func (t *TeamDto) Validate() error {
	if t.TeamName == "" {
		return errors.New("TeamName is required in Team object")
	}

	if len(t.Members) == 0 {
		return errors.New("[]Members object is nil or it's len is 0")
	}
	for _, m := range t.Members {
		if err := m.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (su *SetUserStatusDto) Validate() error {
	if su.UserId == "" {
		return errors.New("UserId is required in body")
	}

	return nil
}

func (pr *CreatePullRequestDto) Validate() error {
	if pr.PullRequestId == "" {
		return errors.New("PullRequestId is required in body")
	}

	if pr.PullRequestName == "" {
		return errors.New("PullRequestName is required in body")
	}

	if pr.AuthorId == "" {
		return errors.New("AuthorId is required in body")
	}
	return nil
}

func (mr *MergePullRequestDto) Validate() error {
	if mr.PullRequestId == "" {
		return errors.New("PullRequestId is required in body")
	}

	return nil
}

func (re *ReassignDto) Validate() error {
	if re.PullRequestId == "" {
		return errors.New("PullRequestId is required in body")
	}

	if re.OldUserId == "" {
		return errors.New("OldUserId is required in body")
	}

	return nil
}
