package test

// SessionKeyGen implements the webwire.SessionKeyGenerator interface
type SessionKeyGen struct {
	OnGenerate func() string
}

// Generate implements the webwire.SessionKeyGenerator interface
func (gen *SessionKeyGen) Generate() string {
	if gen.OnGenerate != nil {
		return gen.OnGenerate()
	}
	return ""
}
