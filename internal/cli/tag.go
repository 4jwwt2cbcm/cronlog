package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"cronlog/internal/logentry"
	"cronlog/internal/storage"
)

// TagSummary holds aggregated stats for a single tag.
type TagSummary struct {
	Tag     string
	Total   int
	Errors  int
	Jobs    []string
}

// ComputeTagSummary groups log entries by tag and returns a summary per tag.
// Tags are derived from the entry's Tags field (comma-separated in Label).
func ComputeTagSummary(entries []logentry.Entry) []TagSummary {
	type bucket struct {
		total  int
		errors int
		jobs   map[string]struct{}
	}

	buckets := map[string]*bucket{}

	for _, e := range entries {
		for _, tag := range parseTags(e.Tags) {
			b, ok := buckets[tag]
			if !ok {
				b = &bucket{jobs: map[string]struct{}{}}
				buckets[tag] = b
			}
			b.total++
			if e.IsError() {
				b.errors++
			}
			b.jobs[e.Job] = struct{}{}
		}
	}

	summaries := make([]TagSummary, 0, len(buckets))
	for tag, b := range buckets {
		jobs := make([]string, 0, len(b.jobs))
		for j := range b.jobs {
			jobs = append(jobs, j)
		}
		sort.Strings(jobs)
		summaries = append(summaries, TagSummary{
			Tag:    tag,
			Total:  b.total,
			Errors: b.errors,
			Jobs:   jobs,
		})
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Tag < summaries[j].Tag
	})
	return summaries
}

func parseTags(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func buildTagCmd(dbPath string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "Show log entry counts grouped by tag",
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := storage.New(dbPath)
			if err != nil {
				return err
			}
			entries, err := st.All()
			if err != nil {
				return err
			}
			summaries := ComputeTagSummary(entries)
			if len(summaries) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no tagged entries found")
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%-20s %8s %8s  %s\n", "TAG", "TOTAL", "ERRORS", "JOBS")
			for _, s := range summaries {
				fmt.Fprintf(cmd.OutOrStdout(), "%-20s %8d %8d  %s\n",
					s.Tag, s.Total, s.Errors, strings.Join(s.Jobs, ", "))
			}
			return nil
		},
	}
	return cmd
}
