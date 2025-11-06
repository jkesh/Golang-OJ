package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	Kafka    KafkaConfig    `json:"kafka"`
	Judge    JudgeConfig    `json:"judge"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int    `json:"port"`
	Host         string `json:"host"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// KafkaConfig Kafka配置
type KafkaConfig struct {
	Brokers     []string `json:"brokers"`
	Topic       string   `json:"topic"`
	ResultTopic string   `json:"result_topic"`
}

// JudgeConfig 判题配置
type JudgeConfig struct {
	Timeout      int      `json:"timeout"`
	MaxMemory    int      `json:"max_memory"`
	AllowedLangs []string `json:"allowed_langs"`
}

var (
	config *Config
	once   sync.Once
)

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	once.Do(func() {
		config = &Config{}

		// 如果配置文件不存在，使用默认配置
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			config = getDefaultConfig()
			return
		}

		file, err := os.Open(configPath)
		if err != nil {
			fmt.Printf("无法打开配置文件: %v\n", err)
			config = getDefaultConfig()
			return
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(config); err != nil {
			fmt.Printf("解析配置文件失败: %v\n", err)
			config = getDefaultConfig()
			return
		}
	})

	return config, nil
}

// GetConfig 获取配置实例
func GetConfig() *Config {
	if config == nil {
		_, _ = LoadConfig("./config.json")
	}
	return config
}

// 默认配置
func getDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         8080,
			Host:         "0.0.0.0",
			ReadTimeout:  60,
			WriteTimeout: 60,
		},
		Database: DatabaseConfig{
			Host:     "43.131.41.101",
			Port:     5432,
			User:     "postgres",
			Password: "jkesh1024",
			DBName:   "ojplus",
			SSLMode:  "disable",
		},
		Redis: RedisConfig{
			Host:     "111.228.56.17",
			Port:     6379,
			Password: "jkesh1024",
			DB:       0,
		},
		Kafka: KafkaConfig{
			Brokers:     []string{"localhost:9092"},
			Topic:       "judge-tasks",
			ResultTopic: "judge-results",
		},
		Judge: JudgeConfig{
			Timeout:      10000,
			MaxMemory:    256,
			AllowedLangs: []string{"go", "cpp", "java", "python"},
		},
	}
}
