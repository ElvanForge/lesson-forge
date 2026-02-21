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
		client: {start:"_app/immutable/entry/start.CH-Bdzki.js",app:"_app/immutable/entry/app.CZ53k567.js",imports:["_app/immutable/entry/start.CH-Bdzki.js","_app/immutable/chunks/DbUhHK1O.js","_app/immutable/chunks/DGkkJ6Jd.js","_app/immutable/chunks/BsAtrqI-.js","_app/immutable/entry/app.CZ53k567.js","_app/immutable/chunks/DGkkJ6Jd.js","_app/immutable/chunks/859_MDgp.js","_app/immutable/chunks/C3rKXLWF.js","_app/immutable/chunks/BsAtrqI-.js","_app/immutable/chunks/BBOZHhJG.js","_app/immutable/chunks/CWKISfjT.js"],stylesheets:[],fonts:[],uses_env_dynamic_public:false},
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
