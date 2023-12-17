import type { ClassValue } from "clsx";
import { clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/**
 * Generates Avatar URL using DiceBear API
 * @param options
 * @returns Avatar URL in SVG format
 */

interface AvatarOptions {
  backgroundColors?: string[];
  flip?: boolean;
  name: string;
  style?:
    | "adventurer"
    | "adventurer-neutral"
    | "avataaars"
    | "avataaars-neutral"
    | "big-ears"
    | "big-ears-neutral"
    | "big-smile"
    | "bottts"
    | "bottts-neutral"
    | "croodles"
    | "croodles-neutral"
    | "fun-emoji"
    | "icons"
    | "identicon"
    | "initials"
    | "lorelei"
    | "lorelei-neutral"
    | "micah"
    | "miniavs"
    | "notionists"
    | "notionists-neutral"
    | "open-peeps"
    | "personas"
    | "pixel-art"
    | "pixel-art-neutral"
    | "rings"
    | "shapes"
    | "thumbs";
}

export function generateAvatarUrl(options: AvatarOptions) {
  const params = new URLSearchParams();
  params.set("seed", btoa(options.name));
  if (options.backgroundColors) {
    params.set("backgroundColor", options.backgroundColors.join(",").replace(/#/g, ""));
  }
  if (options.flip) {
    params.set("flip", options.flip.toString());
  }
  return `https://api.dicebear.com/7.x/${options.style || "notionists"}/svg?${params.toString()}`;
}

export function formatDate(date: Date | string | number) {
  return new Date(date).toLocaleDateString("en-US", {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    hour12: true,
  });
}
