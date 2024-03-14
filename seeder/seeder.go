package seeder

import (
	"context"
	"log"

	"github.com/mikestefanello/pagoda/config"
	"github.com/mikestefanello/pagoda/ent"
	"golang.org/x/crypto/bcrypt"
)

// RunIdempotentSeeder seeds all regular objects that the app needs to run smoothly
func RunIdempotentSeeder(cfg *config.Config, client *ent.Client) error {
	ctx := context.Background()

	SeedUser(client, ctx, "Alice Bonjovi", "alice@test.com", "password")
	SeedUser(client, ctx, "Bob Lacoste", "bob@test.com", "password")
	return nil
}

func SeedUser(client *ent.Client, ctx context.Context, name, email, password string) error {
	hashPassword := func(password string) string {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			panic("failed to hash password")
		}
		return string(hash)
	}

	_, err := client.User.
		Create().
		SetName(name).
		SetEmail(email).
		SetPassword(hashPassword(password)).
		SetVerified(true).
		Save(ctx)
	if err != nil {
		log.Printf("User with email %s already exists. Skipping.", email)
		return nil
	}
	return nil
}
