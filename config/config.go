package config

import "os"

type Config struct {
	App      AppConfig
	Gemini   GeminiConfig
	Weaviate WeaviateConfig
}

type AppConfig struct {
	Port string
	IP   string
}

type GeminiConfig struct {
	APIKey         string
	EmbeddingModel string
	LLMModel       string
}

type WeaviateConfig struct {
	Host       string
	Port       string
	Collection string
	Scheme     string
	APIKey     string
}

func Load() *Config {
	return &Config{
		App: AppConfig{
			Port: getEnv("SERVER_PORT", "3000"),
			IP:   getEnv("SERVER_IP", "localhost"),
		},
		Gemini: GeminiConfig{
			APIKey:         getEnv("GEMINI_API_KEY", ""),
			EmbeddingModel: getEnv("GEMINI_EMBEDDING_MODEL", "gemini-embedding-2"),
			LLMModel:       getEnv("GEMINI_LLM_MODEL", "gemini-2.5-flash-lite"),
		},
		Weaviate: WeaviateConfig{
			Host:       getEnv("WEAVIATE_HOST", "localhost"),
			Port:       getEnv("WEAVIATE_PORT", "8080"),
			Scheme:     getEnv("WEAVIATE_SCHEME", "http"),
			Collection: getEnv("WEAVIATE_COLLECTION", "Documents"),
			APIKey:     getEnv("WEAVIATE_API_KEY", ""),
		},
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
