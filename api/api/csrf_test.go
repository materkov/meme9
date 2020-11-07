package api

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateCSRFAndValidate(t *testing.T) {
	hash := GenerateCSRFToken(151)
	require.True(t, strings.HasPrefix(hash, "151-"))
	require.Len(t, hash, 68)

	require.True(t, ValidateCSRFToken(151, "151-ee6c789031a8738ce6e741e4869a65828ee4e643aaf2be0ec2b699f20812f503"))

	require.False(t, ValidateCSRFToken(151, "151-ee6c789031a8738ce6e741e4869a65828ee4e643aaf2be0ec2b699f20812f50311111111"))
	require.False(t, ValidateCSRFToken(151, "152-ee6c789031a8738ce6e741e4869a65828ee4e643aaf2be0ec2b699f20812f503"))
	require.False(t, ValidateCSRFToken(152, "151-ee6c789031a8738ce6e741e4869a65828ee4e643aaf2be0ec2b699f20812f503"))
}
