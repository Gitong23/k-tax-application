package tax

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) CalTax() {
	// Implement the logic here
}
