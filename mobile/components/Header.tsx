import React from 'react';
import { View, Text, TouchableOpacity } from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';

interface HeaderProps {
  appId?: string;
  reviewCount?: number;
  onRefresh: () => void;
  loading: boolean;
}

export const Header: React.FC<HeaderProps> = ({ appId, reviewCount, onRefresh, loading }) => {
  return (
    <View className="border-b border-gray-200 bg-white p-4 shadow-sm">
      <View className="flex-row items-center justify-between">
        <View className="flex-1">
          <Text className="text-2xl font-bold text-gray-900">App Reviews</Text>
          <Text className="text-gray-600">
            {appId && reviewCount ? (
              <>
                App ID: {appId} • {reviewCount} reviews (96h)
              </>
            ) : (
              <>App ID: - </>
            )}
          </Text>
        </View>
        <TouchableOpacity
          onPress={onRefresh}
          disabled={loading}
          testID="refresh-button"
          className={`rounded-lg p-2 ${loading ? 'bg-gray-300' : 'bg-blue-500'}`}>
          <MaterialIcons name="refresh" size={24} color="white" />
        </TouchableOpacity>
      </View>
    </View>
  );
};
