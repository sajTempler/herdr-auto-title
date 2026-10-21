package resolver

import (
	"testing"

	"github.com/kryptamine/herdr-auto-title/internal/herdr"
	"github.com/kryptamine/herdr-auto-title/internal/state"
)

// hiddenDirectoryChain is the chain as the plugin ships it, agent name
// hidden, with the directory hidden too.
func hiddenDirectoryChain() *Deterministic {
	return Default(Options{
		MaxLength:     DefaultMaxLength,
		BranchMax:     DefaultBranchMaxLength,
		HideAgentName: true,
		HideDirectory: true,
	})
}

func TestAHiddenDirectoryLeavesTheWorkAlone(t *testing.T) {
	t.Parallel()

	got := hiddenDirectoryChain().Resolve(tabWithPane(&state.PaneState{
		Dir:         dashboard,
		Agent:       "claude",
		AgentStatus: herdr.AgentStatusWorking,
		AgentTitle:  "Implement OAuth scopes",
	}))

	if want := "Implement OAuth scopes"; got.Name != want {
		t.Errorf("name = %q, want %q", got.Name, want)
	}
}

func TestAHiddenDirectoryStillNamesAnIdleShell(t *testing.T) {
	t.Parallel()

	// Nothing else is known about the pane, so dropping the directory would
	// leave every such tab reading Shell.
	got := hiddenDirectoryChain().Resolve(tabWithCWD(dashboard))
	if want := "dashboard"; got.Name != want {
		t.Errorf("name = %q, want %q", got.Name, want)
	}
}

func TestAHiddenDirectoryOutlivesAnActivityThatEchoesIt(t *testing.T) {
	t.Parallel()

	got := hiddenDirectoryChain().Resolve(tabWithPane(&state.PaneState{
		Dir:           dashboard,
		TerminalTitle: "dashboard",
	}))

	if want := "dashboard"; got.Name != want {
		t.Errorf("name = %q, want %q", got.Name, want)
	}
}

func TestAHiddenDirectoryKeepsTheBranch(t *testing.T) {
	t.Parallel()

	pane := repoPane("feat/oauth", "main")
	pane.TerminalTitle = "auth.ts - Nvim"
	pane.Processes = []state.Process{{Name: "nvim"}}

	got := hiddenDirectoryChain().Resolve(tabWithPane(pane))
	if want := "feat/oauth › nvim › auth.ts"; got.Name != want {
		t.Errorf("name = %q, want %q", got.Name, want)
	}
}

func TestAHiddenDirectoryKeepsTheRemoteHost(t *testing.T) {
	t.Parallel()

	// The host is the task rather than where the user stands, so it is not
	// the directory the setting hides.
	pane := sshPane("ssh", "prod-01")
	pane.Dir = dashboard
	pane.TerminalTitle = "htop"

	want := defaultChain().Resolve(tabWithPane(pane)).Name

	pane = sshPane("ssh", "prod-01")
	pane.Dir = dashboard
	pane.TerminalTitle = "htop"

	if got := hiddenDirectoryChain().Resolve(tabWithPane(pane)); got.Name != want {
		t.Errorf("name = %q, want %q", got.Name, want)
	}
}
