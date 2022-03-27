import { writable } from "svelte/store";
import type { Celebration } from "../types/Celebration";

export const celebrations = writable<Celebration[]>([]);
