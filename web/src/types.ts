import { JSX } from "solid-js";

export interface Menu {
  name: string;
  isGroup?: boolean;
  show?: boolean;
  children?: Menu[];
  icon?: JSX.Element;
  href?: string;
}
