package github

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-github/v39/github"
	"github.com/openshift-pipelines/pipelines-as-code/pkg/params"
	"github.com/openshift-pipelines/pipelines-as-code/pkg/params/clients"
	"github.com/openshift-pipelines/pipelines-as-code/pkg/params/info"
	ghtesthelper "github.com/openshift-pipelines/pipelines-as-code/pkg/test/github"
	"gotest.tools/v3/assert"
	rtesting "knative.dev/pkg/reconciler/testing"
)

func TestPayLoadFix(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/pull_request_with_newlines.json")
	assert.NilError(t, err)
	ctx, _ := rtesting.SetupFakeContext(t)
	fakeclient, mux, _, teardown := ghtesthelper.SetupGH()
	defer teardown()
	mux.HandleFunc("/repos/repo/owner/commits/SHA", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"commit": {"message": "HELLO"}}`)
	})
	mux.HandleFunc("/repos/repo/owner/git/commits/SHA", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"commit": {"message": "HELLO"}}`)
	})

	gvcs := VCS{
		Client: fakeclient,
	}

	logger, _ := getLogger()

	event := &info.Event{
		EventType:     "pull_request",
		TriggerTarget: "pull_request",
	}
	run := &params.Run{
		Clients: clients.Clients{
			Log: logger,
		},
		Info: info.Info{
			Event: event,
		},
	}
	_, err = gvcs.ParsePayload(ctx, run, string(b))
	assert.NilError(t, err)

	// would bomb out on "assertion failed: error is not nil: invalid character
	// '\n' in string literal" if we don't fix the payload
	assert.NilError(t, err)
}

func TestParsePayloadRerequestFromPullRequest(t *testing.T) {
	checkrunSender := "jean-pierre"
	prOwner := "owner"
	repoName := "repo"
	prNumber := "123"
	sha := "ParsePayloadRerequestFromPullRequestSHA"
	checkrunEvent := fmt.Sprintf(`{"action": "rerequested",
	"sender": {"login": "%s"},
	"check_run": {"check_suite": {"pull_requests": [{"number": %s}]}},
	"repository": {"name": "%s", "owner": {"login": "%s"}}}`,
		checkrunSender, prNumber, repoName, prOwner)
	fakeclient, mux, _, teardown := ghtesthelper.SetupGH()
	defer teardown()
	mux.HandleFunc("/repos/"+prOwner+"/"+repoName+"/pulls/"+prNumber, func(rw http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(rw, `{"head": {"sha": "%s", "ref": "123"}, "user": {"login": "%s"}}`, sha, prOwner)
	})
	mux.HandleFunc("/repos/owner/repo/commits/"+sha, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"commit": {"message": "HELLO"}}`)
	})
	mux.HandleFunc("/repos/owner/repo/git/commits/"+sha, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"commit": {"message": "HELLO"}}`)
	})
	ctx, _ := rtesting.SetupFakeContext(t)
	gvcs := VCS{
		Client: fakeclient,
	}
	logger, observer := getLogger()
	event := &info.Event{
		EventType:     "check_run",
		TriggerTarget: "issue-recheck",
	}
	run := &params.Run{
		Clients: clients.Clients{
			Log: logger,
		},
		Info: info.Info{
			Event: event,
		},
	}
	runinfo, err := gvcs.ParsePayload(ctx, run, checkrunEvent)
	assert.NilError(t, err)

	assert.Equal(t, prOwner, runinfo.Owner)
	assert.Equal(t, repoName, runinfo.Repository)
	assert.Assert(t, checkrunSender != runinfo.Sender)
	assert.Equal(t, runinfo.EventType, "pull_request")
	assert.Assert(t, strings.Contains(observer.TakeAll()[0].Message, "Recheck of PR"))
}

func TestParsePayloadRerequestFromPush(t *testing.T) {
	sender := "jean-pierre"
	headBranch := "tartonpion"
	headSHA := "TestParsePayloadRerequestFromPushSHA"
	owner := "owner"
	repository := "repo"
	url := fmt.Sprintf("https://github.com/%s/%s", owner, repository)
	checkrunEvent := fmt.Sprintf(`{
  "action": "rerequested",
  "check_run": {
    "check_suite": {
      "head_branch": "%s",
      "head_sha": "%s",
      "pull_requests": []
    }
  },
  "repository": {
    "default_branch": "main",
    "html_url": "%s",
    "name": "%s",
    "owner": {
      "login": "%s"
    }
  },
  "sender": {
    "login": "%s"
  }
}`,
		headBranch, headSHA, url, repository, owner, sender)
	ctx, _ := rtesting.SetupFakeContext(t)
	fakeclient, mux, _, teardown := ghtesthelper.SetupGH()
	defer teardown()
	gvcs := VCS{
		Client: fakeclient,
	}
	mux.HandleFunc("/repos/owner/repo/commits/"+headSHA, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"commit": {"message": "HELLO"}}`)
	})
	mux.HandleFunc("/repos/owner/repo/git/commits/"+headSHA, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"commit": {"message": "HELLO"}}`)
	})

	event := &info.Event{
		EventType:     "check_run",
		TriggerTarget: "issue-recheck",
	}
	logger, _ := getLogger()

	run := &params.Run{
		Clients: clients.Clients{
			Log: logger,
		},
		Info: info.Info{
			Event: event,
		},
	}
	runinfo, err := gvcs.ParsePayload(ctx, run, checkrunEvent)
	assert.NilError(t, err)

	assert.Equal(t, runinfo.EventType, "push")
	assert.Equal(t, runinfo.HeadBranch, headBranch)
	assert.Equal(t, runinfo.Owner, owner)
	assert.Equal(t, runinfo.Repository, repository)
	assert.Equal(t, runinfo.URL, url)
	assert.Assert(t, sender == runinfo.Sender) // TODO: should it be set to the push sender?
}

func TestParsePayLoadRetest(t *testing.T) {
	issueSender := "tartanpion"
	prOwner := "user1"
	repoOwner := "openshift"
	repoName := "pipelines"
	prNumber := "123"
	sha := "TestParsePayLoadRetestSHA"

	fakeclient, mux, _, teardown := ghtesthelper.SetupGH()
	defer teardown()
	mux.HandleFunc("/repos/"+prOwner+"/"+repoName+"/pulls/"+prNumber, func(rw http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(rw, `{"head": {"sha": "%s", "ref": "123"}, "user": {"login": "%s"}}`, sha, prOwner)
	})
	mux.HandleFunc("/repos/user1/pipelines/commits/"+sha, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"commit": {"message": "HELLO"}}`)
	})
	mux.HandleFunc("/repos/user1/pipelines/git/commits/"+sha, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"commit": {"message": "HELLO"}}`)
	})
	issueEvent := fmt.Sprintf(`{
  "sender": {
	"login": "%s"
  },
  "repository": {
	"name": "%s",
	"owner": {
	  "login": "%s"
	}
  },
  "issue": {
	"pull_request": {
	  "html_url": "https://github.com/%s/%s/pull/%s"
	}
  }
}`, issueSender, repoName, prOwner, repoName, repoOwner, prNumber)

	ctx, _ := rtesting.SetupFakeContext(t)
	logger, observer := getLogger()
	gvcs := VCS{
		Client: fakeclient,
	}
	// TODO
	event := &info.Event{
		EventType:     "issue_comment",
		TriggerTarget: "issue_comment",
	}
	run := &params.Run{
		Clients: clients.Clients{
			Log: logger,
		},
		Info: info.Info{
			Event: event,
			Pac: info.PacOpts{
				VCSToken: "TOKENSET",
			},
		},
	}

	runinfo, err := gvcs.ParsePayload(ctx, run, issueEvent)
	assert.NilError(t, err)
	assert.Equal(t, prOwner, runinfo.Owner)
	// Make sure the PR owner is the runinfo.Owner and not the issueSender
	assert.Assert(t, issueSender != runinfo.Owner)
	firstObservedMessage := observer.TakeAll()[0].Message
	assert.Assert(t, strings.Contains(firstObservedMessage, "recheck"))
	assert.Equal(t, runinfo.EventType, "pull_request")

	// We cannot parse payload without a token when getting from issue_comment
	run.Info.Pac = info.PacOpts{VCSToken: ""}
	runinfo, err = gvcs.ParsePayload(ctx, run, issueEvent)
	assert.ErrorContains(t, err, "gitops style comments operation is only supported with github apps")
	assert.Assert(t, runinfo.Event == nil)
}

func TestParsePayload(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/pull_request.json")
	assert.NilError(t, err)
	fakeclient, mux, _, teardown := ghtesthelper.SetupGH()
	defer teardown()
	mux.HandleFunc("/repos/chmouel/scratchpad/commits/cc8334de8e056317d18bd00c2588c3f7c95af294", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"commit": {"message": "HELLO"}}`)
	})
	mux.HandleFunc("/repos/chmouel/scratchpad/git/commits/cc8334de8e056317d18bd00c2588c3f7c95af294", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"commit": {"message": "HELLO"}}`)
	})
	gvcs := VCS{
		Client: fakeclient,
	}
	ctx, _ := rtesting.SetupFakeContext(t)
	logger, _ := getLogger()

	event := &info.Event{
		EventType:     "pull_request",
		TriggerTarget: "pull_request",
	}
	run := &params.Run{
		Clients: clients.Clients{
			Log: logger,
		},
		Info: info.Info{
			Event: event,
		},
	}
	runinfo, err := gvcs.ParsePayload(ctx, run, string(b))
	assert.NilError(t, err)
	assert.Assert(t, runinfo.BaseBranch == "master")
	assert.Assert(t, runinfo.Owner == "chmouel")
	assert.Assert(t, runinfo.Repository == "scratchpad")
	assert.Equal(t, runinfo.EventType, "pull_request")
	assert.Assert(t, runinfo.URL == "https://github.com/chmouel/scratchpad")
}

func TestParsePayloadInvalid(t *testing.T) {
	ctx, _ := rtesting.SetupFakeContext(t)
	gvcs := VCS{
		Token:  github.String("none"),
		APIURL: github.String("nothing"),
	}
	logger, _ := getLogger()
	event := &info.Event{
		EventType:     "pull_request",
		TriggerTarget: "pull_request",
	}
	run := &params.Run{
		Clients: clients.Clients{
			Log: logger,
		},
		Info: info.Info{
			Event: event,
		},
	}
	_, err := gvcs.ParsePayload(ctx, run, "hello moto")
	assert.ErrorContains(t, err, "invalid character")
}

func TestParsePayloadUnkownEvent(t *testing.T) {
	ctx, _ := rtesting.SetupFakeContext(t)
	gvcs := VCS{
		Token:  github.String("none"),
		APIURL: github.String(""),
	}
	logger, _ := getLogger()
	event := &info.Event{
		EventType:     "foo",
		TriggerTarget: "foo",
	}
	run := &params.Run{
		Clients: clients.Clients{
			Log: logger,
		},
		Info: info.Info{
			Event: event,
		},
	}
	_, err := gvcs.ParsePayload(ctx, run, "{\"hello\": \"moto\"}")
	assert.ErrorContains(t, err, "unknown X-Github-Event")
}

func TestParsePayCannotParse(t *testing.T) {
	ctx, _ := rtesting.SetupFakeContext(t)
	gvcs := VCS{
		Token:  github.String("none"),
		APIURL: github.String(""),
	}
	logger, _ := getLogger()
	event := &info.Event{
		EventType:     "gollum",
		TriggerTarget: "gollum",
	}
	run := &params.Run{
		Clients: clients.Clients{
			Log: logger,
		},
		Info: info.Info{
			Event: event,
		},
	}
	_, err := gvcs.ParsePayload(ctx, run, "{}")
	assert.Error(t, err, "this event is not supported")
}
