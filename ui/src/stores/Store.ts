import { writable } from "svelte/store";
import type { Celebration } from "../types/Celebration";
import type { CelebrationOverview } from "../types/CelebrationOverview";

export const celebrations = writable<Celebration[]>([]);
export const eventOverview = writable<CelebrationOverview>(null);
