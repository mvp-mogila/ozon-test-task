package usecase

type CommentsRepo interface {
	CreateComment()
	GetComments()
}

type CommentUsecase struct {
	repo CommentsRepo
}

func NewCommentUsecase(r CommentsRepo) *CommentUsecase {
	return &CommentUsecase{
		repo: r,
	}
}

func (u *CommentUsecase) CreateComment() {

}

func (u *CommentUsecase) GetComments() {

}
