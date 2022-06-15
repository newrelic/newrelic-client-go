package notifications

const (
	graphqlDestinationFields = `
						      id
							  type
							  accountId
							  createdAt
							  updatedAt
							  updatedBy
							  name
							  active
							  status
							  lastSent
							  properties {
								label
								key
								value
								displayValue
							  }
							  auth {
								  ... on AiNotificationsTokenAuth {
									authType
									prefix
								  }
								  ... on AiNotificationsBasicAuth {
									authType
									user
								  }
							  }
	`

	notificationsCreateDestination = `mutation CreateDestination($accountID: Int!, $destinationInput: DestinationInput!) {
		aiNotificationsCreateDestination(accountId: $accountID, destination: $destinationInput) {
			destination {
				` + graphqlDestinationFields + `
			}
		} 
	}`

	notificationsDeleteDestination = `mutation DeleteDestinationMutation($accountID: Int!, $destinationID: ID!) {
		aiNotificationsDeleteDestination(accountId: $accountID, destinationId: $destinationID) {
			ids
		}
	}`

	notificationsGetDestinations = `query ($accountID: Int!) {
			actor {
			  account(id: $accountID) {
				aiNotifications {
				  destinations {
					entities {
						` + graphqlDestinationFields + `
					}
				  }
				}
			  }
			}
		}`

	notificationsGetDestination = `query ($accountID: Int!, $filters: AiNotificationsDestinationFilter!) {
			actor {
			  account(id: $accountID) {
				aiNotifications {
				  destinations(filters: $filters) {
					entities {
						` + graphqlDestinationFields + `
					}
				  }
				}
			  }
			}
		}`
)
