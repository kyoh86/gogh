package repo

// Branch will get current branch
func (r *Repository) Branch() (string, error) {
	head, err := r.repository.Head()
	if err != nil {
		return "", err
	}
	return head.Name().Short(), nil
}
