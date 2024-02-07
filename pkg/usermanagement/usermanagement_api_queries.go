package usermanagement

import (
	"context"
)

// GET User(s)

func (a *Usermanagement) UserManagementGetUsers(
	authenticationDomainIDs []string,
	userIDs []string,
	nameContains string,
	emailContains string,
) (*UserManagementAuthenticationDomains, error) {
	return a.UserManagementGetUsersWithContext(context.Background(),
		authenticationDomainIDs,
		userIDs,
		nameContains,
		emailContains,
	)
}

func (a *Usermanagement) UserManagementGetUsersWithContext(
	ctx context.Context,
	authenticationDomainIDs []string,
	userIDs []string,
	nameContains string,
	emailContains string,
) (*UserManagementAuthenticationDomains, error) {

	resp := authenticationDomainsResponse{}
	vars := map[string]interface{}{
		"authenticationDomainIDs": authenticationDomainIDs,
		"userIDs":                 userIDs,
		"nameContains":            nameContains,
		"emailContains":           emailContains,
	}

	if len(authenticationDomainIDs) == 0 {
		delete(vars, "authenticationDomainIDs")
	}

	if len(userIDs) == 0 {
		delete(vars, "userIDs")
	}

	if nameContains == "" {
		delete(vars, "nameContains")
	}

	if emailContains == "" {
		delete(vars, "emailContains")
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getUsersQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.Actor.Organization.UserManagement.AuthenticationDomains, nil
}

const getUsersQuery = `query(
  $authenticationDomainIDs: [ID!],
  $userIDs: [ID!],
  $nameContains: String,
  $emailContains: String
)
{
  actor {
    organization {
      userManagement {
        authenticationDomains(id: $authenticationDomainIDs) {
          authenticationDomains {
            users(
              filter: {
                email: { contains: $emailContains }
                name: { contains: $nameContains }
                id: { in: $userIDs }
              }
            ) {
              users {
                id
                emailVerificationState
                email
                lastActive
                name
                timeZone
                pendingUpgradeRequest {
                  id
                  message
                  requestedUserType {
                    displayName
                    id
                  }
                }
                type {
                  displayName
                  id
                }
                groups {
                  groups {
                    id
                    displayName
                  }
                }
              }
            }
			id
			name
          }
		nextCursor
		totalCount
        }
      }
    }
  }
}`

// GET Group(s) With User(s)

func (a *Usermanagement) UserManagementGetGroupsWithUsers(
	authenticationDomainIDs []string,
	groupIDs []string,
	nameContains string,
) (*UserManagementAuthenticationDomains, error) {
	return a.UserManagementGetGroupsWithUsersWithContext(context.Background(),
		authenticationDomainIDs,
		groupIDs,
		nameContains,
	)
}

func (a *Usermanagement) UserManagementGetGroupsWithUsersWithContext(
	ctx context.Context,
	authenticationDomainIDs []string,
	groupIDs []string,
	nameContains string,
) (*UserManagementAuthenticationDomains, error) {

	resp := authenticationDomainsResponse{}
	vars := map[string]interface{}{
		"authenticationDomainIDs": authenticationDomainIDs,
		"groupIDs":                groupIDs,
		"nameContains":            nameContains,
	}

	if len(authenticationDomainIDs) == 0 {
		delete(vars, "authenticationDomainIDs")
	}

	if len(groupIDs) == 0 {
		delete(vars, "groupIDs")
	}

	if nameContains == "" {
		delete(vars, "nameContains")
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getGroupsWithUsersQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.Actor.Organization.UserManagement.AuthenticationDomains, nil
}

const getGroupsWithUsersQuery = `query (
  $authenticationDomainIDs: [ID!]
  $groupIDs: [ID!]
  $nameContains: String
) {
  actor {
    organization {
      userManagement {
        authenticationDomains(id: $authenticationDomainIDs) {
          authenticationDomains {
            groups(
              filter: {
                displayName: { contains: $nameContains }
                id: { in: $groupIDs }
              }
            ) {
              groups {
                displayName
                id
                users {
                  users {
                    email
                    id
                    name
                    timeZone
                  }
                }
              }
            }
            id
            name
          }
          nextCursor
          totalCount
        }
      }
    }
  }
}
`
