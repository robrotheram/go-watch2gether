import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import * as child from "child_process";

const manifestForPlugin = {
	registerType: "autoUpdate",
	injectRegister: 'auto',
	includeAssets: ["favicon.ico", "apple-touch-icon.png", "logo.svg"],
	manifest: {
		name: "Watch2Gether",
		short_name: "Watch2Gether",
		description: "",
		icons: [
			{
				src: "/android-chrome-192x192.png",
				sizes: "192x192",
				type: "image/png",
			},
			{
				src: "/android-chrome-512x512.png",
				sizes: "512x512",
				type: "image/png",
			},
			{
				src: "/apple-touch-icon.png",
				sizes: "180x180",
				type: "image/png",
				purpose: "apple touch icon",
			},
			{
				src: "/maskable_icon.png",
				sizes: "225x225",
				type: "image/png",
				purpose: "any maskable",
			},
		],
		workbox: {
			navigateFallback: "/",
      		navigateFallbackDenylist: [/auth/],
		},
		theme_color: "#171717",
		background_color: "#e8ebf2",
		display: "standalone",
		scope: "/",
		start_url: "/",
		orientation: "portrait",
	},
};



const commitHash = child.execSync("git rev-parse --short HEAD").toString()
const gitTag = child.execSync("git describe --tags --abbrev=0").toString()
// https://vitejs.dev/config/
export default defineConfig({
	define: {
		'import.meta.env.VITE_APP_VERSION': JSON.stringify(gitTag),
		'import.meta.env.VITE_APP_COMMIT': JSON.stringify(commitHash)
	},
	plugins: [react()],
})
