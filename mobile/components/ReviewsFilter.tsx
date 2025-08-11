import React, { useState } from 'react';
import { View, TouchableOpacity, Text, Modal, ScrollView } from 'react-native';
import { Rating, TimeRange } from 'types/reviews';
import { MaterialIcons } from '@expo/vector-icons';

const validRatings: Rating[] = [1, 2, 3, 4, 5];
const validTimeRanges: TimeRange[] = [24, 48, 72, 96];

interface ReviewsFilterProps {
  selectedRating?: Rating;
  onSelectRating: (rating?: Rating) => void;
  selectedTimeRange: TimeRange;
  onSelectTimeRange: (timeRange: TimeRange) => void;
}

export const ReviewsFilter = ({
  selectedRating,
  onSelectRating,
  selectedTimeRange,
  onSelectTimeRange,
}: ReviewsFilterProps) => {
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);

  const handleTimeRangeSelect = (timeRange: TimeRange) => {
    onSelectTimeRange(timeRange);
    setIsDropdownOpen(false);
  };

  return (
    <View className="relative mb-2 items-center gap-3 border-b border-gray-200 bg-white px-4 py-3">
      <MaterialIcons
        name="filter-list"
        size={24}
        color="gray"
        className="absolute left-4 top-1/2"
      />

      {/* Rating Filter */}
      <View className="flex-row items-center justify-center space-x-2">
        {validRatings.map((rating) => (
          <TouchableOpacity
            key={rating}
            testID={`star-button-${rating}`}
            onPress={() => onSelectRating(selectedRating === rating ? undefined : rating)}
            className="rounded-full px-2 py-1">
            <MaterialIcons
              name={selectedRating && selectedRating >= rating ? 'star' : 'star-border'}
              size={24}
              color={'#EAB308'}
            />
          </TouchableOpacity>
        ))}
      </View>

      {/* Time Range Dropdown */}
      <View className="relative">
        <TouchableOpacity
          testID="time-range-dropdown-trigger"
          onPress={() => setIsDropdownOpen(!isDropdownOpen)}
          className="flex-row items-center rounded-lg border border-gray-300 bg-white px-4 py-2 shadow-sm">
          <MaterialIcons name="schedule" size={16} color="gray" className="mr-2" />
          <Text className="mr-2 text-sm font-medium text-gray-700">Last {selectedTimeRange}h</Text>
          <MaterialIcons
            name={isDropdownOpen ? 'keyboard-arrow-up' : 'keyboard-arrow-down'}
            size={20}
            color="gray"
          />
        </TouchableOpacity>

        <Modal
          visible={isDropdownOpen}
          transparent={true}
          animationType="fade"
          onRequestClose={() => setIsDropdownOpen(false)}>
          <TouchableOpacity
            className="flex-1 bg-black/20"
            activeOpacity={1}
            onPress={() => setIsDropdownOpen(false)}>
            <View className="absolute left-4 right-4 top-32 rounded-lg border border-gray-200 bg-white shadow-lg">
              <ScrollView>
                {validTimeRanges.map((timeRange) => (
                  <TouchableOpacity
                    key={timeRange}
                    testID={`time-range-option-${timeRange}`}
                    onPress={() => handleTimeRangeSelect(timeRange)}
                    className={`flex-row items-center justify-between px-4 py-3 ${
                      selectedTimeRange === timeRange ? 'bg-blue-50' : ''
                    }`}>
                    <Text
                      className={`text-base ${
                        selectedTimeRange === timeRange
                          ? 'font-semibold text-blue-600'
                          : 'text-gray-700'
                      }`}>
                      Last {timeRange}h
                    </Text>
                    {selectedTimeRange === timeRange && (
                      <MaterialIcons name="check" size={20} color="#2563EB" />
                    )}
                  </TouchableOpacity>
                ))}
              </ScrollView>
            </View>
          </TouchableOpacity>
        </Modal>
      </View>
    </View>
  );
};
