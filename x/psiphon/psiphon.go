// Copyright 2024 Jigsaw Operations LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package psiphon

import (
	"context"
	"net"

	"github.com/Jigsaw-Code/outline-sdk/transport"
	psi "github.com/Psiphon-Labs/psiphon-tunnel-core/psiphon"
)

type PsiphonDialer struct {
	cancel     context.CancelFunc
	controller *psi.Controller
}

func (d *PsiphonDialer) DialStream(ctx context.Context, addr string) (transport.StreamConn, error) {
	netConn, err := d.controller.Dial(addr, nil)
	if err != nil {
		return nil, err
	}
	return streamConn{netConn}, nil
}

func (d *PsiphonDialer) Close() {
	d.cancel()
}

func NewStreamDialer(configJSON string) (*PsiphonDialer, error) {
	config, err := psi.LoadConfig([]byte(configJSON))
	if err != nil {
		return nil, err
	}
	controller, err := psi.NewController(config)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	controller.Run(ctx)
	return &PsiphonDialer{cancel, controller}, nil
}

var _ transport.StreamDialer = (*PsiphonDialer)(nil)

// streamConn wraps a [net.Conn] to provide a [transport.StreamConn] interface.
type streamConn struct {
	net.Conn
}

var _ transport.StreamConn = (*streamConn)(nil)

func (c streamConn) CloseWrite() error {
	return nil
}

func (c streamConn) CloseRead() error {
	return nil
}