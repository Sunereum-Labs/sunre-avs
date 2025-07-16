import React, { useState, useEffect } from 'react';
import { Location, WeatherDataPoint } from '../types/insurance';

interface WeatherSimulatorProps {
  location: Location | null;
  scenario: 'normal' | 'heat_wave' | 'cold_snap' | 'storm';
  onDataUpdate: (data: WeatherDataPoint[]) => void;
}

const WeatherSimulator: React.FC<WeatherSimulatorProps> = ({ location, scenario, onDataUpdate }) => {
  const [currentData, setCurrentData] = useState<WeatherDataPoint[]>([]);
  const [dayIndex, setDayIndex] = useState(0);
  const [isSimulating, setIsSimulating] = useState(false);

  // Weather data sources
  const sources = ['OpenMeteo', 'WeatherAPI', 'Tomorrow.io', 'VisualCrossing', 'OpenWeatherMap'];

  const generateWeatherData = (day: number): WeatherDataPoint[] => {
    const baseTime = new Date();
    baseTime.setDate(baseTime.getDate() - (7 - day));

    return sources.map(source => {
      let temperature: number;
      let humidity: number;
      let windSpeed: number;
      let precipitation: number;

      switch (scenario) {
        case 'heat_wave':
          // Simulate extreme heat for days 3-7
          if (day >= 3 && day <= 7) {
            temperature = 36 + (day - 3) * 1.5 + Math.random() * 2;
            humidity = 30 + Math.random() * 10;
            windSpeed = 5 + Math.random() * 10;
            precipitation = 0;
          } else {
            temperature = 28 + Math.random() * 3;
            humidity = 50 + Math.random() * 15;
            windSpeed = 10 + Math.random() * 5;
            precipitation = Math.random() * 5;
          }
          break;

        case 'cold_snap':
          if (day >= 2 && day <= 5) {
            temperature = -12 - (day - 2) * 2 + Math.random() * 2;
            humidity = 70 + Math.random() * 10;
            windSpeed = 20 + Math.random() * 15;
            precipitation = 10 + Math.random() * 20;
          } else {
            temperature = 5 + Math.random() * 5;
            humidity = 60 + Math.random() * 10;
            windSpeed = 15 + Math.random() * 10;
            precipitation = Math.random() * 10;
          }
          break;

        case 'storm':
          if (day === 4 || day === 5) {
            temperature = 15 + Math.random() * 5;
            humidity = 90 + Math.random() * 10;
            windSpeed = 60 + Math.random() * 20;
            precipitation = 50 + Math.random() * 30;
          } else {
            temperature = 20 + Math.random() * 5;
            humidity = 70 + Math.random() * 10;
            windSpeed = 15 + Math.random() * 10;
            precipitation = 5 + Math.random() * 10;
          }
          break;

        default: // normal
          temperature = 20 + Math.sin(day * 0.5) * 5 + Math.random() * 2;
          humidity = 50 + Math.random() * 20;
          windSpeed = 10 + Math.random() * 10;
          precipitation = Math.random() * 15;
      }

      // Add small variations between sources
      const variation = (Math.random() - 0.5) * 0.5;
      
      return {
        source,
        temperature: temperature + variation,
        humidity,
        windSpeed,
        precipitation,
        timestamp: baseTime,
        confidence: 0.85 + Math.random() * 0.15
      };
    });
  };

  useEffect(() => {
    if (isSimulating && dayIndex < 10) {
      const timer = setTimeout(() => {
        const newData = generateWeatherData(dayIndex);
        setCurrentData(prev => [...prev, ...newData]);
        onDataUpdate([...currentData, ...newData]);
        setDayIndex(dayIndex + 1);
      }, 1000);

      return () => clearTimeout(timer);
    } else if (dayIndex >= 10) {
      setIsSimulating(false);
    }
  }, [isSimulating, dayIndex]);

  const startSimulation = () => {
    setCurrentData([]);
    setDayIndex(0);
    setIsSimulating(true);
  };

  const getLatestMetrics = () => {
    if (currentData.length === 0) return null;

    const latestData = currentData.slice(-sources.length);
    const avgTemp = latestData.reduce((sum, d) => sum + d.temperature, 0) / latestData.length;
    const avgHumidity = latestData.reduce((sum, d) => sum + (d.humidity || 0), 0) / latestData.length;
    const avgWind = latestData.reduce((sum, d) => sum + (d.windSpeed || 0), 0) / latestData.length;
    const avgPrecip = latestData.reduce((sum, d) => sum + (d.precipitation || 0), 0) / latestData.length;

    return { avgTemp, avgHumidity, avgWind, avgPrecip };
  };

  const metrics = getLatestMetrics();
  const isExtreme = scenario !== 'normal' && dayIndex >= 3 && dayIndex <= 7;

  if (!location) {
    return <div className="placeholder">Select a policy to view weather simulation</div>;
  }

  return (
    <div className="weather-simulator">
      <div className="weather-header">
        <h3>📍 {location.city}, {location.country}</h3>
        <p>Simulating: {scenario.replace('_', ' ').toUpperCase()} pattern</p>
      </div>

      {!isSimulating && currentData.length === 0 && (
        <div className="simulation-start">
          <button className="process-button" onClick={startSimulation}>
            Start Weather Simulation
          </button>
          <p>Simulate 10 days of weather data from multiple sources</p>
        </div>
      )}

      {(isSimulating || currentData.length > 0) && (
        <div className="weather-display">
          <div className="simulation-progress">
            <span>Day {dayIndex} of 10</span>
            <div className="progress-bar">
              <div className="progress-fill" style={{ width: `${(dayIndex / 10) * 100}%` }}></div>
            </div>
          </div>

          {metrics && (
            <div className="weather-metrics">
              <div className={`metric-card ${isExtreme && scenario === 'heat_wave' ? 'alert' : ''}`}>
                <h4>Temperature</h4>
                <div className="value">{metrics.avgTemp.toFixed(1)}°C</div>
              </div>
              <div className="metric-card">
                <h4>Humidity</h4>
                <div className="value">{metrics.avgHumidity.toFixed(0)}%</div>
              </div>
              <div className={`metric-card ${isExtreme && scenario === 'storm' ? 'alert' : ''}`}>
                <h4>Wind Speed</h4>
                <div className="value">{metrics.avgWind.toFixed(0)} km/h</div>
              </div>
              <div className={`metric-card ${isExtreme && scenario === 'storm' ? 'alert' : ''}`}>
                <h4>Precipitation</h4>
                <div className="value">{metrics.avgPrecip.toFixed(0)} mm</div>
              </div>
            </div>
          )}

          <div className="data-sources">
            <h4>Data Sources ({sources.length} active)</h4>
            <div className="source-list">
              {sources.map(source => (
                <div key={source} className="source-badge">
                  <span className="source-indicator"></span>
                  {source}
                </div>
              ))}
            </div>
          </div>

          {currentData.length > 0 && (
            <div className="consensus-info">
              <p>📊 Consensus Algorithm: MAD (Median Absolute Deviation)</p>
              <p>✅ Data Points Collected: {currentData.length}</p>
              <p>🔒 Average Confidence: {(currentData.reduce((sum, d) => sum + d.confidence, 0) / currentData.length * 100).toFixed(1)}%</p>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default WeatherSimulator;