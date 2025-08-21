package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config represents the complete application configuration
type Config struct {
	Kafka   *KafkaConfig   `mapstructure:"kafka"`
	Service *ServiceConfig `mapstructure:"service"`
	App     *AppConfig     `mapstructure:"app"`
	Logging *LoggingConfig `mapstructure:"logging"`
	Health  *HealthConfig  `mapstructure:"health"`
	Secrets *SecretConfig  `json:"-"` // Don't serialize secrets
}

// KafkaConfig holds Kafka-related configuration
type KafkaConfig struct {
	Brokers         []string `mapstructure:"brokers"`
	GroupID         string   `mapstructure:"group_id"`
	Topic           string   `mapstructure:"topic"`
	AutoOffsetReset string   `mapstructure:"auto_offset_reset"`
}

// ServiceConfig holds service-specific configuration
type ServiceConfig struct {
	Port        int    `mapstructure:"port"`
	Environment string `mapstructure:"environment"`
	Host        string `mapstructure:"host"`
}

// AppConfig holds general application configuration
type AppConfig struct {
	Name           string `mapstructure:"name"`
	Version        string `mapstructure:"version"`
	Timeout        int    `mapstructure:"timeout"`
	MaxConnections int    `mapstructure:"max_connections"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// HealthConfig holds health check configuration
type HealthConfig struct {
	CheckInterval int `mapstructure:"check_interval"`
	Timeout       int `mapstructure:"timeout"`
}

// SecretConfig holds sensitive configuration data loaded from environment variables
type SecretConfig struct {
	// Kafka configuration
	KafkaBrokers          []string
	KafkaUsername         string
	KafkaPassword         string
	KafkaSASLMechanism    string
	KafkaSecurityProtocol string

	// Database credentials
	DatabaseURL      string
	DatabaseUsername string
	DatabasePassword string

	// API security
	APIKey    string
	JWTSecret string

	// MF API Integration
	MFAPIBaseURL string
	MFAPIKey     string
	MFAPISecret  string

	// External services
	ExternalAPIKey    string
	ExternalAPISecret string
}

// ConfigBuilder provides a fluent interface for building configuration
type ConfigBuilder struct {
	config *Config
	errors []error
}

// NewConfigBuilder creates a new configuration builder
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: &Config{
			Kafka:   &KafkaConfig{},
			Service: &ServiceConfig{},
			App:     &AppConfig{},
			Logging: &LoggingConfig{},
			Health:  &HealthConfig{},
			Secrets: &SecretConfig{},
		},
		errors: make([]error, 0),
	}
}

// LoadFromFile loads configuration from YAML file
func (cb *ConfigBuilder) LoadFromFile(filename string) *ConfigBuilder {
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		cb.errors = append(cb.errors, fmt.Errorf("failed to read config file: %w", err))
		return cb
	}

	if err := viper.Unmarshal(cb.config); err != nil {
		cb.errors = append(cb.errors, fmt.Errorf("failed to unmarshal config: %w", err))
		return cb
	}

	return cb
}

// LoadFromEnv loads sensitive configuration from environment variables
func (cb *ConfigBuilder) LoadFromEnv() *ConfigBuilder {
	// Load .env file if it exists (ignore errors)
	_ = godotenv.Load()

	// Enable automatic environment variable reading
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Load secrets
	secrets := cb.loadSecrets()
	cb.config.Secrets = secrets

	// Apply environment overrides
	cb.applyEnvironmentOverrides()

	return cb
}

// Build returns the final configuration or error if validation fails
func (cb *ConfigBuilder) Build() (*Config, error) {
	if len(cb.errors) > 0 {
		return nil, fmt.Errorf("configuration errors: %v", cb.errors)
	}

	// Validate required fields
	if err := cb.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cb.config, nil
}

// GetKafkaBrokers returns the effective Kafka brokers (from env or config)
func (c *Config) GetKafkaBrokers() []string {
	if len(c.Secrets.KafkaBrokers) > 0 {
		return c.Secrets.KafkaBrokers
	}
	return c.Kafka.Brokers
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return strings.ToLower(c.Service.Environment) == "production"
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	env := strings.ToLower(c.Service.Environment)
	return env == "development" || env == "dev"
}

// HasKafkaAuth returns true if Kafka authentication is configured
func (c *Config) HasKafkaAuth() bool {
	return c.Secrets.KafkaUsername != "" && c.Secrets.KafkaPassword != ""
}

// HasDatabaseConfig returns true if database configuration is present
func (c *Config) HasDatabaseConfig() bool {
	return c.Secrets.DatabaseURL != ""
}

// loadSecrets loads sensitive configuration from environment variables
func (cb *ConfigBuilder) loadSecrets() *SecretConfig {
	// Parse Kafka brokers from environment
	var kafkaBrokers []string
	if brokersList := os.Getenv("KAFKA_BROKERS"); brokersList != "" {
		brokers := strings.Split(brokersList, ",")
		for _, broker := range brokers {
			if trimmed := strings.TrimSpace(broker); trimmed != "" {
				kafkaBrokers = append(kafkaBrokers, trimmed)
			}
		}
	}

	return &SecretConfig{
		// Kafka configuration
		KafkaBrokers:          kafkaBrokers,
		KafkaUsername:         getEnvWithDefault("KAFKA_USERNAME", ""),
		KafkaPassword:         getEnvWithDefault("KAFKA_PASSWORD", ""),
		KafkaSASLMechanism:    getEnvWithDefault("KAFKA_SASL_MECHANISM", ""),
		KafkaSecurityProtocol: getEnvWithDefault("KAFKA_SECURITY_PROTOCOL", ""),

		// Database credentials
		DatabaseURL:      getEnvWithDefault("DATABASE_URL", ""),
		DatabaseUsername: getEnvWithDefault("DATABASE_USERNAME", ""),
		DatabasePassword: getEnvWithDefault("DATABASE_PASSWORD", ""),

		// API security
		APIKey:    getEnvWithDefault("API_KEY", ""),
		JWTSecret: getEnvWithDefault("JWT_SECRET", ""),

		// MF API Integration
		MFAPIBaseURL: getEnvWithDefault("MF_API_BASE_URL", ""),
		MFAPIKey:     getEnvWithDefault("MF_API_KEY", ""),
		MFAPISecret:  getEnvWithDefault("MF_API_SECRET", ""),

		// External services
		ExternalAPIKey:    getEnvWithDefault("EXTERNAL_API_KEY", ""),
		ExternalAPISecret: getEnvWithDefault("EXTERNAL_API_SECRET", ""),
	}
}

// applyEnvironmentOverrides applies environment variable overrides to configuration
func (cb *ConfigBuilder) applyEnvironmentOverrides() {
	// Override service port from environment if provided
	if portStr := os.Getenv("SERVICE_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cb.config.Service.Port = port
		}
	}

	// Override service environment
	if env := os.Getenv("SERVICE_ENVIRONMENT"); env != "" {
		cb.config.Service.Environment = env
	}

	// Override service host
	if host := os.Getenv("SERVICE_HOST"); host != "" {
		cb.config.Service.Host = host
	}
}

// validate performs configuration validation
func (cb *ConfigBuilder) validate() error {
	config := cb.config

	// Validate service configuration
	if config.Service.Port <= 0 || config.Service.Port > 65535 {
		return fmt.Errorf("invalid service port: %d", config.Service.Port)
	}

	if config.Service.Environment == "" {
		return fmt.Errorf("service environment is required")
	}

	// Validate Kafka configuration
	brokers := config.GetKafkaBrokers()
	if len(brokers) == 0 {
		return fmt.Errorf("at least one Kafka broker must be configured")
	}

	if config.Kafka.GroupID == "" {
		return fmt.Errorf("kafka group ID is required")
	}

	if config.Kafka.Topic == "" {
		return fmt.Errorf("kafka topic is required")
	}

	// Validate production-specific requirements
	if config.IsProduction() {
		if config.Secrets.JWTSecret == "" {
			return fmt.Errorf("JWT secret is required in production")
		}

		if len(config.Secrets.JWTSecret) < 32 {
			return fmt.Errorf("JWT secret must be at least 32 characters in production")
		}
	}

	return nil
}

// getEnvWithDefault gets an environment variable with a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// LoadConfig is a convenience function that creates and builds configuration using the fluent interface
func LoadConfig() (*Config, error) {
	return NewConfigBuilder().
		LoadFromFile("config").
		LoadFromEnv().
		Build()
}