# notifications
--
    import "github.com/newrelic/newrelic-client-go/v3/pkg/notifications"

## Usage

#### type Notifications

```go
type Notifications struct {
}
```

Notifications is used to communicate with New Relic Notifications.

#### func  New

```go
func New(config config.Config) Notifications
```
New is used to create a new Notifications' client instance.

### Destinations

#### func  TestNotificationMutationDestination

```go
func TestNotificationMutationDestination(t *testing.T)
```

#### type AiNotificationsDestination

```go
type AiNotificationsDestination struct {
    AccountID           int                                 `json:"accountId"`
    Active              bool                                `json:"active"`
    Auth                ai.AiNotificationsAuth              `json:"auth,omitempty"`
    CreatedAt           nrtime.DateTime                     `json:"createdAt"`
    ID                  string                              `json:"id"`
    IsUserAuthenticated bool                                `json:"isUserAuthenticated"`
    LastSent            nrtime.DateTime                     `json:"lastSent,omitempty"`
    Name                string                              `json:"name"`
    Properties          []AiNotificationsProperty           `json:"properties"`
    Status              AiNotificationsDestinationStatus    `json:"status"`
    Type                AiNotificationsDestinationType      `json:"type"`
    UpdatedAt           nrtime.DateTime                     `json:"updatedAt"`
    UpdatedBy           int                                 `json:"updatedBy"`
}
```

AiNotificationsDestination represents a New Relic notification destination.

#### func (*Notifications) AiNotificationsCreateDestination

```go
func (a *Notifications) AiNotificationsCreateDestination(accountID int,destination AiNotificationsDestinationInput) (*AiNotificationsDestinationResponse, error)
```
AiNotificationsCreateDestination creates a new notification destination for a given account.

#### func (*Notifications) GetDestinations

```go
func (a *Notifications) GetDestinations(accountID int,cursor string, filters ai.AiNotificationsDestinationFilter, sorter AiNotificationsDestinationSorter) (*AiNotificationsDestinationsResponse, error)
```
GetDestinations returns a list of notifications destinations for a given account. You can filter by ID.

#### func (*Notifications) AiNotificationsUpdateDestination

```go
func (a *Notifications) AiNotificationsUpdateDestination(accountID int,destination AiNotificationsDestinationUpdate, destinationId string) (*AiNotificationsDestinationResponse, error)
```
AiNotificationsUpdateDestination update a notification destination for a given account.

#### type AiNotificationsDeleteDestination

```go
func (a *Notifications) AiNotificationsDeleteDestination(accountID int, destinationId string) (*AiNotificationsDeleteResponse, error)
```

AiNotificationsDeleteDestination delete a notification destination for a given account.

### Channels
#### func  TestNotificationMutationChannel

```go
func TestNotificationMutationChannel(t *testing.T)
```

#### type AiNotificationsChannel

```go
type AiNotificationsChannel struct {
    AccountID       int                             `json:"accountId"`
    Active          bool                            `json:"active"`
    CreatedAt       nrtime.DateTime                 `json:"createdAt"`
    DestinationId   string                          `json:"destinationId"`
    ID              string                          `json:"id"`
    Name            string                          `json:"name"`
    Product         AiNotificationsProduct          `json:"product"`
    Properties      []AiNotificationsProperty       `json:"properties"`
    Status          AiNotificationsChannelStatus    `json:"status"`
    Type            AiNotificationsChannelType      `json:"type"`
    UpdatedAt       nrtime.DateTime                 `json:"updatedAt"`
    UpdatedBy       int                             `json:"updatedBy"`
}
```

AiNotificationsChannel represents a New Relic notification channel.

#### func (*Notifications) AiNotificationsCreateChannel

```go
func (a *Notifications) AiNotificationsCreateChannel(accountID int, destination AiNotificationsChannelInput) (*AiNotificationsChannelResponse, error)
```
AiNotificationsCreateChannel creates a new notification channel for a given account.

#### func (*Notifications) GetChannels

```go
func (a *Notifications) GetChannels(accountID int,cursor string, filters ai.AiNotificationsChannelFilter, sorter AiNotificationsChannelSorter) (*AiNotificationsChannelsResponse, error)
```
GetChannels returns a list of notifications channels for a given account. You can filter by ID.

#### func (*Notifications) AiNotificationsUpdateChannel

```go
func (a *Notifications) AiNotificationsUpdateChannel(accountID int, destination AiNotificationsChannelUpdate, channelId string) (*AiNotificationsChannelResponse, error)
```
AiNotificationsUpdateChannel update a notification channel for a given account.

#### type AiNotificationsDeleteChannel

```go
func (a *Notifications) AiNotificationsDeleteChannel(accountID int, channelId string) (*AiNotificationsDeleteResponse, error)
```

AiNotificationsDeleteChannel delete a notification channel for a given account.
