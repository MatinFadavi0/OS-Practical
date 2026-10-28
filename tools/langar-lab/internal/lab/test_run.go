package lab

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type caseResult struct {
	ID     string `json:"id"`
	Points int    `json:"points"`
	Earned int    `json:"earned"`
	Passed bool   `json:"passed"`
	Error  string `json:"error"`
}

type report struct {
	Assignment       string            `json:"assignment"`
	Status           string            `json:"status"`
	ExitCode         int               `json:"exit_code"`
	Score            *int              `json:"score"`
	Total            *int              `json:"total"`
	Results          []caseResult      `json:"results"`
	Error            string            `json:"error,omitempty"`
	Environment      map[string]string `json:"environment,omitempty"`
	Cleanup          string            `json:"cleanup"`
	Backend          string            `json:"backend"`
	Profile          string            `json:"profile"`
	CLIVersion       string            `json:"cli_version"`
	JudgeVersion     string            `json:"judge_version"`
	Source           string            `json:"source"`
	SourceHash       string            `json:"source_hash"`
	StartedAt        string            `json:"started_at"`
	FinishedAt       string            `json:"finished_at"`
	Workspace        string            `json:"workspace,omitempty"`
	WorkspaceCleanup string            `json:"workspace_cleanup"`
	BuildExitCode    *int              `json:"build_exit_code,omitempty"`
	JudgeExitCode    *int              `json:"judge_exit_code,omitempty"`
}

func validateReport(r report) error {
	expected := map[string]int{"passed": 0, "test_failed": 1, "environment_error": 2, "build_failed": 3, "interrupted": 130}
	code, ok := expected[r.Status]
	if !ok || code != r.ExitCode || CaseCounts[r.Assignment] == 0 {
		return fmt.Errorf("invalid worker status/assignment")
	}
	if r.Status == "passed" || r.Status == "test_failed" {
		if r.Score == nil || r.Total == nil || *r.Total != 100 || len(r.Results) != CaseCounts[r.Assignment] {
			return fmt.Errorf("missing or invalid functional scores")
		}
		total, earned, all := 0, 0, true
		seen := map[string]bool{}
		for _, c := range r.Results {
			if c.ID == "" || seen[c.ID] || c.Points <= 0 || (c.Earned != 0 && c.Earned != c.Points) || c.Passed != (c.Earned == c.Points) {
				return fmt.Errorf("invalid test case result")
			}
			seen[c.ID] = true
			total += c.Points
			earned += c.Earned
			all = all && c.Passed
		}
		if total != *r.Total || earned != *r.Score || (r.Status == "passed") != all {
			return fmt.Errorf("judge totals/status disagree")
		}
	} else if r.Score != nil || r.Total != nil {
		return fmt.Errorf("non-functional failure must not have a score")
	}
	return nil
}

func (a *App) test(ctx context.Context, o options) int {
	if o.Assignment == "" {
		o.Assignment = "hw2"
	}
	source, err := filepath.Abs(o.Source)
	if err == nil {
		_, err = sourceFiles(source)
	}
	if err != nil {
		fmt.Fprintln(a.Err, err)
		return 2
	}
	unlock, err := lock(o.State)
	if err != nil {
		fmt.Fprintln(a.Err, err)
		return 2
	}
	defer unlock()
	id := time.Now().UTC().Format("20060102T150405") + "-" + token()[:8]
	dir := filepath.Join(o.State, "results", o.Assignment, id)
	if err = os.MkdirAll(dir, 0700); err != nil {
		fmt.Fprintln(a.Err, err)
		return 2
	}
	r := report{Assignment: o.Assignment, Backend: o.Backend, Status: "environment_error", ExitCode: 2,
		CLIVersion: Version, JudgeVersion: judgeVersion(o.Assignment), Profile: Profile, Source: source,
		StartedAt: time.Now().UTC().Format(time.RFC3339), Cleanup: "not_started",
		WorkspaceCleanup: "not_started"}
	if o.Backend == "native" {
		r.Profile = NativeProfile
	}
	var cfg vmConfig
	var work string
	workCreated := false
	finish := func(failure error) int {
		if failure != nil {
			r.Status, r.ExitCode, r.Score, r.Total, r.Error = "environment_error", 2, nil, nil, failure.Error()
		}
		if ctx.Err() != nil {
			r.Status, r.ExitCode, r.Score, r.Total, r.Error = "interrupted", 130, nil, nil, "Run cancelled"
		}
		if workCreated && r.WorkspaceCleanup == "pending" && r.Cleanup != "unconfirmed" {
			if o.KeepWork {
				r.WorkspaceCleanup = "retained_for_debug"
			} else {
				var cleanupErr error
				if o.Backend == "native" {
					cleanupErr = os.RemoveAll(work)
				} else {
					cleanCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
					_, cleanupErr = a.capture(cleanCtx, "multipass", "exec", cfg.Name, "--", "rm", "-rf", "--", strings.TrimSuffix(work, "/work"))
					cancel()
				}
				if cleanupErr == nil {
					r.Workspace, r.WorkspaceCleanup = "", "complete"
				} else {
					r.WorkspaceCleanup = "failed"
					r.Error += "; workspace cleanup failed: " + cleanupErr.Error()
				}
			}
		}
		r.FinishedAt = time.Now().UTC().Format(time.RFC3339)
		paths := []string{filepath.Join(dir, "report.json")}
		if o.JSON != "" {
			paths = append(paths, o.JSON)
		}
		for _, path := range paths {
			if err := writeJSON(path, r); err != nil {
				fmt.Fprintln(a.Err, "Could not save report:", err)
				return 2
			}
		}
		if r.Error != "" {
			fmt.Fprintln(a.Err, r.Error)
		}
		if r.Score != nil {
			fmt.Fprintf(a.Out, "[%s] Score: %d/%d\n", r.Status, *r.Score, *r.Total)
		} else {
			fmt.Fprintf(a.Out, "[%s] No functional score\n", r.Status)
		}
		fmt.Fprintln(a.Out, "Results:", dir)
		return r.ExitCode
	}
	archive := filepath.Join(dir, "source.tar.gz")
	r.SourceHash, err = snapshot(source, archive, o.Assignment)
	if err != nil {
		return finish(err)
	}
	fmt.Fprintf(a.Out, "==> %s source snapshot: %s\n==> Backend: %s\n", strings.ToUpper(o.Assignment), source, o.Backend)
	if o.Backend == "native" {
		if err = a.doctor(ctx, o); err != nil {
			return finish(err)
		}
		work = filepath.Join(dir, "work")
		workCreated, r.Workspace, r.WorkspaceCleanup = true, work, "pending"
		if err = unpackSnapshot(archive, work); err != nil {
			return finish(err)
		}
	} else {
		setupOptions := o
		setupOptions.Timeout = 15 * time.Minute
		if cfg, err = a.ensureVM(ctx, setupOptions); err != nil {
			return finish(err)
		}
		work = "/home/ubuntu/langar-lab-runs/" + id + "/work"
		if _, err = a.capture(ctx, "multipass", "exec", cfg.Name, "--", "mkdir", "-p", work); err != nil {
			return finish(err)
		}
		workCreated, r.Workspace, r.WorkspaceCleanup = true, work, "pending"
		remoteArchive := strings.TrimSuffix(work, "/work") + "/source.tar.gz"
		if _, err = a.capture(ctx, "multipass", "transfer", archive, cfg.Name+":"+remoteArchive); err != nil {
			return finish(err)
		}
		if _, err = a.capture(ctx, "multipass", "exec", cfg.Name, "--", "tar", "-xzf", remoteArchive, "-C", work); err != nil {
			return finish(err)
		}
	}
	r.Workspace, r.WorkspaceCleanup = work, "pending"
	worker := work + "/course/agent.py"
	if o.Backend == "native" {
		worker = filepath.Join(work, "course", "agent.py")
	}
	workerArgs := []string{worker, "--workdir", work, "--assignment", o.Assignment, "--cc", o.CC, "--timeout", fmt.Sprint(o.Timeout.Seconds())}
	command, args := "python3", workerArgs
	if o.Backend == "multipass" {
		command = "multipass"
		args = append([]string{"exec", cfg.Name, "--", "python3"}, workerArgs...)
	}
	runCtx, cancelRun := context.WithTimeout(ctx, o.Timeout+15*time.Second)
	code, runErr := a.Exec.Run(runCtx, command, args, nil, a.Out, a.Err)
	cancelled := runCtx.Err() != nil
	cancelRun()
	if cancelled || runErr != nil {
		cancelCtx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
		cancelArgs := []string{worker, "--workdir", work, "--cancel"}
		if o.Backend == "multipass" {
			_, err = a.capture(cancelCtx, "multipass", append([]string{"exec", cfg.Name, "--", "python3"}, cancelArgs...)...)
		} else {
			_, err = a.capture(cancelCtx, "python3", cancelArgs...)
		}
		cancel()
		if err != nil {
			r.Cleanup = "unconfirmed"
			return finish(fmt.Errorf("worker interrupted; could not confirm cleanup: %w", err))
		}
	}
	collectCtx, cancelCollect := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelCollect()
	for _, name := range []string{"report.json", "build.log", "test.log", "judge.json"} {
		local := filepath.Join(dir, name)
		if o.Backend == "native" {
			data, readErr := os.ReadFile(filepath.Join(work, "artifacts", name))
			if readErr == nil {
				err = os.WriteFile(local, data, 0600)
			} else {
				err = readErr
			}
		} else {
			_, err = a.capture(collectCtx, "multipass", "transfer", cfg.Name+":"+work+"/artifacts/"+name, local)
		}
		if err != nil && name == "report.json" {
			return finish(fmt.Errorf("worker returned %d without a report: %v (%v)", code, err, runErr))
		}
	}
	data, err := os.ReadFile(filepath.Join(dir, "report.json"))
	if err != nil {
		return finish(err)
	}
	var workerReport report
	if err = json.Unmarshal(data, &workerReport); err != nil {
		return finish(err)
	}
	if workerReport.Assignment != o.Assignment {
		return finish(fmt.Errorf("worker assignment mismatch"))
	}
	if err = validateReport(workerReport); err != nil {
		return finish(err)
	}
	if !cancelled && runErr == nil && code != workerReport.ExitCode {
		return finish(fmt.Errorf("worker exit %d disagrees with reported %d", code, workerReport.ExitCode))
	}
	workerReport.Backend, workerReport.CLIVersion, workerReport.Profile = o.Backend, Version, r.Profile
	workerReport.JudgeVersion = judgeVersion(o.Assignment)
	workerReport.Source, workerReport.SourceHash, workerReport.StartedAt = source, r.SourceHash, r.StartedAt
	workerReport.Workspace = work
	r = workerReport
	if o.KeepWork {
		r.WorkspaceCleanup = "retained_for_debug"
		fmt.Fprintln(a.Out, "Linux workspace retained:", work)
	} else {
		if o.Backend == "native" {
			err = os.RemoveAll(work)
		} else {
			// work and its parent are generated IDs under our fixed guest root.
			_, err = a.capture(collectCtx, "multipass", "exec", cfg.Name, "--", "rm", "-rf", "--", strings.TrimSuffix(work, "/work"))
		}
		if err != nil {
			r.WorkspaceCleanup = "failed"
			return finish(fmt.Errorf("test completed, but workspace cleanup failed: %w", err))
		}
		r.Workspace, r.WorkspaceCleanup = "", "complete"
	}
	return finish(nil)
}
