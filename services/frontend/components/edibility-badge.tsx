import { Badge } from "@/components/ui/badge";
import { getEdibilityLabel, getEdibilityColor } from "@/lib/types";

interface EdibilityBadgeProps {
  edibility: number;
}

export function EdibilityBadge({ edibility }: EdibilityBadgeProps) {
  const color = getEdibilityColor(edibility);
  const label = getEdibilityLabel(edibility);

  return <Badge variant={color}>{label}</Badge>;
}
