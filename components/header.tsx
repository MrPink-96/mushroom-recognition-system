"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";

export function Header() {
  const pathname = usePathname();

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container mx-auto flex h-16 items-center justify-between px-4">
        <Link href="/" className="flex items-center gap-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
            <svg
              className="h-5 w-5 text-primary-foreground"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 3c-4 0-7 3-7 7 0 2 1 4 2 5v6h10v-6c1-1 2-3 2-5 0-4-3-7-7-7z"
              />
            </svg>
          </div>
          <span className="text-lg font-semibold text-foreground">
            MushroomID
          </span>
        </Link>

        <nav className="flex items-center gap-6">
          <Link
            href="/"
            className={cn(
              "text-sm font-medium transition-colors hover:text-primary",
              pathname === "/" ? "text-primary" : "text-muted-foreground"
            )}
          >
            Главная
          </Link>
          <Link
            href="/recognize"
            className={cn(
              "text-sm font-medium transition-colors hover:text-primary",
              pathname === "/recognize" ? "text-primary" : "text-muted-foreground"
            )}
          >
            Распознать
          </Link>
        </nav>
      </div>
    </header>
  );
}
