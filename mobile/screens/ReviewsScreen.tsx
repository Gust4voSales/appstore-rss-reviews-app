import React, { useMemo, useState, useCallback } from 'react';
import { View, Text, ScrollView, ActivityIndicator, RefreshControl } from 'react-native';
import { fetchReviews, ReviewsResponse } from 'services/reviewsService';
import { useQuery } from 'hooks/useQuery';
import { Header } from 'components/Header';
import { ReviewItem } from 'components/ReviewItem';
import { ReviewsFilter } from 'components/ReviewsFilter';
import { DEFAULT_TIME_RANGE, Rating, TimeRange } from 'types/reviews';

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

  // Memoize the query function to prevent infinite updates
  const queryFn = useCallback(
    () => fetchReviews(selectedRating, selectedTimeRange),
    [selectedRating, selectedTimeRange]
  );

  // Note: I would typically use a library like react-query for this, but as requested I'm avoiding 3rd party libraries
  const { data: reviewsData, loading, error, refetch } = useQuery<ReviewsResponse>(queryFn);

  const content = useMemo(() => {
    if (loading) return <LoadingContent />;
    if (error) return <ErrorContent error={error} />;
    if (!reviewsData || reviewsData.reviews.length === 0) return <NoDataContent />;

    // render content
    return (
      <ScrollView
        className="flex-1"
        contentContainerStyle={{ padding: 16, paddingTop: 0 }}
        showsVerticalScrollIndicator={true}
        refreshControl={<RefreshControl refreshing={loading} onRefresh={refetch} />}>
        {reviewsData.reviews.map((review) => (
          <ReviewItem key={review.id} review={review} />
        ))}
      </ScrollView>
    );
  }, [loading, reviewsData, error, refetch]);

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
    </View>
  );
};
