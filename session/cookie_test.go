package session

import (
	"bytes"
	"testing"

	"github.com/benjamw/golibs/random"
)

var (
	sigKey   = []byte("fake_sig_key_1_2_3_4_5_6_7_8_9_0") // 32 chars
	cryptKey = []byte("fake_crypt_key_1_2_3_4_5_6_7_8_9") // 32 chars
)

func TestEncodeDecode(t *testing.T) {
	var b bytes.Buffer
	b.Write([]byte(random.Stringn(50)))
	msg := b.Bytes()

	out, err := SignAndEncode(b, sigKey, cryptKey)
	if err != nil {
		t.Fatalf("SignAndEncode threw an error: %v", err)
	}

	in, err := DecodeAndCheckSig(out, sigKey, cryptKey)
	if err != nil {
		t.Fatalf("DecodeAndCheckSig threw an error: %v", err)
	}

	if bytes.Compare(in, msg) != 0 {
		t.Fatalf("Encode then Decode did not return original value. Wanted: %v. Got: %v.", msg, in)
	}
}

func TestPlayerCookie(t *testing.T) {
	in := Data{
		IsPlayer: true,
		PlayerID: random.String(),
	}

	encoded, err := in.Serialize()
	if err != nil {
		t.Errorf("Serialize returned an error: %v", err)
	}

	if len(encoded) < 1 {
		t.Error("Encoded string has no bytes")
	}

	var out Data

	err = out.Deserialize(encoded)
	if err != nil {
		t.Errorf("Deserialize returned an error: %v", err)
	}
	if !out.IsPlayer {
		t.Error("Doesn't know that they are still a player")
	}
	if out.IsSuperUser {
		t.Error("Thinks a normal user is a super user")
	}
	if out.PlayerID != in.PlayerID {
		t.Error("Didn't encode and decode the player key to the same thing")
	}

	// test with really long ID string

	in2 := in
	in2.PlayerID = random.Stringn(1025)

	encoded2, err := in2.Serialize()
	if err != nil {
		t.Errorf("The second Serialize returned an error: %v", err)
	}

	if len(encoded2) < 1 {
		t.Error("The second encoded string has no bytes")
	}

	var out2 Data

	err = out2.Deserialize(encoded2)
	if err != nil {
		t.Errorf("The second Deserialize returned an error: %v", err)
	}
	if !out2.IsPlayer {
		t.Error("Doesn't know that the second player is still a player")
	}
	if out2.IsSuperUser {
		t.Error("Thinks the second normal user is a super user")
	}
	if out2.PlayerID != in2.PlayerID {
		t.Error("Didn't encode and decode the second player key to the same thing")
	}
}
