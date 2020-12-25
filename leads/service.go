package leads

// Service interface for leads service
type Service interface {
	Store(email string) error
}
