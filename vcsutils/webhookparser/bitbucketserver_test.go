package webhookparser

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/jfrog/froggit-go/vcsutils"
	"github.com/stretchr/testify/assert"
)

const (
	bitbucketServerPushSha256       = "c7921dbb415bf4bb47098f8612212bdf2d9fb28f4c147a9544a7163f4c65b540"
	bitbucketServerPushExpectedTime = int64(1631178392)

	bitbucketServerPrCreateExpectedTime = int64(1631178661)
	bitbucketServerPrCreatedSha256      = "d3ed184563b8e373b9fd5cfe869d23a5f60fb2095e8254f1c3195a99dfe6fc1a"

	bitbucketServerPrUpdateExpectedTime = int64(1631180185)
	bitbucketServerPrUpdatedSha256      = "c0d3d62cad1fa261b9cfd092bffe49c74971109d9c617a7ee55dd639d59f8377"
)

func TestBitbucketServerParseIncomingPushWebhook(t *testing.T) {
	reader, err := os.Open(filepath.Join("testdata", "bitbucketserver", "pushpayload"))
	assert.NoError(t, err)
	defer reader.Close()

	// Create request
	request := httptest.NewRequest("POST", "https://127.0.0.1", reader)
	request.Header.Add(EventHeaderKey, "repo:refs_changed")
	request.Header.Add(Sha256Signature, "sha256="+bitbucketServerPushSha256)

	// Parse webhook
	actual, err := ParseIncomingWebhook(vcsutils.BitbucketServer, token, request)
	assert.NoError(t, err)

	// Check values
	assert.Equal(t, "~"+expectedRepoName, actual.Repository)
	assert.Equal(t, expectedBranch, actual.Branch)
	assert.Equal(t, bitbucketServerPushExpectedTime, actual.Timestamp)
	assert.Equal(t, vcsutils.Push, actual.Event)
}

func TestBitbucketServerParseIncomingPrCreateWebhook(t *testing.T) {
	reader, err := os.Open(filepath.Join("testdata", "bitbucketserver", "prcreatepayload"))
	assert.NoError(t, err)
	defer reader.Close()

	// Create request
	request := httptest.NewRequest("POST", "https://127.0.0.1?", reader)
	request.Header.Add(EventHeaderKey, "pr:opened")
	request.Header.Add(Sha256Signature, "sha256="+bitbucketServerPrCreatedSha256)

	// Parse webhook
	actual, err := ParseIncomingWebhook(vcsutils.BitbucketServer, token, request)
	assert.NoError(t, err)

	// Check values
	assert.Equal(t, "~"+expectedRepoName, actual.Repository)
	assert.Equal(t, expectedBranch, actual.Branch)
	assert.Equal(t, bitbucketServerPrCreateExpectedTime, actual.Timestamp)
	assert.Equal(t, "~"+expectedRepoName, actual.SourceRepository)
	assert.Equal(t, expectedSourceBranch, actual.SourceBranch)
	assert.Equal(t, vcsutils.PrCreated, actual.Event)
}

func TestBitbucketServerParseIncomingPrUpdateWebhook(t *testing.T) {
	reader, err := os.Open(filepath.Join("testdata", "bitbucketserver", "prupdatepayload"))
	assert.NoError(t, err)
	defer reader.Close()

	// Create request
	request := httptest.NewRequest("POST", "https://127.0.0.1", reader)
	request.Header.Add(EventHeaderKey, "pr:from_ref_updated")
	request.Header.Add(Sha256Signature, "sha256="+bitbucketServerPrUpdatedSha256)

	// Parse webhook
	actual, err := ParseIncomingWebhook(vcsutils.BitbucketServer, token, request)
	assert.NoError(t, err)

	// Check values
	assert.Equal(t, "~"+expectedRepoName, actual.Repository)
	assert.Equal(t, expectedBranch, actual.Branch)
	assert.Equal(t, bitbucketServerPrUpdateExpectedTime, actual.Timestamp)
	assert.Equal(t, "~"+expectedRepoName, actual.SourceRepository)
	assert.Equal(t, expectedSourceBranch, actual.SourceBranch)
	assert.Equal(t, vcsutils.PrEdited, actual.Event)
}

func TestBitbucketServerPayloadMismatchSignature(t *testing.T) {
	reader, err := os.Open(filepath.Join("testdata", "bitbucketserver", "pushpayload"))
	assert.NoError(t, err)
	defer reader.Close()

	// Create request
	request := httptest.NewRequest("POST", "https://127.0.0.1", reader)
	request.Header.Add(EventHeaderKey, "repo:refs_changed")
	request.Header.Add(Sha256Signature, "sha256=wrongsianature")

	// Parse webhook
	_, err = ParseIncomingWebhook(vcsutils.BitbucketServer, token, request)
	assert.EqualError(t, err, "Payload signature mismatch")
}
