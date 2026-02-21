import type { Config } from 'tailwindcss';

export default {
    content: ['./src/**/*.{html,js,svelte,ts}'],

    theme: {
        extend: {
            colors: {
                primary: '#016B61',
                secondary: '#70B2B2',
                accent: '#9ECFD4',
                surface: '#E5E9C5',
            }
        }
    },

    plugins: []
} satisfies Config;
