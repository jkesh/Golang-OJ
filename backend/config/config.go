package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	Judge    JudgeConfig    `mapstructure:"judge"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	Host         string `mapstructure:"host"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// KafkaConfig Kafka配置
type KafkaConfig struct {
	Brokers     []string `mapstructure:"brokers"`
	Topic       string   `mapstructure:"topic"`
	ResultTopic string   `mapstructure:"result_topic"`
}

// JudgeConfig 判题配置
type JudgeConfig struct {
	Timeout      int      `mapstructure:"timeout"`
	MaxMemory    int      `mapstructure:"max_memory"`
	AllowedLangs []string `mapstructure:"allowed_langs"`
}

var (
	config *Config
	once   sync.Once
)

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	once.Do(func() {
		viper.SetConfigFile(configPath)
		viper.SetConfigType("json")

		// 设置默认值
		setDefaults()

		// 如果配置文件存在则读取
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("警告: 无法读取配置文件: %v\n", err)
			fmt.Println("将使用默认配置")
		}

		// 将配置反序列化到结构体
		config = &Config{}
		if err := viper.Unmarshal(config); err != nil {
			fmt.Printf("警告: 无法解析配置: %v\n", err)
			config = getDefaultConfig()
		}
	})

	return config, nil
}

// GetConfig 获取配置实例
func GetConfig() *Config {
	if config == nil {
		_, _ = LoadConfig("./config/config.json")
	}
	return config
}

// 设置默认配置值
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.read_timeout", 60)
	viper.SetDefault("server.write_timeout", 60)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.dbname", "ojplus")
	viper.SetDefault("database.sslmode", "disable")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// Kafka defaults
	viper.SetDefault("kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("kafka.topic", "judge-tasks")
	viper.SetDefault("kafka.result_topic", "judge-results")

	// Judge defaults
	viper.SetDefault("judge.timeout", 10000)
	viper.SetDefault("judge.max_memory", 256)
	viper.SetDefault("judge.allowed_langs", []string{"go", "cpp", "java", "python"})
}

// 默认配置
func getDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         viper.GetInt("server.port"),
			Host:         viper.GetString("server.host"),
			ReadTimeout:  viper.GetInt("server.read_timeout"),
			WriteTimeout: viper.GetInt("server.write_timeout"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetInt("database.port"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			DBName:   viper.GetString("database.dbname"),
			SSLMode:  viper.GetString("database.sslmode"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("redis.host"),
			Port:     viper.GetInt("redis.port"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		},
		Kafka: KafkaConfig{
			Brokers:     viper.GetStringSlice("kafka.brokers"),
			Topic:       viper.GetString("kafka.topic"),
			ResultTopic: viper.GetString("kafka.result_topic"),
		},
		Judge: JudgeConfig{
			Timeout:      viper.GetInt("judge.timeout"),
			MaxMemory:    viper.GetInt("judge.max_memory"),
			AllowedLangs: viper.GetStringSlice("judge.allowed_langs"),
		},
	}
}
