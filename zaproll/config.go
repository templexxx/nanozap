/*
 * Copyright (c) 2019. Temple3x (temple3x@gmail.com)
 * Copyright (c) 2014 Nate Finch
 *
 * Use of this source code is governed by the MIT License
 * that can be found in the LICENSE file.
 */

package zaproll

// Config of zaproll.
type Config struct {
	// OutputPath is the log file path.
	OutputPath string `json:"output_path" toml:"output_path"`
	// MaxSize is the maximum size of a log file before it gets rotated.
	// Unit: MB.
	// Default: 128 (128MB).
	MaxSize int64 `json:"max_size_mb" toml:"max_size_mb"`
	// MaxBackups is the maximum number of backup log files to retain.
	MaxBackups int `json:"max_backups" toml:"max_backups"`
	// LocalTime is the timestamp in backup log file. Default is to use UTC time.
	// If true, use local time.
	LocalTime bool `json:"local_time" toml:"local_time"`

	// PerWriteSize is zaproll's write size,
	// zaproll writes data to page cache every PerWriteSize.
	// Unit: KB.
	// Default: 64 (64KB).
	//
	// It's used for combining writes.
	// The size of it should be aligned to page size,
	// and it shouldn't be too large, because that may block zaproll write.
	PerWriteSize int64 `json:"per_write_size" toml:"per_write_size"`
	// PerSyncSize is zaproll's sync size,
	// zaproll flushes data to storage media(hint) every PerSyncSize.
	// Unit: MB.
	// Default: 16 (16MB).
	//
	// The size of it should be aligned to page size,
	// and it shouldn't be too large, avoiding burst I/O.
	PerSyncSize int64 `json:"per_sync_size" toml:"per_sync_size"`

	// Develop mode. Default is false.
	// It' used for testing, if it's true, the page cache control unit could not be aligned to page cache size.
	Developed bool `json:"developed" toml:"developed"`
}

const (
	kb int64 = 1024
	mb       = 1024 * kb
)

// Default configs.
var (
	defaultPerWriteSize = 64 * kb
	defaultPerSyncSize  = 16 * mb

	// We don't need to keep too many backups,
	// in practice, log shipper will collect the logs.
	defaultMaxSize    = 128 * mb
	defaultMaxBackups = 4
)

func (c *Config) adjust() {

	k, m := kb, mb
	if c.Developed {
		k, m = 1, 1
	}

	if c.MaxSize <= 0 {
		c.MaxSize = defaultMaxSize
	} else {
		c.MaxSize = c.MaxSize * m
	}
	if c.MaxBackups <= 0 {
		c.MaxBackups = defaultMaxBackups
	}

	if c.PerWriteSize <= 0 {
		c.PerWriteSize = defaultPerWriteSize
	} else {
		c.PerWriteSize = c.PerWriteSize * k
	}
	if c.PerSyncSize <= 0 {
		c.PerSyncSize = defaultPerSyncSize
	} else {
		c.PerSyncSize = c.PerSyncSize * m
	}

	if !c.Developed {
		if c.PerSyncSize < 2*c.PerWriteSize {
			c.PerSyncSize = 2 * c.PerWriteSize
		}
		c.MaxSize = alignToPage(c.MaxSize)
		c.PerWriteSize = alignToPage(c.PerWriteSize)
		c.PerSyncSize = alignToPage(c.PerSyncSize)
	}
}

const pageSize = 1 << 12 // 4KB.

func alignToPage(n int64) int64 {
	return (n + pageSize - 1) &^ (pageSize - 1)
}
