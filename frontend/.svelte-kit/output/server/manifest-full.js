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
		client: {start:"_app/immutable/entry/start.D8HKY465.js",app:"_app/immutable/entry/app.DKhofKNd.js",imports:["_app/immutable/entry/start.D8HKY465.js","_app/immutable/chunks/CPGch2rv.js","_app/immutable/chunks/CZ8iVQTt.js","_app/immutable/chunks/DfM9hz9f.js","_app/immutable/entry/app.DKhofKNd.js","_app/immutable/chunks/CZ8iVQTt.js","_app/immutable/chunks/MdFL1Qgq.js","_app/immutable/chunks/SldK7h0U.js","_app/immutable/chunks/DfM9hz9f.js","_app/immutable/chunks/DwbYP2Bb.js","_app/immutable/chunks/iW0WkGsQ.js"],stylesheets:[],fonts:[],uses_env_dynamic_public:false},
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
