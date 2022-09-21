package opensearch

// EnsureIndexTemplate makes sure that an Index Template is present and is configured in a given way.
//
// Note: If the configuration of an exisiting index template doesn't match the given configuration an error
// will be returned. Currently there is no save way for us to update the index template.
func EnsureIndexTemplate() {}
