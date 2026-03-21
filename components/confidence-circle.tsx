"use client";

import { cn } from "@/lib/utils";

interface ConfidenceCircleProps {
  confidence: number;
  size?: "sm" | "md" | "lg";
}

export function ConfidenceCircle({
  confidence,
  size = "md",
}: ConfidenceCircleProps) {
  const percentage = Math.round(confidence * 100);

  const sizeClasses = {
    sm: "h-12 w-12 text-xs",
    md: "h-20 w-20 text-sm",
    lg: "h-28 w-28 text-lg",
  };

  const strokeWidth = size === "sm" ? 3 : size === "md" ? 4 : 5;
  const radius = size === "sm" ? 20 : size === "md" ? 34 : 48;
  const circumference = 2 * Math.PI * radius;
  const strokeDashoffset = circumference - (confidence * circumference);

  const getColor = () => {
    if (percentage >= 80) return "text-safe";
    if (percentage >= 60) return "text-warning";
    return "text-danger";
  };

  return (
    <div
      className={cn(
        "relative flex items-center justify-center",
        sizeClasses[size]
      )}
    >
      <svg className="h-full w-full -rotate-90" viewBox="0 0 100 100">
        {/* Background circle */}
        <circle
          cx="50"
          cy="50"
          r={radius}
          fill="none"
          stroke="currentColor"
          strokeWidth={strokeWidth}
          className="text-muted"
        />
        {/* Progress circle */}
        <circle
          cx="50"
          cy="50"
          r={radius}
          fill="none"
          stroke="currentColor"
          strokeWidth={strokeWidth}
          strokeLinecap="round"
          strokeDasharray={circumference}
          strokeDashoffset={strokeDashoffset}
          className={cn("transition-all duration-500", getColor())}
        />
      </svg>
      <div className="absolute flex flex-col items-center">
        <span className={cn("font-bold", getColor())}>{percentage}%</span>
      </div>
    </div>
  );
}
