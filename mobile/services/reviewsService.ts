import { DEFAULT_TIME_RANGE, Rating, Review, TimeRange } from 'types/reviews';
import config from 'config/env';

export interface ReviewsResponse {
  appId: string;
  count: number;
  reviews: Review[];
  lastHours: TimeRange;
}

export const fetchReviews = async (rating?: Rating, timeRange: TimeRange = DEFAULT_TIME_RANGE): Promise<ReviewsResponse> => {
  try {
    const url = new URL(`${config.API_URL}/reviews`);
    if (rating !== undefined) {
      url.searchParams.append('rating', rating.toString());
    }
    if (timeRange !== undefined) {
      url.searchParams.append('hours', timeRange.toString());
    }
    const response = await fetch(url.toString(), {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data: ReviewsResponse = await response.json();
    return data;
  } catch (error) {
    console.error('Error fetching reviews:', error);
    throw error;
  }
};
