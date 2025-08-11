import React, { useMemo, useState, useCallback, useRef } from 'react';
import {
  View,
  Text,
  ScrollView,
  ActivityIndicator,
  RefreshControl,
  TouchableOpacity,
  Animated,
} from 'react-native';
import { fetchReviews, ReviewsResponse } from 'services/reviewsService';
import { useQuery } from 'hooks/useQuery';
import { Header } from 'components/Header';
import { ReviewItem } from 'components/ReviewItem';
import { ReviewsFilter } from 'components/ReviewsFilter';
import { DEFAULT_TIME_RANGE, Rating, TimeRange } from 'types/reviews';
import { MaterialIcons } from '@expo/vector-icons';

const LoadingContent = () => (
  <View className="flex-1 items-center justify-center">
    <ActivityIndicator size="large" color="#3B82F6" testID="activity-indicator" />
    <Text className="mt-4 text-gray-600">Loading reviews...</Text>
  </View>
);

const ErrorContent = ({ error }: { error: string }) => (
  <View className="flex-1 items-center justify-center">
    <Text className="text-center text-red-500">{error}</Text>
  </View>
);

const NoDataContent = () => (
  <View className="flex-1 items-center justify-center">
    <Text className="text-center text-gray-600">No reviews data available</Text>
  </View>
);

export const ReviewsScreen = () => {
  const [selectedRating, setSelectedRating] = useState<Rating | undefined>();
  const [selectedTimeRange, setSelectedTimeRange] = useState<TimeRange>(DEFAULT_TIME_RANGE);
  const [showBackToTop, setShowBackToTop] = useState(false);

  const scrollViewRef = useRef<ScrollView>(null);
  const fadeAnim = useRef(new Animated.Value(0)).current;

  // Memoize the query function to prevent infinite updates
  const queryFn = useCallback(
    () => fetchReviews(selectedRating, selectedTimeRange),
    [selectedRating, selectedTimeRange]
  );

  // Note: I would typically use a library like react-query for this, but as requested I'm avoiding 3rd party libraries
  const { data: reviewsData, loading, error, refetch } = useQuery<ReviewsResponse>(queryFn);

  // Handle scroll events to show/hide back to top button
  const handleScroll = useCallback(
    (event: any) => {
      const offsetY = event.nativeEvent.contentOffset.y;
      const shouldShow = offsetY > 200; // Show button after scrolling 200px down

      if (shouldShow !== showBackToTop) {
        setShowBackToTop(shouldShow);
        Animated.timing(fadeAnim, {
          toValue: shouldShow ? 1 : 0,
          duration: 300,
          useNativeDriver: true,
        }).start();
      }
    },
    [showBackToTop, fadeAnim]
  );

  // Scroll to top function
  const scrollToTop = useCallback(() => {
    scrollViewRef.current?.scrollTo({ y: 0, animated: true });
  }, []);

  const content = useMemo(() => {
    if (loading) return <LoadingContent />;
    if (error) return <ErrorContent error={error} />;
    if (!reviewsData || reviewsData.reviews.length === 0) return <NoDataContent />;

    // render content
    return (
      <ScrollView
        ref={scrollViewRef}
        className="flex-1"
        contentContainerStyle={{ padding: 16, paddingTop: 0 }}
        showsVerticalScrollIndicator={true}
        onScroll={handleScroll}
        scrollEventThrottle={16}
        refreshControl={<RefreshControl refreshing={loading} onRefresh={refetch} />}>
        {reviewsData.reviews.map((review) => (
          <ReviewItem key={review.id} review={review} />
        ))}
      </ScrollView>
    );
  }, [loading, reviewsData, error, refetch, handleScroll]);

  return (
    <View className="flex-1 bg-gray-50">
      <Header
        appId={reviewsData?.appId}
        reviewCount={reviewsData?.count}
        timeRange={reviewsData?.lastHours}
        onRefresh={refetch}
        loading={loading}
      />
      <ReviewsFilter
        selectedRating={selectedRating}
        onSelectRating={setSelectedRating}
        selectedTimeRange={selectedTimeRange}
        onSelectTimeRange={setSelectedTimeRange}
      />
      {content}

      {/* Floating Back to Top Button */}
      <Animated.View
        style={{
          position: 'absolute',
          bottom: 20,
          right: 20,
          opacity: fadeAnim,
        }}
        pointerEvents={showBackToTop ? 'auto' : 'none'}>
        <TouchableOpacity
          onPress={scrollToTop}
          className="rounded-full bg-blue-500 p-3 shadow-lg"
          testID="back-to-top-button">
          <MaterialIcons name="keyboard-arrow-up" size={36} color="white" />
        </TouchableOpacity>
      </Animated.View>
    </View>
  );
};
