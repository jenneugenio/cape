package capepg

import (
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

func TestNew(t *testing.T) {
	wantPool := &pgxpool.Pool{}
	pg := New(wantPool)

	gotPool := pg.pool
	if gotPool != wantPool {
		t.Errorf("CapePg create with incorrect pool: got %v want %v", gotPool, wantPool)
	}
}

func TestPolicies(t *testing.T) {
	pg := New(&pgxpool.Pool{})
	got := pg.Policies()
	if got == nil {
		t.Errorf("Policies() did not return a PolicyDB")
	}
}

func TestRoles(t *testing.T) {
	pg := New(&pgxpool.Pool{})
	got := pg.Roles()
	if got == nil {
		t.Errorf("Roles() did not return a RoleDB")
	}
}
