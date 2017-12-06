package model

// DeleteToken is an individual delete token for a player
type DeleteToken struct {
	Token
}

const deleteTokenEntityType = "DeleteToken"

// EntityType returns the entity type
func (m *DeleteToken) EntityType() string {
	return deleteTokenEntityType
}
