package profile

type Store interface {
	Create(Profile) error
	Get(name string) (Profile, error)
	List() ([]Profile, error)
	Delete(name string) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Create(p Profile) error {
	if err := p.Validate(); err != nil {
		return err
	}
	return s.store.Create(p)
}

func (s *Service) Get(name string) (Profile, error) {
	if err := ValidateName(name); err != nil {
		return Profile{}, err
	}
	return s.store.Get(name)
}

func (s *Service) List() ([]Profile, error) {
	return s.store.List()
}

func (s *Service) Delete(name string) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	return s.store.Delete(name)
}

