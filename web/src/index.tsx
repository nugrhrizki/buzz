/* @refresh reload */
import { ColorModeProvider, ColorModeScript } from "@kobalte/core";
import { QueryClient, QueryClientProvider } from "@tanstack/solid-query";
import { render } from "solid-js/web";

import { Toaster } from "@/components/ui/toast";

import "./root.css";
import { Routes } from "./routes";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      gcTime: 1000 * 60 * 60 * 24, // 24 hours
    },
  },
});

render(
  () => (
    <QueryClientProvider client={queryClient}>
      <ColorModeScript storageType="localStorage" />
      <ColorModeProvider>
        <Routes />
      </ColorModeProvider>
      <Toaster />
    </QueryClientProvider>
  ),
  document.getElementById("root")!,
);
