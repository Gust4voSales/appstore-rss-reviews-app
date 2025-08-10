import React from 'react';
import { View, TouchableOpacity } from 'react-native';
import { Rating } from 'types/reviews';
import { MaterialIcons } from '@expo/vector-icons';

const validRatings: Rating[] = [1, 2, 3, 4, 5];

interface ReviewsFilterProps {
  selectedRating?: Rating;
  onSelectRating: (rating?: Rating) => void;
}

export const ReviewsFilter = ({ selectedRating, onSelectRating }: ReviewsFilterProps) => (
  <View className="relative flex-row items-center justify-center px-4 py-2">
    <MaterialIcons name="filter-list" size={24} color="gray" className="absolute left-4" />
    <View className="flex-row justify-center space-x-2">
      {validRatings.map((rating) => (
        <TouchableOpacity
          key={rating}
          testID={`star-button-${rating}`}
          onPress={() => onSelectRating(selectedRating === rating ? undefined : rating)}
          className={`rounded-full px-4 py-2`}>
          <MaterialIcons
            name={selectedRating && selectedRating >= rating ? 'star' : 'star-border'}
            size={24}
            color={'#EAB308'}
          />
        </TouchableOpacity>
      ))}
    </View>
  </View>
);
