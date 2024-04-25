package github

type PushEvent struct {
	Repository string
	Branch     string
	AuthorName string
}
