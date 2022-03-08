package repository

import (
	"context"
	"fmt"
	"gqlgen-ent/ent"
	"gqlgen-ent/ent/car"
	"gqlgen-ent/ent/group"
	"gqlgen-ent/ent/user"
	"log"
	"time"
)

func CreateGraph(ctx context.Context, client *ent.Client) error {
	a8m, err := client.User.Create().SetAge(30).SetName("Ariel").Save(ctx)
	if err != nil {
		return err
	}
	neta, err := client.User.Create().SetAge(30).SetName("Neta").Save(ctx)
	if err != nil {
		return err
	}
	err = client.Car.Create().SetModel("Tesla").SetRegisteredAt(time.Now()).SetOwner(a8m).Exec(ctx)
	if err != nil {
		return err
	}
	err = client.Car.Create().SetModel("Mazta").SetRegisteredAt(time.Now()).SetOwner(a8m).Exec(ctx)
	if err != nil {
		return err
	}
	err = client.Car.Create().SetModel("Ford").SetRegisteredAt(time.Now()).SetOwner(neta).Exec(ctx)
	if err != nil {
		return err
	}
	err = client.Group.Create().SetName("GitHub").AddUsers(a8m).Exec(ctx)
	if err != nil {
		return err
	}
	err = client.Group.Create().SetName("GitLab").AddUsers(a8m, neta).Exec(ctx)
	if err != nil {
		return err
	}
	log.Println("The graph was created successfully.")
	return nil
}

func QueryGitHub(ctx context.Context, client *ent.Client) error {
	cars, err := client.Group.Query().Where(group.Name("GitHub")).QueryUsers().QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed getting cars: %w", err)
	}
	log.Println("car returned:", cars)
	return nil
}

func QueryArielCars(ctx context.Context, client *ent.Client) error {
	a8m := client.User.
		Query().
		Where(
			user.HasCars(),
			user.Name("Ariel"),
		).
		OnlyX(ctx)
	cars, err := a8m.
		QueryGroups().
		QueryUsers().
		QueryCars().
		Where(
			car.Not(
				car.Model("Mazda"),
			),
		).
		All(ctx)
	if err != nil {
		return fmt.Errorf("failed getting cars: %w", err)
	}
	log.Println("cars returned:", cars)
	return nil
}

func QueryGroupWithUsers(ctx context.Context, client *ent.Client) error {
	groups, err := client.Group.Query().Where(group.HasUsers()).All(ctx)
	if err != nil {
		return fmt.Errorf("failed getting groups: %w", err)
	}
	log.Println("groups returned:", groups)
	return nil
}
