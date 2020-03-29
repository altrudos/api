package main

func (c *RouteContext) Authenticate() error {
	// TODO: authentication
	c.UserId = "12345"
	return nil
}