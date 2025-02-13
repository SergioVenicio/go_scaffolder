package repositories

type Repository[T any] interface {
	Insert(T) error
}
