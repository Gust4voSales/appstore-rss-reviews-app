import React from "react";
import { Star } from "lucide-react";
import { Review } from "../types/reviews";

interface ReviewItemProps {
  review: Review;
}

export const ReviewItem: React.FC<ReviewItemProps> = ({ review }) => {
  const renderStars = (rating: number) => {
    return (
      <div className="flex flex-row">
        {[1, 2, 3, 4, 5].map((index) => (
          <Star key={index} size={18} fill={index <= rating ? "#EAB308" : "transparent"} color="#EAB308" />
        ))}
      </div>
    );
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  return (
    <div className="mb-4 rounded-lg border border-gray-100 bg-white p-4 shadow-sm">
      <div className="mb-2 flex flex-row items-start justify-between gap-2">
        <h3 className="flex-1 text-lg font-bold text-gray-900 break-words">{review.title}</h3>
        {renderStars(review.rating)}
      </div>

      <p className="mb-2 text-gray-700">{review.content}</p>

      <div className="flex flex-row items-start justify-between gap-2">
        <p className="flex-1 text-sm text-gray-500">By {review.author}</p>
        <p className="text-sm text-gray-400">{formatDate(review.updatedAt)}</p>
      </div>
    </div>
  );
};
