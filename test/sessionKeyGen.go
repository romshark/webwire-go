package test

// sessionKeyGen implements the webwire.SessionKeyGenerator interface
type sessionKeyGen struct {
	generate func() string
}

// Generate implements the webwire.SessionKeyGenerator interface
func (gen *sessionKeyGen) Generate() string {
	return gen.generate()
}
