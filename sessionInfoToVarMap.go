package webwire

// SessionInfoToVarMap is a utility function that turns a
// session info compliant object into a map of variants.
// This is helpful for serialization of session info objects.
func SessionInfoToVarMap(info SessionInfo) map[string]interface{} {
	if info == nil {
		return nil
	}
	varMap := make(map[string]interface{})
	for _, field := range info.Fields() {
		varMap[field] = info.Value(field)
	}
	return varMap
}
