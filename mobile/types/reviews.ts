export type Rating = 1 | 2 | 3 | 4 | 5;

export interface Review {
  id: string;
  title: string;
  content: string;
  author: string;
  rating: Rating;
  updatedAt: string;
}
