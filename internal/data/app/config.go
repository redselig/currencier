package app

type Config struct {
	Log   Log   `yaml:"log"`
	DB    DB    `yaml:"db"`
	Update Update `yaml:"update"`
	API API `yaml:"api"`

}

type Log struct {
	File string `yaml:"file"`
}

type API struct {
	HTTPPort string `yaml:"httpport"`
}


type DB struct {
	DSN     string `yaml:"dsn"`
	Dialect string `yaml:"dialect"`
}

type Update struct {
	time string `yaml:"time"`
	source string `yaml:"source"`
}
