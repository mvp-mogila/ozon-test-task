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

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, createPostInput dto.CreatePostInput) (*dto.PostData, error) {
	newPost, err := r.PostUsecase.CreatePost(ctx, dto.ConvertCreatePostInputDTOToModel(createPostInput))
	if err != nil {
		return nil, gqlerror.Wrap(err)
	}

	postDTO := dto.ConvertPostModelDataToDTO(newPost)
	return &postDTO, nil
}

// Comments is the resolver for the comments field.
func (r *postResolver) Comments(ctx context.Context, obj *dto.Post, page *int, commentsCount *int) ([]*dto.Comment, error) {
	if page == nil {
		page = new(int)
	}
	if commentsCount == nil {
		commentsCount = new(int)
	}

	comments, err := r.CommentUsecase.GetCommentsForPost(ctx, obj.PostData.ID, *page, *commentsCount)
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

// GetPost is the resolver for the getPost field.
func (r *queryResolver) GetPost(ctx context.Context, id int) (*dto.Post, error) {
	post, err := r.PostUsecase.GetPost(ctx, id)
	if err != nil {
		return nil, gqlerror.Wrap(err)
	}

	postDataDTO := dto.ConvertPostModelDataToDTO(post)
	postDTO := &dto.Post{
		PostData: &postDataDTO,
	}
	return postDTO, nil
}

// GetPosts is the resolver for the getPosts field.
func (r *queryResolver) GetPosts(ctx context.Context, page *int, postsCount *int) ([]*dto.PostData, error) {
	posts, err := r.PostUsecase.GetPosts(ctx, *page, *postsCount)
	if err != nil {
		return nil, gqlerror.Wrap(err)
	}

	postDTOs := make([]*dto.PostData, 0, len(posts))
	for _, p := range posts {
		postDTO := dto.ConvertPostModelDataToDTO(p)
		postDTOs = append(postDTOs, &postDTO)
	}
	return postDTOs, nil
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Post returns graph.PostResolver implementation.
func (r *Resolver) Post() graph.PostResolver { return &postResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type postResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
