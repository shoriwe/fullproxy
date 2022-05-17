package listeners

type NoFilter struct {
}

func (nf *NoFilter) Outbound(_ string) error {
	return nil
}

func (nf *NoFilter) Inbound(_ string) error {
	return nil
}

func (nf *NoFilter) Listen(_ string) error {
	return nil
}

func (nf *NoFilter) Accept(_ string) error {
	return nil
}
