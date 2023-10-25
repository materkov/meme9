#!/usr/bin/env node

const cssModulesPlugin = require("esbuild-css-modules-plugin");
const esbuild = require("esbuild")

async function build() {
    await esbuild.build({
        logLevel: "debug",
        entryPoints: ["src/index.tsx"],
        bundle: true,
        minify: true,
        //outfile: "dist/bundle/bundle.js",
        outdir: "dist/bundle",
        treeShaking: true,
        plugins: [cssModulesPlugin({
            v2: true,
            localsConvention: 'camelCase',
        })],
        define: {
            'process.env.NODE_ENV': '"production"',
        }
    });
}
build();
