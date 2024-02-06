package usermanagement

import (
	"context"

	"github.com/newrelic/newrelic-client-go/v2/pkg/errors"
)

// An "authentication domain" is a grouping of New Relic users governed by the same user management settings, like how they're provisioned (added and updated), how they're authenticated (logged in), session settings, and how user upgrades are managed.
func (a *Usermanagement) GetAuthenticationDomains(
	cursor string,
	iD []string,
) (*UserManagementAuthenticationDomains, error) {
	return a.GetAuthenticationDomainsWithContext(context.Background(),
		cursor,
		iD,
	)
}

// An "authentication domain" is a grouping of New Relic users governed by the same user management settings, like how they're provisioned (added and updated), how they're authenticated (logged in), session settings, and how user upgrades are managed.
func (a *Usermanagement) GetAuthenticationDomainsWithContext(
	ctx context.Context,
	cursor string,
	iD []string,
) (*UserManagementAuthenticationDomains, error) {

	resp := authenticationDomainsResponse{}
	vars := map[string]interface{}{
		"cursor": cursor,
		"id":     iD,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getAuthenticationDomainsQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.Actor.Organization.UserManagement.AuthenticationDomains, nil
}

const getAuthenticationDomainsQuery = `query(
	$id: [ID!],
) { actor { organization { userManagement { authenticationDomains(
	id: $id,
) {
	authenticationDomains {
		id
		name
		provisioningType
	}
	nextCursor
	totalCount
} } } } }`

// GET Groups

func (a *Usermanagement) GetGroups(
	authenticationDomainIDs []string,
	groupIDs []string,
	name string,
) (*UserManagementAuthenticationDomains, error) {
	return a.GetGroupsWithContext(context.Background(),
		authenticationDomainIDs,
		groupIDs,
		name,
	)
}

func (a *Usermanagement) GetGroupsWithContext(
	ctx context.Context,
	authenticationDomainIDs []string,
	groupIDs []string,
	name string,
) (*UserManagementAuthenticationDomains, error) {

	resp := authenticationDomainsResponse{}
	vars := map[string]interface{}{
		"authenticationDomainIDs": authenticationDomainIDs,
	}

	if len(groupIDs) != 0 {
		vars["groupIDs"] = groupIDs
	}

	if name != "" {
		vars["name"] = name
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getGroupsQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.Actor.Organization.UserManagement.AuthenticationDomains, nil
}

const getGroupsQuery = `query(
	$authenticationDomainIDs: [ID!],
	$groupIDs: [ID!],
	$name: String,
) {
  actor {
    organization {
      userManagement {
        authenticationDomains(
          id: $authenticationDomainIDs
        ) {
          authenticationDomains {
            groups(filter: {displayName: {contains: $name}, id: {in: $groupIDs}}) {
              groups {
                id
                displayName
              }
            }
            id
            name
            provisioningType
          }
        }
      }
    }
  }
}`

// GET Groups and Users belonging to them

func (a *Usermanagement) GetGroupsWithUsers(
	authenticationDomainIDs []string,
	groupIDs []string,
) (*UserManagementAuthenticationDomains, error) {
	return a.GetGroupsWithUsersWithContext(context.Background(),
		authenticationDomainIDs,
		groupIDs,
	)
}

func (a *Usermanagement) GetGroupsWithUsersWithContext(
	ctx context.Context,
	authenticationDomainIDs []string,
	groupIDs []string,
) (*UserManagementAuthenticationDomains, error) {

	resp := authenticationDomainsResponse{}
	vars := map[string]interface{}{
		"authenticationDomainIDs": authenticationDomainIDs,
		"groupIDs":                groupIDs,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getGroupsWithUsersQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.Actor.Organization.UserManagement.AuthenticationDomains, nil
}

const getGroupsWithUsersQuery = `query(
	$authenticationDomainIDs: [ID!],
	$groupIDs: [ID!],
) {
  actor {
    organization {
      userManagement {
        authenticationDomains(id: $authenticationDomainIDs) {
          authenticationDomains {
            groups(filter: {id: {in: $groupIDs}}) {
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
          }
        }
      }
    }
  }
}`

// GET User(s)

func (a *Usermanagement) GetUsers(
	authenticationDomainIDs []string,
	userIDs []string,
	nameContains string,
	emailContains string,
) (*UserManagementAuthenticationDomains, error) {
	return a.GetUsersWithContext(context.Background(),
		authenticationDomainIDs,
		userIDs,
		nameContains,
		emailContains,
	)
}

func (a *Usermanagement) GetUsersWithContext(
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

// container for authentication_domains enabling cursor based pagination
// ------ TO BE REMOVED -------
// commented out for now

//func (a *Usermanagement) GetAuthenticationDomains(
//	authenticationDomainsID []string,
//) (*[]UserManagementAuthenticationDomain, error) {
//	return a.GetAuthenticationDomainsWithContext(context.Background(),
//		authenticationDomainsID,
//	)
//}

// container for authentication_domains enabling cursor based pagination
func (a *Usermanagement) GetAllAuthenticationDomains() (*[]UserManagementAuthenticationDomain, error) {
	return a.GetAllAuthenticationDomainsWithContext(context.Background())
}

// container for authentication_domains enabling cursor based pagination
// ------ TO BE REMOVED -------
// commented out for now

//func (a *Usermanagement) GetAuthenticationDomainsWithContext(
//	ctx context.Context,
//	authenticationDomainsID []string,
//) (*[]UserManagementAuthenticationDomain, error) {
//
//	resp := authenticationDomainsResponse{}
//	vars := map[string]interface{}{
//		"authenticationDomainsID": authenticationDomainsID,
//	}
//
//	if err := a.client.NerdGraphQueryWithContext(ctx, getAuthenticationDomainsQuery, vars, &resp); err != nil {
//		return nil, err
//	}
//
//	if len(resp.Actor.Organization.UserManagement.AuthenticationDomains.AuthenticationDomains) == 0 {
//		return nil, errors.NewNotFound("")
//	}
//
//	return &resp.Actor.Organization.UserManagement.AuthenticationDomains.AuthenticationDomains, nil
//}

// GetAllAuthenticationDomainsWithContext is a modified function that uses a modified query to fetch all authentication domains (not query by ID)
func (a *Usermanagement) GetAllAuthenticationDomainsWithContext(
	ctx context.Context,
) (*[]UserManagementAuthenticationDomain, error) {

	resp := authenticationDomainsResponse{}

	vars := map[string]interface{}{}

	if err := a.client.NerdGraphQueryWithContext(ctx, getAllAuthenticationDomainsQuery, vars, &resp); err != nil {
		return nil, err
	}

	if len(resp.Actor.Organization.UserManagement.AuthenticationDomains.AuthenticationDomains) == 0 {
		return nil, errors.NewNotFound("")
	}

	return &resp.Actor.Organization.UserManagement.AuthenticationDomains.AuthenticationDomains, nil
}

//const getAuthenticationDomainsQuery = `query(
//	$authenticationDomainsID: ID!,
//) { actor { organization { userManagement { authenticationDomains(cursor: $authenticationDomainsCursor) { authenticationDomains(id: $authenticationDomainsID) { authenticationDomains {
//	id
//	name
//	provisioningType
//} } } } } } }`

// The following query is out of the scope of Tutone. DO NOT DELETE THIS.
// To be split into two queries based on user/group management usage (e.g. users need not be fetched if groups are needed)

// ------ TO BE REMOVED -------
// commented out for now

//const getAuthenticationDomainsQuery = `query ($authenticationDomainsID: [ID!]) {
//  actor {
//    organization {
//      userManagement {
//        authenticationDomains(id: $authenticationDomainsID) {
//          authenticationDomains {
//            groups {
//              groups {
//                displayName
//                id
//                users {
//                  users {
//                    email
//                    id
//                    name
//                    timeZone
//                  }
//                  nextCursor
//                  totalCount
//                }
//              }
//              nextCursor
//              totalCount
//            }
//            id
//            name
//            provisioningType
//            users {
//              users {
//                email
//                emailVerificationState
//                groups {
//                  groups {
//                    displayName
//                    id
//                  }
//                  nextCursor
//                  totalCount
//                }
//                id
//                lastActive
//                name
//                pendingUpgradeRequest {
//                  id
//                  message
//                  requestedUserType {
//                    displayName
//                    id
//                  }
//                }
//                timeZone
//                type {
//                  displayName
//                  id
//                }
//              }
//              nextCursor
//              totalCount
//            }
//          }
//        }
//      }
//    }
//  }
//}`

// The following query is out of the scope of Tutone. DO NOT DELETE THIS.
// To be split into two queries based on user/group management usage (e.g. users need not be fetched if groups are needed)
const getAllAuthenticationDomainsQuery = `query {
  actor {
    organization {
      userManagement {
        authenticationDomains {
          authenticationDomains {
            id
            name
            provisioningType
          }
        }
      }
    }
  }
}`

// GET Users in Groups

func (a *Usermanagement) GetUsersInGroups(
	authenticationDomainIDs []string,
	groupIDs []string,
) (*UserManagementAuthenticationDomains, error) {
	return a.GetUsersInGroupsWithContext(context.Background(),
		authenticationDomainIDs,
		groupIDs,
	)
}

func (a *Usermanagement) GetUsersInGroupsWithContext(
	ctx context.Context,
	authenticationDomainIDs []string,
	groupIDs []string,
) (*UserManagementAuthenticationDomains, error) {

	resp := authenticationDomainsResponse{}
	vars := map[string]interface{}{
		"authenticationDomainIDs": authenticationDomainIDs,
		"groupIDs":                groupIDs,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getUsersInGroupsQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.Actor.Organization.UserManagement.AuthenticationDomains, nil
}

const getUsersInGroupsQuery = `query ($authenticationDomainIDs: [ID!], $groupIDs: [ID!]) {
  actor {
    organization {
      userManagement {
        authenticationDomains(id: $authenticationDomainIDs) {
          authenticationDomains {
            groups(filter: {id: {in: $groupIDs}}) {
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
        }
      }
    }
  }
}`
