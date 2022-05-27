package http

import "fmt"

/*
This file is an attempt to represent resource using JSON Hypertext Application Language
*/

type ResourceLink struct {
	Href string `json:"href"`
} // @name Link

type RedirectionLinks struct {
	Self    ResourceLink `json:"self"`
	Count   ResourceLink `json:"count"`
	Consume ResourceLink `json:"consume"`
} // @name RedirectionLinks

type RedirectionRepresentation struct {
	Key   string           `json:"key"`
	Links RedirectionLinks `json:"_links"`
} // @name RedirectionRepresentation

func getRedirectionRepresentation(key string) RedirectionRepresentation {
	return RedirectionRepresentation{
		Key: key,
		Links: RedirectionLinks{
			Self: ResourceLink{
				Href: fmt.Sprintf("/api/redirection/%s", key),
			},
			Count: ResourceLink{
				Href: fmt.Sprintf("/api/redirection/%s/count", key),
			},
			Consume: ResourceLink{
				Href: fmt.Sprintf("/%s", key),
			},
		},
	}
}

type RedirectionListLinks struct {
	Self ResourceLink `json:"self"`
} // @name RedirectionListLinks

type RedirectionListRepresentation struct {
	Items []RedirectionRepresentation `json:"items"`
	Links RedirectionListLinks        `json:"_links"`
} // @name RedirectionListRepresentation

func getRedirectionListRepresentation(keys []string) RedirectionListRepresentation {
	items := []RedirectionRepresentation{}
	for _, key := range keys {
		items = append(items, getRedirectionRepresentation(key))
	}

	return RedirectionListRepresentation{
		Items: items,
		Links: RedirectionListLinks{
			Self: ResourceLink{
				Href: "/api/redirections",
			},
		},
	}
}

type RedirectionCountLinks struct {
	Self   ResourceLink `json:"self"`
	Parent ResourceLink `json:"parent"`
} // @name RedirectionCountLinks

type RedirectionCountRepresentation struct {
	Count int                   `json:"count"`
	Links RedirectionCountLinks `json:"_links"`
} // @name RedirectionCountRepresentation

func getRedirectionCountRepresentation(key string, count int) RedirectionCountRepresentation {
	return RedirectionCountRepresentation{
		Count: count,
		Links: RedirectionCountLinks{
			Self: ResourceLink{
				Href: fmt.Sprintf("/api/redirection/%s/count", key),
			},
			Parent: ResourceLink{
				Href: fmt.Sprintf("/api/redirection/%s", key),
			},
		},
	}
}
