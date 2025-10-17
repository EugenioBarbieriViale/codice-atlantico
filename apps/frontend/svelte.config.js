import adapter from "@sveltejs/adapter-node";
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";

const config = {
  preprocess: vitePreprocess(),
  kit: {
    adapter: adapter({
      out: "build",
    }),
    alias: {
      $components: "src/components",
      $lib: "src/lib",
    },
  },
};

export default config;
