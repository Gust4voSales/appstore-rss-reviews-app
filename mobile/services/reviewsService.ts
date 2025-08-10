import { Rating, Review } from 'types/reviews';
import config from 'config/env';

export interface ReviewsResponse {
  appId: string;
  count: number;
  reviews: Review[];
}

export const fetchReviews = async (rating?: Rating): Promise<ReviewsResponse> => {
  try {
    const url = new URL(`${config.API_URL}/reviews/96h`);
    if (rating !== undefined) {
      url.searchParams.append('rating', rating.toString());
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
