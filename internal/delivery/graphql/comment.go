package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.55

import (
	"context"
	"fmt"

	"github.com/mvp-mogila/ozon-test-task/gen/graph"
	"github.com/mvp-mogila/ozon-test-task/internal/models"
)

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, createCommentInput models.CreateCommentInput) (*models.Comment, error) {
	panic(fmt.Errorf("not implemented: CreateComment - createComment"))
}

// CommentSubscription is the resolver for the commentSubscription field.
func (r *subscriptionResolver) CommentSubscription(ctx context.Context, postID string) (<-chan *models.Comment, error) {
	panic(fmt.Errorf("not implemented: CommentSubscription - commentSubscription"))
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Subscription returns graph.SubscriptionResolver implementation.
func (r *Resolver) Subscription() graph.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }