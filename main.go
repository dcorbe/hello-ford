package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// These are our command-line flags
var humanFlag bool
var recursiveFlag bool

/* Convert size to human-readable format
 * Parameters:
 * 	- size: Size in bytes
 * Returns:
 * 	- string: Human-readable size string
 */
func humanReadableSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

/* Calculate the size of a directory or file
 * Parameters:
 *  - path: Path to the directory or file
 * Returns:
 *  - (int64, error): Size of the directory or file, or an error if one occured
 */
func dirSize(path string) (int64, error) {
	var totalSize int64
	walkFunc := func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// If it's a directory and recursion is not enabled, skip subdirectories
		if info.IsDir() && p != path && !recursiveFlag {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	}
	err := filepath.Walk(path, walkFunc)
	return totalSize, err
}

/* Print the size of each directory and calculate the cumulative size
 * Parameters:
 *	- dirs: List of directories to process
 */
func processDirectories(dirs []string) {
	var cumulativeSize int64

	for _, dir := range dirs {
		size, err := dirSize(dir)
		if err != nil {
			fmt.Printf("Error processing directory %s: %v\n", dir, err)
			continue
		}

		cumulativeSize += size
		if humanFlag {
			fmt.Printf("%s: %s\n", dir, humanReadableSize(size))
		} else {
			fmt.Printf("%s: %d bytes\n", dir, size)
		}
	}

	// Output cumulative size
	if humanFlag {
		fmt.Printf("Total: %s\n", humanReadableSize(cumulativeSize))
	} else {
		fmt.Printf("Total: %d bytes\n", cumulativeSize)
	}
}

func main() {
	// Parse command-line flags
	flag.BoolVar(&humanFlag, "human", false, "Display sizes in human-readable format (e.g., 1K, 234M, 2G)")
	flag.BoolVar(&recursiveFlag, "recursive", false, "Recursively calculate the sizes of directories and subdirectories")
	flag.Parse()

	// Remaining command-line arguments are the directories
	dirs := flag.Args()

	if len(dirs) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	processDirectories(dirs)
}
