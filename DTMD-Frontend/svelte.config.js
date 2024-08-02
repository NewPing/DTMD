import preprocess from 'svelte-preprocess';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import adapterAuto from '@sveltejs/adapter-auto';
import adapterNode from '@sveltejs/adapter-node';
import { defineConfig } from 'vite';

// Determine the environment
const isProduction = process.env.NODE_ENV === 'production';

/** @type {import('@sveltejs/kit').Config} */
const config = {
    extensions: ['.svelte'],
    preprocess: [vitePreprocess()],

    kit: {
        adapter: isProduction ? adapterNode() : adapterAuto(),
    },

    // Vite configuration should be directly here
    vite: defineConfig({
        define: {
            'import.meta.env.VITE_API_BASE_URL': JSON.stringify(process.env.VITE_API_BASE_URL)
        }
    })
};

export default config;