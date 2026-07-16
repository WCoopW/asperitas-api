package user

import domain "reddit/internal/domain/user"

type UserMapper struct{}

func (m *UserMapper) SchemaToEntity(schema *UserSchema) domain.User {
	return domain.User{
		ID:           schema.ID,
		Username:     schema.Username,
		PasswordHash: schema.PasswordHash,
	}
}

func (m *UserMapper) EntityToSchema(entity domain.User) UserSchema {
	return UserSchema{
		ID:           entity.ID,
		Username:     entity.Username,
		PasswordHash: entity.PasswordHash,
	}
}
