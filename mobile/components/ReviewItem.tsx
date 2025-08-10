import React from 'react';
import { View, Text } from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { Review } from 'types/reviews';

interface ReviewItemProps {
  review: Review;
}

export const ReviewItem: React.FC<ReviewItemProps> = ({ review }) => {
  const renderStars = (rating: number) => {
    return (
      <View className="flex-row">
        {[1, 2, 3, 4, 5].map((index) => (
          <MaterialIcons
            key={index}
            name={index <= rating ? 'star' : 'star-border'}
            size={18}
            color="#EAB308"
          />
        ))}
      </View>
    );
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  return (
    <View className="mb-4 rounded-lg border border-gray-100 bg-white p-4 shadow-sm">
      <View className="mb-2 flex-row items-start justify-between gap-2">
        <Text className="flex-1 text-lg font-bold text-gray-900" numberOfLines={0}>
          {review.title}
        </Text>
        {renderStars(review.rating)}
      </View>

      <Text className="mb-2 text-gray-700">{review.content}</Text>

      <View className="flex-row items-start justify-between gap-2">
        <Text className="flex-1 text-sm text-gray-500">By {review.author}</Text>
        <Text className="text-sm text-gray-400">{formatDate(review.updatedAt)}</Text>
      </View>
    </View>
  );
};
