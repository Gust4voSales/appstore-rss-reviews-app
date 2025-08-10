import React, { useMemo } from 'react';
import { View, Text, ScrollView, ActivityIndicator, RefreshControl } from 'react-native';
import { fetchReviews, ReviewsResponse } from 'services/reviewsService';
import { useQuery } from 'hooks/useQuery';
import { Header } from 'components/Header';
import { ReviewItem } from 'components/ReviewItem';

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
  // Note: I would typically use a library like react-query for this, but as requested I'm avoiding 3rd party libraries
  const { data: reviewsData, loading, error, refetch } = useQuery<ReviewsResponse>(fetchReviews);

  const content = useMemo(() => {
    if (loading) return <LoadingContent />;
    if (error) return <ErrorContent error={error} />;
    if (!reviewsData || reviewsData.reviews.length === 0) return <NoDataContent />;

    // render content
    return (
      <ScrollView
        className="flex-1"
        contentContainerStyle={{ padding: 16 }}
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
        onRefresh={refetch}
        loading={loading}
      />
      {content}
    </View>
  );
};
