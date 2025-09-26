package earlydecoder

import "fmt"

type ephemeralResource struct {
	resource
}

func (r *ephemeralResource) MapKey() string {
	return fmt.Sprintf("ephemeral.%s.%s", r.Type, r.Name)
}
