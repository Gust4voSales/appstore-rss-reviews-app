import { Review } from '../types/reviews';
import config from 'config/env';

export interface ReviewsResponse {
  appId: string;
  count: number;
  reviews: Review[];
}

export const fetchReviews = async (): Promise<ReviewsResponse> => {
  try {
    const response = await fetch(`${config.API_URL}/reviews/96h`);

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
