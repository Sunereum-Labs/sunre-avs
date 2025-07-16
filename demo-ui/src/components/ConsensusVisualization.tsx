import React from 'react';
import { Activity, CheckCircle, Info } from 'lucide-react';

interface ConsensusVisualizationProps {
  data: {
    consensus: {
      algorithm: string;
      confidence: number;
      totalSources: number;
    };
    sources: Array<{
      name: string;
      temperature: number;
      confidence: number;
    }>;
  };
  detailed?: boolean;
}

const ConsensusVisualization: React.FC<ConsensusVisualizationProps> = ({ data, detailed = false }) => {
  const getSourceIcon = (name: string) => {
    if (name.toLowerCase().includes('openmeteo')) return '🌍';
    if (name.toLowerCase().includes('tomorrow')) return '🔮';
    if (name.toLowerCase().includes('weather')) return '☁️';
    return '📡';
  };

  const averageTemp = data.sources.reduce((sum, s) => sum + s.temperature, 0) / data.sources.length;
  const variance = data.sources.reduce((sum, s) => sum + Math.pow(s.temperature - averageTemp, 2), 0) / data.sources.length;
  const stdDev = Math.sqrt(variance);

  return (
    <div className={`bg-gradient-to-br from-green-900/20 to-emerald-900/20 backdrop-blur-sm rounded-2xl p-6 border border-green-800/50 ${detailed ? '' : ''}`}>
      <h3 className="text-xl font-semibold mb-4 flex items-center">
        <Activity className="w-6 h-6 mr-2 text-green-400" />
        Consensus Algorithm
      </h3>

      <div className="grid grid-cols-2 gap-4 mb-6">
        <div className="bg-dark-500/50 rounded-lg p-4">
          <p className="text-sm text-gray-400 mb-1">Algorithm</p>
          <p className="font-semibold">{data.consensus.algorithm}</p>
        </div>
        <div className="bg-dark-500/50 rounded-lg p-4">
          <p className="text-sm text-gray-400 mb-1">Confidence</p>
          <p className="font-semibold text-green-400">{(data.consensus.confidence * 100).toFixed(1)}%</p>
        </div>
      </div>

      <div className="space-y-3">
        {data.sources.map((source) => {
          const deviation = Math.abs(source.temperature - averageTemp);
          const isOutlier = deviation > stdDev * 2;
          
          return (
            <div key={source.name} className="flex items-center justify-between p-3 bg-dark-500/30 rounded-lg">
              <div className="flex items-center space-x-3">
                <span className="text-2xl">{getSourceIcon(source.name)}</span>
                <div>
                  <p className="font-medium">{source.name}</p>
                  <p className="text-sm text-gray-400">{source.temperature.toFixed(1)}°C</p>
                </div>
              </div>
              <div className="flex items-center space-x-2">
                {isOutlier ? (
                  <span className="px-2 py-1 bg-yellow-900/50 text-yellow-400 text-xs rounded">Outlier</span>
                ) : (
                  <CheckCircle className="w-5 h-5 text-green-500" />
                )}
                <div className="text-right">
                  <p className="text-sm font-medium">{(source.confidence * 100).toFixed(0)}%</p>
                  <p className="text-xs text-gray-500">confidence</p>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {detailed && (
        <div className="mt-6 p-4 bg-dark-500/30 rounded-lg">
          <div className="flex items-start space-x-2">
            <Info className="w-5 h-5 text-blue-400 mt-0.5" />
            <div className="text-sm text-gray-300">
              <p className="font-medium mb-1">How MAD Consensus Works:</p>
              <ul className="space-y-1 text-gray-400">
                <li>• Collects temperature from {data.consensus.totalSources} independent sources</li>
                <li>• Calculates median absolute deviation to identify outliers</li>
                <li>• Removes data points beyond 2.5 standard deviations</li>
                <li>• Produces consensus temperature with {(data.consensus.confidence * 100).toFixed(0)}% confidence</li>
              </ul>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default ConsensusVisualization;