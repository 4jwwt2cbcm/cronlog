package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"cronlog/internal/storage"
)

// buildAnnotateCmd returns a cobra command that adds or replaces tags on a
// stored log entry identified by its ID.
func buildAnnotateCmd(store *storage.Store) *cobra.Command {
	var tagsFlag string

	cmd := &cobra.Command{
		Use:   "annotate <id>",
		Short: "Add or replace tags on a log entry",
		Long:  annotateDoc,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			if id == "" {
				return errors.New("entry id must not be empty")
			}

			newTags := parseFlagTags(tagsFlag)

			all, err := store.All()
			if err != nil {
				return fmt.Errorf("loading entries: %w", err)
			}

			for i, e := range all {
				if e.ID == id {
					all[i].Tags = mergeTags(e.Tags, newTags)
					if err := store.Replace(all); err != nil {
						return fmt.Errorf("saving entries: %w", err)
					}
					fmt.Fprintf(cmd.OutOrStdout(), "annotated entry %s\n", id)
					return nil
				}
			}

			return fmt.Errorf("entry %s not found", id)
		},
	}

	cmd.Flags().StringVar(&tagsFlag, "tags", "", "comma-separated key=value tags to add (e.g. env=prod,team=ops)")
	_ = cmd.MarkFlagRequired("tags")
	return cmd
}

// parseFlagTags splits a comma-separated "key=value" string into a map.
func parseFlagTags(raw string) map[string]string {
	out := make(map[string]string)
	for _, pair := range strings.Split(raw, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			out[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return out
}

// mergeTags overlays newTags onto existing tags, returning the combined map.
func mergeTags(existing, newTags map[string]string) map[string]string {
	result := make(map[string]string, len(existing)+len(newTags))
	for k, v := range existing {
		result[k] = v
	}
	for k, v := range newTags {
		result[k] = v
	}
	return result
}
