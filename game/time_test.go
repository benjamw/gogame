package game

import (
	"os"
	"testing"
	"time"

	"github.com/benjamw/golibs/test"
)

func TestMain(m *testing.M) {
	test.InitCtx()
	runVal := m.Run()
	test.ReleaseCtx()
	os.Exit(runVal)
}

func TestSetNow(t *testing.T) {
	ctx := test.GetCtx()

	diff := 500 // 0.5 seconds
	now := time.Now()
	add := 800 // ~ one month
	newNow := now.Add(time.Duration(add) * time.Hour)
	ctx = SetNow(ctx, newNow)

	resp := Now(ctx)

	if resp.Sub(now) < time.Duration(diff)*time.Microsecond {
		t.Fatalf("game.Now (%v) returned time.Now (%v) when SetNow was used", resp, now)
	}

	if resp.Sub(newNow) > time.Duration(0)*time.Microsecond {
		t.Fatalf("game.Now (%v) did not return time.Now (%v) when SetNow was used", resp, newNow)
	}
}

func TestNow(t *testing.T) {
	ctx := test.GetCtx()

	diff := 500 // 0.5 seconds
	now := time.Now()
	resp := Now(ctx)

	if resp.Sub(now) > time.Duration(diff)*time.Microsecond {
		t.Fatalf("game.Now (%v) did not roughly return actual time.Now (%v)", resp, now)
	}

	add := 800 // ~ one month
	newNow := now.Add(time.Duration(add) * time.Hour)
	ctx = SetNow(ctx, newNow)

	resp = Now(ctx)

	if resp.Sub(now) < time.Duration(diff)*time.Microsecond {
		t.Fatalf("game.Now (%v) returned time.Now (%v) when SetNow was used", resp, now)
	}

	if resp.Sub(newNow) > time.Duration(0)*time.Microsecond {
		t.Fatalf("game.Now (%v) did not return exact set time (%v) when SetNow was used", resp, newNow)
	}
}
