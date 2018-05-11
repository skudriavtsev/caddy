// Copyright 2015 Light Code Labs, LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package caddyhttp

import (
	// plug in the server
	_ "github.com/skudriavtsev/caddy/caddyhttp/httpserver"

	// plug in the standard directives
	_ "github.com/skudriavtsev/caddy/caddyhttp/basicauth"
	_ "github.com/skudriavtsev/caddy/caddyhttp/bind"
	_ "github.com/skudriavtsev/caddy/caddyhttp/browse"
	_ "github.com/skudriavtsev/caddy/caddyhttp/errors"
	_ "github.com/skudriavtsev/caddy/caddyhttp/expvar"
	_ "github.com/skudriavtsev/caddy/caddyhttp/extensions"
	_ "github.com/skudriavtsev/caddy/caddyhttp/fastcgi"
	_ "github.com/skudriavtsev/caddy/caddyhttp/gzip"
	_ "github.com/skudriavtsev/caddy/caddyhttp/header"
	_ "github.com/skudriavtsev/caddy/caddyhttp/index"
	_ "github.com/skudriavtsev/caddy/caddyhttp/internalsrv"
	_ "github.com/skudriavtsev/caddy/caddyhttp/limits"
	_ "github.com/skudriavtsev/caddy/caddyhttp/log"
	_ "github.com/skudriavtsev/caddy/caddyhttp/markdown"
	_ "github.com/skudriavtsev/caddy/caddyhttp/mime"
	_ "github.com/skudriavtsev/caddy/caddyhttp/pprof"
	_ "github.com/skudriavtsev/caddy/caddyhttp/proxy"
	_ "github.com/skudriavtsev/caddy/caddyhttp/push"
	_ "github.com/skudriavtsev/caddy/caddyhttp/redirect"
	_ "github.com/skudriavtsev/caddy/caddyhttp/ibredirect"
	_ "github.com/skudriavtsev/caddy/caddyhttp/ibproxy"
	_ "github.com/skudriavtsev/caddy/caddyhttp/requestid"
	_ "github.com/skudriavtsev/caddy/caddyhttp/rewrite"
	_ "github.com/skudriavtsev/caddy/caddyhttp/root"
	_ "github.com/skudriavtsev/caddy/caddyhttp/status"
	_ "github.com/skudriavtsev/caddy/caddyhttp/templates"
	_ "github.com/skudriavtsev/caddy/caddyhttp/timeouts"
	_ "github.com/skudriavtsev/caddy/caddyhttp/websocket"
	_ "github.com/skudriavtsev/caddy/onevent"
)
