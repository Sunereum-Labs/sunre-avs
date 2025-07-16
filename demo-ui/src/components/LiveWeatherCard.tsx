import React from 'react';
import { Thermometer, TrendingUp, TrendingDown, Activity, MapPin, Clock } from 'lucide-react';

interface LiveWeatherCardProps {
  data: {
    temperature: number;
    trend: 'warming' | 'cooling' | 'stable';
    historicalData: {
      temperature4hAgo: number;
      change: number;
    };
  };
  detailed?: boolean;
}

const LiveWeatherCard: React.FC<LiveWeatherCardProps> = ({ data, detailed = false }) => {
  const getTrendIcon = () => {
    switch (data.trend) {
      case 'warming':
        return <TrendingUp className="w-5 h-5 text-red-500" />;
      case 'cooling':
        return <TrendingDown className="w-5 h-5 text-blue-500" />;
      default:
        return <Activity className="w-5 h-5 text-gray-500" />;
    }
  };

  const getTrendColor = () => {
    switch (data.trend) {
      case 'warming':
        return 'text-red-500';
      case 'cooling':
        return 'text-blue-500';
      default:
        return 'text-gray-500';
    }
  };

  return (
    <div className={`bg-gradient-to-br from-blue-900/20 to-cyan-900/20 backdrop-blur-sm rounded-2xl p-6 border border-blue-800/50 ${detailed ? 'h-full' : ''}`}>
      <div className="flex items-start justify-between mb-4">
        <div>
          <h3 className="text-xl font-semibold mb-2 flex items-center">
            <Thermometer className="w-6 h-6 mr-2 text-blue-400" />
            Current Temperature
          </h3>
          <div className="flex items-center text-sm text-gray-400 space-x-4">
            <span className="flex items-center">
              <MapPin className="w-4 h-4 mr-1" />
              New York City
            </span>
            <span className="flex items-center">
              <Clock className="w-4 h-4 mr-1" />
              Live
            </span>
          </div>
        </div>
        <div className="text-right">
          <div className="text-4xl font-bold">{data.temperature.toFixed(1)}°C</div>
          <div className="text-sm text-gray-400 mt-1">
            {(data.temperature * 9/5 + 32).toFixed(1)}°F
          </div>
        </div>
      </div>

      <div className="border-t border-gray-700 pt-4">
        <div className="grid grid-cols-2 gap-4">
          <div className="bg-dark-500/50 rounded-lg p-3">
            <p className="text-sm text-gray-400 mb-1">4 hours ago</p>
            <p className="text-xl font-semibold">{data.historicalData.temperature4hAgo.toFixed(1)}°C</p>
          </div>
          <div className="bg-dark-500/50 rounded-lg p-3">
            <p className="text-sm text-gray-400 mb-1">Change</p>
            <div className="flex items-center space-x-2">
              {getTrendIcon()}
              <p className={`text-xl font-semibold ${getTrendColor()}`}>
                {data.historicalData.change > 0 ? '+' : ''}{data.historicalData.change.toFixed(1)}°C
              </p>
            </div>
          </div>
        </div>
      </div>

      {detailed && (
        <div className="mt-4 p-4 bg-dark-500/30 rounded-lg">
          <p className="text-sm text-gray-400">Temperature Trend</p>
          <p className="text-lg font-medium capitalize mt-1">{data.trend}</p>
          <div className="mt-2 h-24 flex items-end justify-between">
            {[...Array(8)].map((_, i) => {
              const height = Math.random() * 60 + 20;
              return (
                <div
                  key={i}
                  className="w-8 bg-gradient-to-t from-blue-500 to-cyan-500 rounded-t"
                  style={{ height: `${height}%` }}
                />
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
};

export default LiveWeatherCard;