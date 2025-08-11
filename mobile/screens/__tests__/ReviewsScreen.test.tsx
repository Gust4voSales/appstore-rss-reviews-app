import React from 'react';
import { render, screen, waitFor, act, fireEvent } from '@testing-library/react-native';
import { ReviewsScreen } from '../ReviewsScreen';
import { fetchReviews, ReviewsResponse } from 'services/reviewsService';
import { DEFAULT_TIME_RANGE, Review } from 'types/reviews';

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
    lastHours: 48,
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
      const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {}); // avoid console error
      // Arrange
      const errorMessage = 'Network error';
      mockFetchReviews.mockRejectedValue(new Error(errorMessage));

      // Act
      render(<ReviewsScreen />);

      // Assert
      await waitFor(() => {
        expect(screen.getByText('An error occurred while fetching data')).toBeTruthy();
      });
      // Clean up
      consoleSpy.mockRestore();
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
        expect(
          screen.getByText(`App ID: test-app-123 • 2 reviews (${DEFAULT_TIME_RANGE}h)`)
        ).toBeTruthy();

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
        lastHours: 48,
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
      expect(mockFetchReviews).toHaveBeenLastCalledWith(3, DEFAULT_TIME_RANGE);
      expect(screen.getByText('Could be better')).toBeTruthy();
    });

    it('filters reviews by date range', async () => {
      // Arrange
      mockFetchReviews.mockResolvedValue(mockReviewsResponse);

      // Act
      const { getByTestId } = render(<ReviewsScreen />);

      // Wait for initial load
      await waitFor(() => {
        expect(mockFetchReviews).toHaveBeenCalledTimes(1);
      });

      // Open the time range dropdown
      const timeRangeDropdownTrigger = getByTestId('time-range-dropdown-trigger');
      act(() => {
        fireEvent.press(timeRangeDropdownTrigger);
      });

      // Wait for dropdown to open and select 72h option
      await waitFor(() => {
        expect(getByTestId('time-range-option-72')).toBeTruthy();
      });

      const timeRange72Option = getByTestId('time-range-option-72');
      act(() => {
        fireEvent.press(timeRange72Option);
      });

      // Assert
      await waitFor(() => {
        expect(mockFetchReviews).toHaveBeenCalledTimes(2);
      });
      expect(mockFetchReviews).toHaveBeenLastCalledWith(undefined, 72);
    });

    it('combines rating and time range filters', async () => {
      // Arrange
      mockFetchReviews.mockResolvedValue(mockReviewsResponse);

      // Act
      const { getByTestId } = render(<ReviewsScreen />);

      // Wait for initial load
      await waitFor(() => {
        expect(mockFetchReviews).toHaveBeenCalledTimes(1);
      });

      // First select a rating filter
      const star4Button = getByTestId('star-button-4');
      act(() => {
        fireEvent.press(star4Button);
      });

      await waitFor(() => {
        expect(mockFetchReviews).toHaveBeenCalledTimes(2);
      });

      // Then select a time range
      const timeRangeDropdownTrigger = getByTestId('time-range-dropdown-trigger');
      act(() => {
        fireEvent.press(timeRangeDropdownTrigger);
      });

      await waitFor(() => {
        expect(getByTestId('time-range-option-24')).toBeTruthy();
      });

      const timeRange24Option = getByTestId('time-range-option-24');
      act(() => {
        fireEvent.press(timeRange24Option);
      });

      // Assert both filters are applied
      await waitFor(() => {
        expect(mockFetchReviews).toHaveBeenCalledTimes(3);
      });
      expect(mockFetchReviews).toHaveBeenLastCalledWith(4, 24);
    });
  });
});
