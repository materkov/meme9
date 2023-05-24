#!/usr/bin/env node

const cssModulesPlugin = require("esbuild-css-modules-plugin");
const esbuild = require("esbuild")

async function watch() {
    let ctx = await esbuild.context({
        logLevel: "debug",
        entryPoints: ["src/index.tsx"],
        bundle: true,
        //outfile: "dist/bundle/bundle.js",
        outdir: "dist/bundle",
        plugins: [cssModulesPlugin({
            v2: true,
            localsConvention: 'camelCase',
        })],
    });
    await ctx.watch();
}
watch();
