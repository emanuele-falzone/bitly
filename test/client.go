package test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/slices"
)

type Client struct {
	driver   Driver
	ctx      context.Context
	key      string
	location string
	count    int
	err      error
	keys     []string
	gotKeys  []string
}

const (
	delay = 50 * time.Millisecond
)

func NewClient(ctx context.Context, driver Driver) *Client {
	return &Client{driver: driver, ctx: ctx}
}

func (c *Client) CreateRedirection(location string) error {
	c.key, c.err = c.driver.CreateRedirection(c.ctx, location)
	c.keys = append(c.keys, c.key)

	time.Sleep(delay)

	return nil
}

func (c *Client) DeleteRedirection() error {
	c.err = c.driver.DeleteRedirection(c.ctx, c.key)

	time.Sleep(delay)

	return nil
}

func (c *Client) GetNonExistingKey() error {
	c.key = "not_found"

	time.Sleep(delay)

	return nil
}

func (c *Client) GetRedirectionLocation() error {
	c.location, c.err = c.driver.GetRedirectionLocation(c.ctx, c.key)

	time.Sleep(delay)

	return nil
}

func (c *Client) GetRedirectionList() error {
	c.gotKeys, c.err = c.driver.GetRedirectionList(c.ctx)

	time.Sleep(delay)

	return nil
}

func (c *Client) GetRedirectionLocationTimes(times int) error {
	for i := 1; i <= times; i++ {
		c.location, c.err = c.driver.GetRedirectionLocation(c.ctx, c.key)

		time.Sleep(delay)
	}

	return nil
}

func (c *Client) GetRedirectionCount() error {
	c.count, c.err = c.driver.GetRedirectionCount(c.ctx, c.key)

	time.Sleep(delay)

	return nil
}

func (c *Client) ConfirmLocationToBe(location string) error {
	if location != c.location {
		return fmt.Errorf("the locations are not equal, expected %s, got %s", location, c.location)
	}

	return nil
}

func (c *Client) ConfirmCorrectList() error {
	for _, mine := range c.keys {
		if !slices.Contains(c.gotKeys, mine) {
			return fmt.Errorf("the key %s was not found in %s", mine, c.gotKeys)
		}
	}

	return nil
}

func (c *Client) ConfirmCountToBe(count int) error {
	if count != c.count {
		return fmt.Errorf("the counts are not equal, expected %d, got %d", count, c.count)
	}

	return nil
}

func (c *Client) ConfirmHasKey() error {
	if c.key == "" {
		return fmt.Errorf("the client does not have any key")
	}

	return nil
}

func (c *Client) ConfirmError() error {
	if c.err == nil {
		return errors.New("expecting an error")
	}

	return nil
}

func (c *Client) ConfirmNoError() error {
	if c.err != nil {
		return errors.New("not expecting an error")
	}

	return nil
}
