import { Badge } from "@/components/ui/badge";
import { getToxicityLabel, getToxicityColor } from "@/lib/types";

interface ToxicityBadgeProps {
  level: number;
}

export function ToxicityBadge({ level }: ToxicityBadgeProps) {
  const color = getToxicityColor(level);
  const label = getToxicityLabel(level);

  return <Badge variant={color}>{label}</Badge>;
}
