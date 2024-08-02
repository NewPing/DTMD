import preprocess from 'svelte-preprocess';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import dotenv from 'dotenv';

// Load environment variables
const mode = process.env.NODE_ENV || 'development';
dotenv.config({
    path: mode === 'production' ? '.env.production' : '.env'
});

let adapter;
if (mode === 'production') {
    adapter = require('@sveltejs/adapter-node').default();
} else {
    adapter = require('@sveltejs/adapter-auto').default();
}

/** @type {import('@sveltejs/kit').Config} */
const config = {
    extensions: ['.svelte'],
    preprocess: [vitePreprocess()],

    kit: {
        adapter,
        vite: {
            define: {
                'process.env.VITE_API_BASE_URL': JSON.stringify(process.env.VITE_API_BASE_URL)
            }
        }
    }
};

export default config;