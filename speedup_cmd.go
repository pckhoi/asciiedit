package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func speedupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "speedup FILE TIMES",
		Long: "Speedup the section ",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			file := args[0]
			times, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return err
			}
			r, err := cmd.Flags().GetStringSlice("range")
			if err != nil {
				return err
			}
			if len(r) == 0 {
				return fmt.Errorf("must specify at least 1 range")
			}
			ranges, err := processRanges(r)
			if err != nil {
				return err
			}
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			line := 1
			out := cmd.OutOrStdout()
			var prevTime float64
			var accumDur float64
			var ri int
			var firstFrom = ranges[0][0]
			for scanner.Scan() {
				if line == 1 {
					fmt.Fprintln(out, scanner.Text())
				} else {
					b := scanner.Bytes()
					arr := []interface{}{}
					if err := json.Unmarshal(b, &arr); err != nil {
						return err
					}
					curTime := arr[0].(float64)
					if line > firstFrom {
						for ri < len(ranges) && line > ranges[ri][1] {
							ri += 1
						}
						if ri < len(ranges) && line > ranges[ri][0] && line <= ranges[ri][1]+1 {
							duration := curTime - prevTime
							arr[0] = (prevTime - accumDur) + duration/times
							accumDur += duration - duration/times
						} else {
							arr[0] = curTime - accumDur
						}
						b, err = json.Marshal(arr)
						if err != nil {
							return err
						}
					}
					prevTime = curTime
					fmt.Fprintf(out, "%s\n", string(b))
				}
				line += 1
			}
			return scanner.Err()
		},
	}
	cmd.Flags().StringSliceP("range", "r", nil, "specify line range FROM:TO (inclusive) to speed up. Line 1 is the line that starts with '{\"version\":'")
	return cmd
}

func processRanges(r []string) ([][2]int, error) {
	sl := make([][2]int, 0, len(r))
	for _, s := range r {
		p := strings.Split(s, ":")
		from, err := strconv.Atoi(p[0])
		if err != nil {
			return nil, err
		}
		to, err := strconv.Atoi(p[1])
		if err != nil {
			return nil, err
		}
		sl = append(sl, [2]int{from, to})
	}
	sort.Slice(sl, func(i, j int) bool {
		return sl[i][0] < sl[j][0]
	})
	return sl, nil
}
