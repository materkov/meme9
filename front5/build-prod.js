#!/usr/bin/env node

const cssModulesPlugin = require("esbuild-css-modules-plugin");

require("esbuild")
    .build({
        logLevel: "debug",
        entryPoints: ["src/index.tsx"],
        bundle: true,
        outfile: "dist/bundle.js",
        plugins: [cssModulesPlugin({
            localsConvention: 'camelCase',
        })],
    })
    .catch(() => process.exit(1));
