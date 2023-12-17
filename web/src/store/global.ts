import { createSignal } from "solid-js";

const [sidebarPin, setSidebarPin] = createSignal<boolean>(localStorage.getItem("sidebar-pin") === "true" || false);

function toggleSidebarPin() {
  setSidebarPin(!sidebarPin());
  localStorage.setItem("sidebar-pin", sidebarPin().toString());
}

export { sidebarPin, toggleSidebarPin };
