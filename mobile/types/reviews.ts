export type Rating = 1 | 2 | 3 | 4 | 5;

export type TimeRange = 24 | 48 | 72 | 96;
export const DEFAULT_TIME_RANGE: TimeRange = 48;

export interface Review {
  id: string;
  title: string;
  content: string;
  author: string;
  rating: Rating;
  updatedAt: string;
}
