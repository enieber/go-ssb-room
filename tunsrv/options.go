package tunsrv

import (
	"context"
	"fmt"
	"strings"

	"go.mindeco.de/ssb-tunnel/internal/keys"
	"go.mindeco.de/ssb-tunnel/internal/repo"
)

type Option func(srv *Server) error

// WithRepoPath changes where the replication database and blobs are stored.
func WithRepoPath(path string) Option {
	return func(s *Server) error {
		s.repoPath = path
		return nil
	}
}

// WithAppKey changes the appkey (aka secret-handshake network cap).
// See https://ssbc.github.io/scuttlebutt-protocol-guide/#handshake for more.
func WithAppKey(k []byte) Option {
	return func(s *Server) error {
		if n := len(k); n != 32 {
			return fmt.Errorf("appKey: need 32 bytes got %d", n)
		}
		s.appKey = k
		return nil
	}
}

// WithNamedKeyPair changes from the default `secret` file, useful for testing.
func WithNamedKeyPair(name string) Option {
	return func(s *Server) error {
		r := repo.New(s.repoPath)
		var err error
		s.KeyPair, err = repo.LoadKeyPair(r, name)
		if err != nil {
			return fmt.Errorf("loading named key-pair %q failed: %w", name, err)
		}
		return nil
	}
}

// WithJSONKeyPair expectes a JSON-string as blob and calls ssb.ParseKeyPair on it.
// This is useful if you dont't want to place the keypair on the filesystem.
func WithJSONKeyPair(blob string) Option {
	return func(s *Server) error {
		var err error
		s.KeyPair, err = keys.ParseKeyPair(strings.NewReader(blob))
		if err != nil {
			return fmt.Errorf("JSON KeyPair decode failed: %w", err)
		}
		return nil
	}
}

// WithKeyPair exepect a initialized ssb.KeyPair. Useful for testing.
func WithKeyPair(kp *keys.KeyPair) Option {
	return func(s *Server) error {
		s.KeyPair = kp
		return nil
	}
}

// WithInfo changes the info/warn/debug loging output.
func WithInfo(log kitlog.Logger) Option {
	return func(s *Server) error {
		s.info = log
		return nil
	}
}

// WithContext changes the context that is context.Background() by default.
// Handy to setup cancelation against a interup signal like ctrl+c.
// Canceling the context also shuts down indexing. If no context is passed sbot.Shutdown() can be used.
func WithContext(ctx context.Context) Option {
	return func(s *Server) error {
		s.rootCtx, s.Shutdown = ShutdownContext(ctx)
		return nil
	}
}
