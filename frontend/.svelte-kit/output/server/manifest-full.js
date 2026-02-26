export const manifest = (() => {
function __memo(fn) {
	let value;
	return () => value ??= (value = fn());
}

return {
	appDir: "_app",
	appPath: "_app",
	assets: new Set(["robots.txt"]),
	mimeTypes: {".txt":"text/plain"},
	_: {
		client: {start:"_app/immutable/entry/start.DPuYTsHP.js",app:"_app/immutable/entry/app.DKiAiHbN.js",imports:["_app/immutable/entry/start.DPuYTsHP.js","_app/immutable/chunks/BULynGLC.js","_app/immutable/chunks/Cc-YLYwh.js","_app/immutable/chunks/DjQGscSg.js","_app/immutable/entry/app.DKiAiHbN.js","_app/immutable/chunks/Cc-YLYwh.js","_app/immutable/chunks/CapGnVIO.js","_app/immutable/chunks/CFHtMSnc.js","_app/immutable/chunks/DjQGscSg.js","_app/immutable/chunks/D_9wn8cK.js","_app/immutable/chunks/qnO9cl5L.js"],stylesheets:[],fonts:[],uses_env_dynamic_public:false},
		nodes: [
			__memo(() => import('./nodes/0.js')),
			__memo(() => import('./nodes/1.js')),
			__memo(() => import('./nodes/2.js')),
			__memo(() => import('./nodes/3.js')),
			__memo(() => import('./nodes/4.js'))
		],
		remotes: {
			
		},
		routes: [
			{
				id: "/",
				pattern: /^\/$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 2 },
				endpoint: null
			},
			{
				id: "/auth/callback",
				pattern: /^\/auth\/callback\/?$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 3 },
				endpoint: null
			},
			{
				id: "/login",
				pattern: /^\/login\/?$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 4 },
				endpoint: null
			}
		],
		prerendered_routes: new Set([]),
		matchers: async () => {
			
			return {  };
		},
		server_assets: {}
	}
}
})();
