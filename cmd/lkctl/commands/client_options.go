package commands

type clientOptions struct {
	server       string
	authURL      string
	clientID     string
	clientSecret string
	scope        []string
	boostrap     bool
}
