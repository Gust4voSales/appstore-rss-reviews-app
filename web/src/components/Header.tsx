import React from "react";
import { RefreshCw } from "lucide-react";
import { DEFAULT_TIME_RANGE, TimeRange } from "../types/reviews";
import { DefaultButton } from "./DefaultButton";

interface HeaderProps {
  appId?: string;
  reviewCount?: number;
  timeRange?: TimeRange;
  onRefresh: () => void;
  loading: boolean;
}

export const Header: React.FC<HeaderProps> = ({ appId, reviewCount, timeRange, onRefresh, loading }) => {
  return (
    <div className="border-b border-gray-200 bg-white p-4 shadow-sm">
      <div className="flex flex-row items-center justify-between">
        <div className="flex-1">
          <h1 className="text-2xl font-bold text-gray-900">App Reviews</h1>
          <p className="text-gray-600">
            App ID: {appId || "-"} • {reviewCount || "0"} reviews ({timeRange || DEFAULT_TIME_RANGE}
            h)
          </p>
        </div>

        <DefaultButton onClick={onRefresh} className="!rounded-lg !p-2" data-testid="refresh-button">
          <RefreshCw size={24} color="white" className={`size-6 ${loading ? "animate-spin" : ""}`} />
        </DefaultButton>
      </div>
    </div>
  );
};
