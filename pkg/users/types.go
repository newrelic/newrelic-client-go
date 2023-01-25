// Code generated by tutone: DO NOT EDIT
package users

// Actor - The `Actor` object contains fields that are scoped to the API user's access level.
type Actor struct {
	// The authenticated `User` who made this request.
	User User `json:"user,omitempty"`
}

// User - The `User` object provides general data about the user.
type User struct {
	//
	Email string `json:"email,omitempty"`
	//
	ID int `json:"id,omitempty"`
	//
	Name string `json:"name,omitempty"`
}

// UserManagementUser - A user of New Relic scoped to an authentication domain.
type UserManagementUser struct {
	// Email address of the user.
	Email string `json:"email,omitempty"`
	// One of: "Not Verifiable", "Verified", and "Pending".
	//
	// Not Verifiable: the user's email does not require verification.
	//
	// Verified: the user's email requires verification and has been.
	//
	// Pending: the user's email requires verification and has not been.
	EmailVerificationState string `json:"emailVerificationState"`
	// container for groups enabling cursor based pagination
	Groups UserManagementUserGroups `json:"groups,omitempty"`
	// a value that uniquely identifies this object
	ID int `json:"id"`
	// The last active date of the user.
	LastActive DateTime `json:"lastActive,omitempty"`
	// The full name of the user.
	Name string `json:"name,omitempty"`
	// Time zone of the user in IANA Time Zone database format, also known as the "Olson" time zone database format (for exmaple, "America/Los_Angeles").
	TimeZone string `json:"timeZone,omitempty"`
	// A "user type" is what determines the set of New Relic capabilities a user can theoretically access.
	Type UserManagementUserType `json:"type"`
}

// UserManagementUserGroup - For users on our New Relic One user model, a "group" represents a group of users. Putting users in a group allows the managing of permissions for multiple users at the same time.
type UserManagementUserGroup struct {
	// the name of the object
	DisplayName string `json:"displayName"`
	// a value that uniquely identifies this object
	ID int `json:"id"`
}

// UserManagementUserGroups - container for groups enabling cursor based pagination
type UserManagementUserGroups struct {
	// container for groups enabling cursor based pagination
	Groups []UserManagementUserGroup `json:"groups"`
	// an opaque cursor to supply with subsequent   requests to get the next page of results, null if there are no more pages
	NextCursor string `json:"nextCursor,omitempty"`
	// the total number of results
	TotalCount int `json:"totalCount"`
}

// UserManagementUserType - A "user type" is what determines the set of New Relic capabilities a user can theoretically access.
type UserManagementUserType struct {
	// the name of the object
	DisplayName string `json:"displayName"`
	// a value that uniquely identifies this object
	ID int `json:"id"`
}

// UserReference - The `UserReference` object provides basic identifying information about the user.
type UserReference struct {
	//
	Email string `json:"email,omitempty"`
	//
	Gravatar string `json:"gravatar,omitempty"`
	//
	ID int `json:"id,omitempty"`
	//
	Name string `json:"name,omitempty"`
}

type userResponse struct {
	Actor Actor `json:"actor"`
}

// DateTime - The `DateTime` scalar represents a date and time. The `DateTime` appears as an ISO8601 formatted string.
type DateTime string

// ID - The `ID` scalar type represents a unique identifier, often used to
// refetch an object or as key for a cache. The ID type appears in a JSON
// response as a String; however, it is not intended to be human-readable.
// When expected as an input type, any string (such as `"4"`) or integer
// (such as `4`) input value will be accepted as an ID.
type ID string
