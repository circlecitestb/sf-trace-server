// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package metadata

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/v7/libbeat/common"
)

func TestUserFields(t *testing.T) {
	id := "1234"
	ip := net.ParseIP("127.0.0.1")
	email := "test@mail.co"
	name := "user123"
	userAgent := "rum-1.0"

	tests := []struct {
		User   User
		Output common.MapStr
	}{
		{
			User:   User{},
			Output: nil,
		},
		{
			User: User{
				ID:        id,
				IP:        ip,
				Email:     email,
				Name:      name,
				UserAgent: userAgent,
			},
			Output: common.MapStr{
				"id":    "1234",
				"email": "test@mail.co",
				"name":  "user123",
			},
		},
	}

	for _, test := range tests {
		output := test.User.Fields()
		assert.Equal(t, test.Output, output)
	}
}

func TestUserClientFields(t *testing.T) {
	id := "1234"
	email := "test@mail.co"
	userName := "user123"
	userAgent := "rum-1.0"

	for name, tc := range map[string]struct {
		ip  string
		out common.MapStr
	}{
		"Empty":   {ip: "", out: nil},
		"IPv4":    {ip: "192.0.0.1", out: common.MapStr{"ip": "192.0.0.1"}},
		"IPv6":    {ip: "2001:db8::68", out: common.MapStr{"ip": "2001:db8::68"}},
		"Invalid": {ip: "192.0.1", out: nil},
	} {
		t.Run(name, func(t *testing.T) {
			u := User{IP: net.ParseIP(tc.ip), ID: id, Email: email, Name: userName, UserAgent: userAgent}
			assert.Equal(t, tc.out, u.ClientFields())
		})
	}

	t.Run("NilValues", func(t *testing.T) {
		var u User
		assert.Nil(t, u.ClientFields())
		u = User{}
		assert.Nil(t, u.ClientFields())
	})
}

func TestUserAgentFields(t *testing.T) {
	id := "1234"
	email := "test@mail.co"
	name := "user123"
	userAgent := "rum-1.0"

	tests := []struct {
		User   User
		Output common.MapStr
	}{
		{
			User:   User{},
			Output: nil,
		},
		{
			User: User{
				ID:        id,
				IP:        net.ParseIP("127.0.0.1"),
				Email:     email,
				Name:      name,
				UserAgent: userAgent,
			},
			Output: common.MapStr{"original": "rum-1.0"},
		},
	}

	for _, test := range tests {
		output := test.User.UserAgentFields()
		assert.Equal(t, test.Output, output)
	}
}
