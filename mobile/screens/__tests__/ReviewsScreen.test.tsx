import React from 'react';
import { render, screen, waitFor, act, fireEvent } from '@testing-library/react-native';
import { ReviewsScreen } from '../ReviewsScreen';
import { fetchReviews, ReviewsResponse } from 'services/reviewsService';
import { Review } from 'types/reviews';

// Mock the services
jest.mock('services/reviewsService');
const mockFetchReviews = fetchReviews as jest.MockedFunction<typeof fetchReviews>;

describe('ReviewsScreen', () => {
  // Test data
  const mockReviews: Review[] = [
    {
      id: '1',
      title: 'Great app!',
      content: 'This app is amazing and works perfectly.',
      author: 'John Doe',
      rating: 5,
      updatedAt: '2024-01-15T10:00:00Z',
    },
    {
      id: '2',
      title: 'Could be better',
      content: 'The app has some issues but overall decent.',
      author: 'Jane Smith',
      rating: 3,
      updatedAt: '2024-01-14T15:30:00Z',
    },
  ];

  const mockReviewsResponse: ReviewsResponse = {
    appId: 'test-app-123',
    count: 2,
    reviews: mockReviews,
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('when rendered', () => {
    it('fetches data when rendered', async () => {
      // Arrange
      mockFetchReviews.mockResolvedValue(mockReviewsResponse);

      // Act
      render(<ReviewsScreen />);

      // Assert
      await waitFor(() => {
        expect(mockFetchReviews).toHaveBeenCalledTimes(1);
      });
    });
  });

  describe('loading state', () => {
    it('shows loading indicator when loading', async () => {
      // Arrange
      let resolvePromise: (value: ReviewsResponse) => void;
      const promise = new Promise<ReviewsResponse>((resolve) => {
        resolvePromise = resolve;
      });
      mockFetchReviews.mockReturnValue(promise);

      // Act
      render(<ReviewsScreen />);

      // Assert - should show loading indicator immediately
      expect(screen.getByText('Loading reviews...')).toBeTruthy();
      expect(screen.getByTestId('activity-indicator')).toBeTruthy();

      // Clean up - resolve the promise to avoid hanging
      await act(async () => {
        resolvePromise!(mockReviewsResponse);
        await promise;
      });
    });
  });

  describe('error state', () => {
    it('shows error message when failing to get data', async () => {
      // Arrange
      const errorMessage = 'Network error';
      mockFetchReviews.mockRejectedValue(new Error(errorMessage));

      // Act
      render(<ReviewsScreen />);

      // Assert
      await waitFor(() => {
        expect(screen.getByText('An error occurred while fetching data')).toBeTruthy();
      });
    });
  });

  describe('success state', () => {
    it('renders fetched data', async () => {
      // Arrange
      mockFetchReviews.mockResolvedValue(mockReviewsResponse);

      // Act
      render(<ReviewsScreen />);

      // Assert - Wait for data to load and verify all components are rendered
      await waitFor(() => {
        // Header should show app info
        expect(screen.getByText('App Reviews')).toBeTruthy();
        expect(screen.getByText('App ID: test-app-123 • 2 reviews (96h)')).toBeTruthy();

        // Reviews should be rendered
        expect(screen.getByText('Great app!')).toBeTruthy();
        expect(screen.getByText('This app is amazing and works perfectly.')).toBeTruthy();
        expect(screen.getByText('By John Doe')).toBeTruthy();

        expect(screen.getByText('Could be better')).toBeTruthy();
        expect(screen.getByText('The app has some issues but overall decent.')).toBeTruthy();
        expect(screen.getByText('By Jane Smith')).toBeTruthy();
      });

      // Should not show loading indicator after success
      expect(screen.queryByText('Loading reviews...')).toBeNull();
      expect(screen.queryByText('An error occurred while fetching data')).toBeNull();
    });

    it('renders no data message when reviews array is empty', async () => {
      // Arrange
      const emptyResponse: ReviewsResponse = {
        appId: 'test-app-123',
        count: 0,
        reviews: [],
      };
      mockFetchReviews.mockResolvedValue(emptyResponse);

      // Act
      render(<ReviewsScreen />);

      // Assert
      await waitFor(() => {
        expect(screen.getByText('No reviews data available')).toBeTruthy();
      });
    });

    it('renders no data message when response is null', async () => {
      // Arrange
      mockFetchReviews.mockResolvedValue(null as any);

      // Act
      render(<ReviewsScreen />);

      // Assert
      await waitFor(() => {
        expect(screen.getByText('No reviews data available')).toBeTruthy();
      });
    });
  });

  describe('refresh functionality', () => {
    it('calls fetchReviews again when refresh is triggered', async () => {
      // Arrange
      mockFetchReviews.mockResolvedValue(mockReviewsResponse);

      // Act
      const { getByTestId } = render(<ReviewsScreen />);

      // Wait for initial load
      await waitFor(() => {
        expect(mockFetchReviews).toHaveBeenCalledTimes(1);
      });

      // Trigger refresh by pressing the refresh button
      const refreshButton = getByTestId('refresh-button');
      act(() => {
        fireEvent.press(refreshButton);
      });

      // Assert
      await waitFor(() => {
        expect(mockFetchReviews).toHaveBeenCalledTimes(2);
      });
    });
  });

  describe('filter functionality', () => {
    it('filters reviews by rating', async () => {
      // Arrange
      mockFetchReviews.mockResolvedValue(mockReviewsResponse);

      // Act
      const { getByTestId } = render(<ReviewsScreen />);

      // Wait for initial load
      await waitFor(() => {
        expect(mockFetchReviews).toHaveBeenCalledTimes(1);
      });

      // Trigger filter by pressing the filter button
      const star3Button = getByTestId('star-button-3');
      act(() => {
        fireEvent.press(star3Button);
      });

      // Assert
      await waitFor(() => {
        expect(mockFetchReviews).toHaveBeenCalledTimes(2);
        expect(screen.queryByText('Great app!')).toBeNull();
      });
      expect(mockFetchReviews).toHaveBeenLastCalledWith(3);
      expect(screen.getByText('Could be better')).toBeTruthy();
    });
  });
});
