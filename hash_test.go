package bf

import (
	"fmt"
	"testing"
)

func TestShaHash_doHash(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		n        byte
		expected string
	}{
		{
			name:     "1 time",
			input:    "hello",
			n:        1,
			expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			name:  "2 times",
			input: "hello",
			n:     2,
			expected: "8a2a5c9b768827de5a9552c38a044c66959c68f6d2f21b5260af54d2f87db827" +
				"cceeb7a985ecc3dabcb4c8f666cd637f16f008e3c963db6aa6f83a7b288c54ef",
		},
		{
			name:  "3 times",
			input: "hello",
			n:     3,
			expected: "8a2a5c9b768827de5a9552c38a044c66959c68f6d2f21b5260af54d2f87db827" +
				"cceeb7a985ecc3dabcb4c8f666cd637f16f008e3c963db6aa6f83a7b288c54ef" +
				"29f3ced0b171e52626c66bedaf76469f1efda5c110b47ea24228ef25e61859cc",
		},
		{
			name:  "4 times",
			input: "hello",
			n:     4,
			expected: "8a2a5c9b768827de5a9552c38a044c66959c68f6d2f21b5260af54d2f87db827" +
				"cceeb7a985ecc3dabcb4c8f666cd637f16f008e3c963db6aa6f83a7b288c54ef" +
				"29f3ced0b171e52626c66bedaf76469f1efda5c110b47ea24228ef25e61859cc" +
				"0b4d354d56ea9a985571a56b1829f33d072e7902c1afaf981381089b9eb00ffe",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := shaHash{}
			input := []byte(tc.input)
			r := h.doHash(tc.n, &input)

			result := fmt.Sprintf("%x", r)
			assertStringEqual(t, result, tc.expected)
		})
	}
}
