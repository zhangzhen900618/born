package logger

func getDefaultOptions() options {
	return options{
		level:         "debug",
		fileName:      "./logs/game.log",
		maxSize:       200,
		maxAge:        30,
		maxBackups:    20,
		enableConsole: true,
		enableWrite:   false,
	}
}

type (
	Option func(o *options)

	options struct {
		level         string // 日志级别
		fileName      string // 日志文件的位置
		maxSize       int    // 在进行切割之前，日志文件的最大大小（以 MB 为单位）
		maxAge        int    // 保留旧文件的最大天数
		maxBackups    int    // 保留旧文件的最大个数
		enableConsole bool   // 输出到控制台
		enableWrite   bool   // 输出到文件内
	}
)

func WithLevel(level string) Option {
	return func(o *options) {
		o.level = level
	}
}

func WithFileName(fileName string) Option {
	return func(o *options) {
		o.fileName = fileName
	}
}

func WithMaxSize(size int) Option {
	return func(o *options) {
		o.maxSize = size
	}
}

func WithMaxAge(age int) Option {
	return func(o *options) {
		o.maxAge = age
	}
}

func WithMaxBackups(backups int) Option {
	return func(o *options) {
		o.maxBackups = backups
	}
}

func WithEnableConsole(enable bool) Option {
	return func(o *options) {
		o.enableConsole = enable
	}
}

func WithEnableWrite(enable bool) Option {
	return func(o *options) {
		o.enableWrite = enable
	}
}
