import React, { useMemo, useState, useCallback, useRef } from "react";
import { ArrowUp, Loader2 } from "lucide-react";
import { fetchReviews, ReviewsResponse } from "../services/reviewsService";
import { useQuery } from "../hooks/useQuery";
import { Header } from "../components/Header";
import { ReviewItem } from "../components/ReviewItem";
import { ReviewsFilter } from "../components/ReviewsFilter";
import { DEFAULT_TIME_RANGE, Rating, TimeRange } from "../types/reviews";
import { DefaultButton } from "../components/DefaultButton";

const LoadingContent = () => (
  <div className="flex-1 flex items-center justify-center">
    <div className="text-center">
      <Loader2 className="animate-spin rounded-full h-12 w-12 mx-auto text-blue-600" />
      <p className="mt-4 text-gray-600">Loading reviews...</p>
    </div>
  </div>
);

const ErrorContent = ({ error }: { error: string }) => (
  <div className="flex-1 flex items-center justify-center">
    <p className="text-center text-red-500">{error}</p>
  </div>
);

const NoDataContent = () => (
  <div className="flex-1 flex items-center justify-center">
    <p className="text-center text-gray-600">No reviews data available</p>
  </div>
);

export const ReviewsPage = () => {
  const [selectedRating, setSelectedRating] = useState<Rating | undefined>();
  const [selectedTimeRange, setSelectedTimeRange] = useState<TimeRange>(DEFAULT_TIME_RANGE);
  const [showBackToTop, setShowBackToTop] = useState(false);

  const scrollContainerRef = useRef<HTMLDivElement>(null);

  // Memoize the query function to prevent infinite updates
  const queryFn = useCallback(
    () => fetchReviews(selectedRating, selectedTimeRange),
    [selectedRating, selectedTimeRange]
  );

  // Note: I would typically use a library like react-query for this, but as requested I'm avoiding 3rd party libraries
  const { data: reviewsData, loading, error, refetch } = useQuery<ReviewsResponse>(queryFn);

  // Handle scroll events to show/hide back to top button
  const handleScroll = useCallback(
    (event: React.UIEvent<HTMLDivElement>) => {
      const scrollTop = event.currentTarget.scrollTop;
      const shouldShow = scrollTop > 200; // Show button after scrolling 200px down

      if (shouldShow !== showBackToTop) {
        setShowBackToTop(shouldShow);
      }
    },
    [showBackToTop]
  );

  const scrollToTop = () => {
    scrollContainerRef.current?.scrollTo({ top: 0, behavior: "smooth" });
  };

  const content = useMemo(() => {
    if (loading) return <LoadingContent />;
    if (error) return <ErrorContent error={error} />;
    if (!reviewsData || reviewsData.reviews.length === 0) return <NoDataContent />;

    // render content
    return (
      <div
        ref={scrollContainerRef}
        className="flex-1 overflow-y-auto"
        style={{ padding: "16px", paddingTop: "0" }}
        onScroll={handleScroll}
      >
        {reviewsData.reviews.map((review) => (
          <ReviewItem key={review.id} review={review} />
        ))}
      </div>
    );
  }, [loading, reviewsData, error, handleScroll]);

  return (
    <div className="flex flex-col h-screen bg-gray-50">
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
      {showBackToTop && (
        <DefaultButton
          onClick={scrollToTop}
          className="fixed bottom-5 right-5 !rounded-full h-10 w-10 shadow-lg"
          data-testid="back-to-top-button"
        >
          <ArrowUp className="size-6" />
        </DefaultButton>
      )}
    </div>
  );
};
