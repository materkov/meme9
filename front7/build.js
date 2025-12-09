const cssModulesPlugin = require("esbuild-css-modules-plugin");
const esbuild = require("esbuild");

const isProduction = process.env.NODE_ENV === "production";

esbuild.build({
  entryPoints: ["src/index.tsx"],
  bundle: true,
  outfile: "dist/index.js",
  format: "iife",
  globalName: "app",
  jsx: "automatic",
  minify: isProduction,
  sourcemap: !isProduction,
  treeShaking: true,
  target: ["es2020"],
  plugins: [
    cssModulesPlugin({
      v2: true,
      localsConvention: "camelCase",
    }),
  ],
}).then(() => {
  if (isProduction) {
    console.log("Production build completed!");
  } else {
    console.log("Build completed!");
  }
}).catch((error) => {
  console.error("Build failed:", error);
  process.exit(1);
});

