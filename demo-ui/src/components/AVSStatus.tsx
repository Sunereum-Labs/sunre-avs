import React from 'react';
import { Wifi, WifiOff, Loader2 } from 'lucide-react';
import { AVSConnection } from '../types/insurance';

interface AVSStatusProps {
  connection: AVSConnection;
}

const AVSStatus: React.FC<AVSStatusProps> = ({ connection }) => {
  const getStatusColor = () => {
    switch (connection.status) {
      case 'connected':
        return 'bg-green-500';
      case 'connecting':
        return 'bg-yellow-500';
      default:
        return 'bg-red-500';
    }
  };

  const getStatusText = () => {
    switch (connection.status) {
      case 'connected':
        return 'Connected';
      case 'connecting':
        return 'Connecting...';
      default:
        return 'Disconnected';
    }
  };

  const getStatusIcon = () => {
    switch (connection.status) {
      case 'connected':
        return <Wifi className="w-4 h-4" />;
      case 'connecting':
        return <Loader2 className="w-4 h-4 animate-spin" />;
      default:
        return <WifiOff className="w-4 h-4" />;
    }
  };

  return (
    <div className="flex items-center space-x-3 bg-dark-500/50 backdrop-blur-sm px-4 py-2 rounded-lg">
      <div className="flex items-center space-x-2">
        <div className={`w-2 h-2 rounded-full ${getStatusColor()} animate-pulse`} />
        {getStatusIcon()}
        <span className="text-sm font-medium">{getStatusText()}</span>
      </div>
      <div className="text-xs text-gray-400 border-l border-gray-700 pl-3">
        {connection.networkType}
      </div>
    </div>
  );
};

export default AVSStatus;