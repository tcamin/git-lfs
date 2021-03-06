package lfsapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndpointDefaultsToOrigin(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"remote.origin.lfsurl": "abc",
	}))

	e := finder.Endpoint("download", "")
	assert.Equal(t, "abc", e.Url)
	assert.Equal(t, "", e.SshUserAndHost)
	assert.Equal(t, "", e.SshPath)
}

func TestEndpointOverridesOrigin(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"lfs.url":              "abc",
		"remote.origin.lfsurl": "def",
	}))

	e := finder.Endpoint("download", "")
	assert.Equal(t, "abc", e.Url)
	assert.Equal(t, "", e.SshUserAndHost)
	assert.Equal(t, "", e.SshPath)
}

func TestEndpointNoOverrideDefaultRemote(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"remote.origin.lfsurl": "abc",
		"remote.other.lfsurl":  "def",
	}))

	e := finder.Endpoint("download", "")
	assert.Equal(t, "abc", e.Url)
	assert.Equal(t, "", e.SshUserAndHost)
	assert.Equal(t, "", e.SshPath)
}

func TestEndpointUseAlternateRemote(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"remote.origin.lfsurl": "abc",
		"remote.other.lfsurl":  "def",
	}))

	e := finder.Endpoint("download", "other")
	assert.Equal(t, "def", e.Url)
	assert.Equal(t, "", e.SshUserAndHost)
	assert.Equal(t, "", e.SshPath)
}

func TestEndpointAddsLfsSuffix(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"remote.origin.url": "https://example.com/foo/bar",
	}))

	e := finder.Endpoint("download", "")
	assert.Equal(t, "https://example.com/foo/bar.git/info/lfs", e.Url)
	assert.Equal(t, "", e.SshUserAndHost)
	assert.Equal(t, "", e.SshPath)
}

func TestBareEndpointAddsLfsSuffix(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"remote.origin.url": "https://example.com/foo/bar.git",
	}))

	e := finder.Endpoint("download", "")
	assert.Equal(t, "https://example.com/foo/bar.git/info/lfs", e.Url)
	assert.Equal(t, "", e.SshUserAndHost)
	assert.Equal(t, "", e.SshPath)
}

func TestEndpointSeparateClonePushUrl(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"remote.origin.url":     "https://example.com/foo/bar.git",
		"remote.origin.pushurl": "https://readwrite.com/foo/bar.git",
	}))

	e := finder.Endpoint("download", "")
	assert.Equal(t, "https://example.com/foo/bar.git/info/lfs", e.Url)
	assert.Equal(t, "", e.SshUserAndHost)
	assert.Equal(t, "", e.SshPath)

	e = finder.Endpoint("upload", "")
	assert.Equal(t, "https://readwrite.com/foo/bar.git/info/lfs", e.Url)
	assert.Equal(t, "", e.SshUserAndHost)
	assert.Equal(t, "", e.SshPath)
}

func TestEndpointOverriddenSeparateClonePushLfsUrl(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"remote.origin.url":        "https://example.com/foo/bar.git",
		"remote.origin.pushurl":    "https://readwrite.com/foo/bar.git",
		"remote.origin.lfsurl":     "https://examplelfs.com/foo/bar",
		"remote.origin.lfspushurl": "https://readwritelfs.com/foo/bar",
	}))

	e := finder.Endpoint("download", "")
	assert.Equal(t, "https://examplelfs.com/foo/bar", e.Url)
	assert.Equal(t, "", e.SshUserAndHost)
	assert.Equal(t, "", e.SshPath)

	e = finder.Endpoint("upload", "")
	assert.Equal(t, "https://readwritelfs.com/foo/bar", e.Url)
	assert.Equal(t, "", e.SshUserAndHost)
	assert.Equal(t, "", e.SshPath)
}

func TestEndpointGlobalSeparateLfsPush(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"lfs.url":     "https://readonly.com/foo/bar",
		"lfs.pushurl": "https://write.com/foo/bar",
	}))

	e := finder.Endpoint("download", "")
	assert.Equal(t, "https://readonly.com/foo/bar", e.Url)
	assert.Equal(t, "", e.SshUserAndHost)
	assert.Equal(t, "", e.SshPath)

	e = finder.Endpoint("upload", "")
	assert.Equal(t, "https://write.com/foo/bar", e.Url)
	assert.Equal(t, "", e.SshUserAndHost)
	assert.Equal(t, "", e.SshPath)
}

func TestSSHEndpointOverridden(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"remote.origin.url":    "git@example.com:foo/bar",
		"remote.origin.lfsurl": "lfs",
	}))

	e := finder.Endpoint("download", "")
	assert.Equal(t, "lfs", e.Url)
	assert.Equal(t, "", e.SshUserAndHost)
	assert.Equal(t, "", e.SshPath)
	assert.Equal(t, "", e.SshPort)
}

func TestSSHEndpointAddsLfsSuffix(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"remote.origin.url": "ssh://git@example.com/foo/bar",
	}))

	e := finder.Endpoint("download", "")
	assert.Equal(t, "https://example.com/foo/bar.git/info/lfs", e.Url)
	assert.Equal(t, "git@example.com", e.SshUserAndHost)
	assert.Equal(t, "foo/bar", e.SshPath)
	assert.Equal(t, "", e.SshPort)
}

func TestSSHCustomPortEndpointAddsLfsSuffix(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"remote.origin.url": "ssh://git@example.com:9000/foo/bar",
	}))

	e := finder.Endpoint("download", "")
	assert.Equal(t, "https://example.com/foo/bar.git/info/lfs", e.Url)
	assert.Equal(t, "git@example.com", e.SshUserAndHost)
	assert.Equal(t, "foo/bar", e.SshPath)
	assert.Equal(t, "9000", e.SshPort)
}

func TestBareSSHEndpointAddsLfsSuffix(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"remote.origin.url": "git@example.com:foo/bar.git",
	}))

	e := finder.Endpoint("download", "")
	assert.Equal(t, "https://example.com/foo/bar.git/info/lfs", e.Url)
	assert.Equal(t, "git@example.com", e.SshUserAndHost)
	assert.Equal(t, "foo/bar.git", e.SshPath)
	assert.Equal(t, "", e.SshPort)
}

func TestBareSSSHEndpointWithCustomPortInBrackets(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"remote.origin.url": "[git@example.com:2222]:foo/bar.git",
	}))

	e := finder.Endpoint("download", "")
	assert.Equal(t, "https://example.com/foo/bar.git/info/lfs", e.Url)
	assert.Equal(t, "git@example.com", e.SshUserAndHost)
	assert.Equal(t, "foo/bar.git", e.SshPath)
	assert.Equal(t, "2222", e.SshPort)
}

func TestSSHEndpointFromGlobalLfsUrl(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"lfs.url": "git@example.com:foo/bar.git",
	}))

	e := finder.Endpoint("download", "")
	assert.Equal(t, "https://example.com/foo/bar.git", e.Url)
	assert.Equal(t, "git@example.com", e.SshUserAndHost)
	assert.Equal(t, "foo/bar.git", e.SshPath)
	assert.Equal(t, "", e.SshPort)
}

func TestHTTPEndpointAddsLfsSuffix(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"remote.origin.url": "http://example.com/foo/bar",
	}))

	e := finder.Endpoint("download", "")
	assert.Equal(t, "http://example.com/foo/bar.git/info/lfs", e.Url)
	assert.Equal(t, "", e.SshUserAndHost)
	assert.Equal(t, "", e.SshPath)
	assert.Equal(t, "", e.SshPort)
}

func TestBareHTTPEndpointAddsLfsSuffix(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"remote.origin.url": "http://example.com/foo/bar.git",
	}))

	e := finder.Endpoint("download", "")
	assert.Equal(t, "http://example.com/foo/bar.git/info/lfs", e.Url)
	assert.Equal(t, "", e.SshUserAndHost)
	assert.Equal(t, "", e.SshPath)
	assert.Equal(t, "", e.SshPort)
}

func TestGitEndpointAddsLfsSuffix(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"remote.origin.url": "git://example.com/foo/bar",
	}))

	e := finder.Endpoint("download", "")
	assert.Equal(t, "https://example.com/foo/bar.git/info/lfs", e.Url)
	assert.Equal(t, "", e.SshUserAndHost)
	assert.Equal(t, "", e.SshPath)
	assert.Equal(t, "", e.SshPort)
}

func TestGitEndpointAddsLfsSuffixWithCustomProtocol(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"remote.origin.url": "git://example.com/foo/bar",
		"lfs.gitprotocol":   "http",
	}))

	e := finder.Endpoint("download", "")
	assert.Equal(t, "http://example.com/foo/bar.git/info/lfs", e.Url)
	assert.Equal(t, "", e.SshUserAndHost)
	assert.Equal(t, "", e.SshPath)
	assert.Equal(t, "", e.SshPort)
}

func TestBareGitEndpointAddsLfsSuffix(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"remote.origin.url": "git://example.com/foo/bar.git",
	}))

	e := finder.Endpoint("download", "")
	assert.Equal(t, "https://example.com/foo/bar.git/info/lfs", e.Url)
	assert.Equal(t, "", e.SshUserAndHost)
	assert.Equal(t, "", e.SshPath)
	assert.Equal(t, "", e.SshPort)
}

func TestLocalPathEndpointAddsDotGitDir(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"remote.origin.url": "/local/path",
	}))
	e := finder.Endpoint("download", "")
	assert.Equal(t, "file:///local/path/.git/info/lfs", e.Url)
}

func TestLocalPathEndpointPreservesDotGit(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"remote.origin.url": "/local/path.git",
	}))
	e := finder.Endpoint("download", "")
	assert.Equal(t, "file:///local/path.git/info/lfs", e.Url)
}

func TestAccessConfig(t *testing.T) {
	type accessTest struct {
		Access        string
		PrivateAccess bool
	}

	tests := map[string]accessTest{
		"":            {"none", false},
		"basic":       {"basic", true},
		"BASIC":       {"basic", true},
		"private":     {"basic", true},
		"PRIVATE":     {"basic", true},
		"invalidauth": {"invalidauth", true},
	}

	for value, expected := range tests {
		finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
			"lfs.url":                        "http://example.com",
			"lfs.http://example.com.access":  value,
			"lfs.https://example.com.access": "bad",
		}))

		dl := finder.Endpoint("upload", "")
		ul := finder.Endpoint("download", "")

		if access := finder.AccessFor(dl.Url); access != Access(expected.Access) {
			t.Errorf("Expected Access() with value %q to be %v, got %v", value, expected.Access, access)
		}
		if access := finder.AccessFor(ul.Url); access != Access(expected.Access) {
			t.Errorf("Expected Access() with value %q to be %v, got %v", value, expected.Access, access)
		}
	}

	// Test again but with separate push url
	for value, expected := range tests {
		finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
			"lfs.url":                           "http://example.com",
			"lfs.pushurl":                       "http://examplepush.com",
			"lfs.http://example.com.access":     value,
			"lfs.http://examplepush.com.access": value,
			"lfs.https://example.com.access":    "bad",
		}))

		dl := finder.Endpoint("upload", "")
		ul := finder.Endpoint("download", "")

		if access := finder.AccessFor(dl.Url); access != Access(expected.Access) {
			t.Errorf("Expected Access() with value %q to be %v, got %v", value, expected.Access, access)
		}
		if access := finder.AccessFor(ul.Url); access != Access(expected.Access) {
			t.Errorf("Expected Access() with value %q to be %v, got %v", value, expected.Access, access)
		}
	}
}

func TestAccessAbsentConfig(t *testing.T) {
	finder := NewEndpointFinder(nil)
	assert.Equal(t, NoneAccess, finder.AccessFor(finder.Endpoint("download", "").Url))
	assert.Equal(t, NoneAccess, finder.AccessFor(finder.Endpoint("upload", "").Url))
}

func TestSetAccess(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{}))

	assert.Equal(t, NoneAccess, finder.AccessFor("http://example.com"))
	finder.SetAccess("http://example.com", NTLMAccess)
	assert.Equal(t, NTLMAccess, finder.AccessFor("http://example.com"))
}

func TestChangeAccess(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"lfs.http://example.com.access": "basic",
	}))

	assert.Equal(t, BasicAccess, finder.AccessFor("http://example.com"))
	finder.SetAccess("http://example.com", NTLMAccess)
	assert.Equal(t, NTLMAccess, finder.AccessFor("http://example.com"))
}

func TestDeleteAccessWithNone(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"lfs.http://example.com.access": "basic",
	}))

	assert.Equal(t, BasicAccess, finder.AccessFor("http://example.com"))
	finder.SetAccess("http://example.com", NoneAccess)
	assert.Equal(t, NoneAccess, finder.AccessFor("http://example.com"))
}

func TestDeleteAccessWithEmptyString(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"lfs.http://example.com.access": "basic",
	}))

	assert.Equal(t, BasicAccess, finder.AccessFor("http://example.com"))
	finder.SetAccess("http://example.com", Access(""))
	assert.Equal(t, NoneAccess, finder.AccessFor("http://example.com"))
}

type EndpointParsingTestCase struct {
	Given    string
	Expected Endpoint
}

func (c *EndpointParsingTestCase) Assert(t *testing.T) {
	finder := NewEndpointFinder(NewContext(nil, nil, map[string]string{
		"url.https://github.com/.insteadof": "gh:",
	}))
	actual := finder.NewEndpoint(c.Given)
	assert.Equal(t, c.Expected, actual, "lfsapi: expected endpoint for %q to be %#v (was %#v)", c.Given, c.Expected, actual)
}

func TestEndpointParsing(t *testing.T) {
	// Note that many of these tests will produce silly or completely broken
	// values for the Url, and that's okay: they work nevertheless.
	for desc, c := range map[string]EndpointParsingTestCase{
		"simple bare ssh": {
			"git@github.com:git-lfs/git-lfs.git",
			Endpoint{
				Url:            "https://github.com/git-lfs/git-lfs.git",
				SshUserAndHost: "git@github.com",
				SshPath:        "git-lfs/git-lfs.git",
				SshPort:        "",
				Operation:      "",
			},
		},
		"port bare ssh": {
			"[git@ssh.github.com:443]:git-lfs/git-lfs.git",
			Endpoint{
				Url:            "https://ssh.github.com/git-lfs/git-lfs.git",
				SshUserAndHost: "git@ssh.github.com",
				SshPath:        "git-lfs/git-lfs.git",
				SshPort:        "443",
				Operation:      "",
			},
		},
		"no user bare ssh": {
			"github.com:git-lfs/git-lfs.git",
			Endpoint{
				Url:            "https://github.com/git-lfs/git-lfs.git",
				SshUserAndHost: "github.com",
				SshPath:        "git-lfs/git-lfs.git",
				SshPort:        "",
				Operation:      "",
			},
		},
		"bare word bare ssh": {
			"github:git-lfs/git-lfs.git",
			Endpoint{
				Url:            "https://github/git-lfs/git-lfs.git",
				SshUserAndHost: "github",
				SshPath:        "git-lfs/git-lfs.git",
				SshPort:        "",
				Operation:      "",
			},
		},
		"insteadof alias": {
			"gh:git-lfs/git-lfs.git",
			Endpoint{
				Url:            "https://github.com/git-lfs/git-lfs.git",
				SshUserAndHost: "",
				SshPath:        "",
				SshPort:        "",
				Operation:      "",
			},
		},
		"remote helper": {
			"remote::git-lfs/git-lfs.git",
			Endpoint{
				Url:            "remote::git-lfs/git-lfs.git",
				SshUserAndHost: "",
				SshPath:        "",
				SshPort:        "",
				Operation:      "",
			},
		},
	} {
		t.Run(desc, c.Assert)
	}
}
