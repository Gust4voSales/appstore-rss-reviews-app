import React, { useState } from "react";
import { Star, Filter, Clock, ChevronDown, ChevronUp, Check } from "lucide-react";
import { Rating, TimeRange } from "../types/reviews";

const validRatings: Rating[] = [1, 2, 3, 4, 5];
const validTimeRanges: TimeRange[] = [24, 48, 72, 96];

interface ReviewsFilterProps {
  selectedRating?: Rating;
  onSelectRating: (rating?: Rating) => void;
  selectedTimeRange: TimeRange;
  onSelectTimeRange: (timeRange: TimeRange) => void;
}

export const ReviewsFilter = ({
  selectedRating,
  onSelectRating,
  selectedTimeRange,
  onSelectTimeRange,
}: ReviewsFilterProps) => {
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const [hoveredRating, setHoveredRating] = useState<Rating | undefined>(undefined);

  const handleTimeRangeSelect = (timeRange: TimeRange) => {
    onSelectTimeRange(timeRange);
    setIsDropdownOpen(false);
  };

  return (
    <div className="relative mb-2 items-center gap-3 border-b border-gray-200 bg-white px-4 py-3">
      <Filter size={24} color="gray" className="absolute left-4 top-1/2 transform -translate-y-1/2" />

      {/* Rating Filter */}
      <div className="flex flex-row items-center justify-center mb-3">
        {validRatings.map((rating) => (
          <button
            key={rating}
            data-testid={`star-button-${rating}`}
            onClick={() => onSelectRating(selectedRating === rating ? undefined : rating)}
            onMouseEnter={() => setHoveredRating(rating)}
            onMouseLeave={() => setHoveredRating(undefined)}
            className="cursor-pointer rounded-full px-2 py-1 transition-colors"
          >
            <Star
              size={24}
              fill={
                (selectedRating && selectedRating >= rating) || (hoveredRating && hoveredRating >= rating)
                  ? "#EAB308"
                  : "transparent"
              }
              color={"#EAB308"}
            />
          </button>
        ))}
      </div>

      {/* Time Range Dropdown */}
      <div className="relative flex justify-center">
        <button
          data-testid="time-range-dropdown-trigger"
          onClick={() => setIsDropdownOpen(!isDropdownOpen)}
          className="flex flex-row items-center rounded-lg border border-gray-300 bg-white px-4 py-2 shadow-sm hover:bg-gray-50 transition-colors"
        >
          <Clock size={16} color="gray" className="mr-2" />
          <span className="mr-2 text-sm font-medium text-gray-700">Last {selectedTimeRange}h</span>
          {isDropdownOpen ? <ChevronUp size={20} color="gray" /> : <ChevronDown size={20} color="gray" />}
        </button>

        {isDropdownOpen && (
          <>
            <div className="fixed inset-0 z-10 bg-black/20" onClick={() => setIsDropdownOpen(false)} />
            <div className="absolute top-full mt-1 z-20 w-32 rounded-lg border border-gray-200 bg-white shadow-lg">
              {validTimeRanges.map((timeRange) => (
                <button
                  key={timeRange}
                  data-testid={`time-range-option-${timeRange}`}
                  onClick={() => handleTimeRangeSelect(timeRange)}
                  className={`flex flex-row items-center justify-between w-full px-4 py-3 text-left hover:bg-gray-50 transition-colors ${
                    selectedTimeRange === timeRange ? "bg-blue-50" : ""
                  }`}
                >
                  <span
                    className={`text-base ${
                      selectedTimeRange === timeRange ? "font-semibold text-blue-600" : "text-gray-700"
                    }`}
                  >
                    Last {timeRange}h
                  </span>
                  {selectedTimeRange === timeRange && <Check size={20} color="#2563EB" />}
                </button>
              ))}
            </div>
          </>
        )}
      </div>
    </div>
  );
};
