package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Layr-Labs/hourglass-avs-template/internal/types"

	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type WeatherDataSource interface {
	GetTemperature(ctx context.Context, location types.Location) (*types.WeatherResponse, error)
	GetName() string
}

type BaseWeatherSource struct {
	Name        string
	BaseURL     string
	APIKey      string
	RateLimiter *rate.Limiter
	Client      *http.Client
	Cache       *Cache
}

type Cache struct {
	mu      sync.RWMutex
	entries map[string]types.CacheEntry
	ttl     time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	cache := &Cache{
		entries: make(map[string]types.CacheEntry),
		ttl:     ttl,
	}
	
	go cache.cleanup()
	
	return cache
}

func (c *Cache) Get(key string) (*types.WeatherResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	entry, ok := c.entries[key]
	if !ok || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	
	return &entry.Data, true
}

func (c *Cache) Set(key string, data types.WeatherResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.entries[key] = types.CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

func (c *Cache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.entries {
			if now.After(entry.ExpiresAt) {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

func (b *BaseWeatherSource) GetName() string {
	return b.Name
}

func (b *BaseWeatherSource) makeRequest(ctx context.Context, url string) (*http.Response, error) {
	if err := b.RateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, err
	}
	
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	
	return resp, nil
}

type OpenWeatherMapSource struct {
	BaseWeatherSource
}

func NewOpenWeatherMapSource(apiKey string, rateLimit int, cache *Cache) *OpenWeatherMapSource {
	return &OpenWeatherMapSource{
		BaseWeatherSource: BaseWeatherSource{
			Name:        "OpenWeatherMap",
			BaseURL:     "https://api.openweathermap.org/data/2.5",
			APIKey:      apiKey,
			RateLimiter: rate.NewLimiter(rate.Every(time.Minute/time.Duration(rateLimit)), 1),
			Client:      &http.Client{Timeout: 30 * time.Second},
			Cache:       cache,
		},
	}
}

func (o *OpenWeatherMapSource) GetTemperature(ctx context.Context, location types.Location) (*types.WeatherResponse, error) {
	cacheKey := fmt.Sprintf("%s:%f:%f", o.Name, location.Latitude, location.Longitude)
	if cached, ok := o.Cache.Get(cacheKey); ok {
		log.Debugf("Cache hit for %s", o.Name)
		return cached, nil
	}
	
	url := fmt.Sprintf("%s/weather?lat=%f&lon=%f&appid=%s&units=metric",
		o.BaseURL, location.Latitude, location.Longitude, o.APIKey)
	
	resp, err := o.makeRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var data struct {
		Main struct {
			Temp     float64 `json:"temp"`
			Humidity float64 `json:"humidity"`
			Pressure float64 `json:"pressure"`
		} `json:"main"`
		Wind struct {
			Speed float64 `json:"speed"`
		} `json:"wind"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	
	result := &types.WeatherResponse{
		Temperature: data.Main.Temp,
		Humidity:    data.Main.Humidity,
		Pressure:    data.Main.Pressure,
		WindSpeed:   data.Wind.Speed,
		Timestamp:   time.Now(),
		Source:      o.Name,
	}
	
	o.Cache.Set(cacheKey, *result)
	return result, nil
}

type WeatherAPISource struct {
	BaseWeatherSource
}

func NewWeatherAPISource(apiKey string, rateLimit int, cache *Cache) *WeatherAPISource {
	return &WeatherAPISource{
		BaseWeatherSource: BaseWeatherSource{
			Name:        "WeatherAPI",
			BaseURL:     "https://api.weatherapi.com/v1",
			APIKey:      apiKey,
			RateLimiter: rate.NewLimiter(rate.Every(time.Minute/time.Duration(rateLimit)), 1),
			Client:      &http.Client{Timeout: 30 * time.Second},
			Cache:       cache,
		},
	}
}

func (w *WeatherAPISource) GetTemperature(ctx context.Context, location types.Location) (*types.WeatherResponse, error) {
	cacheKey := fmt.Sprintf("%s:%f:%f", w.Name, location.Latitude, location.Longitude)
	if cached, ok := w.Cache.Get(cacheKey); ok {
		log.Debugf("Cache hit for %s", w.Name)
		return cached, nil
	}
	
	url := fmt.Sprintf("%s/current.json?key=%s&q=%f,%f",
		w.BaseURL, w.APIKey, location.Latitude, location.Longitude)
	
	resp, err := w.makeRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var data struct {
		Current struct {
			TempC      float64 `json:"temp_c"`
			Humidity   float64 `json:"humidity"`
			PressureMb float64 `json:"pressure_mb"`
			WindKph    float64 `json:"wind_kph"`
		} `json:"current"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	
	result := &types.WeatherResponse{
		Temperature: data.Current.TempC,
		Humidity:    data.Current.Humidity,
		Pressure:    data.Current.PressureMb,
		WindSpeed:   data.Current.WindKph / 3.6,
		Timestamp:   time.Now(),
		Source:      w.Name,
	}
	
	w.Cache.Set(cacheKey, *result)
	return result, nil
}

type TomorrowIOSource struct {
	BaseWeatherSource
}

func NewTomorrowIOSource(apiKey string, rateLimit int, cache *Cache) *TomorrowIOSource {
	return &TomorrowIOSource{
		BaseWeatherSource: BaseWeatherSource{
			Name:        "TomorrowIO",
			BaseURL:     "https://api.tomorrow.io/v4",
			APIKey:      apiKey,
			RateLimiter: rate.NewLimiter(rate.Every(time.Minute/time.Duration(rateLimit)), 1),
			Client:      &http.Client{Timeout: 30 * time.Second},
			Cache:       cache,
		},
	}
}

func (t *TomorrowIOSource) GetTemperature(ctx context.Context, location types.Location) (*types.WeatherResponse, error) {
	cacheKey := fmt.Sprintf("%s:%f:%f", t.Name, location.Latitude, location.Longitude)
	if cached, ok := t.Cache.Get(cacheKey); ok {
		log.Debugf("Cache hit for %s", t.Name)
		return cached, nil
	}
	
	url := fmt.Sprintf("%s/weather/realtime?location=%f,%f&apikey=%s&units=metric",
		t.BaseURL, location.Latitude, location.Longitude, t.APIKey)
	
	resp, err := t.makeRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var data struct {
		Data struct {
			Values struct {
				Temperature float64 `json:"temperature"`
				Humidity    float64 `json:"humidity"`
				Pressure    float64 `json:"pressureSurfaceLevel"`
				WindSpeed   float64 `json:"windSpeed"`
			} `json:"values"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	
	result := &types.WeatherResponse{
		Temperature: data.Data.Values.Temperature,
		Humidity:    data.Data.Values.Humidity,
		Pressure:    data.Data.Values.Pressure,
		WindSpeed:   data.Data.Values.WindSpeed,
		Timestamp:   time.Now(),
		Source:      t.Name,
	}
	
	t.Cache.Set(cacheKey, *result)
	return result, nil
}

type VisualCrossingSource struct {
	BaseWeatherSource
}

func NewVisualCrossingSource(apiKey string, rateLimit int, cache *Cache) *VisualCrossingSource {
	return &VisualCrossingSource{
		BaseWeatherSource: BaseWeatherSource{
			Name:        "VisualCrossing",
			BaseURL:     "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline",
			APIKey:      apiKey,
			RateLimiter: rate.NewLimiter(rate.Every(time.Minute/time.Duration(rateLimit)), 1),
			Client:      &http.Client{Timeout: 30 * time.Second},
			Cache:       cache,
		},
	}
}

func (v *VisualCrossingSource) GetTemperature(ctx context.Context, location types.Location) (*types.WeatherResponse, error) {
	cacheKey := fmt.Sprintf("%s:%f:%f", v.Name, location.Latitude, location.Longitude)
	if cached, ok := v.Cache.Get(cacheKey); ok {
		log.Debugf("Cache hit for %s", v.Name)
		return cached, nil
	}
	
	url := fmt.Sprintf("%s/%f,%f/today?key=%s&unitGroup=metric&include=current",
		v.BaseURL, location.Latitude, location.Longitude, v.APIKey)
	
	resp, err := v.makeRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var data struct {
		CurrentConditions struct {
			Temp      float64 `json:"temp"`
			Humidity  float64 `json:"humidity"`
			Pressure  float64 `json:"pressure"`
			WindSpeed float64 `json:"windspeed"`
		} `json:"currentConditions"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	
	result := &types.WeatherResponse{
		Temperature: data.CurrentConditions.Temp,
		Humidity:    data.CurrentConditions.Humidity,
		Pressure:    data.CurrentConditions.Pressure,
		WindSpeed:   data.CurrentConditions.WindSpeed * 0.277778,
		Timestamp:   time.Now(),
		Source:      v.Name,
	}
	
	v.Cache.Set(cacheKey, *result)
	return result, nil
}

type OpenMeteoSource struct {
	BaseWeatherSource
}

func NewOpenMeteoSource(rateLimit int, cache *Cache) *OpenMeteoSource {
	return &OpenMeteoSource{
		BaseWeatherSource: BaseWeatherSource{
			Name:        "OpenMeteo",
			BaseURL:     "https://api.open-meteo.com/v1",
			RateLimiter: rate.NewLimiter(rate.Every(time.Minute/time.Duration(rateLimit)), 1),
			Client:      &http.Client{Timeout: 30 * time.Second},
			Cache:       cache,
		},
	}
}

func (o *OpenMeteoSource) GetTemperature(ctx context.Context, location types.Location) (*types.WeatherResponse, error) {
	cacheKey := fmt.Sprintf("%s:%f:%f", o.Name, location.Latitude, location.Longitude)
	if cached, ok := o.Cache.Get(cacheKey); ok {
		log.Debugf("Cache hit for %s", o.Name)
		return cached, nil
	}
	
	url := fmt.Sprintf("%s/forecast?latitude=%f&longitude=%f&current=temperature_2m,relative_humidity_2m,surface_pressure,wind_speed_10m",
		o.BaseURL, location.Latitude, location.Longitude)
	
	resp, err := o.makeRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var data struct {
		Current struct {
			Temperature      float64 `json:"temperature_2m"`
			Humidity         float64 `json:"relative_humidity_2m"`
			Pressure         float64 `json:"surface_pressure"`
			WindSpeed        float64 `json:"wind_speed_10m"`
		} `json:"current"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	
	result := &types.WeatherResponse{
		Temperature: data.Current.Temperature,
		Humidity:    data.Current.Humidity,
		Pressure:    data.Current.Pressure,
		WindSpeed:   data.Current.WindSpeed / 3.6,
		Timestamp:   time.Now(),
		Source:      o.Name,
	}
	
	o.Cache.Set(cacheKey, *result)
	return result, nil
}

type DataSourceManager struct {
	sources map[string]WeatherDataSource
	mu      sync.RWMutex
}

func NewDataSourceManager(config map[string]struct {
	BaseURL   string `yaml:"base_url"`
	RateLimit int    `yaml:"rate_limit"`
	APIKey    string `yaml:"api_key,omitempty"`
}, cacheTTL time.Duration) *DataSourceManager {
	cache := NewCache(cacheTTL)
	manager := &DataSourceManager{
		sources: make(map[string]WeatherDataSource),
	}
	
	for name, cfg := range config {
		switch name {
		case "openweathermap":
			if cfg.APIKey != "" {
				manager.AddSource(NewOpenWeatherMapSource(cfg.APIKey, cfg.RateLimit, cache))
			}
		case "weatherapi":
			if cfg.APIKey != "" {
				manager.AddSource(NewWeatherAPISource(cfg.APIKey, cfg.RateLimit, cache))
			}
		case "tomorrowio":
			if cfg.APIKey != "" {
				manager.AddSource(NewTomorrowIOSource(cfg.APIKey, cfg.RateLimit, cache))
			}
		case "visualcrossing":
			if cfg.APIKey != "" {
				manager.AddSource(NewVisualCrossingSource(cfg.APIKey, cfg.RateLimit, cache))
			}
		case "openmeteo":
			manager.AddSource(NewOpenMeteoSource(cfg.RateLimit, cache))
		}
	}
	
	return manager
}

func (m *DataSourceManager) AddSource(source WeatherDataSource) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sources[source.GetName()] = source
}

func (m *DataSourceManager) GetSource(name string) (WeatherDataSource, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	source, ok := m.sources[name]
	return source, ok
}

func (m *DataSourceManager) GetAllSources() []WeatherDataSource {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	sources := make([]WeatherDataSource, 0, len(m.sources))
	for _, source := range m.sources {
		sources = append(sources, source)
	}
	return sources
}

func (m *DataSourceManager) GetSourceNames() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	names := make([]string, 0, len(m.sources))
	for name := range m.sources {
		names = append(names, name)
	}
	return names
}