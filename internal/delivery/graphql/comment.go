package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.55

import (
	"context"

	"github.com/mvp-mogila/ozon-test-task/gen/graph"
	"github.com/mvp-mogila/ozon-test-task/internal/delivery/graphql/dto"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Comments is the resolver for the comments field.
func (r *commentResolver) Comments(ctx context.Context, obj *dto.Comment, commentsCount *int) ([]*dto.Comment, error) {
	if commentsCount == nil {
		commentsCount = new(int)
	}

	comments, err := r.CommentUsecase.GetCommentsForComment(ctx, obj.ID, *commentsCount)
	if err != nil {
		return nil, gqlerror.Wrap(err)
	}

	commentDTOs := make([]*dto.Comment, 0, len(comments))
	for _, c := range comments {
		commentDTO := dto.ConvertCommentModelToDTO(c)
		commentDTOs = append(commentDTOs, &commentDTO)
	}
	return commentDTOs, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, createCommentInput dto.CreateCommentInput) (*dto.Comment, error) {
	newComment, err := r.CommentUsecase.CreateComment(ctx, dto.ConvertCreateCommentInputDTOToModel(createCommentInput))
	if err != nil {
		return nil, gqlerror.Wrap(err)
	}

	commentDTO := dto.ConvertCommentModelToDTO(newComment)
	return &commentDTO, nil
}

// CommentSubscription is the resolver for the commentSubscription field.
func (r *subscriptionResolver) CommentSubscription(ctx context.Context, postID int) (<-chan *dto.Comment, error) {
	commentDTOCh, err := r.CommentUsecase.CreateSubscriber(ctx, postID)
	if err != nil {
		return nil, gqlerror.Wrap(err)
	}

	return commentDTOCh, nil
}

// Comment returns graph.CommentResolver implementation.
func (r *Resolver) Comment() graph.CommentResolver { return &commentResolver{r} }

// Subscription returns graph.SubscriptionResolver implementation.
func (r *Resolver) Subscription() graph.SubscriptionResolver { return &subscriptionResolver{r} }

type commentResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
