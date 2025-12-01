package weather

// NewTestClient creates a client with a custom base URL for testing
func NewTestClient(apiKey, baseURL string) *Client {
	c := NewClient(apiKey)
	c.baseURL = baseURL
	return c
}

// BuildURL exposes buildURL for testing
func (c *Client) BuildURL(opts FetchOptions) string {
	return c.buildURL(opts)
}
