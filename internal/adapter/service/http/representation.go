package http

import "fmt"

/*
This file is an attempt to represent resource using JSON Hypertext Application Language
*/

type resourceLink struct {
	Href string `json:"href"`
} // @name Link

type redirectionLinks struct {
	Self    resourceLink `json:"self"`
	Count   resourceLink `json:"count"`
	Consume resourceLink `json:"consume"`
} // @name RedirectionLinks

type redirectionRepresentation struct {
	Key   string           `json:"key"`
	Links redirectionLinks `json:"_links"`
} // @name RedirectionRepresentation

func getRedirectionRepresentation(key string) redirectionRepresentation {
	return redirectionRepresentation{
		Key: key,
		Links: redirectionLinks{
			Self: resourceLink{
				Href: fmt.Sprintf("/api/redirection/%s", key),
			},
			Count: resourceLink{
				Href: fmt.Sprintf("/api/redirection/%s/count", key),
			},
			Consume: resourceLink{
				Href: fmt.Sprintf("/%s", key),
			},
		},
	}
}

type redirectionListLinks struct {
	Self resourceLink `json:"self"`
} // @name RedirectionListLinks

type redirectionListRepresentation struct {
	Items []redirectionRepresentation `json:"items"`
	Links redirectionListLinks        `json:"_links"`
} // @name RedirectionListRepresentation

func getRedirectionListRepresentation(keys []string) redirectionListRepresentation {
	items := []redirectionRepresentation{}
	for _, key := range keys {
		items = append(items, getRedirectionRepresentation(key))
	}

	return redirectionListRepresentation{
		Items: items,
		Links: redirectionListLinks{
			Self: resourceLink{
				Href: "/api/redirections",
			},
		},
	}
}

type redirectionCountLinks struct {
	Self   resourceLink `json:"self"`
	Parent resourceLink `json:"parent"`
} // @name RedirectionCountLinks

type redirectionCountRepresentation struct {
	Count int                   `json:"count"`
	Links redirectionCountLinks `json:"_links"`
} // @name RedirectionCountRepresentation

func getRedirectionCountRepresentation(key string, count int) redirectionCountRepresentation {
	return redirectionCountRepresentation{
		Count: count,
		Links: redirectionCountLinks{
			Self: resourceLink{
				Href: fmt.Sprintf("/api/redirection/%s/count", key),
			},
			Parent: resourceLink{
				Href: fmt.Sprintf("/api/redirection/%s", key),
			},
		},
	}
}
