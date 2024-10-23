package usecase

type PostsRepo interface {
	CreatePost()
	GetPost()
	GetPosts()
}

type PostUsecase struct {
	repo PostsRepo
}

func NewPostUsecase(r PostsRepo) *PostUsecase {
	return &PostUsecase{
		repo: r,
	}
}

func (u *PostUsecase) CreatePost() {

}

func (u *PostUsecase) GetPost() {

}

func (u *PostUsecase) GetPosts() {

}
